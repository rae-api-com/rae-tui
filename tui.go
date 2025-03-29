package main

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"

	rae "github.com/rae-api-com/go-rae"
	"github.com/rivo/tview"
	"github.com/sonirico/stadio/fp"
)

type State struct {
	searching bool
}

type Tui struct {
	cli *rae.Client

	app *tview.Application

	mainLayout *tview.Flex

	header *tview.TextView

	footer *tview.TextView

	resultsView *tview.List

	modalContainer *tview.Flex

	inputField *tview.InputField

	form *tview.Form

	pages *tview.Pages

	state *State
}

func (t *Tui) exit() {
	t.app.Stop()
}

func (t *Tui) goBack() {
	switch {

	case t.state.searching:
		t.state.searching = false
		t.pages.SwitchToPage("main")

	default:
		t.exit()
	}

}

func (t *Tui) search(ctx context.Context, w string) {
	defer func() {
		t.state.searching = false
		t.pages.SwitchToPage("main")
	}()

	res, err := t.cli.Word(ctx, w)
	if err != nil {
		t.modalContainer.SetBackgroundColor(tcell.ColorRed)
		t.inputField.SetText("Palabra no encontrada")
		return
	}

	t.resultsView.Clear()

	for _, meaning := range res.Meanings {
		for _, def := range meaning.Definitions {
			// Formato: "1. Definición: [Descripción]"
			// resultText := fmt.Sprintf("%d. %s: %s", def.MeaningNumber, def.Category, def.Raw)
			resultText := def.Raw
			t.resultsView.AddItem(resultText, "", 0, nil)
		}
		if meaning.Conjugations != nil {
			t.resultsView.AddItem("", "", 0, nil) // Espacio en blanco
			t.resultsView.AddItem(
				"[::b]Conjugaciones:",
				"",
				0,
				nil,
			)

			t.resultsView.AddItem("", "", 0, nil) // Espacio en blanco
			renderConjugations(meaning.Conjugations, t.resultsView)
		}

	}
}

func (t *Tui) handleEvent(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		t.goBack()

	case tcell.KeyRune:
		if !t.state.searching && event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'q':
				t.exit()
			case 'n':
				t.state.searching = true
				t.pages.ShowPage("modal")
				//app.SetFocus(modalContainer)
			case 'j', 'k':
				if event.Rune() == 'j' || event.Key() == tcell.KeyDown {
					t.resultsView.SetCurrentItem(t.resultsView.GetCurrentItem() + 1)
				}
				if event.Rune() == 'k' || event.Key() == tcell.KeyUp {
					t.resultsView.SetCurrentItem(t.resultsView.GetCurrentItem() - 1)
				}
			}
		}
	case tcell.KeyUp:
		if !t.state.searching {
			t.resultsView.SetCurrentItem(t.resultsView.GetCurrentItem() - 1)
		}
	case tcell.KeyDown:
		if !t.state.searching {
			t.resultsView.SetCurrentItem(t.resultsView.GetCurrentItem() + 1)
		}
	default:
	}
	return event
}

func NewTUI(cli *rae.Client) *Tui {
	return &Tui{
		app:            tview.NewApplication(),
		mainLayout:     tview.NewFlex(),
		header:         tview.NewTextView(),
		footer:         tview.NewTextView(),
		resultsView:    tview.NewList(),
		modalContainer: tview.NewFlex(),
		inputField:     tview.NewInputField(),
		form:           tview.NewForm(),
		pages:          tview.NewPages(),
		cli:            cli,
		state:          new(State),
	}
}

