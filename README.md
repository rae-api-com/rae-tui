# 📖 RAE-TUI

[![Go Version](https://img.shields.io/github/go-mod/go-version/rae-api-com/rae-tui)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/rae-api-com/rae-tui)](https://goreportcard.com/report/github.com/rae-api-com/rae-tui)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/rae-api-com/rae-tui)](https://github.com/rae-api-com/rae-tui/releases)
[![Build Status](https://github.com/rae-api-com/rae-tui/workflows/CI/badge.svg)](https://github.com/rae-api-com/rae-tui/actions)

Un cliente de diccionario español elegante y rápido para la terminal, conectado a la Real Academia Española (RAE). 

![RAE-TUI Demo](demo.gif)

> 🎯 **Perfecto para**: Escritores, estudiantes, desarrolladores y cualquiera que necesite consultas rápidas del diccionario sin salir de la terminal.

## ✨ Características

- 🔍 **Búsqueda Instantánea** - Consulta palabras directamente desde tu terminal
- 📚 **Definiciones Completas** - Visualiza significados detallados y acepciones múltiples
- 🖥️ **Interfaz Interactiva** - TUI moderna con navegación por teclado intuitiva
- 🔄 **Conjugaciones Verbales** - Muestra todas las conjugaciones en todos los tiempos
- ⚡ **Modo CLI Rápido** - Consultas no-interactivas para uso en scripts
- 🔓 **Sin API Key** - Funciona inmediatamente sin configuración
- 🎨 **Colores y Formato** - Salida elegante y fácil de leer
- 📱 **Multiplataforma** - Funciona en Linux, macOS y Windows

## 📦 Instalación

### 🚀 Instalación Rápida (Recomendada)

```bash
go install github.com/rae-api-com/rae-tui@latest
```

### 📥 Binarios Precompilados

Descarga el binario apropiado para tu plataforma desde la página de [Releases](https://github.com/rae-api-com/rae-tui/releases).

#### Linux/macOS
```bash
# Descargar y instalar automáticamente
curl -sf https://gobinaries.com/rae-api-com/rae-tui | sh
```

#### Windows
```powershell
# Usando Scoop
scoop bucket add rae-api-com https://github.com/rae-api-com/scoop-bucket
scoop install rae-tui
```

### 🔨 Desde Código Fuente

```bash
git clone https://github.com/rae-api-com/rae-tui.git
cd rae-tui
go build -o rae-tui
```

## 🎮 Uso

### 🖥️ Modo Interactivo (TUI)

Lanza la interfaz interactiva completa:

```bash
rae-tui
# o explícitamente
rae-tui tui
```

**Búsqueda directa con una palabra:**
```bash
rae-tui tui hola
```

### ⚡ Modo CLI (Rápido)

Consulta directa para scripts o uso rápido:

```bash
rae-tui palabra
# Ejemplo:
rae-tui computadora
```

**Salida en formato JSON:**
```bash
rae-tui --json palabra
```

### 🎯 Ejemplos Prácticos

```bash
# Buscar un verbo y ver conjugaciones
rae-tui tui comer

# Consulta rápida en scripts
if rae-tui existir > /dev/null; then
    echo "La palabra existe en el diccionario"
fi

# Integración con otros comandos
echo "palabras\nque\nbuscar" | xargs -I {} rae-tui {}
```

### ⌨️ Atajos de Teclado (Modo TUI)

| Tecla          | Acción                         |
| -------------- | ------------------------------ |
| `↑` / `k`      | ⬆️ Mover selección hacia arriba |
| `↓` / `j`      | ⬇️ Mover selección hacia abajo  |
| `n` / `Ctrl+N` | 🔍 Buscar nueva palabra         |
| `q` / `ESC`    | ❌ Salir o volver atrás         |
| `Enter`        | ✅ Seleccionar elemento         |
| `Tab`          | 🔄 Cambiar entre paneles        |
| `?` / `h`      | ❓ Mostrar ayuda                |
| `Ctrl+C`       | 🚪 Salir inmediatamente         |

### 🎨 Personalización

**Variables de entorno:**
```bash
# Personalizar colores
export RAE_TUI_THEME="dark"  # dark, light, auto
export RAE_TUI_ACCENT="blue" # blue, green, red, purple

# Configurar timeout
export RAE_TUI_TIMEOUT="10s"
```

**Archivo de configuración (`~/.config/rae-tui/config.yaml`):**
```yaml
theme: "dark"
accent_color: "blue"
timeout: "10s"
cache_enabled: true
cache_duration: "24h"
```

## 🏗️ Desarrollo

### 🔧 Configuración del Entorno

```bash
# Clonar el repositorio
git clone https://github.com/rae-api-com/rae-tui.git
cd rae-tui

# Instalar dependencias
go mod download

# Ejecutar en modo desarrollo
go run . tui
```

### 🧪 Testing

```bash
# Ejecutar tests
make test

# Tests con coverage
make test-coverage

# Ejecutar todos los checks
make ci
```

### 📦 Build

```bash
# Build local
make build

# Build para múltiples plataformas
make build-all

# Build con información de versión
make build VERSION=v1.0.0
```

## 🤝 Contribuir

¡Las contribuciones son súper bienvenidas! 

1. 🍴 Haz fork del proyecto
2. 🌿 Crea tu rama de feature (`git checkout -b feature/AmazingFeature`)
3. 💾 Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. 📤 Push a la rama (`git push origin feature/AmazingFeature`)
5. 🔁 Abre un Pull Request

### 💡 Ideas para Contribuir

- 🎨 Mejoras en la UI/UX
- 🔍 Nuevas funcionalidades de búsqueda
- 🌍 Soporte para más idiomas
- 📱 Integración con otras APIs
- 🐛 Corrección de bugs
- 📚 Mejoras en documentación

## 📈 Roadmap

- [ ] 🌙 Tema oscuro/claro automático
- [ ] 💾 Sistema de cache local
- [ ] 📖 Historial de búsquedas
- [ ] 🔖 Sistema de favoritos
- [ ] 🎵 Pronunciación de palabras
- [ ] 📱 Versión móvil/web
- [ ] 🤖 Integración con ChatGPT/AI

## 📝 Licencia

Este proyecto está bajo la [Licencia MIT](LICENSE).

## 🙏 Reconocimientos

- 🚀 **[go-rae](https://github.com/rae-api-com/go-rae)** - Cliente Go para la API de la RAE
- 🖥️ **[tview](https://github.com/rivo/tview)** - Biblioteca para interfaces de terminal
- 🎬 **[Terminalizer](https://terminalizer.com/)** - Usado para crear los GIFs de demostración
- 📚 **[Real Academia Española](https://www.rae.es/)** - Por mantener el diccionario más completo del español

## 📧 Soporte

¿Tienes algún problema o sugerencia?

- 🐛 [Reportar un bug](https://github.com/rae-api-com/rae-tui/issues/new?template=bug_report.md)
- 💡 [Solicitar una funcionalidad](https://github.com/rae-api-com/rae-tui/issues/new?template=feature_request.md)
- 💬 [Iniciar una discusión](https://github.com/rae-api-com/rae-tui/discussions)
- 📖 [Documentación completa](https://pkg.go.dev/github.com/rae-api-com/rae-tui)

---

<div align="center">

**Hecho con ❤️ para los amantes del español y la terminal**

[⭐ Dale una estrella si te gusta el proyecto](https://github.com/rae-api-com/rae-tui/stargazers) • [🐛 Reportar un problema](https://github.com/rae-api-com/rae-tui/issues) • [💬 Unirse a la discusión](https://github.com/rae-api-com/rae-tui/discussions)

</div>
