# 目的

[golangci-lint](https://github.com/golangci/golangci-lint)についてまとめCI, VSCodeへの導入時の推奨設定の検討をサポートします。

## Linters

* 参考: [Linters](https://golangci-lint.run/usage/linters/#asasalint)

✅: 推奨  
🔧: オプション  
🔗: リンク  
⚙: 設定  

### デフォルト
#### errcheck ✅

プログラム上でチェックされていない `error` を検出。  
チェックされていないエラーはクリティカルなバグになるため推奨。

[🔗](https://github.com/kisielk/errcheck), [⚙](https://golangci-lint.run/usage/linters/#errcheck)

#### gosimple ✅

コードをシンプルに書くためのチェックを実行。  
必要ない処理を記述している箇所などを指摘してくれます。  
余計なコードを減らせるため推奨。

[参考: golangci-lintに搭載されている linter を学ぶ-gosimple](https://zenn.dev/sanpo_shiho/books/61bc1e1a30bf27/viewer/642fe9#gosimple)

[🔗](https://github.com/dominikh/go-tools/tree/master/simple), [⚙](https://golangci-lint.run/usage/linters/#gosimple)

#### govet ✅

公式が提供している `go vet`。  
引数と書式が一致しない`Printf`呼び出しなどの構造をチェックする。  
コンパイラが検出できないエラーを検出できるため、推奨。

[🔗](https://pkg.go.dev/cmd/vet), [⚙](https://golangci-lint.run/usage/linters/#govet)

#### ineffassign ✅

変数への不要な代入が行われている箇所を検出する。  
[参考: golangci-lintに搭載されている linter を学ぶ-ineffassign](https://zenn.dev/sanpo_shiho/books/61bc1e1a30bf27/viewer/642fe9#gosimple)  
`ineffectual`(効果がない)と`assigned`(割り当て)から命名されていると思われる。  
不要な割り当てを検出できるため、推奨。

[🔗](https://github.com/gordonklaus/ineffassign), ⚙: なし

#### staticcheck ✅

複数の静的解析を実施します。

[チェック項目一覧](https://staticcheck.io/docs/checks/)
- SA系: staticcheck
  - SA1: 標準ライブラリの誤用検出
  - SA2: 同時実行の問題検出
  - SA3: テスト問題の検出
  - SA4: 何もしていない不要なコード検出
  - SA5: 正確性の問題
    - ex) nil mapへの代入の検出, 無限ループ内の`defer` 定義...etc
  - SA6: パフォーマンス問題の検出
  - SA9: 間違っている可能性があるコード構造の検出
- S系: コード簡略化
- ST系: stylecheck

※`staticcheck`に承認されていないため、`staticcheck`と同じバイナリではない。

[🔗](https://staticcheck.io/), [⚙](https://golangci-lint.run/usage/linters/#staticcheck)

#### typecheck 🔧

コンパイルが通るかのチェック  
CI上でビルドステージの配置によってはなくても良いため、オプショナル  

#### unused ✅

未使用の変数、定数、関数、型を検出する。  
未使用の定義を残すことは言語仕様で推奨されていないため(ToDo: 出典探す)、推奨。

[🔗](https://github.com/dominikh/go-tools/tree/master/unused), ⚙: なし

### デフォルトは無効

#### asasalint ✅

anyの可変長(...any)を引数で受け取る関数に[]anyを渡している箇所を検出する。
バグを検出できるため、推奨。

[fmt.Printf](https://pkg.go.dev/fmt#Printf)

```golang
package main

import "fmt"

func main() {
	msg := []any{"hello", "world"}
	// p1 [hello world] p2 %!s(MISSING)
	fmt.Printf("p1 %s p2 %s\n", msg)
	// p1 hello p2 world
	fmt.Printf("p1 %s p2 %s\n", msg...)
}
```

[🔗](https://github.com/alingse/asasalint), [⚙](https://golangci-lint.run/usage/linters/#asasalint)

#### asciicheck ✅

コードに非 ASCII 識別子が含まれている箇所を検出する。  
非推奨文字コードチェックのため、推奨。

[🔗](https://github.com/tdakkota/asciicheck), ⚙: なし

#### bidichk ✅

危険と定義されている Unicode 文字シーケンスを検出する。  
非推奨文字コードチェックのため、推奨。

[🔗](https://github.com/alingse/asasalint), [⚙](https://golangci-lint.run/usage/linters/#asasalint)

#### bodyclose ✅

httpのレスポンスbodyがCloseされていない箇所を検出する。  
レスポンスbodyをCloseしない場合、コネクションがブロックされ再利用されず残ったままとなる。  
新しい接続を行う度に新しいGoroutineとファイルディスクリプタを生成してしまい、メモリリークの要因となり得る。  
[参考: Goのnet/httpのclientでなぜresponseBodyをClose、読み切らなくてはいけないのか](https://zenn.dev/cube/articles/4ce18a672fc991#responsebody%E3%82%92close%E3%81%97%E3%81%AA%E3%81%84%E3%81%A8%E3%81%84%E3%81%91%E3%81%AA%E3%81%84%E3%82%8F%E3%81%91)  
バグの原因になるため、推奨。

[🔗](https://github.com/timakin/bodyclose), ⚙: なし

#### containedctx ✅

golangで推奨されていない、contextを構造体のフィールドに持っている箇所を検出する。  
公式非推奨の構文のため、推奨。contextは第一引数に渡す形式で統一しておくのが良いと思われる。  

[🔗](https://github.com/sivchari/containedctx), ⚙: なし

#### contextcheck ✅

contextが適切に継承されていない関数の呼び出しを検出します。  
contextを意図通りに渡せていない箇所を検出できるため、推奨。

[🔗](https://github.com/kkHAIKE/contextcheck), ⚙: なし

#### cyclop 🔧

関数とパッケージの循環的複雑度(Cyclomatic complexity)をチェックし、設定した閾値を超えたことを検出します。  
max-complexityとpackage-averageを設定可能です。  
リファクタリングの際などの参考にできるため、オプショナル

[🔗](https://github.com/bkielbasa/cyclop), [⚙](https://golangci-lint.run/usage/linters/#cyclop)

#### deadcode

現在保守されていないので、非推奨

[🔗](https://github.com/remyoudompheng/go-misc/tree/master/deadcode), ⚙: なし

## 参考文献

* [golangci-lint](https://github.com/golangci/golangci-lint)
* [golangci-lintに搭載されている linter を学ぶ](https://zenn.dev/sanpo_shiho/books/61bc1e1a30bf27/viewer/642fe9)