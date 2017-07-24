\set ECHO queries
\set ON_ERROR_STOP on

DROP SCHEMA IF EXISTS oms CASCADE;
CREATE SCHEMA oms;
SET search_path=oms;

CREATE TABLE order_events (
  id serial PRIMARY KEY,
  order_id int NOT NULL,
  event text NOT NULL,
  time timestamp DEFAULT now() NOT NULL
);

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

SELECT state, event, order_events_transition(state, event)
FROM (VALUES
  ('start', 'create'),
  ('awaiting_payment', 'pay'),
  ('awaiting_payment', 'cancel'),
  ('awaiting_payment', 'ship')
) AS examples(state, event);

CREATE AGGREGATE order_events_fsm(text) (
  SFUNC = order_events_transition,
  STYPE = text,
  INITCOND = 'start'
);


SELECT order_events_fsm(event ORDER BY id)
FROM (VALUES
  (1, 'create'),
  (2, 'pay'),
  (3, 'cancel')
) examples(id, event);

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

CREATE TRIGGER order_events_trigger BEFORE INSERT ON order_events FOR EACH ROW EXECUTE PROCEDURE order_events_tigger_func();

INSERT INTO order_events (order_id, event) VALUES
  (1, 'create'),
  (1, 'pay'),
  (1, 'ship');

\set ON_ERROR_STOP off
INSERT INTO order_events (order_id, event) VALUES
  (2, 'create'),
  (2, 'ship');
\set ON_ERROR_STOP on

SELECT id, order_id, event FROM order_events;

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

SELECT time, order_events_fsm(event) OVER (ORDER BY id)
FROM order_events
WHERE order_id = 3;


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
