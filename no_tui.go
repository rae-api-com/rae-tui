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
	fmt.Printf("  %s0%s. Cancel\n", Yellow, Reset)
	fmt.Printf("\n%sSelect a word (1-%d) or 0 to cancel: %s", Cyan, len(suggestions), Reset)

	// Read user input
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	// Parse selection
	choice, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("%sInvalid input. Please enter a number.%s\n", Red, Reset)
		return ""
	}

	// Validate choice
	if choice == 0 {
		fmt.Printf("%sCancelled.%s\n", Yellow, Reset)
		return ""
	}

	if choice < 1 || choice > len(suggestions) {
		fmt.Printf(
			"%sInvalid choice. Please select a number between 1 and %d.%s\n",
			Red,
			len(suggestions),
			Reset,
		)
		return ""
	}

	return suggestions[choice-1]
}

func renderNoTUI(ctx context.Context, cli *rae.Client, word string) {
	res, err := cli.Word(ctx, word)
	if err != nil {
		if len(res.Suggestions) > 0 {
			fmt.Printf("Did you mean:\n")
			selectedWord := selectWordFromSuggestions(res.Suggestions)
			if selectedWord != "" {
				fmt.Printf("\n%sSearching for: %s%s\n", Bold, selectedWord, Reset)
				renderNoTUI(ctx, cli, selectedWord) // Recursively search with selected word
			}
		} else {
			fmt.Printf("%sNo word found and no suggestions available for: %s%s\n", Red, word, Reset)
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
