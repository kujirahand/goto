package main

import (
	"os"
	"strings"
)

// Language represents supported languages
type Language string

const (
	English  Language = "en"
	Japanese Language = "ja"
	Chinese  Language = "zh"
	Korean   Language = "ko"
	Spanish  Language = "es"
)

// Messages contains all user-facing messages
type Messages struct {
	// Interactive mode messages
	AvailableDestinations string
	AddCurrentDirectory   string
	EnterChoice           string
	EnterChoicePrompt     string

	// Navigation messages
	OpeningShell     string
	Destination      string
	TypeExitToReturn string
	YouAreNowIn      string
	FoundDestination string

	// Add directory messages
	CurrentDirectory      string
	EnterLabel            string
	EnterShortcut         string
	EnterShortcutOptional string
	LabelCannotBeEmpty    string
	ShortcutAlreadyExists string
	Added                 string
	Shortcut              string

	// Error messages
	ErrorGettingUser         string
	ErrorReadingConfig       string
	NoDestinationsConfigured string
	DestinationNotFound      string
	DirectoryNotExist        string
	ErrorOpeningShell        string
	ErrorCreatingTempFile    string
	ErrorWritingTempScript   string
	ErrorMakingExecutable    string
	ErrorOpeningConfigFile   string
	ErrorWritingConfigFile   string
	ErrorGettingCurrentDir   string
	OperationCancelled       string
	InvalidInput             string

	// History messages
	RecentUsageHistory           string
	NoUsageHistoryFound          string
	WarningFailedToUpdateHistory string

	// Command messages
	WillExecute      string
	ExecutingCommand string
	CommandCompleted string

	// Help messages
	NavigateDirectoriesQuickly  string
	ConfigurationFile           string
	Usage                       string
	ShowInteractiveMenu         string
	GoToDestinationByNumber     string
	GoToDestinationByLabel      string
	GoToDestinationByShortcut   string
	ShowHelpMessage             string
	ShowVersionInfo             string
	ShowCompletionCandidates    string
	ShowRecentUsageHistory      string
	AddCurrentDirectoryToConfig string
	Examples                    string
	NavigateToFirstDest         string
	NavigateToHomeDest          string
	NavigateUsingShortcut       string
	ShowInteractiveMenuExample  string

	// Other messages
	NoDirectorySelected  string
	CreatedDefaultConfig string
}

// detectLanguage detects the system language from environment variables
func detectLanguage() Language {
	// Check environment variables in order of preference
	envVars := []string{"LANG", "LANGUAGE", "LC_ALL", "LC_MESSAGES"}

	for _, env := range envVars {
		if lang := os.Getenv(env); lang != "" {
			// Extract language code (e.g., ja_JP.UTF-8 -> ja)
			langCode := strings.Split(lang, "_")[0]
			langCode = strings.Split(langCode, ".")[0]
			langCode = strings.ToLower(langCode)

			switch langCode {
			case "ja":
				return Japanese
			case "zh":
				return Chinese
			case "ko":
				return Korean
			case "es":
				return Spanish
			default:
				return English
			}
		}
	}

	return English // Default to English
}

