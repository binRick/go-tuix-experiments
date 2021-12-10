package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/k0kubun/pp"
	"github.com/millerlogic/tuix"
	"github.com/rivo/tview"
	logrus "github.com/sirupsen/logrus"

	//ac "github.com/binRick/abduco-dev/go/abduco"
	//ac "github.com/binRick/abduco-dev/go/abducoctl"
	ac "github.com/binRick/abduco-dev/go/abducoctl"
	//ac1 "github.com/binRick/abduco-dev/go/abducoctl"
	//ac1 "github.com/binRick/abduco-dev/go/abducoctl"
)

var l = logrus.New()
var LOG_FILE = `/tmp/logrus.log`
var cmd_mutex sync.Mutex
var radio_mutex sync.Mutex
var DEFAULT_CMD = fmt.Sprintf(`date >&1 ; date >&2`)
var cur_cmd = DEFAULT_CMD

var (
	app     = tview.NewApplication()
	desktop = tuix.NewDesktop()
	wform2  = tview.NewForm()
	wform2r = tview.NewForm()
	tv      = tview.NewTextView()
	tvr     = tview.NewTextView()
)
var sessions = []ac.AbducoSession{}
var sessions_mt sync.Mutex

func get_sessions() []ac.AbducoSession {
	sessions_mt.Lock()
	defer sessions_mt.Unlock()
	sess := sessions
	return sess
}
func monitor_sessions() {
	sessions_mt.Lock()
	started := time.Now()
	SESSIONS, err := ac.List()
	if err != nil {
		panic(err)
	}
	sessions = SESSIONS
	update_items(sessions)
	sessions_mt.Unlock()
	l.WithFields(logrus.Fields{
		"qty": len(sessions),
		"dur": time.Since(started),
	}).Info(fmt.Sprintf("%d Sessions Loaded", len(sessions)))
}

var monitor_once sync.Once

func init() {
	app.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		if false {
			l.WithFields(logrus.Fields{
				"event": pp.Sprintf(`%s`, event),
			}).Info("Mouse Capture")
		}
		return event, action
	})
	l.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
		FullTimestamp: false,
	})
	l.SetReportCaller(false)
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		l.Out = file
	} else {
		panic(err)
	}

	l.WithFields(logrus.Fields{
		"sessions": sessions,
	}).Info("Terminal UI Started")
	monitor_once.Do(func() {
		go func() {
			for {
				monitor_sessions()
				time.Sleep(5000 * time.Millisecond)
			}
		}()
	})
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

	win2 := tuix.NewWindow().SetAutoPosition(false).SetResizable(true)
	win2.SetTitle("Session Selection")
	win2.SetBorder(true).SetRect(68, 1, 40, 10)

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

	tv.SetWordWrap(true).SetDynamicColors(true).SetScrollable(true).SetBorderPadding(1, 1, 1, 1)
	tvr.SetWordWrap(true).SetDynamicColors(true).SetScrollable(true).SetBorderPadding(1, 1, 1, 1)

	win3.SetClient(_list, true)
	win2.SetClient(radioButtons, true)
	win2l.SetClient(tv, true)
	win2r.SetClient(tvr, true)

	desktop.AddWindow(menu_bar).AddWindow(lw).AddWindow(win2l).AddWindow(win2r).AddWindow(win2).AddWindow(win3).SetBackgroundColor(tcell.ColorBlack).SetTitle("Abduco Sessions").SetBorder(true)
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
