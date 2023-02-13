package main

import (
	"context"
	"testing"
)

func BenchmarkUseIncrementOperator(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		useIncrementOperator(ctx)
	}
}

func BenchmarkUseAtomicAddUint64(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		useAtomicAddUint32(ctx)
	}
}

func BenchmarkUseSyncMutexLock(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		useSyncMutexLock(ctx)
	}
}