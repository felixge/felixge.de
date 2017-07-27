---
layout: post
title: "Implementing State Machines in PostgreSQL"
---

[Finite-state machine](https://en.wikipedia.org/wiki/Finite-state_machine)
(FSM) is a wonderful model of computation that has many practical
applications. One of my favorites is turning business logic into simple FSMs.
For example consider the following requirements for an order management system:

* Orders must not ship before they are paid.
* Orders can be canceled, but only if they haven't shipped yet.
* An order that has already been paid, needs to be refunded after canceling.

Informally, turning such a list of requirements into an FSM is simply the
process of coming up with: a set of states an order can be in, a set of events,
a transition function that maps the combination of each state and event to a
new state, and an initial state.

In practice this is usually done by drawing up a directed graph that shows how
the different states are connected to each other via events:

<img src="{% asset_path orders.svg %}" alt="Finite-state machine for an order management system in PostgreSQL" />

Besides guiding our implementation, these graphs can be very useful for
discussions and analysis. For example you can immediately see how there is no
way for a customer to return an order after it has shipped. Additionally the
process of naming the involved states and events automatically creates a
precise language that can be adopted by all stakeholders.

You might have noticed that the graph above does not show all possible
combinations of all states and all events. The reason is that the missing
combinations implicitly lead to an error state. The [complete FSM]({%
asset_path orders-full.svg %}) can be shown as a graph as well, but I don't
think it's very useful in practice.

Anyway, as promised in the title of this article we're going to implement this
FSM in PostgreSQL. This may not be everybody's cup of tea, but you might like
how this approach us gives advanced analytical powers as a free by-product.
Embedding this kind of logic into the database can also help protect against
race conditions, but this will perhaps be the topic of a future post [^race].

Let's begin by creating an `order_events` table which keeps track of all events
for a given `order_id`.

```sql
CREATE TABLE order_events (
  id serial PRIMARY KEY,
  order_id int NOT NULL,
  event text NOT NULL,
  time timestamp DEFAULT now() NOT NULL
);
```

Next, let's implement the transition function, which is the heart of every
FSM:

```sql
CREATE FUNCTION order_events_transition(state text, event text) RETURNS text
LANGUAGE sql AS
$$
  SELECT CASE state
    WHEN 'start' THEN
      CASE event
        WHEN 'create' THEN 'awaiting_payment'
        ELSE 'error'
      END
    WHEN 'awaiting_payment' THEN
      CASE event
        WHEN 'pay' THEN 'awaiting_shipment'
        WHEN 'cancel' THEN 'canceled'
        ELSE 'error'
      END
    WHEN 'awaiting_shipment' THEN
      CASE event
        WHEN 'cancel' THEN 'awaiting_refund'
        WHEN 'ship' THEN 'shipped'
        ELSE 'error'
      END
    WHEN 'awaiting_refund' THEN
      CASE event
        WHEN 'refund' THEN 'canceled'
        ELSE 'error'
      END
    ELSE 'error'
  END
$$;
```

And before we proceed, let's test with a few examples to make sure the function
is working:

```sql
SELECT state, event, order_events_transition(state, event)
FROM (VALUES
  ('start', 'create'),
  ('awaiting_payment', 'pay'),
  ('awaiting_payment', 'cancel'),
  ('awaiting_payment', 'ship')
) AS examples(state, event);
```
```
      state       | event  | order_events_transition
------------------+--------+-------------------------
 start            | create | awaiting_payment
 awaiting_payment | pay    | awaiting_shipment
 awaiting_payment | cancel | canceled
 awaiting_payment | ship   | error
```

The above looks correct, but it's not immediately clear how we could use this
function to enforce our FSM on all rows for the same `order_id` in the
`order_events` table.

The first thing we need is a good way to take a list of events and call our
transition function on them recursively to determine the resulting state. There
are multiple ways to accomplish this in PostgreSQL, but perhaps the most
elegant is a [user-defined aggregate](https://www.postgresql.org/docs/current/static/xaggr.html).

In a nutshell, a user defined aggregate is a function that has an internal
value (state) that is updated for each input value that is being passed into it
using a state transition function. This means it fits our FSM model like a
glove:

```sql
CREATE AGGREGATE order_events_fsm(text) (
  SFUNC = order_events_transition,
  STYPE = text,
  INITCOND = 'start'
);
```

The above defines a new aggregate called `order_events_fsm` which takes a
`text` input (one of our events) and calls the `order_events_transition` state
transition function (`FSUNC`) for each input along with the current state
(`STYPE`) which is also of type `text`. The initial state is `start`
(`INITCOND`).

A quick test shows that it works as expected:

```sql
SELECT order_events_fsm(event ORDER BY id)
FROM (VALUES
  (1, 'create'),
  (2, 'pay'),
  (3, 'cancel')
) examples(id, event);
```
```
 order_events_fsm
------------------
 awaiting_refund
```

Now let's use our `order_events_fsm` to create a `BEFORE INSERT` trigger for
our `order_events` table that makes sure all events of a given `order_id`
are valid and don't lead to an error state. We do this using a simple `plpgsql`
function that executes our `order_events_fsm` against all existing events
of the current `order_id` plus the new event. If the final state is `error`,
we raise an exception which causes the current transaction to roll back. It
looks like this:


```sql
CREATE FUNCTION order_events_tigger_func() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  new_state text;
BEGIN
  SELECT order_events_fsm(event ORDER BY id)
  FROM (
    SELECT id, event FROM order_events WHERE order_id = new.order_id
    UNION
    SELECT new.id, new.event
  ) s
  INTO new_state;
  
  IF new_state = 'error' THEN
    RAISE EXCEPTION 'invalid event';
  END IF;

  RETURN new;
END
$$;

CREATE TRIGGER order_events_trigger BEFORE INSERT ON order_events
FOR EACH ROW EXECUTE PROCEDURE order_events_tigger_func();
```

We can verify that this works by inserting a valid sequence of order events:

```sql
INSERT INTO order_events (order_id, event) VALUES
  (1, 'create'),
  (1, 'pay'),
  (1, 'ship');
```
```
INSERT 0 3
```

As well as an invalid sequence of events:

```sql
INSERT INTO order_events (order_id, event) VALUES
  (2, 'create'),
  (2, 'ship');
```
```
psql:orders.sql:95: ERROR:  invalid event
CONTEXT:  PL/pgSQL function order_events_tigger_func() line 14 at RAISE
```

As expected, only the first series of inserts made it into our table:

```sql
SELECT id, order_id, event FROM order_events;
```
```
 id | order_id | event
----+----------+--------
  1 |        1 | create
  2 |        1 | pay
  3 |        1 | ship
```

If you're still on the fence about embedding this kind of logic into your
database, let's see how our approach gives us advanced analytical powers as a
free by-product. Let's consider a new data set of 3 orders:


```sql
TRUNCATE order_events;
INSERT INTO order_events (order_id, event, time) VALUES
  (1, 'create', '2017-07-23 00:00:00'),
  (1, 'pay', '2017-07-23 12:00:00'),
  (1, 'ship', '2017-07-24 00:00:00'),

  (2, 'create', '2017-07-23 00:00:00'),
  (2, 'cancel', '2017-07-24 00:00:00'),

  (3, 'create', '2017-07-23 00:00:00'),
  (3, 'pay', '2017-07-24 00:00:00'),
  (3, 'cancel', '2017-07-25 00:00:00'),
  (3, 'refund', '2017-07-26 00:00:00');
```

Using our `order_events_fsm` as a [window
function](https://www.postgresql.org/docs/current/static/tutorial-window.html), we
can easily get the state history of a given order:

```sql
SELECT time, order_events_fsm(event) OVER (ORDER BY id)
FROM order_events
WHERE order_id = 3;
```
```
        time         | order_events_fsm
---------------------+-------------------
 2017-07-23 00:00:00 | awaiting_payment
 2017-07-24 00:00:00 | awaiting_shipment
 2017-07-25 00:00:00 | awaiting_refund
 2017-07-26 00:00:00 | canceled
```

But we can go even further and apply our state machines to multiple orders,
e.g. by using the
[generate\_series](https://www.postgresql.org/docs/current/static/functions-srf.html)
function and a
[Lateral](https://www.postgresql.org/docs/current/static/queries-table-expressions.html#QUERIES-LATERAL)
sub-query to break down the number of orders per state for each day of a given
date range:

```sql
SELECT date::date, state, count(1)
FROM
  generate_series('2017-07-23'::date, '2017-07-26', '1 day') date,
  LATERAL (
    SELECT order_id, order_events_fsm(event ORDER BY id) AS state
    FROM order_events
    WHERE time < date + '1 day'::interval
    GROUP BY 1
  ) orders
GROUP BY 1, 2
ORDER BY 1, 2;
```
```
    date    |       state       | count
------------+-------------------+-------
 2017-07-23 | awaiting_payment  |     2
 2017-07-23 | awaiting_shipment |     1
 2017-07-24 | awaiting_shipment |     1
 2017-07-24 | canceled          |     1
 2017-07-24 | shipped           |     1
 2017-07-25 | awaiting_refund   |     1
 2017-07-25 | canceled          |     1
 2017-07-25 | shipped           |     1
 2017-07-26 | canceled          |     2
 2017-07-26 | shipped           |     1
```

There you have it, a FSM implemented as a user defined aggregate in PostgreSQL
providing data integrity and advanced analytics.

That being said, your milage may vary, and embedding your business logic into
your database is always a tradeoff. But if you want some reassurance: I've had
great success in applying this approach in combination with [eager
materialization](https://hashrocket.com/blog/posts/materialized-view-strategies-using-postgresql#eager-materialized-view)
to implement a realtime analytics dashboard for an application with over a
billion rows.

Anyway, I'm really looking forward to feedback on this, and am more than happy
to answer any questions, so please comment.

<small>Thanks to Thorsten Ball and Johannes Boyne for reviewing.</small>

[^race]: The code in this article immune to concurrency anomalies when using the `SERIALIZABLE` transaction isolation level. Alternative you could modify the trigger to aquire an exclusive lock on the `order_events` table. But as mentioned, this topic deserves a separate post.
