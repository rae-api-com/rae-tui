package main

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"

	rae "github.com/rae-api-com/go-rae"
	"github.com/rivo/tview"
	"github.com/sonirico/vago/fp"
)

type State struct {
	searching   bool
	suggestions bool
	fuzzySearch bool
}

type Tui struct {
	cli *rae.Client
	app *tview.Application

	// Layout components
	mainLayout      *tview.Flex
	header          *tview.TextView
	footer          *tview.TextView
	resultsView     *tview.List
	suggestionsList *tview.List

	// Search modal
	modalContainer *tview.Flex
	inputField     *tview.InputField
	form           *tview.Form

	// Pages
	pages *tview.Pages

	// State
	state *State
}

func NewTUI(cli *rae.Client) *Tui {
	return &Tui{
		cli:             cli,
		app:             tview.NewApplication(),
		mainLayout:      tview.NewFlex(),
		header:          tview.NewTextView(),
		footer:          tview.NewTextView(),
		resultsView:     tview.NewList(),
		suggestionsList: tview.NewList(),
		modalContainer:  tview.NewFlex(),
		inputField:      tview.NewInputField(),
		form:            tview.NewForm(),
		pages:           tview.NewPages(),
		state:           &State{},
	}
}

func (t *Tui) Run(ctx context.Context, word fp.Option[string]) {
	t.state.searching = word.IsNone()
	t.setupUI()
	t.setupPages()
	t.setupEventHandlers()

	t.pages.SwitchToPage("main")

	if t.state.searching {
		t.pages.SwitchToPage("modal")
	} else {
		t.search(ctx, word.UnwrapUnsafe())
	}

	if err := t.app.SetRoot(t.pages, true).Run(); err != nil {
		panic(err)
	}
}

func (t *Tui) setupUI() {
	// Header
	t.header.
		SetTextStyle(tcell.StyleDefault.Bold(true)).
		SetText("Diccionario RAE").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorGreen)

	// Footer
	t.updateFooter()
	t.footer.
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorDarkCyan)

	// Results view
	t.resultsView.ShowSecondaryText(false)
	t.resultsView.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		// Optional action when selecting a result
	})

	// Suggestions/Fuzzy search list
	t.suggestionsList.ShowSecondaryText(false)
	t.suggestionsList.SetSelectedStyle(tcell.StyleDefault.
		Foreground(tcell.ColorYellow).
		Background(tcell.ColorDarkBlue).
		Bold(true))
	t.suggestionsList.SetSelectedFunc(
		func(index int, mainText, secondaryText string, shortcut rune) {
			if secondaryText != "" {
				t.selectWord(secondaryText)
			}
		},
	)

	// Main layout
	t.mainLayout.
		SetDirection(tview.FlexRow).
		AddItem(t.header, 1, 1, false).
		AddItem(t.resultsView, 0, 10, true).
		AddItem(t.footer, 1, 1, false)

	// Search modal
	t.setupSearchModal()
}

func (t *Tui) setupSearchModal() {
	t.modalContainer.SetDirection(tview.FlexRow)

	t.inputField.
		SetLabel("Buscar: ").
		SetFieldWidth(20).
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEscape:
				t.goBack()
			case tcell.KeyEnter:
				t.search(context.Background(), t.inputField.GetText())
			}
		})

	t.form.
		AddFormItem(t.inputField).
		AddButton("Buscar", func() {
			t.search(context.Background(), t.inputField.GetText())
		}).
		AddButton("Limpiar", func() {
			t.inputField.SetText("")
		})

	t.modalContainer.AddItem(t.form, 0, 1, true)
}

func (t *Tui) setupPages() {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	// Full-page layout for suggestions/fuzzy search
	listLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.header, 1, 1, false).
		AddItem(t.suggestionsList, 0, 10, true).
		AddItem(t.footer, 1, 1, false)

	t.pages.
		AddPage("main", t.mainLayout, true, true).
		AddPage("modal", modal(t.modalContainer, 40, 10), true, false).
		AddPage("list", listLayout, true, false)
}

