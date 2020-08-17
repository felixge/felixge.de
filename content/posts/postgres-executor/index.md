---
title: "Writing Your Own PostgreSQL Query Executor"
date: 2020-08-16T10:32:14+02:00
draft: true
---

<script src="./main.js"></script>

<blockquote>
<p> Your father's relational database. This is the weapon of a data engineer. Not as clumsy or random as a document store. An elegant weapon ... for a more civilized age.</p>
<footer><cite><a href="https://en.wikipedia.org/wiki/Michael_Stonebraker">Obi-Wan Stonebreaker</a></cite></footer>
</blockquote>

Would you like to improve your intution about SQL performance? Have you ever wondered how queries are actually executed? Are you learning SQL and would like to see how it translates to the code you normally write? Or do you think SQL is an arcane legacy language that should be avoided in favor of NoSQL databases?

Then you have come to the right place! Let's go on a little adventure and write a toy query executor similar to the one found in PostgreSQL and other relational databases. We'll use JavaScript due to its popularity, but try to keep things easy to follow even if prefer Go, Ruby, Python, Java or similar. Once we're done, you'll hopefully understand and appreciate SQL a lot more than before.

Let's start with a simple `JOIN` query between two tables:

```sql
CREATE TABLE users (user_id, email_id) AS VALUES
	(1, 'user1@example.org'),
	(2, 'user2@example.org'),
	(3, 'user3@example.org');

CREATE TABLE orders (order_id, user_id, amount) AS VALUES
	(1, 1, 5),
	(2, 1, 13),
	(3, 2, 7);
```
```sql
SELECT users.email, orders.amount
FROM users
JOIN orders USING (user_id);
```

Running the above in PostgeSQL should give us the following results:

<table>
	<tr>
		<th>email</th>
		<th>amount</th>
	</tr>
	<tr>
		<td>user1@example.org</td>
		<td>5</td>
	</tr>
	<tr>
		<td>user1@example.org</td>
		<td>13</td>
	</tr>
	<tr>
		<td>user2@example.org</td>
		<td>7</td>
	</tr>
</table>

How did PostgreSQL do that? Let's assume we don't know, how could we do it in JS? Let's start by defining our data as arrays of objects:

```js
var users = [
  {user_id: 1, email: 'user1@example.org'},
  {user_id: 2, email: 'user2@example.org'},
  {user_id: 3, email: 'user3@example.org'},
];

var orders = [
  {order_id: 1, user_id: 1, amount: 5},
  {order_id: 2, user_id: 1, amount: 13},
  {order_id: 3, user_id: 2, amount: 7},
];
```

And then let's write down the simplest code that could produce the same results:

```js
var results = [];
for (let order of orders) {
  for (let user of users) {
    if (order.user_id === user.user_id) {
      results.push({email: user.email, amount: order.amount});
    }
  }
}
```

PostgreSQL calls the above a `NestedLoop` JOIN, and it's indeed one of the ways it might execute a query like this. However, is this an efficient solution? Having to iterate over every user for every order seems inefficient. Let's try to fix that:

```js
var results = [];
for (let order of orders) {
  for (let user of users) {
    if (order.user_id === user.user_id) {
      results.push({email: user.email, amount: order.amount});
    }
  }
}
```

Outline:

- Node Types:
  - Nested Loop
  - Hash Join
  - Sequential Scan
  - Index Scan
  - Limit
  - Sort
