#!/usr/bin/env node
var less = require('less');
var input = '';
process.stdin.setEncoding('utf8');
process.stdin
  .on('data', function(chunk) {
    input += chunk
  })
  .on('end', function() {
    var options = {};
    less.render(input, options, function(err, output) {
      if (err) {
        throw err;
      }

      process.stdout.write(output)
    });
  })
  .resume();
