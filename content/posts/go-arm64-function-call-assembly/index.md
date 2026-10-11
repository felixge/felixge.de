---
title: "Go arm64 Function Call Assembly"
date: 2023-02-05T09:17:02.000Z
slug: "go-arm64-function-call-assembly"
cover: cover.png
aliases:
  - /go-arm64-function-call-assembly/
---

I am currently working on implementing [frame pointer unwinding](/2023/01/31/reducing-gos-execution-tracer-overhead-with-frame-pointer-unwinding/) for the Go execution tracer. This involves debugging various problems and crashes caught by [TestTraceSymbolize](https://github.com/golang/go/blob/25f5d9d4a253cb880bff909ebbdc05c8941c4a48/src/runtime/trace/trace_stack_test.go#L26) on arm64.

One of these crashes seems to be caused by a goroutine overwriting the frame pointer on the stack of another goroutine. To debug this, I'm trying to study the assembly output of the compiler. (Un)fortunately I don't read or write Go assembly every day, so this process usually involves chaotic jumping between the following resources:

1.  [A Quick Guide to Go's Assembler](https://go.dev/doc/asm) (Go's Assembly Language is quite quirky)
2.  [Go ARM64 Assembly Instructions Reference Manual](https://pkg.go.dev/cmd/internal/obj/arm64) (arm64 specific quirks)
3.  [Go internal ABI specification](https://go.dev/s/regabi) (Go's register based calling convention)
4.  [Arm Architecture Reference Manual](https://developer.arm.com/documentation/ddi0487/latest/) (11.952 pages PDF)
5.  [Arm A64 Instruction Set Architecture](https://developer.arm.com/documentation/ddi0602/2022-12/?lang=en) (HTML subset of the above)
6.  [ChatGPT](https://chat.openai.com/chat) (with a very healthy amount of mistrust,- it's often wrong)
7.  Various Google Search Results

To make this process a little easier for my future self, I decided it's time to create some high quality notes for my own needs. In particular, I'll try to explain every assembly instruction the Go compiler (`go1.20 darwin/arm64`) emits for the following code in great detail:

``` go
//go:noinline
func foo() { bar() }

//go:noinline
func bar() {}
```

To get the assembly, I use the steps shown below.

``` bash
$ go build main.go
$ go tool objdump -gnu -s 'main.(foo|bar)' ./main
```

The `-s 'main.(foo|bar)'` flag filters the output down the two functions we are interested in, and the `-gnu` flag adds GNU assembly comments which make it easier to lookup instructions in the official arm documentation. Normally I'd also include the `-S` flag to see each line of source code above the instructions it generated, but that's not needed here.

After some manual trimming, the `objdump` output looks like this:

```
TEXT main.foo(SB) /Users/felixge/Desktop/main.go
  MOVD 16(R28), R16                    // ldr x16, [x28,#16]
  CMP R16, RSP                         // cmp sp, x16
  BLS 8(PC)                            // b.ls .+0x20
  MOVD.W R30, -16(RSP)                 // str x30, [sp,#-16]!
  MOVD R29, -8(RSP)                    // stur x29, [sp,#-8]
  SUB $8, RSP, R29                     // sub x29, sp, #0x8
  CALL main.bar(SB)                    // bl .+0x28
  LDP -8(RSP), (R29, R30)              // ldp x29, x30, [sp,#-8]
  ADD $16, RSP, RSP                    // add sp, sp, #0x10
  RET                                  // ret
  MOVD R30, R3                         // mov x3, x30
  CALL runtime.morestack_noctxt.abi0(SB) // bl .+0xffffffffffffdd54
  JMP main.foo(SB)                     // b .+0xffffffffffffffd0
  ?
  ?
  ?

TEXT main.bar(SB) /Users/felixge/Desktop/main.go
  RET                                  // ret
  ?
  ?
  ?
```

Before we dig into this output I want to mention that I'm not using Compiler Explorer for this, because its [output](https://godbolt.org/z/Pr7q619a4) seems to be [missing important instructions](https://a.cl.ly/mXubO0ol?method=GET) for some reason. I'm also not using `go tool compile -S main.go` or `go build -gcflags=-S ./main.go` because those output additional [pseudo instructions](https://gist.github.com/felixge/7f8c8b9adc0e8e26be5209a26e33fcd6#file-compile-txt-L5) that are used by the compiler and linker, but don't remain in the final assembly.

Ok, let's dig in!

## Prologue 1: Check if the stack needs growing

```
  MOVD 16(R28), R16                    // ldr x16, [x28,#16]
  CMP R16, RSP                         // cmp sp, x16
  BLS 8(PC)                            // b.ls .+0x20
```

1.  Store the value of `g.stackguard0` in the `R16` register.
2.  Compare the value of `R16` with the value of the stack pointer `RSP`.
3.  If `RSP <= R16`, jump `8` instructions forward  to grow the stack

**In Detail**

1.  `MOVD`, aka [LDR (immediate)](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/LDR--immediate---Load-Register--immediate--?lang=en), loads the data from the memory address \[`R28`  + `16` \] into `R16`. The `R28` register points to the [current goroutine](https://github.com/golang/go/blob/de4748c47c67392a57f250714509f590f68ad395/src/cmd/compile/abi-internal.md?plain=1#L549), and the `stackguard0` [field](https://github.com/golang/go/blob/de4748c47c67392a57f250714509f590f68ad395/src/runtime/runtime2.go#L416) is found at offset `16` in the `g` struct because the `stack` field above it is `16` bytes wide.
2.  `CMP`, aka [CMP (extended register)](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/CMP--extended-register---Compare--extended-register---an-alias-of-SUBS--extended-register--?lang=en), compares `R16` with `RSP` and stores the results in the condition flags.
3.  `BLS`, aka [B.cond](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/B-cond--Branch-conditionally-?lang=en), checks if the condition flags match the [condition code](https://developer.arm.com/documentation/dui0473/m/condition-codes/comparison-of-condition-code-meanings-in-integer-and-floating-point-code?lang=en) `ls` (lower or same). If yes, the CPU is instructed to jump `8` instructions (`0x20` = `32` bytes) forward. We'll cover those instructions in the [Prologue 3: Growing the stack](#prologue-3-growing-the-stack) section. However, usually this is not the case, and execution continues with the instructions below.

## Prologue 2: Setup foo's frame

```
  MOVD.W R30, -16(RSP)                 // str x30, [sp,#-16]!
  MOVD R29, -8(RSP)                    // stur x29, [sp,#-8]
  SUB $8, RSP, R29                     // sub x29, sp, #0x8
```

1.  Store the return address from the link register `R30` at `16` bytes below the stack pointer register `RSP` and update `RSP` to that location.
2.  Store the frame pointer register `R29` of the caller `8` bytes below the stack pointer.
3.  Set the frame pointer register `R29` to point to the the previous frame pointer we just pushed onto the stack.

**In Detail**

1.  `MOVD.W`, aka [STR (immediate)](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/STR--immediate---Store-Register--immediate--?lang=en), stores the value of the `R30` register at the memory address of \[`RSP` - `16`\]. `R30` is the [link register](https://github.com/golang/go/blob/de4748c47c67392a57f250714509f590f68ad395/src/cmd/compile/abi-internal.md?plain=1#L547) which holds the return address of the caller. The operation is [pre-indexed](https://developer.arm.com/documentation/102374/0101/Loads-and-stores---addressing), which means that the memory address is computed before the memory is accessed and that the `RSP` pointer is updated to the computed address afterwards (i.e. `RSP = RSP - 16`). Technically this frame only uses `8` bytes above `RSP` (for storing the return address), but the architecture [requires](https://github.com/golang/go/blob/de4748c47c67392a57f250714509f590f68ad395/src/cmd/compile/abi-internal.md?plain=1#L572-L573) the stack pointer to be `16`-byte aligned, so `8` bytes of memory are wasted here.
2.  `MOVD`, aka [STUR](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/STUR--Store-Register--unscaled--?lang=en), stores the value of the `R29` register at the memory address of \[`RSP` - `8`\]. `R29` is the [frame pointer register](https://github.com/golang/go/blob/de4748c47c67392a57f250714509f590f68ad395/src/cmd/compile/abi-internal.md?plain=1#L548) that is holding the caller's frame pointer.
3.  `SUB`, aka [SUB (immediate)](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/SUB--immediate---Subtract--immediate--?lang=en), subtracts `8` from the stack pointer register `RSP` and stores the result in the frame pointer register `R29`. In other words, it points the frame pointer register to the caller's frame pointer that was just pushed onto the stack by the previous operation.

## Body: Call bar

```
  CALL main.bar(SB)                    // bl .+0x28
```

1.  Call `main.bar`.

**In Detail**

1.  `CALL`, aka [BL](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/BL--Branch-with-Link-?lang=en), implicitly adds `4` bytes (the size of an instruction) to the value of the current program counter and stores it in the link register `R30`. After that it jumps to the first instruction of the `main.bar` function which is located `10` instructions (`0x28` = `40` bytes) below this instruction.

## Epilogue: Return from foo

```
  LDP -8(RSP), (R29, R30)              // ldp x29, x30, [sp,#-8]
  ADD $16, RSP, RSP                    // add sp, sp, #0x10
  RET                                  // ret
```

1.  Restore the frame pointer register `R29` and the link register `R30` from the values stored `8` bytes below the stack pointer `RSP`.
2.  Restore the stack pointer `RSP` to its original value by adding `16` to it.
3.  Return to the caller.

**In Detail**

1.  `LDP`, aka [LDP](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/LDP--Load-Pair-of-Registers-?lang=en),  loads `16` bytes from the memory address \[`RSP` - `8`\] into a pair of registers. The first `8` bytes of the data are loaded into the frame pointer register `R29`, and the second `8` bytes are loaded into the link register `R30`. This is needed because these registers are [call-preserved](https://stackoverflow.com/a/56178078), so we need to restore them to their original values. Remember: The link register `R30` was overwritten by [Body: Call bar](#body-call-bar) and the frame pointer register `R29` was overwritten in [Prologue 2: Setup foo's frame](#prologue-2-setup-foos-frame).
2.  `ADD`, aka [ADD (immediate)](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/ADD--immediate---Add--immediate--?lang=en), adds `16` to the stack pointer register `RSP` and stores the result in the stack pointer register `RSP`. This is needed because `RSP` is another call-preserved register that we overwrote in [Prologue 2: Setup foo's frame](#prologue-2-setup-foos-frame).
3.  `RET`, aka [RET](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/RET--Return-from-subroutine-?lang=en) returns to the caller. The implied jump location is the value of the link register `R30`.

## Prologue 3: Growing the stack

```
  MOVD R30, R3                         // mov x3, x30
  CALL runtime.morestack_noctxt.abi0(SB) // bl .+0xffffffffffffdd54
  JMP main.foo(SB)                     // b .+0xffffffffffffffd0
  ?
  ?
  ?
```

This is the continuation of [Prologue 1: Check if the stack needs growing](#prologue-1-check-if-the-stack-needs-growing).

1.  Store the value of the `R30` link register in the `R3` register.
2.  Call the `runtime.morestack_noctxt` function to grow the stack.
3.  Jump back to the beginning of the `main.foo` function.
4.  `?` indicates zero padding.

**In Detail**

1.  `MOVD`, aka [MOV (register)](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/MOV--register---Move--register---an-alias-of-ORR--shifted-register--?lang=en), copies the value of the link register `R30` to the `R3` register. This is done to make it available to `runtime.morestack_noctxt`. As far as I can tell the value is only used to help with printing debug information, especially when a function is tyring to split the stack [when it's not supposed to](https://github.com/golang/go/blob/de4748c47c67392a57f250714509f590f68ad395/src/runtime/stack.go#L979-L998).
2.  `CALL`, aka [BL](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/BL--Branch-with-Link-?lang=en), implicitly adds `4` bytes (the size of an instruction) to the value of the current program counter and stores it in the link register `R30`. After that it jumps to the first instruction of the `runtime.morestack_noctxt` [function](https://github.com/golang/go/blob/de4748c47c67392a57f250714509f590f68ad395/src/runtime/asm_arm64.s#L322-L331) which is located `2219` instructions (`0xffffffffffffdd54` = `-8876` bytes in [two's complement](https://en.wikipedia.org/wiki/Two%27s_complement)) above this instruction. The called function usually grows the stack of the goroutine before it [implicitly returns](https://github.com/golang/go/blob/de4748c47c67392a57f250714509f590f68ad395/src/runtime/stack.go#L1119). Another use case is the [preemption](https://github.com/golang/go/blob/de4748c47c67392a57f250714509f590f68ad395/src/runtime/stack.go#L1052-L1073) of goroutines.
3.  `JMP`, aka [B](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/B--Branch-?lang=en), jumps jumps to the first instruction of the current function (`main.foo`) which is located `12` instructions (`0xffffffffffffffd0` = `-48` bytes in [two's complement](https://en.wikipedia.org/wiki/Two%27s_complement)) above this instruction. This causes the function to be executed again, now with a big enough stack to hold all of its values.
4.  The three `?` each indicate `4` bytes of zero padding which are emitted by the compiler to `16`-byte align all function entry points. I'm unable to find an official reference that mentions this value. But since it was [provided](https://go-review.googlesource.com/c/go/+/225397) by an arm employee, and another arm employee added a [similar patch](https://patchwork.ozlabs.org/project/gcc/patch/VI1PR0801MB21279519E21E8E2AEED6C6B883190@VI1PR0801MB2127.eurprd08.prod.outlook.com/#2183053) to gcc citing performance benefits, this alignment is probably a good idea.

## Function bar

```
  RET                                  // ret
  ?
  ?
  ?
```

1.  Return to the caller.
2.  `?` indicates zero padding.

**In Detail**

1.  `RET`, aka [RET](https://developer.arm.com/documentation/ddi0596/2021-12/Base-Instructions/RET--Return-from-subroutine-?lang=en) returns to the caller. The implied jump location is the value of the link register `R30`.
2.  The three `?` indicate zero padding as explained above.

## Final Thoughts

I only covered `arm64` in this article because that's the architecture used by my laptop and perhaps also the future of the cloud. Writing a similar article for `amd64` would be nice, but I'm not sure if I'll find the time.

Anyway, I hope that this information will be useful to others as well as to my future self. Please let me know if you spot any mistakes or have any questions!
