# goto 명령어

`goto` 명령어는 디렉토리를 빠르게 탐색하기 위한 도구입니다.

이는 빠르고 의존성이 없는 디렉토리 탐색을 제공하는 Go 구현입니다.

## 빠른 시작

1. **다운로드** [Releases](https://github.com/kujirahand/goto/releases)에서 플랫폼에 맞는 최신 바이너리를 다운로드합니다
2. **실행 권한 설정** 후 PATH에 배치합니다
3. **실행** `goto`를 실행하여 대화형 메뉴를 확인합니다

## 주요 기능

- **빠른 디렉토리 탐색**: 자주 사용하는 디렉토리로 즉시 이동
- **스마트 히스토리**: 최근 사용 순서로 목적지를 자동 정렬
- **다양한 입력 방법**: 숫자, 라벨 또는 단축키 사용
- **Tab 자동완성**: Bash와 Zsh 자동완성 지원
- **크로스 플랫폼**: Linux, macOS, Windows에서 작동
- **다국어 지원**: 자동 언어 감지(영어, 일본어, 중국어, 한국어)
- **제로 의존성**: 외부 의존성이 없는 단일 바이너리

## 설치

다음 단계에 따라 `goto` 명령어를 설치하세요.

### 사전 빌드된 바이너리 다운로드 (권장)

`goto`를 설치하는 가장 쉬운 방법은 GitHub 릴리스 페이지에서 사전 빌드된 바이너리를 다운로드하는 것입니다:

1. **릴리스 페이지 방문**: <https://github.com/kujirahand/goto/releases>
2. **플랫폼용 바이너리 다운로드**:
   - **Linux amd64**: `goto-linux-amd64`
   - **Linux arm64**: `goto-linux-arm64`
   - **macOS Intel**: `goto-darwin-amd64`
   - **macOS Apple Silicon**: `goto-darwin-arm64`
   - **Windows amd64**: `goto-windows-amd64.exe`
   - **Windows arm64**: `goto-windows-arm64.exe`

3. **실행 권한 설정 후 PATH에 배치**:

   **Linux/macOS의 경우**:

   ```sh
   # 다운로드 후 실행 권한 설정
   chmod +x goto-*
   
   # PATH의 디렉토리로 이동
   sudo mv goto-* /usr/local/bin/goto
   
   # 또는 로컬 bin 디렉토리 생성 (존재하지 않는 경우)
   mkdir -p ~/bin
   mv goto-* ~/bin/goto
   export PATH="$PATH:$HOME/bin"  # 이를 셸 구성에 추가
   ```

   **Windows의 경우**:
   - 다운로드한 파일을 `goto.exe`로 이름 변경
   - PATH에 있는 디렉토리에 배치하거나 새 디렉토리를 생성하여 PATH에 추가

4. **설치 확인**:

   ```sh
   goto --version
   ```

### 소스에서 클론 및 빌드

```sh
# 저장소 클론
git clone https://github.com/kujirahand/goto.git
# 빌드
cd goto
make
```

### PATH에 추가

셸 구성 파일(`.bashrc`, `.zshrc` 등)에 다음 줄을 추가하여 빌드된 `goto` 실행 파일을 PATH에 추가하세요:

```sh
export PATH="$PATH:/path/to/goto"
```

예를 들어, 홈 디렉토리에 클론한 경우:

```sh
export PATH="$PATH:$HOME/goto"
```

PATH에 추가한 후 셸 구성을 다시 로드하세요:

```sh
# zsh의 경우
source ~/.zshrc

# bash의 경우
source ~/.bashrc
```

### Tab 자동완성과 함께 설치 (소스 빌드)

소스에서 빌드한 경우 바이너리와 자동완성 스크립트를 모두 설치할 수 있습니다:

```sh
# 모든 것을 빌드하고 설치 (소스 코드 필요)
make install-all
```

### 수동 Tab 자동완성 설정 (사전 빌드된 바이너리)

사전 빌드된 바이너리를 다운로드한 경우에도 tab 자동완성을 수동으로 설정할 수 있습니다:

1. **자동완성 스크립트 다운로드**:

   ```sh
   # 자동완성 디렉토리 생성
   mkdir -p ~/.bash_completion.d ~/.zsh/completions
   
   # bash 자동완성 스크립트 다운로드
   curl -o ~/.bash_completion.d/goto-completion.bash \
     https://raw.githubusercontent.com/kujirahand/goto/main/completion/goto-completion.bash
   
   # zsh 자동완성 스크립트 다운로드
   curl -o ~/.zsh/completions/_goto \
     https://raw.githubusercontent.com/kujirahand/goto/main/completion/_goto
   ```

2. **셸 구성에 추가**:

   **bash의 경우** (`~/.bashrc` 또는 `~/.bash_profile`):

   ```sh
   source ~/.bash_completion.d/goto-completion.bash
   ```

   **zsh의 경우** (`~/.zshrc`):

   ```sh
   fpath=(~/.zsh/completions $fpath)
   autoload -U compinit && compinit
   ```

3. **셸 재시작 또는 구성 다시 로드**:

   ```sh
   source ~/.bashrc   # bash의 경우
   source ~/.zshrc    # zsh의 경우
   ```

### Tab 자동완성을 통한 고급 설치 (소스 빌드)

소스에서 빌드할 때 최상의 경험을 위해 바이너리와 자동완성 스크립트를 모두 설치하세요:

```sh
# 모든 것을 빌드하고 설치
make install-all
```

이렇게 하면:

1. `goto` 바이너리를 `/usr/local/bin/`에 설치
2. 셸 자동완성 스크립트 설치
3. 자동완성 활성화 지침 표시

#### 대안: 수동 자동완성 설정 (소스 빌드)

소스에서 빌드했지만 자동완성을 수동으로 설치하고 싶다면:

1. 자동완성 스크립트 설치:

   ```sh
   make install-completion
   ```

2. 셸 구성에 다음 내용 추가:

   **bash의 경우** (`~/.bashrc` 또는 `~/.bash_profile`):

   ```sh
   source ~/.bash_completion.d/goto-completion.bash
   ```

   **zsh의 경우** (`~/.zshrc`):

   ```sh
   fpath=(~/.zsh/completions $fpath)
   autoload -U compinit && compinit
   ```

3. 셸 재시작 또는 구성 다시 로드:

   ```sh
   source ~/.bashrc   # bash의 경우
   source ~/.zshrc    # zsh의 경우
   ```

#### Tab 자동완성 사용하기

활성화된 후 `goto` 명령어에서 tab 자동완성을 사용할 수 있습니다:

```sh
goto <TAB>        # 모든 사용 가능한 목적지 표시
goto h<TAB>       # 'h'로 시작하는 단축키 자동완성
goto Home<TAB>    # 'Home'으로 시작하는 라벨 자동완성
goto 1<TAB>       # '1'로 시작하는 숫자 목적지 표시
```

## 구성

### 구성 파일 - `~/.goto.toml`

`goto` 명령어는 `~/.goto.toml`에 위치한 TOML 구성 파일을 사용합니다. 처음 `goto`를 실행할 때 샘플 목적지가 포함된 기본 구성 파일을 자동으로 생성합니다.

구성 예시:

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

각 목적지는 다음을 가질 수 있습니다:

- `path` (필수): 디렉토리 경로 (홈 디렉토리의 경우 `~` 지원)
- `shortcut` (선택): 단일 문자 단축키
- `command` (선택): 디렉토리 변경 후 실행할 명령어

## 사용법

### 기본 사용법

사용 가능한 목적지를 보려면 `goto` 명령어를 실행하세요:

```sh
goto
```

### 명령줄 인수

목적지를 명령줄 인수로 직접 지정할 수도 있습니다:

```sh
# 숫자 사용
goto 1
goto 4

# 라벨명 사용
goto Home
goto MyProject

# 단축키 사용
goto h
goto p

# 사용 이력 보기
goto --history

# 도움말 표시
goto --help

# 버전 표시
goto --version
```

이는 스크립팅이나 정확히 어디로 가고 싶은지 아는 경우에 유용합니다.

### 대화형 모드

인수 없이 실행하면 `goto`는 대화형 메뉴를 표시합니다:

예시 출력:

```text
👉 사용 가능한 디렉토리:
1. Home → /Users/username/ (shortcut: h)
2. Desktop → /Users/username/Desktop (shortcut: d)
3. Downloads → /Users/username/Downloads (shortcut: b)
4. MyProject → /Users/username/workspace/my-project (shortcut: p)

➕ [+] 현재 디렉토리 추가

번호, 단축키, 라벨명 또는 [+]를 입력하세요:
번호, 단축키, 라벨명 또는 [+] 입력:
```

다음 방법으로 탐색할 수 있습니다:

- **숫자**: `1`, `2`, `3` 등 입력
- **단축키**: `h`, `d`, `b` 등 입력
- **현재 디렉토리 추가**: `+` 입력하여 현재 디렉토리 추가

### 현재 디렉토리 추가

`[+]`를 선택하여 현재 디렉토리를 goto 목적지에 추가할 수 있습니다:

```sh
goto
# 메뉴에서 [+] 선택
# 현재 디렉토리의 라벨 입력
# 선택적으로 단축키 입력
```

예시:

```text
번호, 단축키, 라벨명 또는 [+] 입력: +
📍 현재 디렉토리: /Users/username/workspace/new-project
이 디렉토리의 라벨을 입력하세요: NewProject
단축키를 입력하세요 (선택사항, Enter로 건너뛰기): n
✅ 추가되었습니다 'NewProject' → /Users/username/workspace/new-project
🔑 단축키: n
```

이 기능을 통해 자주 사용하는 디렉토리를 goto 목록에 빠르게 추가할 수 있습니다.

### 새 셸 기능

목적지를 선택하면 `goto`는 대상 디렉토리에서 새 셸 세션을 엽니다. 이는 다음을 의미합니다:

- 현재 셸 세션은 그대로 유지
- 새 위치에서 새로운 셸 환경 제공
- 이전 셸로 돌아가려면 `exit` 입력
- 구성에 `command`가 지정된 경우 자동으로 실행

### 사용 이력

`goto`는 자동으로 사용 이력을 추적하고 최근 사용 순서대로 목적지를 표시합니다. 이를 통해 자주 접근하는 디렉토리가 대화형 메뉴 상단에 나타납니다.

#### 사용 이력 보기

최근 사용 이력을 볼 수 있습니다:

```sh
goto --history
```

예시 출력:

```text
📈 최근 사용 기록:
==================================================
 1. Home → /Users/username
    📅 2025-07-18 16:08:38

 2. Desktop → /Users/username/Desktop
    📅 2025-07-18 16:04:40

 3. MyProject → /Users/username/workspace/my-project
    📅 2025-07-18 15:30:15
```

#### 이력 작동 방식

- **자동 추적**: 목적지로 이동할 때마다 타임스탬프가 기록됩니다
- **스마트 정렬**: 대화형 모드에서 목적지는 최근 사용 순서로 정렬됩니다
- **영구 저장**: 이력은 `~/.goto.toml` 구성 파일에 저장됩니다
- **수동 관리 불필요**: 이력은 자동으로 업데이트되므로 수동 관리가 필요하지 않습니다

#### 이력 저장

사용 이력은 다음 형식으로 `~/.goto.toml` 파일에 저장됩니다:

```toml
[[history]]
  label = "Home"
  last_used = "2025-07-18T16:08:38+09:00"

[[history]]
  label = "Desktop"
  last_used = "2025-07-18T16:04:40+09:00"

# ... 목적지들 ...
[Home]
path = "~/"
shortcut = "h"

[Desktop]
path = "~/Desktop"
shortcut = "d"
```

이 지능적인 정렬은 가장 자주 사용하는 디렉토리가 항상 쉽게 접근할 수 있도록 보장합니다.

### 다국어 지원

`goto`는 자동으로 시스템 언어를 감지하고 선호하는 언어로 메시지를 표시합니다. 현재 지원되는 언어:

- **English** (en) - 영어 (기본값)
- **Japanese** (ja) - 일본어
- **Chinese** (zh) - 중국어
- **Korean** (ko) - 한국어

#### 언어 감지 작동 방식

애플리케이션은 다음 환경 변수를 순서대로 확인하여 시스템 언어를 자동으로 감지합니다:

1. `LANG`
2. `LANGUAGE`
3. `LC_ALL`
4. `LC_MESSAGES`

예를 들어, 시스템이 한국어로 설정된 경우(`LANG=ko_KR.UTF-8`), `goto`는 자동으로 모든 메시지를 한국어로 표시합니다.

#### 다양한 언어의 출력 예시

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

#### 언어 재정의

시스템 설정에 관계없이 특정 언어를 사용하려면 `LANG` 환경 변수를 설정할 수 있습니다:

```sh
# 일본어 인터페이스 사용
LANG=ja_JP.UTF-8 goto

# 영어 인터페이스 사용
LANG=en_US.UTF-8 goto

# 중국어 인터페이스 사용
LANG=zh_CN.UTF-8 goto

# 한국어 인터페이스 사용
LANG=ko_KR.UTF-8 goto
```

#### 지원되는 언어

다국어 지원은 다음을 포함한 모든 사용자 인터페이스 요소를 포함합니다:

- 대화형 메뉴 메시지
- 탐색 확인
- 오류 메시지
- 도움말 텍스트
- 이력 표시
- 구성 메시지

모든 메시지는 시스템 언어 설정에 따라 자동으로 현지화되어 국제 사용자에게 네이티브한 경험을 제공합니다.

### 예시

1. **명령줄 인수를 사용한 탐색 (숫자):**

   ```sh
   goto 1
   goto 4
   ```

2. **명령줄 인수를 사용한 탐색 (라벨):**

   ```sh
   goto Home
   goto MyProject
   ```

3. **명령줄 인수를 사용한 탐색 (단축키):**

   ```sh
   goto h
   goto p
   ```

4. **대화형 탐색:**

   ```sh
   goto
   # 그 다음 입력: h (단축키), 1 (숫자), 또는 Home (라벨)
   ```

5. **현재 디렉토리 추가:**

   ```sh
   cd /path/to/important/project
   goto
   # 입력: +
   # 라벨: ImportantProject
   # 단축키: i
   ```

6. **사용 이력 보기:**

   ```sh
   goto --history
   ```