func (t *Tui) setupEventHandlers() {
	t.app.SetInputCapture(t.handleEvent)
}

func (t *Tui) updateFooter() {
	var text string
	switch {
	case t.state.suggestions:
		text = "[yellow]↑/k[:] Subir  ↓/j[:] Bajar  Enter/1-9[:] Seleccionar  q/ESC[:] Volver"
	case t.state.fuzzySearch:
		text = "[yellow]↑/k[:] Subir  ↓/j[:] Bajar  Enter[:] Seleccionar  q/ESC[:] Volver"
	case t.state.searching:
		text = "[yellow]Enter[:] Buscar  ESC[:] Cancelar"
	default:
		text = "[yellow]↑/k[:] Subir  ↓/j[:] Bajar  n[:] Nueva búsqueda  q/ESC[:] Salir"
	}
	t.footer.SetText(text)
	t.footer.SetTextStyle(tcell.StyleDefault.Bold(true))
}

func (t *Tui) resetState() {
	t.state.searching = false
	t.state.suggestions = false
	t.state.fuzzySearch = false
	t.updateFooter()
}

func (t *Tui) goBack() {
	switch {
	case t.state.suggestions, t.state.fuzzySearch:
		t.resetState()
		t.pages.SwitchToPage("main")
	case t.state.searching:
		t.resetState()
		t.pages.SwitchToPage("main")
	default:
		t.exit()
	}
}

func (t *Tui) exit() {
	t.app.Stop()
}

func (t *Tui) selectWord(word string) {
	t.resetState()
	t.search(context.Background(), word)
}

func (t *Tui) search(ctx context.Context, word string) {
	res, err := t.cli.Word(ctx, word)
	if err != nil {
		if len(res.Suggestions) > 0 {
			t.showSuggestions(res.Suggestions)
			return
		}
		t.showFuzzySearchResults(ctx, word)
		return
	}

	t.displayResults(res)
}

func (t *Tui) displayResults(res rae.WordEntry) {
	t.resetState()
	t.resultsView.Clear()

	for _, meaning := range res.Meanings {
		for _, def := range meaning.Definitions {
			t.resultsView.AddItem(def.Raw, "", 0, nil)
		}

		if meaning.Conjugations != nil {
			t.resultsView.AddItem("", "", 0, nil)
			t.resultsView.AddItem("[::b]Conjugaciones:", "", 0, nil)
			t.resultsView.AddItem("", "", 0, nil)
			renderConjugations(meaning.Conjugations, t.resultsView)
		}
	}

	t.pages.SwitchToPage("main")
}

func (t *Tui) showSuggestions(suggestions []string) {
	t.resetState()
	t.state.suggestions = true
	t.updateFooter()

	t.suggestionsList.Clear()
	t.suggestionsList.AddItem("[yellow]¿Quisiste decir?", "", 0, nil)
	t.suggestionsList.AddItem("", "", 0, nil)

	for i, suggestion := range suggestions {
		text := fmt.Sprintf("%d. %s", i+1, suggestion)
		t.suggestionsList.AddItem(text, suggestion, rune('0'+i+1), nil)
	}

	t.suggestionsList.AddItem("", "", 0, nil)
	t.suggestionsList.AddItem("0. Cancelar", "", '0', nil)

	t.pages.SwitchToPage("list")
}

