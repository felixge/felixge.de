#!/usr/bin/env node
var rs = require('robotskirt');
var input = '';
process.stdin.setEncoding('utf8');
process.stdin
	.on('data', function(chunk) {
		input += chunk
	})
	.on('end', function() {
		var renderer = new rs.HtmlRenderer();
		var parser = new rs.Markdown(renderer);
		process.stdout.write(parser.render(input))
	})
	.resume();
