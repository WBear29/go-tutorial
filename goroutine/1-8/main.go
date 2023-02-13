package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"demo/tracer"

	"go.opentelemetry.io/otel"
)

const (
	chapter = "1-8"
	times = 10000
)

func useIncrementOperator(ctx context.Context) uint32 {
	tr := otel.Tracer("useIncrementOperator")
	ctx, span := tr.Start(ctx, chapter+":useIncrementOperator")
	defer span.End()
	var cnt uint32
	var wg sync.WaitGroup

	for i := 0; i < times; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			_, childSpan := tr.Start(ctx, chapter+":useIncrementOperator:child")
			defer childSpan.End()
			cnt++
			wg.Done()
		}(ctx)
	}

	wg.Wait()

	return cnt
}

func useAtomicAddUint32(ctx context.Context) uint32 {
	tr := otel.Tracer("useAtomicAddUint32")
	ctx, span := tr.Start(ctx, chapter+":useAtomicAddUint32")
	defer span.End()
	var cnt uint32
	var wg sync.WaitGroup

	for i := 0; i < times; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			_, childSpan := tr.Start(ctx, chapter+":useAtomicAddUint32:child")
			defer childSpan.End()
			atomic.AddUint32(&cnt, 1)
			wg.Done()
		}(ctx)
	}

	wg.Wait()

	return cnt
}

func useSyncMutexLock(ctx context.Context) uint32 {
	tr := otel.Tracer("useSyncMutexLock")
	ctx, span := tr.Start(ctx, chapter+":useSyncMutexLock")
	defer span.End()
	var cnt uint32
	var wg sync.WaitGroup
	mu := new(sync.Mutex)

	for i := 0; i < times; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			_, childSpan := tr.Start(ctx, chapter+":useSyncMutexLock:child")
			defer childSpan.End()
			mu.Lock()
			defer mu.Unlock()
			cnt++
			wg.Done()
		}(ctx)
	}

	wg.Wait()

	return cnt
}

func main() {
	tp, err := tracer.NewTracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	tr := tp.Tracer("component-main")
	ctx, span := tr.Start(ctx, chapter+":main")
	defer span.End()

	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))

	// トレーシングを確認する際はコメントアウトしつつ各関数を実行する必要あり
	fmt.Printf("useIncrementOperator(): %d\n", useIncrementOperator(ctx))
	fmt.Printf("useAtomicAddUint32(): %d\n", useAtomicAddUint32(ctx))
	fmt.Printf("useSyncMutexLock(): %d\n", useSyncMutexLock(ctx))
}
