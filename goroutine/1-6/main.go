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
	chapter = "1-6"
)

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
	指定した時間待つ
	link: https://www.spinute.org/go-by-example/timers.html
	**/

	// timer1 := time.NewTimer(2 * time.Second)
	// <-timer1.C
	// fmt.Println("Timer 1 fired")

	// timer2 := time.NewTimer(time.Second)
	// go func() {
	// 	<-timer2.C
	// 	fmt.Println("Timer 2 fired")
	// }()
	// stop2 := timer2.Stop() // timer2を途中で停止
	// if stop2 {
	// 	fmt.Println("Timer 2 stopped")
	// }
	// time.Sleep(2 * time.Second)

	/**
	チャネルの使い方
	定期的に繰り返し実行
	link: https://www.spinute.org/go-by-example/tickers.html
	**/

	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}
