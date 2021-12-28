package main

import (
	"fmt"

	"github.com/millerlogic/tuix"
	"github.com/rivo/tview"
	logrus "github.com/sirupsen/logrus"
)

func get_session_log_path(session_id string) string {
	return fmt.Sprintf(`/tmp/%s-passh-stdout.log`, session_id)
}

func new_log_window(title string, width, height int, x, y int) *tuix.Window {
	lw := tuix.NewWindow().SetAutoPosition(false).SetResizable(true)
	lw.SetBorder(true).SetRect(
		width, height,
		x, y,
	)
	lw.SetTitle(title)
	lw.SetClient(new_log_view(title), true)
	lw.Focus(func(p tview.Primitive) {
		l.WithFields(logrus.Fields{
			"title": title,
		}).Info("Log Window Focused")
	})
	desktop.AddWindow(lw)
	return lw
}

func new_log_view(title string) *tview.TextView {
	lv := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	lv.SetScrollable(true).SetBorderPadding(1, 1, 1, 1)

	lv.Focus(func(p tview.Primitive) {
		r, c := lv.GetScrollOffset()
		f := lv.HasFocus()
		l.WithFields(logrus.Fields{
			"title": title,
			"row":   r, "col": c,
			"focused": f,
		}).Info("Log View Focused")
	})
	lv.SetChangedFunc(func() {
		go func() {
			r, c := lv.GetScrollOffset()
			f := lv.HasFocus()
			l.WithFields(logrus.Fields{
				"title": title,
				"row":   r, "col": c,
				"focused": f,
			}).Info("Log View Updated")
			if r == 0 {
				lv.ScrollToEnd()
			}
			app.Draw()
		}()
	})
	return lv
}

var (
	stdin_window = new_log_window("stdin", 1, 25, 155, 20)
	//stdin_window.Focus(func(p tview.Primitive){
	//fmt.Fprintf(os.Stderr, "stdin focus\n", )
	//})
	stdout_window = new_log_window("stdout", 1, 45, 155, 11)
	stderr_window = new_log_window("stderr", 1, 58, 155, 11)
)

func init() {
}
