package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	rae "github.com/rae-api-com/go-rae"
)

// Colores ANSI
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Bold   = "\033[1m"
)

// selectWordFromSuggestions displays a list of suggested words and allows the user to select one
func selectWordFromSuggestions(suggestions []string) string {
	if len(suggestions) == 0 {
		return ""
	}

	// Display numbered list of suggestions
	for i, suggestion := range suggestions {
		fmt.Printf("  %s%d%s. %s\n", Yellow, i+1, Reset, suggestion)
	}
	fmt.Printf("  %s0%s. Cancelar\n", Yellow, Reset)
	fmt.Printf("\n%sSelecciona una palabra (1-%d) o 0 para cancelar: %s", Cyan, len(suggestions), Reset)

	// Read user input
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	// Parse selection
	choice, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("%sEntrada inválida. Por favor ingresa un número.%s\n", Red, Reset)
		return ""
	}

	// Validate choice
	if choice == 0 {
		fmt.Printf("%sCancelado.%s\n", Yellow, Reset)
		return ""
	}

	if choice < 1 || choice > len(suggestions) {
		fmt.Printf(
			"%sOpción inválida. Por favor selecciona un número entre 1 y %d.%s\n",
			Red,
			len(suggestions),
			Reset,
		)
		return ""
	}

	return suggestions[choice-1]
}

// selectWordFromSearchResults displays a list of search results and allows the user to select one
func selectWordFromSearchResults(searchResults []rae.SearchResult) string {
	if len(searchResults) == 0 {
		return ""
	}

	fmt.Printf("\n%sBúsqueda difusa - Resultados encontrados:%s\n", Bold, Reset)

	// Display numbered list of search results
	for i, result := range searchResults {
		wordEntry, err := result.WordEntry()
		if err == nil && wordEntry != nil {
			// Show word and a preview of the first definition if available
			preview := ""
			if len(wordEntry.Meanings) > 0 && len(wordEntry.Meanings[0].Definitions) > 0 {
				def := wordEntry.Meanings[0].Definitions[0].Raw
				if len(def) > 60 {
					preview = def[:60] + "..."
				} else {
					preview = def
				}
			}
			if preview != "" {
				fmt.Printf("  %s%d%s. %s%s%s - %s%s%s\n", Yellow, i+1, Reset, Bold, result.Doc.Word, Reset, Cyan, preview, Reset)
			} else {
				fmt.Printf("  %s%d%s. %s%s%s\n", Yellow, i+1, Reset, Bold, result.Doc.Word, Reset)
			}
		} else {
			fmt.Printf("  %s%d%s. %s%s%s\n", Yellow, i+1, Reset, Bold, result.Doc.Word, Reset)
		}
	}
	fmt.Printf("  %s0%s. Cancelar\n", Yellow, Reset)
	fmt.Printf("\n%sSelecciona una palabra (1-%d) o 0 para cancelar: %s", Cyan, len(searchResults), Reset)

	// Read user input
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	// Parse selection
	choice, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("%sEntrada inválida. Por favor ingresa un número.%s\n", Red, Reset)
		return ""
	}

	// Validate choice
	if choice == 0 {
		fmt.Printf("%sCancelado.%s\n", Yellow, Reset)
		return ""
	}

	if choice < 1 || choice > len(searchResults) {
		fmt.Printf(
			"%sOpción inválida. Por favor selecciona un número entre 1 y %d.%s\n",
			Red,
			len(searchResults),
			Reset,
		)
		return ""
	}

	return searchResults[choice-1].Doc.Word
}

