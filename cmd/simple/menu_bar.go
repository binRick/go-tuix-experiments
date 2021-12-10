package main

import (
	"fmt"

	"github.com/millerlogic/tuix"
	"github.com/rivo/tview"
	logrus "github.com/sirupsen/logrus"
)

var menu_bar = tuix.NewWindow().SetAutoPosition(false).SetTitle("Main Controls").SetBorder(true)
var menu_bar_form = tview.NewForm().SetHorizontal(true)

func init() {
	menu_bar.SetClient(menu_bar_form, true)
	menu_bar.SetRect(
		1,
		1,
		60, 8,
	)

	menu_bar_form.AddButton("Restore", func() {
		menu_bar.SetState(tuix.Restored)
	})

	menu_bar_form.AddButton("Find Command", func() {
		cmd := fmt.Sprintf(`find /etc|head -n 10; echo NO_ERR >&2`)
		cmd = fmt.Sprintf(`journalctl -f`)
		cmd_mutex.Lock()
		defer cmd_mutex.Unlock()
		cur_cmd = cmd
		l.WithFields(logrus.Fields{
			"cmd": cmd,
		}).Info("Find Command")
	})
	menu_bar_form.AddButton("Refresh", func() {
		go monitor_sessions()
	})
	menu_bar_form.AddButton("Maximize", func() {
		menu_bar.SetState(tuix.Maximized)
	})
	menu_bar_form.AddButton("Restart Session", func() {
		menu_bar.SetState(tuix.Maximized)
	})
	menu_bar_form.AddButton("Kill Session", func() {
		menu_bar.SetState(tuix.Maximized)
	})
}
