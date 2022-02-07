package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	file, _ := os.Create("./cpu.pprof")
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	ctx := context.Background()
	span := StartSpan(ctx, "GET /users/:id")
	defer span.Finish()

	ch := time.After(time.Second)
	for {
		select {
		case <-ch:
			return
		default:
		}
	}
}

func StartSpan(ctx context.Context, endpoint string) *span {
	span := &span{id: rand.Uint64(), restoreCtx: ctx}
	labels := pprof.Labels(
		"span_id", fmt.Sprintf("%d", span.id),
		"endpoint", endpoint,
	)
	pprof.SetGoroutineLabels(pprof.WithLabels(ctx, labels))
	return span
}

type span struct {
	id         uint64
	restoreCtx context.Context
}

func (s *span) Finish() {
	pprof.SetGoroutineLabels(s.restoreCtx)
}
