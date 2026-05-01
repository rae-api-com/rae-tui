package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	rae "github.com/rae-api-com/go-rae"
	"github.com/sonirico/vago/fp"
)

// These variables are set during build using ldflags
var (
	// version is set at build time from git tag
	version = "dev"
	// commit is the git commit hash
	commit = "unknown"
	// buildTime is the build timestamp
	buildTime = "unknown"
)

type args struct {
	word fp.Option[string]
	tui  bool
}

func printHelp() {
	fmt.Println("RAE Dictionary CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  rae-tui [WORD]        - Search for a word in RAE dictionary in CLI mode")
	fmt.Println("  rae-tui tui [WORD]    - Open the TUI interface with optional initial word")
	fmt.Println("\nOptions:")
	fmt.Println("  -h, --help            - Display this help message")
	fmt.Println("  -v, --version         - Display version information")
	fmt.Println("\nExamples:")
	fmt.Println("  rae-tui hola          - Show definition of 'hola' in CLI mode")
	fmt.Println("  rae-tui tui           - Open TUI interface")
	fmt.Println("  rae-tui tui casa      - Open TUI interface and search for 'casa'")
	os.Exit(0)
}

func printVersion() {
	fmt.Printf("rae-tui %s\n", version)
	fmt.Printf("commit: %s\n", commit)
	fmt.Printf("built: %s\n", buildTime)
	os.Exit(0)
}

func parseArgs() args {
	// Check for version flag
	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" || arg == "version" {
			printVersion()
			return args{}
		}

		if arg == "-h" || arg == "--help" || arg == "help" {
			printHelp()
			return args{}
		}
	}

	if len(os.Args) > 2 {
		return args{
			word: fp.Some(strings.TrimSpace(os.Args[2])),
			tui:  strings.TrimSpace(os.Args[1]) == "tui",
		}
	}

	if len(os.Args) > 1 {
		uniqArg := strings.TrimSpace(os.Args[1])

		if uniqArg == "tui" {
			return args{
				tui: true,
			}
		}

		return args{
			word: fp.Some(uniqArg),
			tui:  false,
		}
	}

	return args{
		word: fp.None[string](),
		tui:  true,
	}
}

func main() {
	ctx := context.Background()
	cli := rae.New(rae.WithVersion(version))

	arguments := parseArgs()

	if arguments.tui {
		NewTUI(cli).Run(ctx, arguments.word)
	} else {
		renderNoTUI(ctx, cli, arguments.word.UnwrapUnsafe())
	}
}
