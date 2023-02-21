# 目的

[golangci-lint](https://github.com/golangci/golangci-lint)についてまとめCI, VSCodeへの導入時の推奨設定の検討をサポートします。

最終更新日: 2022/02/21

## Linters

* 参考: [Linters](https://golangci-lint.run/usage/linters/#asasalint)

✅: 推奨
🔧: オプション
🔗: リンク
⚙: 設定

### errcheck ✅

プログラム上でチェックされていない `error` を検出。
チェックされていないエラーはクリティカルなバグになるため推奨。
[🔗](https://github.com/kisielk/errcheck), [⚙](https://golangci-lint.run/usage/linters/#errcheck)

### gosimple ✅

コードをシンプルに書くためのチェックを実行。
必要ない処理を記述している箇所などを指摘してくれます。
余計なコードを減らせるため推奨。
[参考: golangci-lintに搭載されている linter を学ぶ-gosimple](https://zenn.dev/sanpo_shiho/books/61bc1e1a30bf27/viewer/642fe9#gosimple)
[🔗](https://github.com/dominikh/go-tools/tree/master/simple), [⚙](https://golangci-lint.run/usage/linters/#gosimple)

### govet ✅

公式が提供している `go vet`。
引数と書式が一致しない`Printf`呼び出しなどの構造をチェックする。
コンパイラが検出できないエラーを検出できるため、推奨。
[🔗](https://pkg.go.dev/cmd/vet), [⚙](https://golangci-lint.run/usage/linters/#govet)

### ineffassign ✅

変数への不要な代入が行われている箇所を検出する。
[参考: golangci-lintに搭載されている linter を学ぶ-ineffassign](https://zenn.dev/sanpo_shiho/books/61bc1e1a30bf27/viewer/642fe9#gosimple)
`ineffectual`(効果がない)と`assigned`(割り当て)から命名されていると思われる。
不要な割り当てを検出できるため、推奨。
[🔗](https://github.com/gordonklaus/ineffassign), ⚙: なし

### staticcheck ✅

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

### typecheck 🔧

コンパイルが通るかのチェック
CI上でビルドステージをどこに配置するかに依存するため、オプショナル

### unused ✅

未使用の変数、定数、関数、型を検出する。
未使用の定義を残すことは言語仕様で推奨されていないため、推奨。
※ToDo: 出典探す
[🔗](https://github.com/dominikh/go-tools/tree/master/unused), ⚙: なし


