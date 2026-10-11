---
title: Killing a child process and all of its children in Go
date: 2015-09-23T14:05:08Z
---

This blog post explores some common unix child process control use cases, and how to perform them in Go. I’m mostly sharing it to improve my own understanding of the subject, but hopefully this will be useful to others as well.

**Problem:** You are using os/exec to run a child process, but need to terminate it under certain conditions, e.g. after a timeout of 3 seconds.

**Solution:** Just call `cmd.Process.Kill()` when your condition triggers:

```go
cmd := exec.Command("sleep", "5")
start := time.Now()
time.AfterFunc(3*time.Second, func() { cmd.Process.Kill() })
err := cmd.Run()
fmt.Printf("pid=%d duration=%s err=%s\n", cmd.Process.Pid, time.Since(start), err)
```

```text
$ go run example1.go
pid=41066 duration=3.00070126s err=signal: killed
$ ps -j
USER    PID  PPID  PGID   SESS JOBC STAT   TT       TIME COMMAND
felix 41008 41002 41008      0    0 Ss   s000    0:00.07 -/bin/bash
```

**Problem:** The child process is executing a child process of its own (they grow up so fast), but unfortunately the solution above only kills the child, but not the grandchild:

```go
cmd := exec.Command("/bin/sh", "-c", "watch date > date.txt")
start := time.Now()
time.AfterFunc(3*time.Second, func() { cmd.Process.Kill() })
err := cmd.Run()
fmt.Printf("pid=%d duration=%s err=%s\n", cmd.Process.Pid, time.Since(start), err)
```

```text
$ go run example2.go
pid=41216 duration=3.00145737s err=signal: killed
$ ps -j
USER    PID  PPID  PGID   SESS JOBC STAT   TT       TIME COMMAND
felix 41008 41002 41008      0    0 Ss   s000    0:00.11 -/bin/bash
felix 41217     1 41212      0    0 S    s000    0:00.00 watch date
$ kill 41217
```

As you can see, the `watch` process is still alive and has been adopted by [init](https://en.wikipedia.org/wiki/Init) (PPID = 1). The reason for this is that, under-the-hood, Go is is using [`kill(2)`](http://linux.die.net/man/2/kill) to send a `KILL` signal to the `PID` of the `sh` process, but not the `watch` process, turning it into an orphan. This is completely normal behavior for a unix programming language, but not desirable for many use cases.

**Solution:** In addition to sending a signal to a single `PID`, `kill(2)` also supports sending a signal to a [Process Group](https://en.wikipedia.org/wiki/Process_group) by passing the process group id (`PGID`) as a negative number. Child processes get the same `PGID` as their parents by default, so in order to avoid suicide, we’ll also need to ask Go to create a new process group for our child by setting the `Setpgid` field of [`syscall.SysProcAttr`](http://golang.org/pkg/syscall/#SysProcAttr). Putting it all together, we get this:

```go
cmd := exec.Command("/bin/sh", "-c", "watch date > date.txt")
cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
start := time.Now()
time.AfterFunc(3*time.Second, func() {
	syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
})
err := cmd.Run()
fmt.Printf("pid=%d duration=%s err=%s\n", cmd.Process.Pid, time.Since(start), err)
```

```text
$ go run example3.go
pid=41352 duration=3.001809822s err=signal: killed
$ ps -j
USER    PID  PPID  PGID   SESS JOBC STAT   TT       TIME COMMAND
felix 41008 41002 41008      0    0 Ss   s000    0:00.14 -/bin/bash
```

Internally this will [cause Go](https://github.com/golang/go/blob/master/src/syscall/exec_linux.go#L156-L163) to call [`setpgid(2)`](http://linux.die.net/man/2/setpgid) between [`fork(2)`](http://linux.die.net/man/2/fork) and [execve(2)](http://linux.die.net/man/2/execve), to assign the child process a new PGID identical to its PID. This allows us to kill all processes in the process group by sending a `KILL` to `-PID` of the process, which is the same as -PGID. Assuming that the child process did not use `setpgid(2)` when spawning its own child, this should kill the child along with all of its children on any *Nix systems.

**Further reading:** Process control in unix a is a complex topic, so I recommend the resources below to learn more.

* [Advanced Programming in the UNIX Environment, 3rd Edition](http://www.amazon.com/Advanced-Programming-UNIX-Environment-Edition/dp/0321637739)  
  *My unix bible.*
* [https://www.win.tue.nl/~aeb/linux/lk/lk-10.html](https://www.win.tue.nl/~aeb/linux/lk/lk-10.html)  
  *Perhaps a bit dated, but free and easy to understand.*

> I remember when it was possible to understand Unix.  
> — [Rob Pike](https://groups.google.com/d/msg/golang-dev/MdNdaAN5ZBY/lE06YS6gk_UJ) (talking about process control)

*Originally published on [Medium](https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773) on September 23, 2015.*
