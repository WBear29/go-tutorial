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
	複数のチャネルからの受信を制御する
	**/

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		tr := otel.Tracer("one")
		_, span := tr.Start(ctx, chapter+":one:goroutine")
		defer span.End()
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		tr := otel.Tracer("two")
		_, span := tr.Start(ctx, chapter+":two:goroutine")
		defer span.End()
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
}
