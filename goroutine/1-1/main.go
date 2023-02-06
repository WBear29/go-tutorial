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

func getLuckyNum(ctx context.Context) {
	tr := otel.Tracer("gorountine-getLuckyNum")
	_, span := tr.Start(ctx, "1-1:goroutine")
	defer span.End()
	fmt.Println("...")

	// 占いにかかる時間はランダム
	rand.Seed(time.Now().Unix())
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)

	num := rand.Intn(10)
	fmt.Printf("Today's your lucky number is %d!\n", num)
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
	ctx, span := tr.Start(ctx, "1-1:main")
	defer span.End()

	fmt.Println("what is today's lucky number?")
	go getLuckyNum(ctx) // goroutineで実行

	time.Sleep(time.Second * 5) // コメントアウトするとmainのgoruntineを先に終了
}
