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
	ConfigFile               string
	ErrorDetails             string
	ConfigFixSuggestion      string
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

	// Interactive cursor mode messages
	CursorModeHint       string
	BackToCursorModeHint string
	CursorNavigationHint string

	// Interactive help message
	InteractiveHelp string

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
			EnterLabel:            "このディレクトリのラベルを入力してください（Enterでデフォルト使用）",
			EnterShortcut:         "ショートカットキーを入力してください:",
			EnterShortcutOptional: "ショートカットキーを入力してください（任意、Enterでスキップ）:",
			LabelCannotBeEmpty:    "❌ ラベルは空にできません。",
			ShortcutAlreadyExists: "❌ ショートカット '%s' は既に使用されています。別のショートカットを入力してください:",
			Added:                 "✅ 追加しました:",
			Shortcut:              "🔑 ショートカット:",

			// Error messages
			ErrorGettingUser:         "❌ 現在のユーザーの取得エラー:",
			ErrorReadingConfig:       "❌ 設定ファイルの読み取りエラーが発生しました",
			ConfigFile:               "設定ファイル",
			ErrorDetails:             "エラー詳細",
			ConfigFixSuggestion:      "💡 設定ファイルを確認し、古い履歴データが含まれている場合は削除してください。または設定ファイルを削除すると、次回実行時に新しい設定ファイルが作成されます。",
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

			// Interactive cursor mode messages
			CursorModeHint:       "💡 ↑↓/j/kキーで移動、Enterで決定、数字・ショートカットで直接選択、ESCで通常モードに切り替え",
			BackToCursorModeHint: "💡 ヒント: Enterキーのみでカーソル移動モードに戻る",
			CursorNavigationHint: "💡 ↑↓キーで移動、Enterで決定、数字・ショートカットで直接選択、ESCで通常モードに戻る",

			// Interactive help message
			InteractiveHelp: "📋 [?]でヘルプ、[0]で終了、[+]で現在のディレクトリを追加",

			// Other messages
			NoDirectorySelected:  "ℹ️  ディレクトリが選択されていないか、操作がキャンセルされました。",
			CreatedDefaultConfig: "デフォルト設定ファイルを作成しました:",
		}
	case Chinese:
		return Messages{
			// Interactive mode messages
			AvailableDestinations: "👉 可用目录:",
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
			EnterLabel:            "请输入此目录的标签（回车使用默认值）:",
			EnterShortcut:         "请输入快捷键:",
			EnterShortcutOptional: "请输入快捷键（可选，按Enter跳过）:",
			LabelCannotBeEmpty:    "❌ 标签不能为空。",
			ShortcutAlreadyExists: "❌ 快捷键 '%s' 已存在。请输入不同的快捷键:",
			Added:                 "✅ 已添加:",
			Shortcut:              "🔑 快捷键:",

			// Error messages
			ErrorGettingUser:         "❌ 获取当前用户错误:",
			ErrorReadingConfig:       "❌ 配置文件读取错误",
			ConfigFile:               "配置文件",
			ErrorDetails:             "错误详情",
			ConfigFixSuggestion:      "💡 请检查配置文件，如果包含旧的历史数据请删除。或者删除配置文件，下次运行时会创建新的配置文件。",
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

			// Interactive cursor mode messages
			CursorModeHint:       "💡 用↑↓/j/k键移动，Enter确认，数字・快捷键直接选择，ESC切换到普通模式",
			BackToCursorModeHint: "💡 提示: 只按Enter键返回光标移动模式",
			CursorNavigationHint: "💡 用↑↓键移动，Enter确认，数字・快捷键直接选择，ESC切换到普通模式",

			// Interactive help message
			InteractiveHelp: "📋 [?]显示帮助，[0]退出，[+]添加当前目录",

			// Other messages
			NoDirectorySelected:  "ℹ️  未选择目录或操作已取消。",
			CreatedDefaultConfig: "已创建默认配置文件:",
		}
	case Korean:
		return Messages{
			// Interactive mode messages
			AvailableDestinations: "👉 사용 가능한 디렉토리:",
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
			EnterLabel:            "이 디렉토리의 라벨을 입력하세요（엔터로 기본값 사용）:",
			EnterShortcut:         "단축키를 입력하세요:",
			EnterShortcutOptional: "단축키를 입력하세요 (선택사항, Enter로 건너뛰기):",
			LabelCannotBeEmpty:    "❌ 라벨은 비워둘 수 없습니다.",
			ShortcutAlreadyExists: "❌ 단축키 '%s'는 이미 존재합니다. 다른 단축키를 입력하세요:",
			Added:                 "✅ 추가되었습니다:",
			Shortcut:              "🔑 단축키:",

			// Error messages
			ErrorGettingUser:         "❌ 현재 사용자 가져오기 오류:",
			ErrorReadingConfig:       "❌ 설정 파일 읽기 오류가 발생했습니다",
			ConfigFile:               "설정 파일",
			ErrorDetails:             "오류 세부사항",
			ConfigFixSuggestion:      "💡 설정 파일을 확인하고 오래된 히스토리 데이터가 포함되어 있으면 삭제하세요. 또는 설정 파일을 삭제하면 다음 실행 시 새 설정 파일이 생성됩니다.",
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

			// Interactive cursor mode messages
			CursorModeHint:       "💡 ↑↓/j/k키로 이동, Enter로 결정, 숫자・단축키로 직접 선택, ESC로 일반 모드 전환",
			BackToCursorModeHint: "💡 팁: Enter키만으로 커서 이동 모드로 돌아가기",
			CursorNavigationHint: "💡 ↑↓키로 이동, Enter로 결정, 숫자・단축키로 직접 선택, ESC로 일반 모드 전환",

			// Interactive help message
			InteractiveHelp: "📋 [?]로 도움말, [0]으로 종료, [+]로 현재 디렉토리 추가",

			// Other messages
			NoDirectorySelected:  "ℹ️  디렉토리가 선택되지 않았거나 작업이 취소되었습니다.",
			CreatedDefaultConfig: "기본 설정 파일을 생성했습니다:",
		}
	case Spanish:
		return Messages{
			// Interactive mode messages
			AvailableDestinations: "👉 Destinos disponibles:",
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
			EnterLabel:            "Ingrese una etiqueta para este directorio (Enter para usar predeterminado):",
			EnterShortcut:         "Ingrese una tecla de acceso rápido:",
			EnterShortcutOptional: "Ingrese una tecla de acceso rápido (opcional, presione Enter para omitir):",
			LabelCannotBeEmpty:    "❌ La etiqueta no puede estar vacía.",
			ShortcutAlreadyExists: "❌ El acceso rápido '%s' ya existe. Ingrese un acceso rápido diferente:",
			Added:                 "✅ Agregado:",
			Shortcut:              "🔑 Acceso rápido:",

			// Error messages
			ErrorGettingUser:         "❌ Error obteniendo usuario actual:",
			ErrorReadingConfig:       "❌ Error de lectura del archivo de configuración",
			ConfigFile:               "Archivo de configuración",
			ErrorDetails:             "Detalles del error",
			ConfigFixSuggestion:      "💡 Verifique el archivo de configuración y elimine los datos de historial antiguos si están incluidos. O elimine el archivo de configuración para crear uno nuevo en la próxima ejecución.",
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

			// Interactive cursor mode messages
			CursorModeHint:       "💡 Mover con ↑↓/j/k, Enter para decidir, números・accesos rápidos para selección directa, ESC para modo normal",
			BackToCursorModeHint: "💡 Consejo: Solo presiona Enter para volver al modo de movimiento del cursor",
			CursorNavigationHint: "💡 Mover con ↑↓, Enter para decidir, números・accesos rápidos para selección directa, ESC para modo normal",

			// Interactive help message
			InteractiveHelp: "📋 [?] para ayuda, [0] para salir, [+] para agregar directorio actual",

			// Other messages
			NoDirectorySelected:  "ℹ️  No se seleccionó directorio o la operación fue cancelada.",
			CreatedDefaultConfig: "Archivo de configuración por defecto creado:",
		}
	default: // English
		return Messages{
			// Interactive mode messages
			AvailableDestinations: "👉 Available destinations:",
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
			EnterLabel:            "Enter a label for this directory (Enter to use default):",
			EnterShortcut:         "Enter a shortcut key:",
			EnterShortcutOptional: "Enter a shortcut key (optional, press Enter to skip):",
			LabelCannotBeEmpty:    "❌ Label cannot be empty.",
			ShortcutAlreadyExists: "❌ Shortcut '%s' already exists. Please enter a different shortcut:",
			Added:                 "✅ Added:",
			Shortcut:              "🔑 Shortcut:",

			// Error messages
			ErrorGettingUser:         "❌ Error getting current user:",
			ErrorReadingConfig:       "❌ Configuration file reading error occurred",
			ConfigFile:               "Configuration file",
			ErrorDetails:             "Error details",
			ConfigFixSuggestion:      "💡 Please check the configuration file and remove any old history data if included. Or delete the configuration file to create a new one on next run.",
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

			// Interactive cursor mode messages
			CursorModeHint:       "💡 Move with ↑↓/j/k keys, Enter to decide, numbers・shortcuts for direct selection, ESC to switch to normal mode",
			BackToCursorModeHint: "💡 Hint: Press Enter only to return to cursor movement mode",
			CursorNavigationHint: "💡 Move with ↑↓ keys, Enter to decide, numbers・shortcuts for direct selection, ESC to switch to normal mode",

			// Interactive help message
			InteractiveHelp: "📋 Press [?] for help, [0] to exit, [+] to add current dir",

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
