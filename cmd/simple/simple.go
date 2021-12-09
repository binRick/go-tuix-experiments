// This demo is public domain, or MPLv2 if you prefer.

package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/millerlogic/tuix"
	"github.com/rivo/tview"

	logrus "github.com/sirupsen/logrus"
)

var l = logrus.New()
var LOG_FILE = `/tmp/logrus.log`
var cmd_mutex sync.Mutex
var DEFAULT_CMD = fmt.Sprintf(`date >&1 ; date >&2`)
var cur_cmd = DEFAULT_CMD
var _list = tview.NewList().
	AddItem("List item 1", "Some explanatory text", 'a', nil).
	AddItem("List item 2", "Some explanatory text", 'b', nil).
	AddItem("List item 3", "Some explanatory text", 'c', nil).
	AddItem("List item 4", "Some explanatory text", 'd', nil).
	AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})

var (
	app      = tview.NewApplication()
	desktop  = tuix.NewDesktop()
	wform    = tview.NewForm()
	wform2   = tview.NewForm()
	wform2r  = tview.NewForm()
	tv       = tview.NewTextView()
	tvr      = tview.NewTextView()
	menu_bar = tuix.NewWindow().SetAutoPosition(false).SetTitle("Buttons").SetBorder(true)
)

func init() {
	l.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
		FullTimestamp: false,
	})
	l.SetReportCaller(true)

	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		l.Out = file
	} else {
		panic(err)
	}

	l.WithFields(logrus.Fields{
		"log": LOG_FILE,
	}).Info("Terminal UI Started")
}

func run() error {
	_screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	err = _screen.Init()
	if err != nil {
		panic(err)
	}
	_screen.EnableMouse()
	app.SetScreen(_screen)

	wform.AddButton("Restore", func() {
		menu_bar.SetState(tuix.Restored)
	})

	wform.AddButton("Find Command", func() {
		cmd := fmt.Sprintf(`find /etc|head -n 10; echo NO_ERR >&2`)
		cmd_mutex.Lock()
		defer cmd_mutex.Unlock()
		cur_cmd = cmd
		l.WithFields(logrus.Fields{
			"cmd": cmd,
		}).Info("Find Command")
	})

	wform.AddButton("New Terminal", func() {
		//		AddNewTerminal(
		menu_bar.SetState(tuix.Restored)
	})
	wform.AddButton("Maximize", func() {
		menu_bar.SetState(tuix.Maximized)
	})

	menu_bar.SetRect(1, 1, 65, 5)

	win2 := tuix.NewWindow().SetAutoPosition(false).SetResizable(true)
	win2.SetTitle("Mode Selection")
	win2.SetBorder(true).SetRect(68, 1, 40, 5)

	win3 := tuix.NewWindow().SetAutoPosition(false).SetResizable(true)
	win3.SetTitle("Mode Selection")
	win3.SetBorder(true).SetRect(115, 1, 40, 12)

	win2l := tuix.NewWindow().SetAutoPosition(false).SetResizable(true)
	win2l.SetTitle("Window 2 Left")
	win2l.SetBorder(true).SetRect(6, 30, 70, 20)

	win2r := tuix.NewWindow().SetAutoPosition(false).SetResizable(true)
	win2r.SetTitle("Window 2 Right")
	win2r.SetBorder(true).SetRect(80, 30, 70, 20)

	win2l.SetClient(wform2, true)
	win2r.SetClient(wform2r, true)

	wform2r.AddButton("Background", func() {
		win2r.SetState(tuix.Restored)
	})
	wform2.AddButton("Background", func() {
		win3.SetState(tuix.Restored)
	})

	tv.SetTextColor(tcell.ColorGreen)
	tvr.SetTextColor(tcell.ColorRed)

	tvrt := ``
	tvr.SetDoneFunc(func(key tcell.Key) {
		tvrt = fmt.Sprintf("cur %d\n%s\nkey:%s\n", len(tvrt), tvrt, key)
		tvr.Clear()
		fmt.Fprintf(tvr, "%s\n", tvrt)
		if key == 27 {

		}
	})

	tvr.SetChangedFunc(func() {
		app.Draw()
	})
	tv.SetChangedFunc(func() {
		app.Draw()
	})

	tv.SetWordWrap(true).SetDynamicColors(true).SetScrollable(false).SetBorderPadding(1, 1, 1, 1)
	tvr.SetWordWrap(true).SetDynamicColors(true).SetScrollable(false).SetBorderPadding(1, 1, 1, 1)

	menu_bar.SetClient(wform, true)
	win3.SetClient(_list, true)
	win2.SetClient(radioButtons, true)
	win2l.SetClient(tv, true)
	win2r.SetClient(tvr, true)

	desktop.AddWindow(menu_bar).AddWindow(win2l).AddWindow(win2r).AddWindow(win3).AddWindow(win2).SetBackgroundColor(tcell.ColorBlack).SetTitle("Desktop").SetBorder(true)
	app.SetRoot(desktop, true)
	qty := 0
	go func() {
		for {
			run_cmd(cur_cmd, tv, tvr)
			qty++
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	return app.Run()
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
