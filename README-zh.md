# goto 命令

用于快速导航目录的 `goto` 命令。

这是一个 Go 实现，提供快速、无依赖的目录导航。

- [English](README.md) / [日本語](README-ja.md) / [한국어](README-ko.md) / [Español](README-es.md)

## 快速开始

### 使用 Homebrew 安装（推荐）

您可以在 macOS 或 Linux 上使用 [Homebrew](https://brew.sh/) 轻松安装 `goto`：

```sh
brew tap kujirahand/goto
brew install kujirahand/goto/goto
```

### 手动安装

1. **下载** 从 [Releases](https://github.com/kujirahand/goto/releases) 下载适用于您平台的最新二进制文件
2. **设置执行权限** 并将其放在 PATH 中
3. **运行** `goto` 查看交互式菜单

## 主要功能

- **快速目录导航**: 立即跳转到常用目录
- **智能历史**: 自动按最近使用的顺序排序目标
- **多种输入方式**: 使用数字、标签或快捷键
- **Tab 补全**: 支持 Bash 和 Zsh 补全
- **跨平台**: 在 Linux、macOS 和 Windows 上运行
- **多语言支持**: 自动语言检测（英语、日语、中文、韩语、西班牙语）
- **零依赖**: 无外部依赖的单一二进制文件

## 安装

请按照以下步骤安装 `goto` 命令。

### 下载预构建二进制文件（推荐）

安装 `goto` 最简单的方法是从 GitHub 发布页面下载预构建的二进制文件：

1. **访问发布页面**: <https://github.com/kujirahand/goto/releases>
2. **下载适用于您平台的二进制文件**:
   - **Linux amd64**: `goto-linux-amd64`
   - **Linux arm64**: `goto-linux-arm64`
   - **macOS Intel**: `goto-darwin-amd64`
   - **macOS Apple Silicon**: `goto-darwin-arm64`
   - **Windows amd64**: `goto-windows-amd64.exe`
   - **Windows arm64**: `goto-windows-arm64.exe`

3. **设置执行权限并放在 PATH 中**:

   **对于 Linux/macOS**:

   ```sh
   # 解压 zip 文件
   unzip goto-*.zip
   # 使二进制文件可执行
   chmod +x goto-*   
   # 移动到 PATH 中的目录
   sudo mv goto-* /usr/local/bin/goto
   ```

   **对于 Windows**:
   - 将下载的文件重命名为 `goto.exe`
   - 将其放在 PATH 中包含的目录中，或创建新目录并将其添加到 PATH

4. **验证安装**:

   ```sh
   goto --version
   ```

### 从源码克隆和构建

```sh
# 克隆仓库
git clone https://github.com/kujirahand/goto.git
# 构建
cd goto
make
```

### 构建发布存档（开发者用）

要为所有平台创建发布存档：

```sh
# 为所有平台创建 ZIP 存档（二进制文件会自动清理）
make build-release-zip
```

### 添加到 PATH

构建后，通过将以下行添加到您的 shell 配置文件（`.bashrc`、`.zshrc` 等）来将编译的 `goto` 可执行文件添加到 PATH：

```sh
export PATH="$PATH:/path/to/goto"
```

例如，如果您克隆到了主目录：

```sh
export PATH="$PATH:$HOME/goto"
```

添加到 PATH 后，重新加载您的 shell 配置：

```sh
# 对于 zsh
source ~/.zshrc

# 对于 bash
source ~/.bashrc
```

### 带 Tab 补全的安装（源码构建）

如果您从源码构建，可以安装二进制文件和补全脚本：

```sh
# 构建并安装所有内容（需要源代码）
make install-all
```

### 手动 Tab 补全设置（预构建二进制文件）

如果您下载了预构建的二进制文件，仍然可以手动设置 tab 补全：

1. **下载补全脚本**:

   ```sh
   # 创建补全目录
   mkdir -p ~/.bash_completion.d ~/.zsh/completions
   
   # 下载 bash 补全脚本
   curl -o ~/.bash_completion.d/goto-completion.bash \
     https://raw.githubusercontent.com/kujirahand/goto/main/completion/goto-completion.bash
   
   # 下载 zsh 补全脚本
   curl -o ~/.zsh/completions/_goto \
     https://raw.githubusercontent.com/kujirahand/goto/main/completion/_goto
   ```

2. **添加到您的 shell 配置**:

   **对于 bash** (`~/.bashrc` 或 `~/.bash_profile`):

   ```sh
   source ~/.bash_completion.d/goto-completion.bash
   ```

   **对于 zsh** (`~/.zshrc`):

   ```sh
   fpath=(~/.zsh/completions $fpath)
   autoload -U compinit && compinit
   ```

3. **重启您的 shell 或重新加载配置**:

   ```sh
   source ~/.bashrc   # 用于 bash
   source ~/.zshrc    # 用于 zsh
   ```

### 带 Tab 补全的高级安装（源码构建）

从源码构建时获得最佳体验，安装二进制文件和补全脚本：

```sh
# 构建并安装所有内容
make install-all
```

这将：

1. 将 `goto` 二进制文件安装到 `/usr/local/bin/`
2. 安装 shell 补全脚本
3. 显示启用补全的说明

#### 替代方案：手动补全设置（源码构建）

如果您从源码构建但希望手动安装补全：

1. 安装补全脚本：

   ```sh
   make install-completion
   ```

2. 将以下内容添加到您的 shell 配置：

   **对于 bash** (`~/.bashrc` 或 `~/.bash_profile`):

   ```sh
   source ~/.bash_completion.d/goto-completion.bash
   ```

   **对于 zsh** (`~/.zshrc`):

   ```sh
   fpath=(~/.zsh/completions $fpath)
   autoload -U compinit && compinit
   ```

3. 重启您的 shell 或重新加载配置：

   ```sh
   source ~/.bashrc   # 用于 bash
   source ~/.zshrc    # 用于 zsh
   ```

#### 使用 Tab 补全

启用后，您可以在 `goto` 命令中使用 tab 补全：

```sh
goto <TAB>        # 显示所有可用目标
goto h<TAB>       # 补全以 'h' 开头的快捷键
goto Home<TAB>    # 补全以 'Home' 开头的标签
goto 1<TAB>       # 显示以 '1' 开头的数字目标
```

## 配置

### 配置文件

`goto` 命令使用以下配置文件：

- **`~/.goto.toml`**: 包含您目标的主配置文件
- **`~/.goto.history.json`**: 存储您最近使用信息的历史数据

首次运行时，`goto` 将自动创建包含示例目标的默认配置文件。

配置示例：

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

每个目标可以有：

- `path`（必需）: 目录路径（支持主目录的 `~`）
- `shortcut`（可选）: 单个字符快捷键
- `command`（可选）: 更改目录后执行的命令

### 注意：包含点的条目

当 TOML 文件中的条目包含点（`.`）时，其含义可能会改变。为了防止这种情况，请用双引号包围条目，如下所示：

```toml
["kujirahand.com"]
path = https://kujirahand.com
```

## 使用方法

### 基本使用

运行 `goto` 命令查看可用目标：

```sh
goto
```

### 命令行参数

您也可以直接将目标指定为命令行参数：

```sh
# 使用数字
goto 1
goto 4

# 使用标签名
goto Home
goto MyProject

# 使用快捷键
goto h
goto p

# 查看使用历史
goto --history

# 显示帮助
goto --help

# 显示版本
goto --version
```

这对于脚本编写或当您确切知道要去哪里时很有用。

### 交互模式

不带参数运行时，`goto` 显示交互式菜单：

示例输出：

```text
👉 可用目标:
1. Home → /Users/username/ (shortcut: h)
2. Desktop → /Users/username/Desktop (shortcut: d)
3. Downloads → /Users/username/Downloads (shortcut: b)
4. MyProject → /Users/username/workspace/my-project (shortcut: p)

➕ [+] 添加当前目录

请输入数字、快捷键或 [+] 添加当前目录:
输入数字、快捷键或 [+]:
```

您可以通过以下方式导航：

- **数字**: 输入 `1`、`2`、`3` 等
- **快捷键**: 输入 `h`、`d`、`b` 等
- **添加当前**: 输入 `+` 添加当前目录

### 添加当前目录

您可以通过选择 `[+]` 将当前目录添加到您的 goto 目标：

```sh
goto
# 从菜单中选择 [+]
# 为当前目录输入标签
# 可选择输入快捷键
```

示例：

```text
输入数字、快捷键或 [+]: +
📍 当前目录: /Users/username/workspace/new-project
为此目录输入标签: NewProject
输入快捷键（可选，按 Enter 跳过）: n
✅ 已添加 'NewProject' → /Users/username/workspace/new-project
🔑 快捷键: n
```

此功能允许您快速将常用目录添加到您的 goto 列表。

### 新 Shell 功能

当您选择目标时，`goto` 在目标目录中打开新的 shell 会话。这意味着：

- 您当前的 shell 会话保持不变
- 您在新位置获得新的 shell 环境
- 输入 `exit` 返回到您之前的 shell
- 如果在配置中指定了 `command`，它将自动执行

### 使用历史

`goto` 自动跟踪使用历史并按最近使用的顺序显示目标。这使得经常访问的目录出现在交互式菜单的顶部。

#### 查看使用历史

您可以查看最近的使用历史：

```sh
goto --history
```

示例输出：

```text
📈 最近使用历史:
==================================================
 1. Home → /Users/username
    📅 2025-07-18 16:08:38

 2. Desktop → /Users/username/Desktop
    📅 2025-07-18 16:04:40

 3. MyProject → /Users/username/workspace/my-project
    📅 2025-07-18 15:30:15
```

#### 历史工作原理

- **自动跟踪**: 每次导航到目标时，都会记录时间戳
- **智能排序**: 在交互模式下，目标按最近使用的顺序排序
- **持久存储**: 历史存储在 `~/.goto.history.json` 文件中
- **无需手动维护**: 历史自动更新 - 无需手动管理

#### 历史存储

使用历史存储在您的 `~/.goto.history.json` 文件中，格式如下：

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

这种智能排序确保您最常用的目录始终易于访问。

## 多语言支持

`goto` 自动检测您的系统语言并以您首选的语言显示消息。目前支持的语言：

- **英语** (en) - 默认
- **日语** (ja) - 日本語
- **中文** (zh) - 中文
- **韩语** (ko) - 한국어
- **西班牙语** (es) - Español

### 语言检测工作原理

应用程序通过按顺序检查以下环境变量来自动检测您的系统语言：

1. `LANG`
2. `LANGUAGE`
3. `LC_ALL`
4. `LC_MESSAGES`

例如，如果您的系统设置为中文（`LANG=zh_CN.UTF-8`），`goto` 将自动以中文显示所有消息。

### 不同语言的示例输出

**英语:**

```text
🚀 goto - Navigate directories quickly
👉 Available destinations:
1. Home → /Users/username/ (shortcut: h)
📈 Recent usage history:
```

**日语:**

```text
🚀 goto - ディレクトリ間を素早く移動
👉 利用可能なディレクトリ:
1. Home → /Users/username/ (shortcut: h)
📈 最近の使用履歴:
```

**中文:**

```text
🚀 goto - 快速导航目录
👉 可用目录:
1. Home → /Users/username/ (shortcut: h)
📈 最近使用历史:
```

**韩语:**

```text
🚀 goto - 디렉토리 빠른 탐색
👉 사용 가능한 디렉토리:
1. Home → /Users/username/ (shortcut: h)
📈 최근 사용 기록:
```

**西班牙语:**

```text
🚀 goto - Navegar directorios rápidamente
👉 Destinos disponibles:
1. Home → /Users/username/ (shortcut: h)
📈 Historial de uso reciente:
```

### 语言覆盖

如果您想使用特定语言而不考虑系统设置，可以设置 `LANG` 环境变量：

```sh
# 使用中文界面
LANG=zh_CN.UTF-8 goto

# 使用英文界面
LANG=en_US.UTF-8 goto

# 使用日文界面
LANG=ja_JP.UTF-8 goto

# 使用韩文界面
LANG=ko_KR.UTF-8 goto

# 使用西班牙语界面
LANG=es_ES.UTF-8 goto
```

### 支持的语言

多语言支持涵盖所有用户界面元素，包括：

- 交互式菜单消息
- 导航确认
- 错误消息
- 帮助文本
- 历史显示
- 配置消息

所有消息都会根据您的系统语言设置自动本地化，为国际用户提供原生体验。

### 示例

1. **使用命令行参数导航（数字）:**

   ```sh
   goto 1
   goto 4
   ```

2. **使用命令行参数导航（标签）:**

   ```sh
   goto Home
   goto MyProject
   ```

3. **使用命令行参数导航（快捷键）:**

   ```sh
   goto h
   goto p
   ```

4. **交互式导航:**

   ```sh
   goto
   # 然后输入: h（快捷键）、1（数字）或 Home（标签）
   ```

5. **添加当前目录:**

   ```sh
   cd /path/to/important/project
   goto
   # 输入: +
   # 标签: ImportantProject
   # 快捷键: i
   ```

6. **查看使用历史:**

   ```sh
   goto --history
   ```
