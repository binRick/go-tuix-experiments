package main

import (
	"github.com/millerlogic/tuix"
	"github.com/rivo/tview"
)

var (
	menu_bar_form = tview.NewForm().SetHorizontal(true)
	menu_bar      = tuix.NewWindow().SetAutoPosition(false).SetTitle("Main Controls").SetBorder(true)
)

func init() {
	menu_bar.SetClient(menu_bar_form, true)
	menu_bar.SetRect(
		1,
		1,
		60, 8,
	)

	menu_bar_form.AddCheckbox("Mode Enabled", true, func(en bool) {
		l.Info("Mode Set!: ", en)
	})

	//	menu_bar_form.AddButton("Restore", func() {
	//		menu_bar.SetState(tuix.Restored)
	//	})

	AddCommandButton(menu_bar_form, `Tail Journal`, `command journalctl -f`, nil)
	AddCommandButton(menu_bar_form, `List etc`, `command ls /etc`, nil)
	AddCommandButton(menu_bar_form, `  w  `, `w`, nil)
	AddCommandButton(menu_bar_form, `  sleep 5  `, `sleep 5`, nil)
	AddCommandButton(menu_bar_form, `  stderr test  `, `>&2 date`, nil)

	menu_bar_form.AddButton("Refresh", func() {
		go monitor_sessions()
	})
	menu_bar_form.AddButton("Maximize", func() {
		menu_bar.SetState(tuix.Maximized)
	})
	//	menu_bar_form.AddButton("Restart Session", func() {
	//		menu_bar.SetState(tuix.Maximized)
	//	})
	//	menu_bar_form.AddButton("Kill Session", func() {
	//		menu_bar.SetState(tuix.Maximized)
	//	})
}
