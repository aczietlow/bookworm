2025.11.16

Trying to work on a refactor

```go 
	app := conf.tui.app
	pages := conf.tui.pages
	registry := conf.registry

	var views []*tview.Primitive
	var registryOrder []string
	for k := range registry {
		registryOrder = append(registryOrder, k)
	}

	commands := tview.NewList().ShowSecondaryText(false)
	commands.SetTitle("Functions").SetBorder(true).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			return tcell.NewEventKey(tcell.KeyDown, rune(0), tcell.ModNone)
		} else if event.Rune() == 'k' {
			return tcell.NewEventKey(tcell.KeyUp, rune(0), tcell.ModNone)
		}
		return event
	})

	results := tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		})
	results.SetTitle("Results").SetBorder(true)

	for i, k := range registryOrder {
		c := registry[k]
		commands.AddItem(c.name, "", 0, nil)
		view := c.view(conf)
		views = append(views, &view)

		switch v := view.(type) {
		case *tview.InputField:
			v.SetDoneFunc(func(key tcell.Key) {
				// TODO: Only using reflection to check if model exists. Will remove this when every review has a model by default
				if !reflect.ValueOf(c.model).IsZero() {
					if key == tcell.KeyEnter {
						searchText := v.GetText()

						result, err := c.callback(conf, searchText)
						if err != nil {
							panic(err)
						}
						results.SetText(fmt.Sprintf("%s", result))
					} else if key == tcell.KeyEsc {
						app.SetFocus(commands)
					}
				}
				// end of TODO: experiment
				if key == tcell.KeyEnter {
					searchText := v.GetText()
					result, err := c.callback(conf, searchText)
					if err != nil {
						panic(err)
					}
					results.SetText(fmt.Sprintf("%s", result))
				} else if key == tcell.KeyEsc {
					app.SetFocus(commands)
				}
			})
		// This is the default primatitive, if there is no interactivity just call the command callback immediately.
		case *tview.Box:
			v.SetFocusFunc(func() {
				result, err := c.callback(conf)
				if err != nil {
					panic(err)
				}
				results.SetText(fmt.Sprintf("%s", result))
				app.SetFocus(commands)
			})

		}

		flexLayout := tview.NewFlex().
			AddItem(commands, 0, 1, true).
			AddItem(view, 0, 1, false).
			AddItem(results, 0, 3, false)

		pages.AddPage(c.name, flexLayout, true, i == 0)
	}

	commands.SetSelectedFunc(func(i int, main string, secondary string, shortcut rune) {
		name := registryOrder[i]
		pages.SwitchToPage(name)

		view := views[i]
		app.SetFocus(*view)
	})

	app.SetRoot(pages, true).SetFocus(pages)
```

I want to be able to define tview components in a way that I can define more "commands". I want the tview components to preform some logic, and allow other tviews to update if its something they care about.

e.g. I want the search command to fetch results, and update have the search results page update, when the command completes.

Current implementation assumes that the type of tview component will all behave the same, and have the same resutls component. e.g. the all `*tview.InputField` assume a callback with a single string input, and passes the whole app state. It also assumes that the updates are meant for a `tview.TextView` component. 

I want to define new commands, each that sets its own view (a tview component), a model that tracks the user input and results of the command

Attempting to start with the search command

What I'm finding is that tview components want to hold the state and view of the component. i.e.

```go
	results := tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		})
```

I want the command model to control that flow.

Trying to implement elm like architecture into this. I want to replace all the custom logic with just command.model.update() which I can use to call command.model.view() with newly painted tview.Primitives.