func (t *Tui) Run(ctx context.Context, word fp.Option[string]) {
	t.state.searching = word.IsNone()

	t.header.
		SetTextStyle(tcell.StyleDefault.Bold(true)).
		SetText("Diccionario RAE").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorGreen)

	t.footer.
		SetTextStyle(tcell.StyleDefault.Bold(true)).
		SetText("[yellow]↑/k[:] Subir  ↓/j[:] Bajar  n[:] Nueva búsqueda  q/ESC[:] Salir").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorDarkCyan)

	t.resultsView.
		ShowSecondaryText(false)

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p,
					height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	t.modalContainer.
		SetDirection(tview.FlexRow)

	t.inputField.
		SetLabel("Buscar: ").
		SetFieldWidth(20).
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEscape:
				t.goBack()
			case tcell.KeyEnter:
				t.search(ctx, t.inputField.GetText())
			}
		})

	t.form.
		AddFormItem(t.inputField).
		AddButton("Buscar", func() {
			t.search(ctx, t.inputField.GetText())
		}).
		AddButton("Limpiar", func() {
			t.inputField.SetText("")
		})

	t.modalContainer.AddItem(t.form, 0, 1, true)

	t.pages.
		AddPage("main", t.mainLayout, true, true).
		AddPage("modal", modal(t.modalContainer, 40, 10), true, word.IsNone())

	t.mainLayout.
		SetDirection(tview.FlexRow).
		AddItem(t.header, 1, 1, false).
		AddItem(t.resultsView, 0, 10, true).
		AddItem(t.footer, 1, 1, false)

	t.resultsView.SetSelectedFunc(
		func(index int, mainText string, secondaryText string, shortcut rune) {
			// Acción al seleccionar un elemento (opcional)
		},
	)
	t.resultsView.SetMainTextStyle(tcell.StyleDefault.Normal())

	t.resultsView.SetChangedFunc(
		func(index int, mainText string, secondaryText string, shortcut rune) {
			// Resaltar el texto seleccionado
			t.resultsView.SetMainTextStyle(tcell.StyleDefault.Normal())
			t.resultsView.SetItemText(index, mainText, secondaryText)
			t.resultsView.SetMainTextStyle(tcell.StyleDefault.Bold(true))
		},
	)

	// Configurar eventos globales
	t.app.SetInputCapture(t.handleEvent)

	t.pages.ShowPage("main")

	if t.state.searching {
		t.pages.ShowPage("modal")
	} else {
		t.search(ctx, word.UnwrapUnsafe())
	}

	if err := t.app.SetRoot(t.pages, true).Run(); err != nil {
		panic(err)
	}
}

func renderConjugations(conjugations *rae.Conjugations, resultsView *tview.List) {
	if np := conjugations.ConjugationNonPersonal; (np != rae.ConjugationNonPersonal{}) {
		resultsView.AddItem("[::b]Formas no personales", "", 0, nil)
		resultsView.AddItem(fmt.Sprintf("  Infinitivo: %s", np.Infinitive), "", 0, nil)
		resultsView.AddItem(fmt.Sprintf("  Participio: %s", np.Participle), "", 0, nil)
		resultsView.AddItem(fmt.Sprintf("  Gerundio: %s", np.Gerund), "", 0, nil)
	}

	if ind := conjugations.ConjugationIndicative; (ind != rae.ConjugationIndicative{}) {
		resultsView.AddItem("", "", 0, nil)
		resultsView.AddItem("[::b]Modo Indicativo", "", 0, nil)
		renderVerbalConjugation("Presente", ind.Present, resultsView)
		renderVerbalConjugation("Pretérito Imperfecto", ind.Imperfect, resultsView)
		renderVerbalConjugation("Pretérito Perfecto Simple", ind.Preterite, resultsView)
		renderVerbalConjugation("Futuro", ind.Future, resultsView)
		renderVerbalConjugation("Condicional", ind.Conditional, resultsView)
	}

	if subj := conjugations.ConjugationSubjunctive; (subj != rae.ConjugationSubjunctive{}) {
		resultsView.AddItem("", "", 0, nil)
		resultsView.AddItem("[::b]Modo Subjuntivo", "", 0, nil)
		renderVerbalConjugation("Presente", subj.Present, resultsView)
		renderVerbalConjugation("Pretérito Imperfecto", subj.Imperfect, resultsView)
		renderVerbalConjugation("Futuro", subj.Future, resultsView)
	}

	if imp := conjugations.ConjugationImperative; (imp != rae.ConjugationImperative{}) {
		resultsView.AddItem("", "", 0, nil)
		resultsView.AddItem("[::b]Modo Imperativo", "", 0, nil)
		resultsView.AddItem(
			fmt.Sprintf("  2da Persona Singular: %s", imp.SingularSecondPerson),
			"",
			0,
			nil,
		)
		resultsView.AddItem(
			fmt.Sprintf("  2da Persona Plural: %s", imp.PluralSecondPerson),
			"",
			0,
			nil,
		)
	}
}

func renderVerbalConjugation(title string, conj rae.Conjugation, resultsView *tview.List) {
	resultsView.AddItem(fmt.Sprintf("[::b]%s", title), "", 0, nil)
	resultsView.AddItem(fmt.Sprintf("  Yo %s", conj.SingularFirstPerson), "", 0, nil)
	resultsView.AddItem(fmt.Sprintf("  Tú %s", conj.SingularSecondPerson), "", 0, nil)
	resultsView.AddItem(fmt.Sprintf("  Él/Ella/Usted %s", conj.SingularThirdPerson), "", 0, nil)
	resultsView.AddItem(fmt.Sprintf("  Nosotros %s", conj.PluralFirstPerson), "", 0, nil)
	resultsView.AddItem(fmt.Sprintf("  Vosotros %s", conj.PluralSecondPerson), "", 0, nil)
	resultsView.AddItem(
		fmt.Sprintf("  Ellos/Ellas/Ustedes: %s", conj.PluralThirdPerson),
		"",
		0,
		nil,
	)
}
