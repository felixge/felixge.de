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

function example1() {
  var results = [];
  for (let order of orders) {
    for (let user of users) {
      if (order.user_id === user.user_id) {
        results.push({email: user.email, amount: order.amount});
      }
    }
  }
  return results;
}

console.log(example1);
