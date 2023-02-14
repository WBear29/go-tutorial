# How to Use

1. 準備

   ```
   $ docker-compose up -d
   http://localhost:16686
   ```

1. 各リポジトリに移動して実行
   ```
   ex)
   $ cd 1-1
   $ go mod tidy
   $ go run .
   ```


## 内容

### 1-1 お試しでGorutine実行

* time.Sleepをコメントアウトすると数値が出力される前にmainが終了

### 1-2 sync.WaitGroupでgoroutineの終了を待つ

* wg.Wait() でgoroutineの終了を待つように記述

### 1-3 チャネルを利用した値の送信

* gorutineで並行処理している関数からチャネルで値を受信

### 1-4 チャネルを利用した値の送信2

1. チャネルのバッファを超えて送信
1. チャネルのバッファを解放しつつ送信
1. 送信専用チャネル
1. 受信専用チャネル

### 1-5 チャネルの使い方: 並行処理の終了を待つ

1. gorutineの終了を待つ
1. 複数のgorutineの終了を待つ
1. sync.WaitGroupで複数のgorutineの終了を待つ

### 1-6 selectで複数のチャネルから値を受け取る

* selectを用いて複数のチャネルからの受信を制御

### 1-7 チャネルの使い方: 指定時間処理を待つ&定期的に処理を実行

1. time.NewTimerで待つ
   * time.Sleepでもできるが途中でキャンセルができる
1. time.NewTickerで処理の定期実行

### 1-8 atomicで並行処理での安全なインクリメント

* アトミック性:
   * それが操作されている特定のコンテキストの中では、分割不能あるいは中断不可であること。
      * 分割不能あるいは中断不可ということは、そのコンテキストの全体で処理をしていて、その他の何かが同時には起きないということ。

atomicを利用することで、並行処理で安全に変数を扱う。


1. GOMAXPROCS=1 go run main.go
   * 全ての関数が同じ値を出力
1. GOMAXPROCS=4 go run main.go
   * useIncrementOperatorのみ値が変わる
1. go run -race main.go 
   * 競合を検知
1. ベンチマーク確認
   * GOMAXPROCS=1 go test -bench . -benchmem
   * GOMAXPROCS=4 go test -bench . -benchmem
   * GOMAXPROCS=10 go test -bench . -benchmem

   ```
   実行関数名 実行回数 1回あたりの実行に掛かる時間 一回あたりのアロケーションで確保した容量(B/op) 一回あたりのアロケーション回数(allocs/op)
   ```

## ToDo
* Tracerの記述方法確認

## 参考文献
- [Go言語による並行処理](https://www.oreilly.co.jp/books/9784873118468/)
- [Go での並行処理を徹底解剖！](https://zenn.dev/hsaki/books/golang-concurrency/viewer/intro)
- [opentelemetry-go/example/jaeger](https://github.com/open-telemetry/opentelemetry-go/blob/main/example/jaeger/main.go)
- [サンプルで学ぶ Go 言語](https://www.spinute.org/go-by-example/)
- [sync/atomic パッケージを使って数値を安全にカウントアップする](https://qiita.com/yuya_takeyama/items/98e5a9fe6786df11c8a7)