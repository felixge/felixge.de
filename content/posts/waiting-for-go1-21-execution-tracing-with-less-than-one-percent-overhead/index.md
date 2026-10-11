---
title: "Waiting for go1.21: Execution Tracing with < 1% Overhead"
date: 2023-06-30T19:15:32.000Z
slug: "waiting-for-go1-21-execution-tracing-with-less-than-one-percent-overhead"
cover: cover.png
aliases:
  - /waiting-for-go1-21-execution-tracing-with-less-than-one-percent-overhead/
---

Go’s execution tracer (aka `go tool trace`) has suffered from high overhead since its [inception in 2014](https://docs.google.com/document/u/1/d/1FP5apqzBgr7ahCCgFO-yoVhk4YZrNIDNf9RybngBc14/pub#h.3uippe21mzdj). Historically this has forced potential users to worry about up to 20% of CPU overhead when turning it on. Due to this, it's mostly been used in test environments or tricky situations rather than gaining adoption as a continuous profiling signal in production.

However, in early 2023 my colleague Nick Ripley and I started to become convinced that we could change this by implementing frame pointer unwinding. You can read more about the details in the blog post below:

[Reducing Go Execution Tracer Overhead With Frame Pointer Unwinding](/2023/01/31/reducing-gos-execution-tracer-overhead-with-frame-pointer-unwinding/)

Since then we’ve been actively collaborating with the Go runtime team at Google, in particular [Michael Knyszek](https://go-review.googlesource.com/q/owner:mknyszek@google.com), [Michael Pratt](https://go-review.googlesource.com/q/owner:mpratt@google.com) and [Cherry Mui](https://go-review.googlesource.com/q/owner:cherryyz@google.com) in order to make this a reality. And at this point it looks like we’ve succeeded: The go1.21 release coming out in ~August will feature a new version of the tracer which should provide **less than 1% overhead** for most applications 🥳 (see upcoming [release notes](https://tip.golang.org/doc/go1.21#runtime/trace)).

Below is a list of directly or indirectly related patches that contributed to this effort:

- 2023-03-13 [runtime: remove NOFRAME from asmcgocall, systemstack and mcall](https://go-review.googlesource.com/c/go/+/472195) (Quim Muntal)
- 2023-03-30 [runtime: fix frame pointer loop on amd64](https://go-review.googlesource.com/c/go/+/476235) (Felix Geisendörfer)
- 2023-03-30 [runtime/trace: enable frame pointer unwinding on amd64](https://go-review.googlesource.com/c/go/+/463835) (Felix Geisendörfer)
- 2023-03-30 [runtime/trace: use regular unwinding for cgo callbacks](https://go-review.googlesource.com/c/go/+/474916) (Felix Geisendörfer)
- 2023-03-30 [runtime: add a benchmark of FPCallers](https://go-review.googlesource.com/c/go/+/475795) (Felix Geisendörfer)
- 2023-04-04 [runtime/trace: Fix TestTraceSymbolize on solaris and illumos](https://go-review.googlesource.com/c/go/+/481336) (Felix Geisendörfer)
- 2023-04-05 [runtime: save frame pointer to the stack in signal handlers for arm64](https://go-review.googlesource.com/c/go/+/481635) (Nick Ripley)
- 2023-04-19 [runtime: adjust frame pointer on stack copy on ARM64](https://go-review.googlesource.com/c/go/+/241158) (Cherry Mui)
- 2023-04-20 [runtime: turn frame pointer unwinding for tracing off by default](https://go-review.googlesource.com/c/go/+/486382) (Nick Ripley)
- 2023-04-28 [runtime/trace: avoid frame pointer unwinding for events during cgocallbackg](https://go-review.googlesource.com/c/go/+/488755) (Nick Ripley)
- 2023-04-30 [runtime: fix systemstack frame pointer adjustment](https://go-review.googlesource.com/c/go/+/489015) (Felix Geisendörfer)
- 2023-05-02 [runtime/trace: enable frame pointer unwinding by default for amd64 and arm64](https://go-review.googlesource.com/c/go/+/490815) (Nick Ripley)
- 2023-05-03 [runtime: add test for systemstack frame pointer adjustment](https://go-review.googlesource.com/c/go/+/489115) (Felix Geisendörfer)
- 2023-05-03 [runtime: remove systemstack logic from adjustframe](https://go-review.googlesource.com/c/go/+/489116) (Felix Geisendörfer)
- 2023-05-12 [runtime: remove unused skip arg from fpTracebackPCs](https://go-review.googlesource.com/c/go/+/489095) (Felix Geisendörfer)
- 2023-05-22 [runtime: rename getcallerfp to getfp](https://go-review.googlesource.com/c/go/+/494857) (Felix Geisendörfer)

So if you're into Go performance and observability, this improvement is definitely something to look forward to.

And this is just the beginning:

1.  Michael Knyszek [has proposed](https://github.com/golang/go/issues/60773) a lot of exciting improvements that will further improve execution tracing in upcoming releases.
2.  Dominik Honnef's [gotraceui](https://github.com/dominikh/gotraceui) viewer keeps improving.
3.  We've also been busy with integrating the execution tracer into our profiling product at Datadog. We'll share more information about this soon.
