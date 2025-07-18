# goto コマンド

`goto` コマンドは、ディレクトリ間を素早く移動するためのツールです。

これは、高速で依存関係のないディレクトリナビゲーションを提供するGo実装です。

## クイックスタート

1. **最新バイナリをダウンロード** - [リリースページ](https://github.com/kujirahand/goto/releases)から、お使いのプラットフォーム用のバイナリをダウンロードします
2. **実行権限を付与してPATHに配置** - ダウンロードしたファイルを実行可能にし、PATH上のディレクトリに配置します
3. **実行** - `goto` を実行してインタラクティブメニューを表示します

## 主な機能

- **高速ディレクトリナビゲーション**: よく使うディレクトリに即座にジャンプ
- **スマート履歴**: 最近使用した順に自動的にディレクトリを並び替え
- **複数の入力方法**: 番号、ラベル、ショートカットキーを使用可能
- **タブ補完**: BashとZshのタブ補完をサポート
- **クロスプラットフォーム**: Linux、macOS、Windowsで動作
- **多言語対応**: 自動言語検出（英語、日本語、中国語、韓国語）
- **依存関係なし**: 外部依存関係のない単一バイナリ

## インストール

以下の手順に従って `goto` コマンドをインストールしてください。

### 事前ビルド済みバイナリのダウンロード（推奨）

最も簡単なインストール方法は、GitHubのリリースページから事前ビルド済みバイナリをダウンロードすることです：

1. **リリースページにアクセス**: <https://github.com/kujirahand/goto/releases>
2. **プラットフォーム用のバイナリをダウンロード**:
   - **Linux amd64**: `goto-linux-amd64`
   - **Linux arm64**: `goto-linux-arm64`
   - **macOS Intel**: `goto-darwin-amd64`
   - **macOS Apple Silicon**: `goto-darwin-arm64`
   - **Windows amd64**: `goto-windows-amd64.exe`
   - **Windows arm64**: `goto-windows-arm64.exe`

3. **実行可能にしてPATHに配置**:

   **Linux/macOSの場合**:

   ```sh
   # ダウンロードして実行可能にする
   chmod +x goto-*
   
   # PATH上のディレクトリに移動
   sudo mv goto-* /usr/local/bin/goto
   
   # またはローカルbinディレクトリを作成（存在しない場合）
   mkdir -p ~/bin
   mv goto-* ~/bin/goto
   export PATH="$PATH:$HOME/bin"  # これをシェル設定に追加
   ```

   **Windowsの場合**:
   - ダウンロードしたファイルを `goto.exe` にリネーム
   - PATH上のディレクトリに配置するか、新しいディレクトリを作成してPATHに追加

4. **インストールの確認**:

   ```sh
   goto --version
   ```

### ソースコードからのクローンとビルド

```sh
# リポジトリをクローン
git clone https://github.com/kujirahand/goto.git
# ビルド
cd goto
make
```

### PATHへの追加

シェル設定ファイル（`.bashrc`、`.zshrc`など）に以下の行を追加して、`goto/go` ディレクトリをPATHに追加してください：

```sh
export PATH="$PATH:/path/to/goto/go"
```

例：ホームディレクトリにクローンした場合：

```sh
export PATH="$PATH:$HOME/goto/go"
```

PATHに追加後、シェル設定を再読み込みしてください：

```sh
# zshの場合
source ~/.zshrc

# bashの場合
source ~/.bashrc
```

### タブ補完付きインストール（ソースビルド）

ソースからビルドした場合、バイナリと補完スクリプトの両方をインストールできます：

```sh
# すべてをビルドしてインストール（ソースコードが必要）
make install-all
```

### 手動タブ補完設定（事前ビルド済みバイナリ用）

事前ビルド済みバイナリをダウンロードした場合でも、タブ補完を手動で設定できます：

1. **補完スクリプトをダウンロード**:

   ```sh
   # 補完ディレクトリを作成
   mkdir -p ~/.bash_completion.d ~/.zsh/completions
   
   # bash補完スクリプトをダウンロード
   curl -o ~/.bash_completion.d/goto-completion.bash \
     https://raw.githubusercontent.com/kujirahand/goto/main/completion/goto-completion.bash
   
   # zsh補完スクリプトをダウンロード
   curl -o ~/.zsh/completions/_goto \
     https://raw.githubusercontent.com/kujirahand/goto/main/completion/_goto
   ```

2. **シェル設定に追加**:

   **bash用** (`~/.bashrc` または `~/.bash_profile`):

   ```sh
   source ~/.bash_completion.d/goto-completion.bash
   ```

   **zsh用** (`~/.zshrc`):

   ```sh
   fpath=(~/.zsh/completions $fpath)
   autoload -U compinit && compinit
   ```

3. **シェルを再起動または設定を再読み込み**:

   ```sh
   source ~/.bashrc   # bashの場合
   source ~/.zshrc    # zshの場合
   ```

### タブ補完付き高度なインストール（ソースビルド）

ソースからビルドする場合の最適な体験のために、バイナリと補完スクリプトの両方をインストールします：

```sh
# すべてをビルドしてインストール
make install-all
```

これにより以下が実行されます：

1. `goto` バイナリを `/usr/local/bin/` にインストール
2. シェル補完スクリプトをインストール
3. 補完を有効にする手順を表示

#### 代替案：手動補完設定（ソースビルド）

ソースからビルドしたが、補完を手動でインストールしたい場合：

1. 補完スクリプトをインストール：

   ```sh
   make install-completion
   ```

2. シェル設定に以下を追加：

   **bash用** (`~/.bashrc` または `~/.bash_profile`):

   ```sh
   source ~/.bash_completion.d/goto-completion.bash
   ```

   **zsh用** (`~/.zshrc`):

   ```sh
   fpath=(~/.zsh/completions $fpath)
   autoload -U compinit && compinit
   ```

3. シェルを再起動または設定を再読み込み：

   ```sh
   source ~/.bashrc   # bashの場合
   source ~/.zshrc    # zshの場合
   ```

#### タブ補完の使用方法

有効にすると、`goto` コマンドでタブ補完を使用できます：

```sh
goto <TAB>        # 利用可能なすべてのディレクトリを表示
goto h<TAB>       # 'h' で始まるショートカットを補完
goto Home<TAB>    # 'Home' で始まるラベルを補完
goto 1<TAB>       # '1' で始まる番号のディレクトリを表示
```

## 設定

### 設定ファイル - `~/.goto.toml`

`goto` コマンドは、`~/.goto.toml` にあるTOML設定ファイルを使用します。初回実行時に、サンプルディレクトリを含むデフォルト設定ファイルが自動的に作成されます。

設定例：

```toml
[Home]
path = "~/"
shortcut = "h"

[Desktop]
path = "~/Desktop"
shortcut = "d"

[Downloads]
path = "~/Downloads"
shortcut = "b"

[MyProject]
path = "~/workspace/my-project"
shortcut = "p"
command = "ls -la && git status"
```

各ディレクトリには以下を設定できます：

- `path`（必須）: ディレクトリパス（ホームディレクトリには `~` を使用可能）
- `shortcut`（任意）: 1文字のショートカットキー
- `command`（任意）: ディレクトリ変更後に実行するコマンド

## 使用方法

### 基本的な使用方法

`goto` コマンドを実行して、利用可能なディレクトリを確認：

```sh
goto
```

### コマンドライン引数

コマンドライン引数として直接ディレクトリを指定することも可能：

```sh
# 番号を使用
goto 1
goto 4

# ラベル名を使用
goto Home
goto MyProject

# ショートカットキーを使用
goto h
goto p

# 使用履歴を表示
goto --history

# ヘルプを表示
goto --help

# バージョンを表示
goto --version
```

これは、スクリプトや移動先が明確な場合に便利です。

### インタラクティブモード

引数なしで実行すると、`goto` はインタラクティブメニューを表示します：

出力例：

```text
👉 利用可能なディレクトリ:
1. Home → /Users/username/ (ショートカット: h)
2. Desktop → /Users/username/Desktop (ショートカット: d)
3. Downloads → /Users/username/Downloads (ショートカット: b)
4. MyProject → /Users/username/workspace/my-project (ショートカット: p)

➕ [+] 現在のディレクトリを追加

番号、ショートカットキー、または [+] を入力してください:
番号、ショートカットキー、または [+] を入力:
```

以下の方法で移動できます：

- **番号**: `1`、`2`、`3` などを入力
- **ショートカット**: `h`、`d`、`b` などを入力
- **現在のディレクトリを追加**: `+` を入力

### 現在のディレクトリの追加

メニューから `[+]` を選択して、現在のディレクトリをgotoディレクトリリストに追加できます：

```sh
goto
# メニューから [+] を選択
# 現在のディレクトリのラベルを入力
# 任意でショートカットキーを入力
```

例：

```text
番号、ショートカットキー、または [+] を入力: +
📍 現在のディレクトリ: /Users/username/workspace/new-project
このディレクトリのラベルを入力してください: NewProject
ショートカットキーを入力してください（任意、Enterでスキップ）: n
✅ 'NewProject' → /Users/username/workspace/new-project を追加しました
🔑 ショートカット: n
```

この機能により、よく使うディレクトリをgotoリストに素早く追加できます。

### 新しいシェル機能

ディレクトリを選択すると、`goto` は対象ディレクトリで新しいシェルセッションを開きます。これは以下を意味します：

- 現在のシェルセッションは変更されません
- 新しい場所で新しいシェル環境を取得
- 前のシェルに戻るには `exit` を入力
- 設定で `command` が指定されている場合、自動的に実行されます

### 使用履歴

`goto` は自動的に使用履歴を追跡し、最近使用した順にディレクトリを表示します。これにより、よくアクセスするディレクトリがインタラクティブメニューの上位に表示されます。

#### 使用履歴の表示

以下のコマンドで最近の使用履歴を表示できます：

```sh
goto --history
```

出力例：

```text
📈 最近の使用履歴:
==================================================
 1. Home → /Users/username
    📅 2025-07-18 16:08:38

 2. Desktop → /Users/username/Desktop
    📅 2025-07-18 16:04:40

 3. MyProject → /Users/username/workspace/my-project
    📅 2025-07-18 15:30:15
```

#### 履歴の動作原理

- **自動追跡**: ディレクトリに移動するたびに、タイムスタンプが記録されます
- **スマートソート**: インタラクティブモードでは、最近使用した順にディレクトリが並びます
- **永続化**: 履歴は `~/.goto.toml` 設定ファイルに保存されます
- **メンテナンス不要**: 履歴は自動的に更新されるため、手動管理は不要です

#### 履歴の保存

使用履歴は `~/.goto.toml` ファイルに以下の形式で保存されます：

```toml
[[history]]
  label = "Home"
  last_used = "2025-07-18T16:08:38+09:00"

[[history]]
  label = "Desktop"
  last_used = "2025-07-18T16:04:40+09:00"

# ... あなたのディレクトリ設定 ...
[Home]
path = "~/"
shortcut = "h"

[Desktop]
path = "~/Desktop"
shortcut = "d"
```

このインテリジェントな順序付けにより、最もよく使うディレクトリが常に簡単にアクセスできるようになります。

### 例

1. **コマンドライン引数を使用したナビゲーション（番号）:**

   ```sh
   goto 1
   goto 4
   ```

2. **コマンドライン引数を使用したナビゲーション（ラベル）:**

   ```sh
   goto Home
   goto MyProject
   ```

3. **コマンドライン引数を使用したナビゲーション（ショートカット）:**

   ```sh
   goto h
   goto p
   ```

4. **インタラクティブナビゲーション:**

   ```sh
   goto
   # その後入力: h（ショートカット）、1（番号）、またはHome（ラベル）
   ```

5. **現在のディレクトリを追加:**

   ```sh
   cd /path/to/important/project
   goto
   # 入力: +
   # ラベル: ImportantProject
   # ショートカット: i
   ```

6. **使用履歴を表示:**

   ```sh
   goto --history
   ```

## 多言語対応

`goto` は自動的にシステムの言語を検出し、お好みの言語でメッセージを表示します。現在サポートされている言語：

- **English** (en) - 英語（デフォルト）
- **Japanese** (ja) - 日本語
- **Chinese** (zh) - 中文
- **Korean** (ko) - 한국어

### 言語検出の仕組み

アプリケーションは以下の環境変数を順番にチェックして、システムの言語を自動的に検出します：

1. `LANG`
2. `LANGUAGE`
3. `LC_ALL`
4. `LC_MESSAGES`

例えば、システムが日本語に設定されている場合（`LANG=ja_JP.UTF-8`）、`goto` は自動的にすべてのメッセージを日本語で表示します。

### 異なる言語での出力例

**English:**

```text
🚀 goto - Navigate directories quickly
👉 Available destinations:
1. Home → /Users/username/ (shortcut: h)
📈 Recent usage history:
```

**Japanese:**

```text
🚀 goto - ディレクトリ間を素早く移動
👉 利用可能なディレクトリ:
1. Home → /Users/username/ (shortcut: h)
📈 最近の使用履歴:
```

**Chinese:**

```text
🚀 goto - 快速导航目录
👉 可用目录:
1. Home → /Users/username/ (shortcut: h)
📈 最近使用历史:
```

**Korean:**

```text
🚀 goto - 디렉토리 빠른 탐색
👉 사용 가능한 디렉토리:
1. Home → /Users/username/ (shortcut: h)
📈 최근 사용 기록:
```

### 言語の上書き

システム設定に関係なく特定の言語を使用したい場合は、`LANG` 環境変数を設定できます：

```sh
# 日本語インターフェースを使用
LANG=ja_JP.UTF-8 goto

# 英語インターフェースを使用
LANG=en_US.UTF-8 goto

# 中国語インターフェースを使用
LANG=zh_CN.UTF-8 goto

# 韓国語インターフェースを使用
LANG=ko_KR.UTF-8 goto
```

### サポートされている言語

多言語対応は以下を含むすべてのユーザーインターフェース要素をカバーしています：

- インタラクティブメニューメッセージ
- ナビゲーション確認
- エラーメッセージ
- ヘルプテキスト
- 履歴表示
- 設定メッセージ

すべてのメッセージは、システムの言語設定に基づいて自動的にローカライズされ、国際的なユーザーにネイティブな体験を提供します。
