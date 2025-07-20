# goto コマンド

ディレクトリ間を素早く移動するための`goto`コマンドです。

これは高速で依存関係のないディレクトリナビゲーションを提供するGo実装です。

- [English](README.md) / [中文](README-zh.md) / [한국어](README-ko.md) / [Español](README-es.md)

## クイックスタート

1. **ダウンロード** [リリースページ](https://github.com/kujirahand/goto/releases)から最新のバイナリをダウンロード
2. **実行可能にして**PATHに配置
3. **実行** `goto`でインタラクティブメニューを表示

## 主な機能

- **高速ディレクトリナビゲーション**: よく使用するディレクトリに瞬時にジャンプ
- **スマート履歴**: 最近使用した順に自動でソート
- **複数の入力方法**: 数字、ラベル、ショートカットキーを使用
- **タブ補完**: BashとZshの補完をサポート
- **クロスプラットフォーム**: Linux、macOS、Windowsで動作
- **多言語サポート**: 自動言語検出（英語、日本語、中国語、韓国語）
- **依存関係なし**: 外部依存関係なしの単一バイナリ

## インストール

以下の手順に従って`goto`コマンドをインストールしてください。

### プリビルトバイナリのダウンロード（推奨）

`goto`をインストールする最も簡単な方法は、GitHubリリースページからプリビルトバイナリをダウンロードすることです：

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
   
   # PATHのディレクトリに移動
   sudo mv goto-* /usr/local/bin/goto
   
   # またはローカルbinディレクトリを作成（存在しない場合）
   mkdir -p ~/bin
   mv goto-* ~/bin/goto
   export PATH="$PATH:$HOME/bin"  # これをシェル設定に追加
   ```

   **Windowsの場合**:
   - ダウンロードしたファイルを`goto.exe`にリネーム
   - PATHにあるディレクトリに配置、または新しいディレクトリを作成してPATHに追加

4. **インストールの確認**:

   ```sh
   goto --version
   ```

### ソースからクローンしてビルド

```sh
# リポジトリをクローン
git clone https://github.com/kujirahand/goto.git
# ビルド
cd goto
make
```

### PATHに追加

ビルド後、コンパイルされた`goto`実行ファイルをPATHに追加するために、シェル設定ファイル（`.bashrc`、`.zshrc`など）に以下の行を追加してください：

```sh
export PATH="$PATH:/path/to/goto"
```

例えば、ホームディレクトリにクローンした場合：

```sh
export PATH="$PATH:$HOME/goto"
```

PATHに追加後、シェル設定をリロードしてください：

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

### 手動タブ補完設定（プリビルトバイナリ用）

プリビルトバイナリをダウンロードした場合でも、タブ補完を手動で設定できます：

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

   **bash**（`~/.bashrc`または`~/.bash_profile`）:

   ```sh
   source ~/.bash_completion.d/goto-completion.bash
   ```

   **zsh**（`~/.zshrc`）:

   ```sh
   fpath=(~/.zsh/completions $fpath)
   autoload -U compinit && compinit
   ```

3. **シェルを再起動または設定をリロード**:

   ```sh
   source ~/.bashrc   # bashの場合
   source ~/.zshrc    # zshの場合
   ```

### 高度なインストール（タブ補完付きソースビルド）

ソースからビルドする際の最高の体験を得るには、バイナリと補完スクリプトの両方をインストールしてください：

```sh
# すべてをビルドしてインストール
make install-all
```

これにより以下が行われます：

1. `goto`バイナリを`/usr/local/bin/`にインストール
2. シェル補完スクリプトをインストール
3. 補完を有効にする手順を表示

#### 代替案：手動補完設定（ソースビルド）

ソースからビルドしたが補完を手動でインストールしたい場合：

1. 補完スクリプトをインストール：

   ```sh
   make install-completion
   ```

2. シェル設定に以下を追加：

   **bash**（`~/.bashrc`または`~/.bash_profile`）:

   ```sh
   source ~/.bash_completion.d/goto-completion.bash
   ```

   **zsh**（`~/.zshrc`）:

   ```sh
   fpath=(~/.zsh/completions $fpath)
   autoload -U compinit && compinit
   ```

3. シェルを再起動または設定をリロード：

   ```sh
   source ~/.bashrc   # bashの場合
   source ~/.zshrc    # zshの場合
   ```

#### タブ補完の使用

有効にすると、`goto`コマンドでタブ補完を使用できます：

```sh
goto <TAB>        # 利用可能なすべての移動先を表示
goto h<TAB>       # 'h'で始まるショートカットを補完
goto Home<TAB>    # 'Home'で始まるラベルを補完
goto 1<TAB>       # '1'で始まる番号の移動先を表示
```

## 設定

### 設定ファイル

`goto`コマンドは以下の設定ファイルを使用します：

- **`~/.goto.toml`**: 移動先を含むメイン設定ファイル
- **`~/.goto.history.json`**: 最近の使用情報を保存する履歴データ

初回実行時に、サンプルの移動先を含むデフォルト設定ファイルが自動的に作成されます。

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

各移動先には以下を設定できます：

- `path`（必須）: ディレクトリパス（ホームディレクトリには`~`をサポート）
- `shortcut`（オプション）: 単一文字のショートカットキー
- `command`（オプション）: ディレクトリ変更後に実行するコマンド

## 使用方法

### 基本的な使用方法

`goto`コマンドを実行して利用可能な移動先を表示：

```sh
goto
```

### コマンドライン引数

移動先をコマンドライン引数として直接指定することもできます：

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

これはスクリプト作成や正確な移動先がわかっている場合に便利です。

### インタラクティブモード

引数なしで実行すると、`goto`はインタラクティブメニューを表示します：

出力例：

```text
👉 利用可能なディレクトリ:
1. Home → /Users/username/ (shortcut: h)
2. Desktop → /Users/username/Desktop (shortcut: d)
3. Downloads → /Users/username/Downloads (shortcut: b)
4. MyProject → /Users/username/workspace/my-project (shortcut: p)

➕ [+] 現在のディレクトリを追加

番号、ショートカットキー、または[+]を入力して現在のディレクトリを追加:
番号、ショートカットキー、または[+]を入力:
```

以下の方法でナビゲートできます：

- **番号**: `1`、`2`、`3`などを入力
- **ショートカット**: `h`、`d`、`b`などを入力
- **現在を追加**: `+`を入力して現在のディレクトリを追加

### 現在のディレクトリを追加

メニューから`[+]`を選択して現在のディレクトリをgoto移動先に追加できます：

```sh
goto
# メニューから[+]を選択
# 現在のディレクトリのラベルを入力
# オプションでショートカットキーを入力
```

例：

```text
番号、ショートカットキー、または[+]を入力: +
📍 現在のディレクトリ: /Users/username/workspace/new-project
このディレクトリのラベルを入力: NewProject
ショートカットキーを入力（オプション、スキップするにはEnterを押下）: n
✅ 'NewProject' → /Users/username/workspace/new-project を追加しました
🔑 ショートカット: n
```

この機能により、よく使用するディレクトリをgotoリストに素早く追加できます。

### 新しいシェル機能

移動先を選択すると、`goto`はターゲットディレクトリで新しいシェルセッションを開きます。これは以下を意味します：

- 現在のシェルセッションは変更されない
- 新しい場所で新鮮なシェル環境を取得
- `exit`と入力して前のシェルに戻る
- 設定で`command`が指定されている場合、自動的に実行される

### 使用履歴

`goto`は使用履歴を自動的に追跡し、最近使用した順に移動先を表示します。これにより、頻繁にアクセスするディレクトリがインタラクティブメニューの上部に表示されます。

#### 使用履歴の表示

最近の使用履歴を表示するには：

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

#### 履歴の仕組み

- **自動追跡**: 移動先にナビゲートするたびにタイムスタンプが記録される
- **スマートソート**: インタラクティブモードで、移動先は最近使用した順にソートされる
- **永続ストレージ**: 履歴は`~/.goto.history.json`ファイルに保存される
- **手動メンテナンス不要**: 履歴は自動的に更新される - 手動管理は不要

#### 履歴の保存

使用履歴は`~/.goto.history.json`ファイルに以下の形式で保存されます：

```json
{
  "entries": [
    {
      "label": "Home",
      "last_used": "2025-07-18T16:08:38+09:00"
    },
    {
      "label": "Desktop",
      "last_used": "2025-07-18T16:04:40+09:00"
    }
  ]
}
```

このインテリジェントな順序により、最も頻繁に使用するディレクトリが常に簡単にアクセスできるようになります。

## 多言語サポート

`goto`はシステム言語を自動検出し、好みの言語でメッセージを表示します。現在サポートされている言語：

- **English** (en) - デフォルト
- **Japanese** (ja) - 日本語
- **Chinese** (zh) - 中文
- **Korean** (ko) - 한국어
- **Spanish** (es) - Español

### 言語検出の仕組み

アプリケーションは以下の環境変数を順番にチェックしてシステム言語を自動検出します：

1. `LANG`
2. `LANGUAGE`
3. `LC_ALL`
4. `LC_MESSAGES`

例えば、システムが日本語に設定されている場合（`LANG=ja_JP.UTF-8`）、`goto`は自動的にすべてのメッセージを日本語で表示します。

### 異なる言語での出力例

**英語:**

```text
🚀 goto - Navigate directories quickly
👉 Available destinations:
1. Home → /Users/username/ (shortcut: h)
📈 Recent usage history:
```

**日本語:**

```text
🚀 goto - ディレクトリ間を素早く移動
👉 利用可能なディレクトリ:
1. Home → /Users/username/ (shortcut: h)
📈 最近の使用履歴:
```

**中国語:**

```text
🚀 goto - 快速导航目录
👉 可用目录:
1. Home → /Users/username/ (shortcut: h)
📈 最近使用历史:
```

**韓国語:**

```text
🚀 goto - 디렉토리 빠른 탐색
👉 사용 가능한 디렉토리:
1. Home → /Users/username/ (shortcut: h)
📈 최근 사용 기록:
```

**スペイン語:**

```text
🚀 goto - Navegar directorios rápidamente
👉 Destinos disponibles:
1. Home → /Users/username/ (shortcut: h)
📈 Historial de uso reciente:
```

### 言語オーバーライド

システム設定に関係なく特定の言語を使用したい場合は、`LANG`環境変数を設定できます：

```sh
# 日本語インターフェースを使用
LANG=ja_JP.UTF-8 goto

# 英語インターフェースを使用
LANG=en_US.UTF-8 goto

# 中国語インターフェースを使用
LANG=zh_CN.UTF-8 goto

# 韓国語インターフェースを使用
LANG=ko_KR.UTF-8 goto

# スペイン語インターフェースを使用
LANG=es_ES.UTF-8 goto
```

### サポート言語

多言語サポートは以下を含むすべてのユーザーインターフェース要素をカバーします：

- インタラクティブメニューメッセージ
- ナビゲーション確認
- エラーメッセージ
- ヘルプテキスト
- 履歴表示
- 設定メッセージ

すべてのメッセージはシステム言語設定に基づいて自動的にローカライズされ、国際ユーザーにネイティブな体験を提供します。

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
