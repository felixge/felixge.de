---
layout: post
title: Node.js vs Go
date: 2013-06-13T16:00:00+01:00
updated: 2013-06-13T10:03:00+01:00
---

If you know me, you're probably aware that I'm heavily invested in node.js. I
was #24 on the node.js mailing list, and contributed over 100 patches to the
early versions of the core. Among other things, I've used node to build a
[product](https://transloadit.com/), created a range of [npm
modules](https://npmjs.org/~felixge), programmed [quad
copters](http://nodecopter.com/), and engaged in [JavaScript vs
C](https://github.com/felixge/faster-than-c) performance battles on Github.

Given this background, many people were surprised by my recent interest in the
[Go](http://golang.org/) (aka *Golang* for better googleability) programming
language. Now that I've been playing with Go for a little over 6 months, I
think I've had enough time to put my experience into words and attempt a
comparison with node.js.

But before we get started ... a quick reality check.

## Use the right tool for the job

If you are paid by somebody to write code for them, and the problem you face
could realistically be solved by employing an established framework for Ruby,
PHP or Python using a relational database, please do that. There are thousands
of man years of work and experience you can built on, and you should not ignore
these options lightly.

Skip to the bottom of this article if you're curious about the problem domains
I'd recommend Node.js and Go for.

## A quick introduction to Go

For the purpose of this article I'm assuming that you're already familar with
Node.js, but are relatively new to Go.

Go is a new programming language that aims to combine the efficiency of a
statically compiled language with the ease of programming of a dynamic
language. A simple web server that responds to "Hello World" for every request
looks like this:

{% highlight go %}
package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World")
	})
	http.ListenAndServe(":8080", nil)
}
{% endhighlight %}


## Maturity


## When to use Node.js or Go

* defer
* panic / error handling
* goroutines
* channels
* go fmt
* godoc
