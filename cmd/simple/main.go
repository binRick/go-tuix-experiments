package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/k0kubun/pp"
	"github.com/millerlogic/tuix"
	"github.com/rivo/tview"
	logrus "github.com/sirupsen/logrus"

	ac "github.com/binRick/abduco-dev/go/abducoctl"
)

var DEBUG_MODE = false
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
var (
	sessions                 = []ac.AbducoSession{}
	last_loaded_session_time = time.Now()
	sessions_mt              sync.Mutex
)

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
	if *debug_mode {
		fmt.Fprintf(os.Stderr, "SESSIONS DEBUG: %s\n", pp.Sprintf(`%s`, SESSIONS))
		os.Exit(1)
	}
	sessions = SESSIONS
	update_items(sessions)
	sessions_mt.Unlock()
	l.WithFields(logrus.Fields{
		"qty":   len(sessions),
		"dur":   time.Since(started),
		"since": time.Since(last_loaded_session_time),
	}).Debug(fmt.Sprintf("%d Sessions Loaded", len(sessions)))
	if DEBUG_MODE {
		pp.Fprintf(os.Stderr, "SESSIONS:         %s\n", sessions)
	}
	last_loaded_session_time = time.Now()
}

var monitor_once sync.Once

func get_fields() *logrus.Fields {
	return &logrus.Fields{
		"pid": os.Getpid(),
	}
}

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

}

type SCREEN struct {
	screen tcell.Screen
}

var Screen *SCREEN

func init() {
	Screen = &SCREEN{}
}
func run() error {
	_screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	Screen.screen = _screen
	err = Screen.screen.Init()
	if err != nil {
		panic(err)
	}
	Screen.screen.EnableMouse()
	app.SetScreen(Screen.screen)
	w, h := Screen.screen.Size()
	l.WithFields(logrus.Fields{
		"sessions": sessions,
		"width":    w, "height": h,
	}).Info("Terminal UI Started")
	monitor_once.Do(func() {
		go func() {
			for {
				monitor_sessions()
				time.Sleep(5000 * time.Millisecond)
			}
		}()
	})

	//w, h := _screen.Size()
	//pp.Println(w, h)
	//os.Exit(1)

	win2 := tuix.NewWindow().SetAutoPosition(false).SetResizable(true)
	go func() {
		on := 0
		for {
			s := string([]rune{clocks[on]})
			win2.SetTitle(fmt.Sprintf("Session Selection  %s|%d  ",
				s,
				on,
			))
			app.Draw()
			time.Sleep(1 * time.Second)
			on++
			if on >= len(clocks) {
				on = 0
			}
		}
	}()
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
	tv.SetWordWrap(true).SetDynamicColors(true).SetScrollable(true)
	tvr.SetTextColor(tcell.ColorRed)
	tvr.SetWordWrap(true).SetDynamicColors(true).SetScrollable(true)

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

	desktop.AddWindow(lw).AddWindow(win2l).AddWindow(win2r).AddWindow(win2).AddWindow(win3).AddWindow(menu_bar).SetBackgroundColor(tcell.ColorBlack).SetTitle("Abduco Sessions").SetBorder(true)
	app.SetRoot(desktop, true)
	return app.Run()
}

var cmds_ran_qty = 0
var cmd_runner_mutex sync.Mutex
var cmd_runner_once sync.Once

func cmd_runner() {
	if false {
		cmd_runner_once.Do(func() {
			for {
				cmd_runner_mutex.Lock()
				l.WithFields(logrus.Fields{
					"cmd": cur_cmd,
				}).Info(fmt.Sprintf("<%d> Running Command #%d", os.Getpid(), cmds_ran_qty))
				run_cmd_in_new_context(cur_cmd, tv, tvr)
				cmds_ran_qty++
				cmd_runner_mutex.Unlock()
				time.Sleep(1000 * time.Millisecond)
			}
		})
	}
}

func main() {
	flag.Parse()
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