func (t *Tui) showFuzzySearchResults(ctx context.Context, word string) {
	t.resetState()
	t.state.fuzzySearch = true
	t.updateFooter()

	searchResults, err := t.cli.Search(ctx, word)
	if err != nil || len(searchResults) == 0 {
		t.resetState()
		t.showError("No se encontraron resultados de búsqueda difusa")
		return
	}

	t.suggestionsList.Clear()
	t.suggestionsList.AddItem("[yellow]Búsqueda difusa - Resultados encontrados:", "", 0, nil)
	t.suggestionsList.AddItem("", "", 0, nil)

	for _, result := range searchResults {
		searchWord := result.Doc.Word
		wordEntry, err := result.WordEntry()

		var text string
		if err == nil && wordEntry != nil && len(wordEntry.Meanings) > 0 &&
			len(wordEntry.Meanings[0].Definitions) > 0 {
			def := wordEntry.Meanings[0].Definitions[0].Raw
			preview := def
			if len(def) > 70 {
				preview = def[:70] + "..."
			}
			text = fmt.Sprintf("[yellow][::b]%s[white] - %s", searchWord, preview)
		} else {
			text = fmt.Sprintf("[yellow][::b]%s", searchWord)
		}

		t.suggestionsList.AddItem(text, searchWord, 0, nil)
	}

	t.pages.SwitchToPage("list")
}

func (t *Tui) showError(message string) {
	t.inputField.SetText(message)
	t.modalContainer.SetBackgroundColor(tcell.ColorRed)
	t.state.searching = true
	t.updateFooter()
	t.pages.SwitchToPage("modal")
}

func (t *Tui) handleEvent(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		t.goBack()
		return nil

	case tcell.KeyEnter:
		if t.state.suggestions || t.state.fuzzySearch {
			_, selectedWord := t.suggestionsList.GetItemText(t.suggestionsList.GetCurrentItem())
			if selectedWord != "" {
				t.selectWord(selectedWord)
			}
			return nil
		}

	case tcell.KeyUp:
		if t.state.suggestions || t.state.fuzzySearch {
			idx := t.suggestionsList.GetCurrentItem()
			if idx > 0 {
				t.suggestionsList.SetCurrentItem(idx - 1)
			}
			return nil
		} else if !t.state.searching {
			idx := t.resultsView.GetCurrentItem()
			if idx > 0 {
				t.resultsView.SetCurrentItem(idx - 1)
			}
			return nil
		}

	case tcell.KeyDown:
		if t.state.suggestions || t.state.fuzzySearch {
			idx := t.suggestionsList.GetCurrentItem()
			if idx < t.suggestionsList.GetItemCount()-1 {
				t.suggestionsList.SetCurrentItem(idx + 1)
			}
			return nil
		} else if !t.state.searching {
			idx := t.resultsView.GetCurrentItem()
			if idx < t.resultsView.GetItemCount()-1 {
				t.resultsView.SetCurrentItem(idx + 1)
			}
			return nil
		}

	case tcell.KeyRune:
		// If searching, let the input field handle all runes (don't intercept 'q')
		if t.state.searching {
			// Let the input field handle all runes normally
			return event
		}

		return t.handleRune(event.Rune())
	}

	return event
}

func (t *Tui) handleRune(r rune) *tcell.EventKey {
	switch r {
	case 'q':
		t.exit()
		return nil

	case 'j':
		if t.state.suggestions || t.state.fuzzySearch {
			idx := t.suggestionsList.GetCurrentItem()
			if idx < t.suggestionsList.GetItemCount()-1 {
				t.suggestionsList.SetCurrentItem(idx + 1)
			}
			return nil
		}
		idx := t.resultsView.GetCurrentItem()
		if idx < t.resultsView.GetItemCount()-1 {
			t.resultsView.SetCurrentItem(idx + 1)
		}
		return nil

	case 'k':
		if t.state.suggestions || t.state.fuzzySearch {
			idx := t.suggestionsList.GetCurrentItem()
			if idx > 0 {
				t.suggestionsList.SetCurrentItem(idx - 1)
			}
			return nil
		}
		idx := t.resultsView.GetCurrentItem()
		if idx > 0 {
			t.resultsView.SetCurrentItem(idx - 1)
		}
		return nil

	case 'n':
		t.state.searching = true
		t.inputField.SetText("") // Clear input
		t.updateFooter()
		t.pages.SwitchToPage("modal")
		t.app.SetFocus(t.inputField) // Set focus to input field
		return nil

	default:
		// Handle number selection for suggestions only (not fuzzy search)
		if t.state.suggestions && r >= '0' && r <= '9' {
			num := int(r - '0')
			if num == 0 {
				t.goBack()
				return nil
			}

			itemCount := t.suggestionsList.GetItemCount()
			for i := 2; i < itemCount-2; i++ {
				text, secondary := t.suggestionsList.GetItemText(i)
				if len(text) > 0 && text[0] == byte('0'+num) && secondary != "" {
					t.selectWord(secondary)
					return nil
				}
			}
		}
	}

	return nil
}

