---
title: "Go Execution Trace Size Breakdown"
date: 2023-01-12T22:34:00.000Z
slug: "go-execution-trace-size-breakdown"
cover: cover.png
aliases:
  - /go-execution-trace-size-breakdown/
---

I'm currently trying to get a better understanding of the data produced by Go's execution tracer (aka [runtime/trace](https://pkg.go.dev/runtime/trace)). One thing I'm interested in is the breakdown of trace sizes.

For this purpose I've just added a new `breakdown` command ([docs](https://github.com/felixge/traceutils#breakdown)) to my [traceutils](https://github.com/felixge/traceutils) project. I haven't had a chance to analyze a lot of different traces yet, but it was already very interesting to see the results for a big 280MiB trace that I had laying around.

I didn't know what to expect, but seeing almost 22% of the trace being made up of [EventGoSysCall](https://github.com/golang/go/blob/245e95dfabd77f337373bf2d6bb47cd353ad8d74/src/runtime/trace.go#L52) events was definitely surprising. Perhaps its indicative of poorly buffered I/O operations?

Another trace I looked at spends only 1% of its size on this event type. It will be interesting to see what patterns emerge when looking at more data.

    $ traceutils breakdown bytes big.trace
    +------------------------+----------+---------+
    |       EVENT TYPE       |  BYTES   |    %    |
    +------------------------+----------+---------+
    | EventGoSysCall         | 63.1 MB  | 21.84%  |
    | EventGoCreate          | 60.2 MB  | 20.86%  |
    | EventGoStartLocal      | 44.6 MB  | 15.44%  |
    | EventGoUnblock         | 38.5 MB  | 13.34%  |
    | EventGoStart           | 19.8 MB  | 6.85%   |
    | EventGoBlockRecv       | 14.7 MB  | 5.09%   |
    | EventHeapAlloc         | 14.1 MB  | 4.90%   |
    | EventGoEnd             | 11.1 MB  | 3.83%   |
    | EventGoSched           | 8.7 MB   | 3.00%   |
    | EventGoUnblockLocal    | 4.8 MB   | 1.67%   |
    | EventGoBlockSync       | 4.7 MB   | 1.63%   |
    | EventGoSysExit         | 1.1 MB   | 0.37%   |
    | EventStack             | 783.3 kB | 0.27%   |
    | EventCPUSample         | 497.4 kB | 0.17%   |
    | EventProcStart         | 445.9 kB | 0.15%   |
    | EventGoBlockSend       | 291.5 kB | 0.10%   |
    | EventGoSysBlock        | 284.1 kB | 0.10%   |
    | EventGCSweepDone       | 275.7 kB | 0.10%   |
    | EventProcStop          | 218.9 kB | 0.08%   |
    | EventGCSweepStart      | 202.3 kB | 0.07%   |
    | EventGoBlockSelect     | 184.8 kB | 0.06%   |
    | EventString            | 72.5 kB  | 0.03%   |
    | EventBatch             | 36.2 kB  | 0.01%   |
    | EventGoPreempt         | 34.4 kB  | 0.01%   |
    | EventGoBlock           | 28.1 kB  | 0.01%   |
    | EventGCMarkAssistStart | 27.3 kB  | 0.01%   |
    | EventGCMarkAssistDone  | 21.5 kB  | 0.01%   |
    | EventGoStartLabel      | 14.8 kB  | 0.01%   |
    | EventGoWaiting         | 9.4 kB   | 0.00%   |
    | EventGoBlockNet        | 2.1 kB   | 0.00%   |
    | EventGoSleep           | 364 B    | 0.00%   |
    | EventGomaxprocs        | 49 B     | 0.00%   |
    | EventGCSTWStart        | 39 B     | 0.00%   |
    | EventHeapGoal          | 28 B     | 0.00%   |
    | EventGCSTWDone         | 24 B     | 0.00%   |
    | EventGCStart           | 20 B     | 0.00%   |
    | EventGoInSyscall       | 16 B     | 0.00%   |
    | EventGCDone            | 8 B      | 0.00%   |
    | EventGoBlockCond       | 5 B      | 0.00%   |
    | EventFrequency         | 5 B      | 0.00%   |
    +------------------------+----------+---------+
    |         TOTAL          | 288.7 MB | 100.00% |
    +------------------------+----------+---------+
