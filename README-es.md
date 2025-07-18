# Comando goto

Comando `goto` para navegar directorios rápidamente.

Esta es una implementación en Go que proporciona navegación de directorios rápida y sin dependencias.

- [English](README.md) / [日本語](README-ja.md) / [中文](README-zh.md) / [한국어](README-ko.md)

## Inicio Rápido

1. **Descargar** el binario más reciente para su plataforma desde [Releases](https://github.com/kujirahand/goto/releases)
2. **Hacerlo ejecutable** y colocarlo en su PATH
3. **Ejecutar** `goto` para ver el menú interactivo

## Características Principales

- **Navegación Rápida de Directorios**: Saltar a directorios frecuentemente usados instantáneamente
- **Historial Inteligente**: Ordena automáticamente destinos por último uso
- **Múltiples Métodos de Entrada**: Usar números, etiquetas o teclas de acceso rápido
- **Completado con Tab**: Soporte de completado para Bash y Zsh
- **Multiplataforma**: Funciona en Linux, macOS y Windows
- **Soporte Multiidioma**: Detección automática de idioma (Inglés, Japonés, Chino, Coreano, Español)
- **Cero Dependencias**: Binario único sin dependencias externas

## Instalación

Por favor instale el comando `goto` siguiendo los pasos a continuación.

### Descargar Binario Pre-construido (Recomendado)

La manera más fácil de instalar `goto` es descargar un binario pre-construido desde la página de releases de GitHub:

1. **Visite la página de releases**: <https://github.com/kujirahand/goto/releases>
2. **Descargue el binario para su plataforma**:
   - **Linux amd64**: `goto-linux-amd64`
   - **Linux arm64**: `goto-linux-arm64`
   - **macOS Intel**: `goto-darwin-amd64`
   - **macOS Apple Silicon**: `goto-darwin-arm64`
   - **Windows amd64**: `goto-windows-amd64.exe`
   - **Windows arm64**: `goto-windows-arm64.exe`

3. **Hacerlo ejecutable y colocarlo en su PATH**:

   **Para Linux/macOS**:

   ```sh
   # Descargar y hacer ejecutable
   chmod +x goto-*
   
   # Mover a un directorio en su PATH
   sudo mv goto-* /usr/local/bin/goto
   
   # O crear un directorio bin local (si no existe)
   mkdir -p ~/bin
   mv goto-* ~/bin/goto
   export PATH="$PATH:$HOME/bin"  # Agregue esto a su configuración de shell
   ```

   **Para Windows**:
   - Renombre el archivo descargado a `goto.exe`
   - Colóquelo en un directorio que esté en su PATH, o cree un nuevo directorio y agréguelo al PATH

4. **Verificar instalación**:

   ```sh
   goto --version
   ```

### Clonar y Construir desde Código Fuente

```sh
# Clonar repositorio
git clone https://github.com/kujirahand/goto.git
# Construir
cd goto
make
```

### Agregar al PATH

Agregue el directorio `goto/go` a su PATH agregando la siguiente línea a su archivo de configuración de shell (`.bashrc`, `.zshrc`, etc.):

```sh
export PATH="$PATH:/ruta/a/goto/go"
```

Por ejemplo, si clonó en su directorio home:

```sh
export PATH="$PATH:$HOME/goto/go"
```

Después de agregar al PATH, recargue su configuración de shell:

```sh
# Para zsh
source ~/.zshrc

# Para bash
source ~/.bashrc
```

### Instalar con Completado con Tab (Construcción desde Código Fuente)

Si construyó desde código fuente, puede instalar tanto el binario como los scripts de completado:

```sh
# Construir e instalar todo (requiere código fuente)
make install-all
```

### Configuración Manual de Completado con Tab (Para Binarios Pre-construidos)

Si descargó un binario pre-construido, aún puede configurar el completado con tab manualmente:

1. **Descargar scripts de completado**:

   ```sh
   # Crear directorios de completado
   mkdir -p ~/.bash_completion.d ~/.zsh/completions
   
   # Descargar script de completado para bash
   curl -o ~/.bash_completion.d/goto-completion.bash \
     https://raw.githubusercontent.com/kujirahand/goto/main/completion/goto-completion.bash
   
   # Descargar script de completado para zsh
   curl -o ~/.zsh/completions/_goto \
     https://raw.githubusercontent.com/kujirahand/goto/main/completion/_goto
   ```

2. **Agregar a su configuración de shell**:

   **Para bash** (`~/.bashrc` o `~/.bash_profile`):

   ```sh
   source ~/.bash_completion.d/goto-completion.bash
   ```

   **Para zsh** (`~/.zshrc`):

   ```sh
   fpath=(~/.zsh/completions $fpath)
   autoload -U compinit && compinit
   ```

3. **Reiniciar su shell o recargar configuración**:

   ```sh
   source ~/.bashrc   # para bash
   source ~/.zshrc    # para zsh
   ```

### Instalación Avanzada con Completado con Tab (Construcción desde Código Fuente)

Para la mejor experiencia al construir desde código fuente, instale tanto el binario como los scripts de completado:

```sh
# Construir e instalar todo
make install-all
```

Esto:

1. Instalará el binario `goto` en `/usr/local/bin/`
2. Instalará scripts de completado de shell
3. Mostrará instrucciones para habilitar el completado

#### Alternativa: Configuración Manual de Completado (Construcción desde Código Fuente)

Si construyó desde código fuente pero prefiere instalar el completado manualmente:

1. Instalar scripts de completado:

   ```sh
   make install-completion
   ```

2. Agregue lo siguiente a su configuración de shell:

   **Para bash** (`~/.bashrc` o `~/.bash_profile`):

   ```sh
   source ~/.bash_completion.d/goto-completion.bash
   ```

   **Para zsh** (`~/.zshrc`):

   ```sh
   fpath=(~/.zsh/completions $fpath)
   autoload -U compinit && compinit
   ```

3. Reinicie su shell o recargue la configuración:

   ```sh
   source ~/.bashrc   # para bash
   source ~/.zshrc    # para zsh
   ```

#### Usar Completado con Tab

Una vez habilitado, puede usar el completado con tab con el comando `goto`:

```sh
goto <TAB>        # Muestra todos los destinos disponibles
goto h<TAB>       # Completa accesos rápidos que comienzan con 'h'
goto Home<TAB>    # Completa etiquetas que comienzan con 'Home'
goto 1<TAB>       # Muestra destinos con números que comienzan con '1'
```

## Configuración

### Archivo de configuración - `~/.goto.toml`

El comando `goto` usa un archivo de configuración TOML ubicado en `~/.goto.toml`. Cuando ejecute `goto` por primera vez, creará automáticamente un archivo de configuración predeterminado con destinos de ejemplo.

Configuración de ejemplo:

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

Cada destino puede tener:

- `path` (requerido): Ruta del directorio (soporta `~` para directorio home)
- `shortcut` (opcional): Tecla de acceso rápido de un solo carácter
- `command` (opcional): Comando a ejecutar después de cambiar de directorio

## Uso

### Uso Básico

Ejecute el comando `goto` para ver destinos disponibles:

```sh
goto
```

### Argumentos de Línea de Comandos

También puede especificar un destino directamente como argumento de línea de comandos:

```sh
# Usando número
goto 1
goto 4

# Usando nombre de etiqueta
goto Home
goto MyProject

# Usando tecla de acceso rápido
goto h
goto p

# Ver historial de uso
goto --history

# Mostrar ayuda
goto --help

# Mostrar versión
goto --version
```

Esto es útil para scripts o cuando sabe exactamente a dónde quiere ir.

### Modo Interactivo

Cuando se ejecuta sin argumentos, `goto` muestra un menú interactivo:

Salida de ejemplo:

```text
👉 Destinos disponibles:
1. Home → /Users/username/ (shortcut: h)
2. Desktop → /Users/username/Desktop (shortcut: d)
3. Downloads → /Users/username/Downloads (shortcut: b)
4. MyProject → /Users/username/workspace/my-project (shortcut: p)

➕ [+] Agregar directorio actual

Ingrese número, tecla de acceso rápido, nombre de etiqueta o [+]:
```

Puede navegar mediante:

- **Número**: Ingrese `1`, `2`, `3`, etc.
- **Acceso rápido**: Ingrese `h`, `d`, `b`, etc.
- **Agregar actual**: Ingrese `+` para agregar el directorio actual

### Agregar Directorio Actual

Puede agregar el directorio actual a sus destinos goto seleccionando `[+]`:

```sh
goto
# Seleccione [+] del menú
# Ingrese una etiqueta para el directorio actual
# Opcionalmente ingrese una tecla de acceso rápido
```

Ejemplo:

```text
Destino: +
📍 Directorio actual: /Users/username/workspace/new-project
Ingrese una etiqueta para este directorio: NewProject
Ingrese una tecla de acceso rápido (opcional, presione Enter para omitir): n
✅ Agregado 'NewProject' → /Users/username/workspace/new-project
🔑 Acceso rápido: n
```

Esta característica le permite agregar rápidamente directorios frecuentemente usados a su lista goto.

### Funcionalidad de Nuevo Shell

Cuando selecciona un destino, `goto` abre una nueva sesión de shell en el directorio objetivo. Esto significa:

- Su sesión de shell actual permanece sin cambios
- Obtiene un ambiente de shell fresco en la nueva ubicación
- Escriba `exit` para regresar a su shell anterior
- Si se especifica un `command` en la configuración, se ejecutará automáticamente

### Historial de Uso

`goto` rastrea automáticamente el historial de uso y muestra destinos en orden del último uso. Esto hace que directorios frecuentemente accedidos aparezcan en la parte superior del menú interactivo.

#### Ver Historial de Uso

Puede ver su historial de uso reciente con:

```sh
goto --history
```

Salida de ejemplo:

```text
📈 Historial de uso reciente:
==================================================
 1. Home → /Users/username
    📅 2025-07-18 16:08:38

 2. Desktop → /Users/username/Desktop
    📅 2025-07-18 16:04:40

 3. MyProject → /Users/username/workspace/my-project
    📅 2025-07-18 15:30:15
```

#### Cómo Funciona el Historial

- **Rastreo automático**: Cada vez que navega a un destino, se registra la marca de tiempo
- **Ordenamiento inteligente**: En modo interactivo, destinos se ordenan por último uso primero
- **Almacenamiento persistente**: El historial se almacena en el archivo de configuración `~/.goto.toml`
- **Sin mantenimiento manual**: El historial se actualiza automáticamente - no necesita gestionarlo manualmente

#### Almacenamiento del Historial

El historial de uso se almacena en su archivo `~/.goto.toml` en el siguiente formato:

```toml
[[history]]
  label = "Home"
  last_used = "2025-07-18T16:08:38+09:00"

[[history]]
  label = "Desktop"
  last_used = "2025-07-18T16:04:40+09:00"

# ... sus destinos ...
[Home]
path = "~/"
shortcut = "h"

[Desktop]
path = "~/Desktop"
shortcut = "d"
```

Este ordenamiento inteligente asegura que sus directorios más frecuentemente usados siempre sean fácilmente accesibles.

## Soporte Multiidioma

`goto` detecta automáticamente el idioma de su sistema y muestra mensajes en su idioma preferido. Idiomas actualmente soportados:

- **Inglés** (en) - Predeterminado
- **Japonés** (ja) - 日本語
- **Chino** (zh) - 中文
- **Coreano** (ko) - 한국어
- **Español** (es) - Español

### Cómo Funciona la Detección de Idioma

La aplicación detecta automáticamente el idioma de su sistema verificando las siguientes variables de entorno en orden:

1. `LANG`
2. `LANGUAGE`
3. `LC_ALL`
4. `LC_MESSAGES`

Por ejemplo, si su sistema está configurado en español (`LANG=es_ES.UTF-8`), `goto` automáticamente mostrará todos los mensajes en español.

### Ejemplo de Salida en Diferentes Idiomas

**Inglés:**

```text
🚀 goto - Navigate directories quickly
👉 Available destinations:
1. Home → /Users/username/ (shortcut: h)
📈 Recent usage history:
```

**Japonés:**

```text
🚀 goto - ディレクトリ間を素早く移動
👉 利用可能なディレクトリ:
1. Home → /Users/username/ (shortcut: h)
📈 最近の使用履歴:
```

**Chino:**

```text
🚀 goto - 快速导航目录
👉 可用目录:
1. Home → /Users/username/ (shortcut: h)
📈 最近使用历史:
```

**Coreano:**

```text
🚀 goto - 디렉토리 빠른 탐색
👉 사용 가능한 디렉토리:
1. Home → /Users/username/ (shortcut: h)
📈 최근 사용 기록:
```

**Español:**

```text
🚀 goto - Navegar directorios rápidamente
👉 Destinos disponibles:
1. Home → /Users/username/ (shortcut: h)
📈 Historial de uso reciente:
```

### Anular Idioma

Si desea usar un idioma específico independientemente de la configuración de su sistema, puede establecer la variable de entorno `LANG`:

```sh
# Usar interfaz en japonés
LANG=ja_JP.UTF-8 goto

# Usar interfaz en inglés
LANG=en_US.UTF-8 goto

# Usar interfaz en chino
LANG=zh_CN.UTF-8 goto

# Usar interfaz en coreano
LANG=ko_KR.UTF-8 goto

# Usar interfaz en español
LANG=es_ES.UTF-8 goto
```

### Idiomas Soportados

El soporte multiidioma cubre todos los elementos de la interfaz de usuario incluyendo:

- Mensajes del menú interactivo
- Confirmaciones de navegación
- Mensajes de error
- Texto de ayuda
- Visualización del historial
- Mensajes de configuración

Todos los mensajes se localizan automáticamente basándose en la configuración de idioma de su sistema, proporcionando una experiencia nativa para usuarios internacionales.

### Ejemplos

1. **Navegar usando argumento de línea de comandos (número):**

   ```sh
   goto 1
   goto 4
   ```

2. **Navegar usando argumento de línea de comandos (etiqueta):**

   ```sh
   goto Home
   goto MyProject
   ```

3. **Navegar usando argumento de línea de comandos (acceso rápido):**

   ```sh
   goto h
   goto p
   ```

4. **Navegación interactiva:**

   ```sh
   goto
   # Luego ingrese: h (acceso rápido), 1 (número), o Home (etiqueta)
   ```

5. **Agregar directorio actual:**

   ```sh
   cd /ruta/a/proyecto/importante
   goto
   # Ingrese: +
   # Etiqueta: ImportantProject
   # Acceso rápido: i
   ```

6. **Ver historial de uso:**

   ```sh
   goto --history
   ```
