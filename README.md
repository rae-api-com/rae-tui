# RAE-TUI

[![Go Version](https://img.shields.io/github/go-mod/go-version/rae-api-com/rae-tui)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/rae-api-com/rae-tui)](https://goreportcard.com/report/github.com/rae-api-com/rae-tui)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/rae-api-com/rae-tui)](https://github.com/rae-api-com/rae-tui/releases)
[![Build Status](https://github.com/rae-api-com/rae-tui/workflows/CI/badge.svg)](https://github.com/rae-api-com/rae-tui/actions)

Cliente de diccionario español para la terminal, conectado a la Real Academia Española (RAE).

![RAE-TUI Demo](demo.gif)

## Características

- **Búsqueda instantánea** - Consulta palabras directamente desde tu terminal
- **Definiciones completas** - Visualiza significados detallados y acepciones múltiples
- **Búsqueda difusa** - Cuando no se encuentra una palabra exacta, busca automáticamente términos similares con resultados relevantes
- **Interfaz interactiva** - TUI moderna con navegación por teclado intuitiva
- **Conjugaciones verbales** - Muestra todas las conjugaciones en todos los tiempos
- **Modo CLI** - Consultas no-interactivas para uso en scripts
- **Sin configuración** - Funciona inmediatamente sin API keys

## Instalación

### Instalación rápida

```bash
go install github.com/rae-api-com/rae-tui@latest
```

### Binarios precompilados

Descarga el binario apropiado para tu plataforma desde la página de [Releases](https://github.com/rae-api-com/rae-tui/releases).

#### Linux/macOS

```bash
curl -sf https://gobinaries.com/rae-api-com/rae-tui | sh
```

### Desde código fuente

```bash
git clone https://github.com/rae-api-com/rae-tui.git
cd rae-tui
go build -o rae-tui
```

## Uso

### Modo interactivo (TUI)

Lanza la interfaz interactiva completa:

```bash
rae-tui
# o explícitamente
rae-tui tui
```

Búsqueda directa con una palabra:

```bash
rae-tui tui hola
```

### Modo CLI

Consulta directa para scripts o uso rápido:

```bash
rae-tui palabra
# Ejemplo:
rae-tui computadora
```

#### Búsqueda difusa automática

Cuando una palabra no se encuentra en el diccionario y no hay sugerencias disponibles, el CLI automáticamente ejecuta una búsqueda difusa para encontrar palabras similares. Puedes seleccionar una de las opciones encontradas:

```bash
rae-tui "persona que programa"
```

Si no hay resultados exactos, se mostrará una lista de palabras similares con una vista previa de sus definiciones:

```
No se encontró la palabra y no hay sugerencias disponibles para: persona que programa
Buscando resultados difusos...

Búsqueda difusa - Resultados encontrados:
  1. programador - 1. adj. Que programa. U. t. c. s.
  2. programadora - 1. adj. Que programa. U. t. c. s.
  3. programadores - 1. adj. Que programa. U. t. c. s.
  4. coguionista - 1. m. y f. Persona que escribe junto con otra u otras el gui...

Selecciona una palabra (1-4) o 0 para cancelar:
```

### Ejemplos prácticos

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

### Atajos de teclado (Modo TUI)

| Tecla          | Acción                         |
| -------------- | ------------------------------ |
| `↑` / `k`      | Mover selección hacia arriba    |
| `↓` / `j`      | Mover selección hacia abajo     |
| `n` / `Ctrl+N` | Buscar nueva palabra            |
| `q` / `ESC`    | Salir o volver atrás            |
| `Enter`        | Seleccionar elemento            |
| `Tab`          | Cambiar entre paneles           |
| `?` / `h`      | Mostrar ayuda                   |
| `Ctrl+C`       | Salir inmediatamente            |

## Desarrollo

### Configuración del entorno

```bash
# Clonar el repositorio
git clone https://github.com/rae-api-com/rae-tui.git
cd rae-tui

# Instalar dependencias
go mod download

# Ejecutar en modo desarrollo
go run . tui
```

### Testing

```bash
# Ejecutar tests
make test

# Tests con coverage
make test-coverage

# Ejecutar todos los checks
make ci
```

### Build

```bash
# Build local
make build

# Build para múltiples plataformas
make build-all

# Build con información de versión
make build VERSION=v1.0.0
```

## Contribuir

Las contribuciones son bienvenidas. 

1. Haz fork del proyecto
2. Crea tu rama de feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

### Ideas para contribuir

- Mejoras en la UI/UX
- Nuevas funcionalidades de búsqueda
- Soporte para más idiomas
- Integración con otras APIs
- Corrección de bugs
- Mejoras en documentación

## Roadmap

- [ ] Tema oscuro/claro automático
- [ ] Sistema de cache local
- [ ] Historial de búsquedas
- [ ] Sistema de favoritos
- [ ] Pronunciación de palabras

## Licencia

Este proyecto está bajo la [Licencia MIT](LICENSE).

## Reconocimientos

- **[go-rae](https://github.com/rae-api-com/go-rae)** - Cliente Go para la API de la RAE
- **[tview](https://github.com/rivo/tview)** - Biblioteca para interfaces de terminal
- **[Real Academia Española](https://www.rae.es/)** - Por mantener el diccionario más completo del español

## Soporte

¿Tienes algún problema o sugerencia?

- [Reportar un bug](https://github.com/rae-api-com/rae-tui/issues/new?template=bug_report.md)
- [Solicitar una funcionalidad](https://github.com/rae-api-com/rae-tui/issues/new?template=feature_request.md)
- [Iniciar una discusión](https://github.com/rae-api-com/rae-tui/discussions)
- [Documentación completa](https://pkg.go.dev/github.com/rae-api-com/rae-tui)

---

<div align="center">

Hecho con ❤️ para los amantes del español y la terminal

[⭐ Dale una estrella si te gusta el proyecto](https://github.com/rae-api-com/rae-tui/stargazers) • [🐛 Reportar un problema](https://github.com/rae-api-com/rae-tui/issues) • [💬 Unirse a la discusión](https://github.com/rae-api-com/rae-tui/discussions)

</div>