// getMessages returns localized messages for the specified language
func getMessages(lang Language) Messages {
	switch lang {
	case Japanese:
		return Messages{
			// Interactive mode messages
			AvailableDestinations: "😊 どこに移動しますか？",
			AddCurrentDirectory:   "🌱 [+] 現在のディレクトリを追加",
			EnterChoice:           "[番号]、(キー)、ラベル、または[+]を入力してください:",
			EnterChoicePrompt:     ">>>",

			// Navigation messages
			OpeningShell:     "🚀 新しいシェルを開いています:",
			Destination:      "📍 ディレクトリ:",
			TypeExitToReturn: "💡 前のシェルに戻るには 'exit' を入力してください",
			YouAreNowIn:      "✅ 現在のディレクトリ:",
			FoundDestination: "🎯 見つかったディレクトリ:",

			// Add directory messages
			CurrentDirectory:      "📍 現在のディレクトリ:",
			EnterLabel:            "このディレクトリのラベルを入力してください:",
			EnterShortcut:         "ショートカットキーを入力してください:",
			EnterShortcutOptional: "ショートカットキーを入力してください（任意、Enterでスキップ）:",
			LabelCannotBeEmpty:    "❌ ラベルは空にできません。",
			ShortcutAlreadyExists: "❌ ショートカット '%s' は既に使用されています。別のショートカットを入力してください:",
			Added:                 "✅ 追加しました:",
			Shortcut:              "🔑 ショートカット:",

			// Error messages
			ErrorGettingUser:         "❌ 現在のユーザーの取得エラー:",
			ErrorReadingConfig:       "❌ 設定ファイルの読み取りエラー:",
			NoDestinationsConfigured: "⚠️  ~/.goto.toml にディレクトリが設定されていません",
			DestinationNotFound:      "❌ ディレクトリ '%s' が見つかりません。",
			DirectoryNotExist:        "❌ ディレクトリが存在しません:",
			ErrorOpeningShell:        "❌ シェルを開くエラー:",
			ErrorCreatingTempFile:    "❌ 一時ファイルの作成エラー:",
			ErrorWritingTempScript:   "❌ 一時スクリプトの書き込みエラー:",
			ErrorMakingExecutable:    "❌ スクリプトを実行可能にするエラー:",
			ErrorOpeningConfigFile:   "❌ 設定ファイルを開くエラー:",
			ErrorWritingConfigFile:   "❌ 設定ファイルの書き込みエラー:",
			ErrorGettingCurrentDir:   "❌ 現在のディレクトリの取得エラー:",
			OperationCancelled:       "❌ 操作がキャンセルされました。",
			InvalidInput:             "無効な入力です。",

			// History messages
			RecentUsageHistory:           "📈 最近の使用履歴:",
			NoUsageHistoryFound:          "📈 使用履歴が見つかりません。",
			WarningFailedToUpdateHistory: "⚠️  警告: 履歴の更新に失敗しました:",

			// Command messages
			WillExecute:      "⚡ 実行します:",
			ExecutingCommand: "⚡ 実行中:",
			CommandCompleted: "✅ コマンドが完了しました。現在のディレクトリ:",

			// Help messages
			NavigateDirectoriesQuickly:  "🚀 goto - ディレクトリ間を素早く移動",
			ConfigurationFile:           "設定ファイル:",
			Usage:                       "使用方法:",
			ShowInteractiveMenu:         "インタラクティブメニューを表示",
			GoToDestinationByNumber:     "番号でディレクトリに移動 (例: goto 1)",
			GoToDestinationByLabel:      "ラベル名でディレクトリに移動",
			GoToDestinationByShortcut:   "ショートカットキーでディレクトリに移動",
			ShowHelpMessage:             "このヘルプメッセージを表示",
			ShowVersionInfo:             "バージョン情報を表示",
			ShowCompletionCandidates:    "補完候補を表示 (シェル補完用)",
			ShowRecentUsageHistory:      "最近の使用履歴を表示",
			AddCurrentDirectoryToConfig: "現在のディレクトリを設定に追加",
			Examples:                    "例:",
			NavigateToFirstDest:         "# 1番目のディレクトリに移動",
			NavigateToHomeDest:          "# 'Home' ディレクトリに移動",
			NavigateUsingShortcut:       "# ショートカット 'h' を使用して移動",
			ShowInteractiveMenuExample:  "# インタラクティブメニューを表示",

			// Other messages
			NoDirectorySelected:  "ℹ️  ディレクトリが選択されていないか、操作がキャンセルされました。",
			CreatedDefaultConfig: "デフォルト設定ファイルを作成しました:",
		}
	case Chinese:
		return Messages{
			// Interactive mode messages
			AvailableDestinations: "👉 可用目录:",
			AddCurrentDirectory:   "🌱 [+] 添加当前目录",
			EnterChoice:           "请输入编号、快捷键、标签名称或 [+] 添加当前目录:",
			EnterChoicePrompt:     "输入编号、快捷键、标签名称或 [+]:",

			// Navigation messages
			OpeningShell:     "🚀 正在打开新Shell:",
			Destination:      "📍 目录:",
			TypeExitToReturn: "💡 输入 'exit' 返回上一个Shell",
			YouAreNowIn:      "✅ 您现在在:",
			FoundDestination: "🎯 找到目录:",

			// Add directory messages
			CurrentDirectory:      "📍 当前目录:",
			EnterLabel:            "请输入此目录的标签:",
			EnterShortcut:         "请输入快捷键:",
			EnterShortcutOptional: "请输入快捷键（可选，按Enter跳过）:",
			LabelCannotBeEmpty:    "❌ 标签不能为空。",
			ShortcutAlreadyExists: "❌ 快捷键 '%s' 已存在。请输入不同的快捷键:",
			Added:                 "✅ 已添加:",
			Shortcut:              "🔑 快捷键:",

			// Error messages
			ErrorGettingUser:         "❌ 获取当前用户错误:",
			ErrorReadingConfig:       "❌ 读取配置文件错误:",
			NoDestinationsConfigured: "⚠️  ~/.goto.toml 中未配置目录",
			DestinationNotFound:      "❌ 未找到目录 '%s'。",
			DirectoryNotExist:        "❌ 目录不存在:",
			ErrorOpeningShell:        "❌ 打开Shell错误:",
			ErrorCreatingTempFile:    "❌ 创建临时文件错误:",
			ErrorWritingTempScript:   "❌ 写入临时脚本错误:",
			ErrorMakingExecutable:    "❌ 设置脚本可执行错误:",
			ErrorOpeningConfigFile:   "❌ 打开配置文件错误:",
			ErrorWritingConfigFile:   "❌ 写入配置文件错误:",
			ErrorGettingCurrentDir:   "❌ 获取当前目录错误:",
			OperationCancelled:       "❌ 操作已取消。",
			InvalidInput:             "无效输入。",

			// History messages
			RecentUsageHistory:           "📈 最近使用历史:",
			NoUsageHistoryFound:          "📈 未找到使用历史。",
			WarningFailedToUpdateHistory: "⚠️  警告: 更新历史失败:",

			// Command messages
			WillExecute:      "⚡ 将执行:",
			ExecutingCommand: "⚡ 执行中:",
			CommandCompleted: "✅ 命令已完成。当前目录:",

			// Help messages
			NavigateDirectoriesQuickly:  "🚀 goto - 快速导航目录",
			ConfigurationFile:           "配置文件:",
			Usage:                       "用法:",
			ShowInteractiveMenu:         "显示交互式菜单",
			GoToDestinationByNumber:     "通过编号转到目录 (例: goto 1)",
			GoToDestinationByLabel:      "通过标签名转到目录",
			GoToDestinationByShortcut:   "通过快捷键转到目录",
			ShowHelpMessage:             "显示此帮助消息",
			ShowVersionInfo:             "显示版本信息",
			ShowCompletionCandidates:    "显示补全候选 (用于Shell补全)",
			ShowRecentUsageHistory:      "显示最近使用历史",
			AddCurrentDirectoryToConfig: "将当前目录添加到配置",
			Examples:                    "示例:",
			NavigateToFirstDest:         "# 导航到第1个目录",
			NavigateToHomeDest:          "# 导航到 'Home' 目录",
			NavigateUsingShortcut:       "# 使用快捷键 'h' 导航",
			ShowInteractiveMenuExample:  "# 显示交互式菜单",

			// Other messages
			NoDirectorySelected:  "ℹ️  未选择目录或操作已取消。",
			CreatedDefaultConfig: "已创建默认配置文件:",
		}
	case Korean:
		return Messages{
			// Interactive mode messages
			AvailableDestinations: "👉 사용 가능한 디렉토리:",
			AddCurrentDirectory:   "🌱 [+] 현재 디렉토리 추가",
			EnterChoice:           "번호, 단축키, 라벨명 또는 [+]를 입력하세요:",
			EnterChoicePrompt:     "번호, 단축키, 라벨명 또는 [+] 입력:",

			// Navigation messages
			OpeningShell:     "🚀 새 셸을 열고 있습니다:",
			Destination:      "📍 디렉토리:",
			TypeExitToReturn: "💡 이전 셸로 돌아가려면 'exit'를 입력하세요",
			YouAreNowIn:      "✅ 현재 위치:",
			FoundDestination: "🎯 디렉토리를 찾았습니다:",

			// Add directory messages
			CurrentDirectory:      "📍 현재 디렉토리:",
			EnterLabel:            "이 디렉토리의 라벨을 입력하세요:",
			EnterShortcut:         "단축키를 입력하세요:",
			EnterShortcutOptional: "단축키를 입력하세요 (선택사항, Enter로 건너뛰기):",
			LabelCannotBeEmpty:    "❌ 라벨은 비워둘 수 없습니다.",
			ShortcutAlreadyExists: "❌ 단축키 '%s'는 이미 존재합니다. 다른 단축키를 입력하세요:",
			Added:                 "✅ 추가되었습니다:",
			Shortcut:              "🔑 단축키:",

			// Error messages
			ErrorGettingUser:         "❌ 현재 사용자 가져오기 오류:",
			ErrorReadingConfig:       "❌ 설정 파일 읽기 오류:",
			NoDestinationsConfigured: "⚠️  ~/.goto.toml에 디렉토리가 설정되지 않았습니다",
			DestinationNotFound:      "❌ 디렉토리 '%s'를 찾을 수 없습니다.",
			DirectoryNotExist:        "❌ 디렉토리가 존재하지 않습니다:",
			ErrorOpeningShell:        "❌ 셸 열기 오류:",
			ErrorCreatingTempFile:    "❌ 임시 파일 생성 오류:",
			ErrorWritingTempScript:   "❌ 임시 스크립트 작성 오류:",
			ErrorMakingExecutable:    "❌ 스크립트 실행 가능 설정 오류:",
			ErrorOpeningConfigFile:   "❌ 설정 파일 열기 오류:",
			ErrorWritingConfigFile:   "❌ 설정 파일 작성 오류:",
			ErrorGettingCurrentDir:   "❌ 현재 디렉토리 가져오기 오류:",
			OperationCancelled:       "❌ 작업이 취소되었습니다.",
			InvalidInput:             "잘못된 입력입니다.",

			// History messages
			RecentUsageHistory:           "📈 최근 사용 기록:",
			NoUsageHistoryFound:          "📈 사용 기록을 찾을 수 없습니다.",
			WarningFailedToUpdateHistory: "⚠️  경고: 기록 업데이트에 실패했습니다:",

			// Command messages
			WillExecute:      "⚡ 실행할 명령:",
			ExecutingCommand: "⚡ 실행 중:",
			CommandCompleted: "✅ 명령이 완료되었습니다. 현재 디렉토리:",

			// Help messages
			NavigateDirectoriesQuickly:  "🚀 goto - 디렉토리 빠른 탐색",
			ConfigurationFile:           "설정 파일:",
			Usage:                       "사용법:",
			ShowInteractiveMenu:         "대화형 메뉴 표시",
			GoToDestinationByNumber:     "번호로 디렉토리 이동 (예: goto 1)",
			GoToDestinationByLabel:      "라벨명으로 디렉토리 이동",
			GoToDestinationByShortcut:   "단축키로 디렉토리 이동",
			ShowHelpMessage:             "이 도움말 메시지 표시",
			ShowVersionInfo:             "버전 정보 표시",
			ShowCompletionCandidates:    "완성 후보 표시 (셸 완성용)",
			ShowRecentUsageHistory:      "최근 사용 기록 표시",
			AddCurrentDirectoryToConfig: "현재 디렉토리를 설정에 추가",
			Examples:                    "예제:",
			NavigateToFirstDest:         "# 첫 번째 디렉토리로 이동",
			NavigateToHomeDest:          "# 'Home' 디렉토리로 이동",
			NavigateUsingShortcut:       "# 단축키 'h' 사용하여 이동",
			ShowInteractiveMenuExample:  "# 대화형 메뉴 표시",

			// Other messages
			NoDirectorySelected:  "ℹ️  디렉토리가 선택되지 않았거나 작업이 취소되었습니다.",
			CreatedDefaultConfig: "기본 설정 파일을 생성했습니다:",
		}
	case Spanish:
		return Messages{
			// Interactive mode messages
			AvailableDestinations: "👉 Destinos disponibles:",
			AddCurrentDirectory:   "🌱 [+] Agregar directorio actual",
			EnterChoice:           "Ingrese número, tecla de acceso rápido, nombre de etiqueta o [+]:",
			EnterChoicePrompt:     "Destino:",

			// Navigation messages
			OpeningShell:     "🚀 Abriendo nuevo shell en:",
			Destination:      "📍 Destino:",
			TypeExitToReturn: "💡 Escriba 'exit' para regresar al shell anterior",
			YouAreNowIn:      "✅ Ahora está en:",
			FoundDestination: "🎯 Destino encontrado:",

			// Add directory messages
			CurrentDirectory:      "📍 Directorio actual:",
			EnterLabel:            "Ingrese una etiqueta para este directorio:",
			EnterShortcut:         "Ingrese una tecla de acceso rápido:",
			EnterShortcutOptional: "Ingrese una tecla de acceso rápido (opcional, presione Enter para omitir):",
			LabelCannotBeEmpty:    "❌ La etiqueta no puede estar vacía.",
			ShortcutAlreadyExists: "❌ El acceso rápido '%s' ya existe. Ingrese un acceso rápido diferente:",
			Added:                 "✅ Agregado:",
			Shortcut:              "🔑 Acceso rápido:",

			// Error messages
			ErrorGettingUser:         "❌ Error obteniendo usuario actual:",
			ErrorReadingConfig:       "❌ Error leyendo archivo de configuración:",
			NoDestinationsConfigured: "⚠️  No hay destinos configurados en ~/.goto.toml",
			DestinationNotFound:      "❌ Destino '%s' no encontrado.",
			DirectoryNotExist:        "❌ El directorio no existe:",
			ErrorOpeningShell:        "❌ Error abriendo shell:",
			ErrorCreatingTempFile:    "❌ Error creando archivo temporal:",
			ErrorWritingTempScript:   "❌ Error escribiendo script temporal:",
			ErrorMakingExecutable:    "❌ Error haciendo ejecutable el script:",
			ErrorOpeningConfigFile:   "❌ Error abriendo archivo de configuración:",
			ErrorWritingConfigFile:   "❌ Error escribiendo archivo de configuración:",
			ErrorGettingCurrentDir:   "❌ Error obteniendo directorio actual:",
			OperationCancelled:       "❌ Operación cancelada.",
			InvalidInput:             "Entrada inválida.",

			// History messages
			RecentUsageHistory:           "📈 Historial de uso reciente:",
			NoUsageHistoryFound:          "📈 No se encontró historial de uso.",
			WarningFailedToUpdateHistory: "⚠️  Advertencia: Falló al actualizar historial:",

			// Command messages
			WillExecute:      "⚡ Ejecutará:",
			ExecutingCommand: "⚡ Ejecutando:",
			CommandCompleted: "✅ Comando completado. Ahora está en:",

			// Help messages
			NavigateDirectoriesQuickly:  "🚀 goto - Navegar directorios rápidamente",
			ConfigurationFile:           "Archivo de configuración:",
			Usage:                       "Uso:",
			ShowInteractiveMenu:         "Mostrar menú interactivo",
			GoToDestinationByNumber:     "Ir al destino por número (ej., goto 1)",
			GoToDestinationByLabel:      "Ir al destino por nombre de etiqueta",
			GoToDestinationByShortcut:   "Ir al destino por tecla de acceso rápido",
			ShowHelpMessage:             "Mostrar este mensaje de ayuda",
			ShowVersionInfo:             "Mostrar información de versión",
			ShowCompletionCandidates:    "Mostrar candidatos de completado (para completado de shell)",
			ShowRecentUsageHistory:      "Mostrar historial de uso reciente",
			AddCurrentDirectoryToConfig: "Agregar directorio actual a la configuración",
			Examples:                    "Ejemplos:",
			NavigateToFirstDest:         "# Navegar al 1er destino",
			NavigateToHomeDest:          "# Navegar al destino 'Home'",
			NavigateUsingShortcut:       "# Navegar usando acceso rápido 'h'",
			ShowInteractiveMenuExample:  "# Mostrar menú interactivo",

			// Other messages
			NoDirectorySelected:  "ℹ️  No se seleccionó directorio o la operación fue cancelada.",
			CreatedDefaultConfig: "Archivo de configuración por defecto creado:",
		}
	default: // English
		return Messages{
			// Interactive mode messages
			AvailableDestinations: "👉 Available destinations:",
			AddCurrentDirectory:   "🌱 [+] Add current directory",
			EnterChoice:           "Please enter the number, shortcut key, label name, or [+]:",
			EnterChoicePrompt:     "Enter choice:",

			// Navigation messages
			OpeningShell:     "🚀 Opening new shell in:",
			Destination:      "📍 Destination:",
			TypeExitToReturn: "💡 Type 'exit' to return to previous shell",
			YouAreNowIn:      "✅ You are now in:",
			FoundDestination: "🎯 Found destination:",

			// Add directory messages
			CurrentDirectory:      "📍 Current directory:",
			EnterLabel:            "Enter a label for this directory:",
			EnterShortcut:         "Enter a shortcut key:",
			EnterShortcutOptional: "Enter a shortcut key (optional, press Enter to skip):",
			LabelCannotBeEmpty:    "❌ Label cannot be empty.",
			ShortcutAlreadyExists: "❌ Shortcut '%s' already exists. Please enter a different shortcut:",
			Added:                 "✅ Added:",
			Shortcut:              "🔑 Shortcut:",

			// Error messages
			ErrorGettingUser:         "❌ Error getting current user:",
			ErrorReadingConfig:       "❌ Error reading configuration file:",
			NoDestinationsConfigured: "⚠️  No destinations configured in ~/.goto.toml",
			DestinationNotFound:      "❌ Destination '%s' not found.",
			DirectoryNotExist:        "❌ Directory does not exist:",
			ErrorOpeningShell:        "❌ Error opening shell:",
			ErrorCreatingTempFile:    "❌ Error creating temp file:",
			ErrorWritingTempScript:   "❌ Error writing temp script:",
			ErrorMakingExecutable:    "❌ Error making script executable:",
			ErrorOpeningConfigFile:   "❌ Error opening config file:",
			ErrorWritingConfigFile:   "❌ Error writing to config file:",
			ErrorGettingCurrentDir:   "❌ Error getting current directory:",
			OperationCancelled:       "❌ Operation cancelled.",
			InvalidInput:             "Invalid input.",

			// History messages
			RecentUsageHistory:           "📈 Recent usage history:",
			NoUsageHistoryFound:          "📈 No usage history found.",
			WarningFailedToUpdateHistory: "⚠️  Warning: Failed to update history:",

			// Command messages
			WillExecute:      "⚡ Will execute:",
			ExecutingCommand: "⚡ Executing:",
			CommandCompleted: "✅ Command completed. You are now in:",

			// Help messages
			NavigateDirectoriesQuickly:  "🚀 goto - Navigate directories quickly",
			ConfigurationFile:           "Configuration file:",
			Usage:                       "Usage:",
			ShowInteractiveMenu:         "Show interactive menu",
			GoToDestinationByNumber:     "Go to destination by number (e.g., goto 1)",
			GoToDestinationByLabel:      "Go to destination by label name",
			GoToDestinationByShortcut:   "Go to destination by shortcut key",
			ShowHelpMessage:             "Show this help message",
			ShowVersionInfo:             "Show version information",
			ShowCompletionCandidates:    "Show completion candidates (for shell completion)",
			ShowRecentUsageHistory:      "Show recent usage history",
			AddCurrentDirectoryToConfig: "Add current directory to configuration",
			Examples:                    "Examples:",
			NavigateToFirstDest:         "# Navigate to 1st destination",
			NavigateToHomeDest:          "# Navigate to 'Home' destination",
			NavigateUsingShortcut:       "# Navigate using shortcut 'h'",
			ShowInteractiveMenuExample:  "# Show interactive menu",

			// Other messages
			NoDirectorySelected:  "ℹ️  No directory selected or operation cancelled.",
			CreatedDefaultConfig: "Created default configuration file:",
		}
	}
}

// Global variables for current language and messages
var (
	currentLanguage Language = detectLanguage()
	messages        Messages = getMessages(currentLanguage)
)