func renderConjugations(conjugations *rae.Conjugations, resultsView *tview.List) {
	if np := conjugations.ConjugationNonPersonal; (np != rae.ConjugationNonPersonal{}) {
		resultsView.AddItem("[yellow][::b]Formas no personales[white]", "", 0, nil)
		resultsView.AddItem(
			fmt.Sprintf(
				"  [cyan]Infinitivo:[white] %s  [cyan]Participio:[white] %s  [cyan]Gerundio:[white] %s",
				np.Infinitive,
				np.Participle,
				np.Gerund,
			),
			"",
			0,
			nil,
		)
		resultsView.AddItem("", "", 0, nil)
	}

	if ind := conjugations.ConjugationIndicative; (ind != rae.ConjugationIndicative{}) {
		resultsView.AddItem("[yellow][::b]Modo Indicativo[white]", "", 0, nil)
		renderVerbalConjugation("Presente", ind.Present, resultsView)
		renderVerbalConjugation("Pretérito Imperfecto", ind.Imperfect, resultsView)
		renderVerbalConjugation("Pretérito Perfecto Simple", ind.Preterite, resultsView)
		renderVerbalConjugation("Futuro", ind.Future, resultsView)
		renderVerbalConjugation("Condicional", ind.Conditional, resultsView)
		resultsView.AddItem("", "", 0, nil)
	}

	if subj := conjugations.ConjugationSubjunctive; (subj != rae.ConjugationSubjunctive{}) {
		resultsView.AddItem("[yellow][::b]Modo Subjuntivo[white]", "", 0, nil)
		renderVerbalConjugation("Presente", subj.Present, resultsView)
		renderVerbalConjugation("Pretérito Imperfecto", subj.Imperfect, resultsView)
		renderVerbalConjugation("Futuro", subj.Future, resultsView)
		resultsView.AddItem("", "", 0, nil)
	}

	if imp := conjugations.ConjugationImperative; (imp != rae.ConjugationImperative{}) {
		resultsView.AddItem("[yellow][::b]Modo Imperativo[white]", "", 0, nil)
		resultsView.AddItem(
			fmt.Sprintf("  [cyan]Tú:[white] %s  [cyan]Vosotros:[white] %s",
				imp.SingularSecondPerson, imp.PluralSecondPerson),
			"",
			0,
			nil,
		)
	}
}

func renderVerbalConjugation(title string, conj rae.Conjugation, resultsView *tview.List) {
	resultsView.AddItem(fmt.Sprintf("  [cyan][::b]%s[white]", title), "", 0, nil)

	// Primera persona singular y plural en una línea
	resultsView.AddItem(
		fmt.Sprintf("    [yellow]Yo:[white] %-25s  [yellow]Nosotros:[white] %s",
			conj.SingularFirstPerson, conj.PluralFirstPerson),
		"",
		0,
		nil,
	)

	// Segunda persona singular y plural
	resultsView.AddItem(
		fmt.Sprintf("    [yellow]Tú:[white] %-25s  [yellow]Vosotros:[white] %s",
			conj.SingularSecondPerson, conj.PluralSecondPerson),
		"",
		0,
		nil,
	)

	// Tercera persona singular y plural
	resultsView.AddItem(
		fmt.Sprintf(
			"    [yellow]Él/Ella/Usted:[white] %-15s  [yellow]Ellos/Ellas/Ustedes:[white] %s",
			conj.SingularThirdPerson,
			conj.PluralThirdPerson,
		),
		"",
		0,
		nil,
	)
}
