---
title: "History of Go's Execution Tracer"
date: 2023-01-13T21:17:52.000Z
slug: "history-of-gos-execution-tracer"
aliases:
  - /history-of-gos-execution-tracer/
---

Today I spent some time studying the early history of Go's execution tracer. Here is a very brief summary of what I learned.

The design and initial implementation was done by Dmitry Vyukov as outlined in his [Go Execution Tracer (Oct 2014)](https://docs.google.com/document/u/1/d/1FP5apqzBgr7ahCCgFO-yoVhk4YZrNIDNf9RybngBc14/pub) design document. This was followed by the [first upstream patch](https://go-review.googlesource.com/c/go/+/1450) being merged the following January, an announcement at GopherCon in July (Fig 1), and the release in [Go 1.5](https://go.dev/doc/devel/release#go1.5) in August 2015.

[GopherCon 2015: Dmitry Vyukov - Go Dynamic Tools](https://www.youtube.com/watch?v=a9xrxRsIbSU&t=1287s)

Fig. 1 *Dmitry briefly announcing the Go Execution Tracer at GopherCon 2015.*

Since then a few [new events](https://github.com/golang/go/blob/16cec4e7a0ed4b1e5cf67a07a1bf24bdf2f6b04c/src/internal/trace/parser.go#L1129-L1141) have been added, along with other improvements, but overall the tracer seems to have retained most of its original design and functionality.
