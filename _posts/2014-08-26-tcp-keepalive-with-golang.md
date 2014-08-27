---
layout: post
title: "Using TCP keepalive with Go"
date: 2014-08-26T15:40:00+01:00
updated: 2014-08-26T15:40:00+01:00
---

If you have ever written some TCP socket code, you may have wondered: "What
will happen to my connection if the network cable is unplugged or the remote
machine crashes?".

The short answer is: nothing. The remote end of the connection won't be able to
send a FIN packet, and the local OS will not detect that the connection is
lost. So it's up to you as the developer to address this scenario.

In Go you have several methods available to you that can help with this.
Perhaps the first one to consider is the `SetReadDeadline` method of the
[net.Conn][] interface. Assuming that your connection is expected to receive
data at a regular interval, you can simply treat a timed out read as equivalent
to an `io.EOF` error and `Close` the connection. Many existing TCP protocols
support this way of error handling by defining some sort of heartbeat mechanism
that requires each endpoint to send PING/PONG probes at a regular interval in
order to detect both networking problems, as well as service
health<sup>1</sup>. Additionally such heartbeats may also help dealing with
proxy servers that look for network activity to determine the health of a
connection.

So if your protocol supports heartbeats, or you have the ability to add
heartbeats to your own protocol, that should be your first choice for
addressing the unplugged network cable scenario.

However, what happens if you have no control over the protocol, and heartbeats
are not supported?

Now it's time to learn about TCP keepalive and how to use it with Go. TCP
keepalive is defined in [RFC 1122][], and is not part of the TCP specification
itself. It can be enabled for individual connections, but MUST be turned off by
default. Enabling it will cause the network stack to probe the health of an
idle connection after a default duration that must be no less than two hours.
The probe packet will contain no data<sup>2</sup>, and failure to reply to an individual
probe MUST NOT be interpreted as a dead connection, as individual probe packets
are not reliably transmitted.

Go allows you to enable TCP keepalive using `net.TCPConn`'s `SetKeepAlive`. On
OSX and Linux this will cause up to 8 TCP keepalive probes to be sent at an
interval of 75 seconds after a connection has been idle for 2 hours. Or in
other words, `Read` will return an `io.EOF` error after 2 hours and 10 minutes
(7200 + 8 * 75).

Depending on your application, that may be too long of a timeout. In this case
you can call `SetKeepAlivePeriod`. However, this method currently behaves
different for different operating systems. On OSX, it will modify the idle time
before probes are being sent. On Linux however, it will modify both the idle
time, as well as the interval that probes are sent at. So calling
`SetKeepAlivePeriod` with an argument of 30 seconds will cause a total timeout
of 10 minutes and 30 seconds for OSX (30 + 8 * 75), but 4 minutes and 30
seconds on Linux (30 + 8 * 30).

I found that situation rather unsatisfying, so I ended up creating a small
package called [tcpkeepalive][] that gives you more control:

```go
kaConn, _ := tcpkeepalive.EnableKeepAlive(conn)
kaConn.SetKeepAliveIdle(30*time.Second)
kaConn.SetKeepAliveCount(4)
kaConn.SetKeepAliveInterval(5*time.Second)
```

Currently only Linux and OSX are supported, but I'd be happy to merge pull
requests for other platforms. If there is interest from the Go core team, I'll
also try to contribute these new methods to Go itself.

Please let me know if you found this article useful, have any questions, or
spotted any errors so I can correct them.

## Appendix

1) Tuning a heartbeat mechanism to detect failures early with a low false
positive rate is tricky business. Checkout the [ϕ Accrual Failure Detector][] for a
statistically sound model, as well as Damian Gryski's [go-failure][]
implementation. Unfortunately there is no way to use it with TCP keepalive that
I can think of.

2) According to RFC 1122 keepalive probes may contain a single garbage octet
for compatibility with broken implementations. However, I'm not sure if this is
filtered out by the OS network stack or not, please comment if you know.

[go-failure]: https://github.com/dgryski/go-failure
[Damian Gryski]: https://github.com/dgryski
[ϕ Accrual Failure Detector]: http://ddg.jaist.ac.jp/pub/HDY+04.pdf
[net.Conn]: http://golang.org/pkg/net/#Conn
[RFC 1122]: http://tools.ietf.org/html/rfc1122#page-101
[tcpkeepalive]: https://github.com/felixge/tcpkeepalive
