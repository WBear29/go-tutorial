package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"demo/tracer"

	"go.opentelemetry.io/otel"
)

const (
	chapter = "1-5"
)

func worker(ctx context.Context, doneCh chan<- bool) {
	tr := otel.Tracer("worker")
	_, span := tr.Start(ctx, chapter+":worker:goroutine")
	defer span.End()

	fmt.Println("working...")
	time.Sleep(time.Duration(1) * time.Second)

	doneCh <- true
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

	/**
	チャネルの使い方:
	goroutineが終了するのを待つ
	**/

	doneCh := make(chan bool, 1) // バッファありチャネルを定義
	go worker(ctx, doneCh)
	<-doneCh
	fmt.Println("Done")

	/**
	複数のgorutine終了を待つ
	**/

	// doneCh2 := make(chan bool, 2)
	// go worker(ctx, doneCh2)
	// go worker(ctx, doneCh2)
	// for i := 0; i < 2; i++ {
	// 	<-doneCh2
	// }
	// fmt.Println("Done")
	// doneCh3 := make(chan int)
	// for i := 0; i < 100; i++ {
	// 	go func(i int) {
	// 		tr := otel.Tracer("channel")
	// 		_, span := tr.Start(ctx, chapter+":channel:goroutine")
	// 		defer span.End()
	// 		time.Sleep(2 * time.Second)
	// 		doneCh3 <- i
	// 	}(i)

	// }
	// for i := 0; i < 100; i++ {
	// 	fmt.Println("End:", <-doneCh3)
	// }
	// fmt.Println("Done")

	/**
	複数のgorutine終了を待つ
	**/

	// var wg sync.WaitGroup
	// for i := 0; i < 100; i++ {
	// 	wg.Add(1)
	// 	go func(i int) {
	// 		defer wg.Done()
	// 		tr := otel.Tracer("waitgroup")
	// 		_, span := tr.Start(ctx, chapter+":waitgroup:goroutine")
	// 		defer span.End()

	// 		time.Sleep(2 * time.Second)
	// 		fmt.Println("End:", i)
	// 	}(i)
	// }
	// wg.Wait()
	// fmt.Println("Done")
}
