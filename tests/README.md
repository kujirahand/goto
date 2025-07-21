# goto テストスイート

このディレクトリには、goto プログラムの包括的なテストスイートが含まれています。

## テストファイル一覧

### メインテストスクリプト

- **`test`** - シンプルなテスト実行スクリプト（CI/CDやクイックテスト用）
- **`run_tests.sh`** - 統合テストスイート実行スクリプト

### 個別テストスクリプト

- **`test.sh`** - 基本機能テスト
- **`functional_test.sh`** - 機能テスト
- **`performance_test.sh`** - パフォーマンステスト

## 使用方法

### クイックテスト（推奨）

```bash
# 最も簡単な方法
./test

# または
./test quick
```

### 完全なテストスイート

```bash
# すべてのテストを実行
./test full

# または
./run_tests.sh
```

### 個別テスト

```bash
# 基本機能テストのみ
./test basic
./run_tests.sh --basic-only

# 機能テストのみ
./test functional
./run_tests.sh --functional-only

# パフォーマンステストのみ
./test performance
./run_tests.sh --performance-only

# パフォーマンステスト以外
./run_tests.sh --no-performance
```

## テスト内容

### 基本機能テスト (`test.sh`)

- バイナリの存在確認
- コマンドラインオプションの動作確認
  - `-h, --help` (ヘルプ表示)
  - `-v, --version` (バージョン表示)
  - `--complete` (補完候補表示)
  - `--history` (履歴表示)
  - `-c` (カーソルモード)
  - `-l` (ラベル入力モード)
- デフォルト設定ファイルの作成
- 無効な引数の処理

### 機能テスト (`functional_test.sh`)

- TOML設定ファイルの解析
- 履歴機能（JSON）
- パス展開（`~/` など）
- URL処理
- ショートカット機能
- インタラクティブモードオプション
- エラーハンドリング

### パフォーマンステスト (`performance_test.sh`)

- 起動時間の測定
- 大きな設定ファイルの読み込み性能
- 大きな履歴ファイルの処理性能
- メモリ使用量の確認
- 並行実行の性能

## 要件

### 必須

- Go (golang) - プログラムのビルドに必要
- Bash - テストスクリプトの実行に必要

### オプション（パフォーマンステスト用）

- `bc` - 計算処理に使用
- `time` - メモリ使用量測定に使用

インストール例：

```bash
# macOS
brew install bc

# Ubuntu/Debian
sudo apt-get install bc time

# CentOS/RHEL
sudo yum install bc time
```

## テスト環境

テストは自動的に以下のように環境を管理します：

1. **バックアップ**: 既存の設定ファイル（`~/.goto.toml`）と履歴ファイル（`~/.goto.history.json`）を一時的にバックアップ
2. **テスト実行**: テスト用の設定でテストを実行
3. **復元**: 元の設定ファイルと履歴ファイルを復元

## CI/CD での使用

CI/CD パイプラインでは以下のように使用できます：

```bash
# GitHub Actions, GitLab CI, etc.
cd test
./test quick  # クイックテストで基本的な動作確認

# より包括的なテスト
./run_tests.sh --no-performance  # パフォーマンステスト以外
```

## トラブルシューティング

### 権限エラー

```bash
chmod +x test/*.sh
```

### Go が見つからない

Go がインストールされていることを確認してください：

```bash
go version
```

### テンポラリファイルの残留

テストが途中で中断された場合、以下で手動クリーンアップできます：

```bash
rm -rf test/temp*
```

## 貢献

新しいテストケースを追加する場合：

1. 適切なテストスクリプト（`test.sh`, `functional_test.sh`, `performance_test.sh`）に追加
2. テスト関数名は `test_` で始める
3. 成功時は `print_success`、失敗時は `print_error` を使用
4. このREADMEを更新

## ライセンス

このテストスイートは goto プログラムと同じライセンスの下で配布されます。
