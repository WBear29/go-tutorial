package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"demo/tracer"

	"go.opentelemetry.io/otel"
)

const (
	chapter = "1-3"
)

func getLuckyNum(ctx context.Context, ch chan<- int) {
	tr := otel.Tracer("gorountine-getLuckyNum")
	_, span := tr.Start(ctx, chapter+":goroutine")
	defer span.End()
	fmt.Println("...")

	// 占いにかかる時間はランダム
	rand.Seed(time.Now().Unix())
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)

	num := rand.Intn(10)
	ch <- num
}

func checkCloseCh(ok bool) {
	if ok {
		fmt.Println("ch is open")
	} else {
		fmt.Println("ch is closed")
	}
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

	fmt.Println("what is today's lucky number?")

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

	ch := make(chan int)    // チャネルを定義
	go getLuckyNum(ctx, ch) // goroutineで実行

	num, ok := <-ch
	checkCloseCh(ok)
	fmt.Printf("Today's your lucky number is %d!\n", num)
	
	// 使い終わったチャネルはcloseする
	close(ch)
	// クローズされたチャネルから受け取る
	num, ok = <-ch
	fmt.Printf("Today's your lucky number is %d!\n", num)
	checkCloseCh(ok)

	// クローズされたチャネルにSend
	// ch <- 3 // panic
	// クローズされたチャネルをClose
	// close(ch) // panic
}
