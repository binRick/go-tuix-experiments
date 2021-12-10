package main

import (
	"fmt"

	"github.com/millerlogic/tuix"
	"github.com/rivo/tview"
)

func get_session_log_path(session_id string) string {
	return fmt.Sprintf(`/tmp/%s-passh-stdout.log`, session_id)
}

var (
	lw       = tuix.NewWindow().SetAutoPosition(false).SetResizable(true)
	log_view = tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetWordWrap(true)
)

func init() {
	lw.SetTitle("Log Window")
	lw.SetBorder(true).SetRect(
		1,
		60,
		155, 11,
	)
	lw.SetClient(tree, true)

	log_view.SetChangedFunc(func() {
		app.Draw()
	})
	log_view.SetWordWrap(true).SetDynamicColors(true).SetScrollable(true).SetBorderPadding(1, 1, 1, 1)
	lw.SetClient(log_view, true)
	/*
		log_view.SetDoneFunc(func(key tcell.Key) {
			currentSelection := log_view.GetHighlights()
			if key == tcell.KeyEnter {
				if len(currentSelection) > 0 {
					log_view.Highlight()
				} else {
					log_view.Highlight("0").ScrollToHighlight()
				}
			} else if len(currentSelection) > 0 {
				index, _ := strconv.Atoi(currentSelection[0])
				if key == tcell.KeyTab {
					index = (index + 1) % numSelections
				} else if key == tcell.KeyBacktab {
					index = (index - 1 + numSelections) % numSelections
				} else {
					return
				}
				log_view.Highlight(strconv.Itoa(index)).ScrollToHighlight()
			}
		})
		log_view.SetBorder(true)
	*/
}
