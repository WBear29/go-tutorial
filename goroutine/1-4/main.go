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
	chapter = "1-4"
)

func overChanelBuf(ctx context.Context, ch chan int) {
	tr := otel.Tracer("overChanelBuf")
	_, span := tr.Start(ctx, chapter+":overChanelBuf:sequential")
	defer span.End()
	ch <- 1
	fmt.Println("send 1")
	ch <- 2
	fmt.Println("send 2")
	ch <- 3
	fmt.Println("send 3")
	ch <- 4 // block
	fmt.Println("send 4")
	fmt.Println("要素数は", len(ch))

	i := <-ch
	fmt.Println(i)
	fmt.Println("要素数は", len(ch))
}
func openChanelBuf(ctx context.Context, ch chan int) {
	tr := otel.Tracer("openChanelBuf")
	_, span := tr.Start(ctx, chapter+":openChanelBuf:sequential")
	defer span.End()
	ch <- 1
	fmt.Println(<-ch)
	fmt.Println("要素数は", len(ch))
	ch <- 2
	ch <- 3
	ch <- 4
	fmt.Println("要素数は", len(ch))

	i := <-ch
	fmt.Println(i)
	fmt.Println("要素数は", len(ch))
}

func onlySendChanel(ctx context.Context, ch chan<- int) {
	tr := otel.Tracer("onlySendChanel")
	_, span := tr.Start(ctx, chapter+":onlySendChanel:goroutine")
	defer span.End()

	fmt.Println("Start onlySendChanel")
	for i := 0; i < cap(ch); i++ {
		ch <- i
	}
	close(ch)
}

func getOnlyReceiveChanel(ctx context.Context) <-chan int {
	tr := otel.Tracer("getOnlyReceiveChanel")
	_, span := tr.Start(ctx, chapter+":getOnlyReceiveChanel:sequential")
	defer span.End()
	var genCh = make(chan int)
	go func() {
		var primes = make([]int, 0)
		for k := 2; ; k++ {
			check := true
			for _, value := range primes {
				if k%value == 0 {
					check = false
					break
				}
			}
			if check {
				primes = append(primes, k)
				genCh <- k
			}
		}
		close(genCh)
	}()
	return genCh
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
	チャネルの作り方
	1. 基本
	ch := make(chan int)
	2. バッファあり
	ch := make(chan string, 2)
	3. 送信専用チャネル(バッファあり)
	make(chan<= int, 10)
	4. 受信専用チャネル(バッファあり)
	make(<-chan int, 10)
	チャネルの状態とoperation: https://docs.google.com/presentation/d/1WDVYRovp4eN_ESUNoZSrS_9WzJGz_-zzvaIF4BgzNws/edit#slide=id.gd0f0d38d56_0_1329
	**/

	ch := make(chan int, 3) // バッファありチャネルを定義
	fmt.Println("現在のバッファサイズは", cap(ch))
	// ex) 1.
	overChanelBuf(ctx, ch)
	// ex) 2.
	// openChanelBuf(ctx, ch)
	// ex) 3.
	// go onlySendChanel(ctx, ch)
	// for i := range ch {
	// 	fmt.Println(i)
	// }
	// ex) 4.
	// receiveCh := getOnlyReceiveChanel(ctx)
	// for i := 0; i < 999; i++ {
	// 	<-receiveCh
	// }
	// fmt.Println(<-receiveCh)
	fmt.Println("Done")
}