func renderNoTUI(ctx context.Context, cli *rae.Client, word string) {
	res, err := cli.Word(ctx, word)
	if err != nil {
		if len(res.Suggestions) > 0 {
			fmt.Printf("¿Quisiste decir:\n")
			selectedWord := selectWordFromSuggestions(res.Suggestions)
			if selectedWord != "" {
				fmt.Printf("\n%sBuscando: %s%s\n", Bold, selectedWord, Reset)
				renderNoTUI(ctx, cli, selectedWord) // Recursively search with selected word
			}
		} else {
			// No word found and no suggestions, try fuzzy search
			fmt.Printf("%sNo se encontró la palabra y no hay sugerencias disponibles para: %s%s\n", Yellow, word, Reset)
			fmt.Printf("%sBuscando resultados difusos...%s\n", Cyan, Reset)

			searchResults, searchErr := cli.Search(ctx, word)
			if searchErr != nil || len(searchResults) == 0 {
				fmt.Printf("%sNo se encontraron resultados de búsqueda difusa para: %s%s\n", Red, word, Reset)
				return
			}

			selectedWord := selectWordFromSearchResults(searchResults)
			if selectedWord != "" {
				fmt.Printf("\n%sBuscando: %s%s\n", Bold, selectedWord, Reset)
				renderNoTUI(ctx, cli, selectedWord) // Recursively search with selected word
			}
		}
		return
	}

	fmt.Printf("\n%sPalabra: %s%s\n\n", Bold, res.Word, Reset)
	for i, meaning := range res.Meanings {
		fmt.Printf("%sSignificado %d:%s\n", Green, i+1, Reset)
		for _, definition := range meaning.Definitions {
			fmt.Printf("  - %s (%s%s%s)\n", definition.Raw, Bold, definition.Category, Reset)
		}

		if meaning.Conjugations != nil {
			fmt.Printf("\n  %sConjugaciones%s\n", Bold, Reset)

			fmt.Printf("    %sFormas no personales%s\n", Bold, Reset)
			fmt.Printf(
				"      %sInfinitivo%s %s\n",
				Cyan,
				Reset,
				meaning.Conjugations.ConjugationNonPersonal.Infinitive,
			)
			fmt.Printf(
				"      %sParticipio%s %s\n",
				Cyan,
				Reset,
				meaning.Conjugations.ConjugationNonPersonal.Participle,
			)
			fmt.Printf(
				"      %sGerundio%s %s\n",
				Cyan,
				Reset,
				meaning.Conjugations.ConjugationNonPersonal.Gerund,
			)

			fmt.Printf("\n    %sModo Indicativo%s\n", Bold, Reset)
			printConjugations(meaning.Conjugations.ConjugationIndicative)

			fmt.Printf("\n    %sModo Subjuntivo%s\n", Bold, Reset)
			printConjugations(meaning.Conjugations.ConjugationSubjunctive)

			fmt.Printf("\n    %sModo Imperativo%s\n", Bold, Reset)
			imperativo := meaning.Conjugations.ConjugationImperative
			fmt.Printf("      %sTú%s %s\n", Cyan, Reset, imperativo.SingularSecondPerson)
			fmt.Printf("      %sUsted%s %s\n", Cyan, Reset, imperativo.SingularFormalSecondPerson)
			fmt.Printf("      %sVosotros%s %s\n", Cyan, Reset, imperativo.PluralSecondPerson)
			fmt.Printf("      %sUstedes%s %s\n", Cyan, Reset, imperativo.PluralFormalSecondPerson)
		}

		fmt.Println() // Separar significados con una línea en blanco
	}
}

func printConjugations(conjugation any) {
	personas := []string{"Yo", "Tú", "Él/Ella", "Nosotros", "Vosotros", "Ellos/Ellas"}

	switch c := conjugation.(type) {
	case rae.ConjugationIndicative:
		printTense("Presente", c.Present, personas)
		printTense("Pretérito perfecto compuesto", c.PresentPerfect, personas)
		printTense("Imperfecto", c.Imperfect, personas)
		printTense("Pretérito pluscuamperfecto", c.PastPerfect, personas)
		printTense("Pretérito perfecto simple", c.Preterite, personas)
		printTense("Pretérito anterior", c.PastAnterior, personas)
		printTense("Futuro simple", c.Future, personas)
		printTense("Futuro compuesto", c.FuturePerfect, personas)
		printTense("Condicional simple", c.Conditional, personas)
		printTense("Condicional compuesto", c.ConditionalPerfect, personas)
	case rae.ConjugationSubjunctive:
		printTense("Presente", c.Present, personas)
		printTense("Pretérito perfecto compuesto", c.PresentPerfect, personas)
		printTense("Imperfecto", c.Imperfect, personas)
		printTense("Pretérito pluscuamperfecto", c.PastPerfect, personas)
		printTense("Futuro simple", c.Future, personas)
		printTense("Futuro compuesto", c.FuturePerfect, personas)
	}
}

func printTense(tenseName string, tense rae.Conjugation, personas []string) {
	fmt.Printf("      %s%s%s", Bold, tenseName, Reset)
	for i, persona := range personas {
		if i%3 == 0 {
			fmt.Println()
		}
		fmt.Printf(
			"        %s%-8s%s %-40s\t",
			Cyan,
			persona,
			Reset,
			getConjugationForPersona(tense, persona),
		)
	}
	fmt.Println()
}

func getConjugationForPersona(conjugation rae.Conjugation, persona string) string {
	switch persona {
	case "Yo":
		return conjugation.SingularFirstPerson
	case "Tú":
		return conjugation.SingularSecondPerson
	case "Él/Ella":
		return conjugation.SingularThirdPerson
	case "Nosotros":
		return conjugation.PluralFirstPerson
	case "Vosotros":
		return conjugation.PluralSecondPerson
	case "Ellos/Ellas":
		return conjugation.PluralThirdPerson
	default:
		return ""
	}
}
