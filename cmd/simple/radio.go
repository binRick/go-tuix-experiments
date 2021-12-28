package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/millerlogic/tuix"
	"github.com/nxadm/tail"
	"github.com/rivo/tview"
	logrus "github.com/sirupsen/logrus"
)

var radioButtons = NewRadioButtons([]string{"📠  🔄   ✅", "🎰   🆘    🛠", "📛  🌀", "🔵   📀", "🔥   🌊", "💣   🌈"})

type RadioButtons struct {
	*tview.Box
	options       []string
	currentOption int
}

func NewRadioButtons(options []string) *RadioButtons {
	return &RadioButtons{
		Box:     tview.NewBox(),
		options: options,
	}
}

func (r *RadioButtons) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyUp:
			r.currentOption--
			if r.currentOption < 0 {
				r.currentOption = 0
			}
		case tcell.KeyDown:
			r.currentOption++
			if r.currentOption >= len(r.options) {
				r.currentOption = len(r.options) - 1
			}
		}
	})
}

var currently_selected_session_index = 0

func (r *RadioButtons) Draw(screen tcell.Screen) {
	radio_mutex.Lock()
	defer radio_mutex.Unlock()
	r.Box.DrawForSubclass(screen, r)
	x, y, width, height := r.GetInnerRect()
	cur_sessions := get_sexpect_sessions()
	for index, option := range r.options {
		if index >= height || len(cur_sessions)-1 < index {
			break
		}
		cs := cur_sessions[index]
		radioButton := " \u25ef" // Unchecked.
		if index == r.currentOption {
			radioButton = " \u25c9" // Checked.
		}
		cs_dur_s := (int64(time.Now().Unix()*1000) - int64(cs.CreateTime)) / 1000
		cs_dur := cs_dur_s
		//time.Duration(time.Second * cs_dur_s))
		line := fmt.Sprintf(`%s[white]  <%d> %dsec %v, [%v] %v> %v`,
			radioButton,
			cs.PID,
			cs_dur_s,
			cs_dur,
			cs.Username,
			option,
		)
		tview.Print(screen, line, x, y+index, width, tview.AlignLeft, tcell.ColorGreen)
	}
	if false {
		l.WithFields(logrus.Fields{
			"qty":     len(r.options),
			"cur":     r.currentOption,
			"options": fmt.Sprintf(`%s`, r.options),
		}).Info("RadioButtons Drawn")
	}
}

func (r *RadioButtons) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return r.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, y := event.Position()
		_, rectY, _, _ := r.GetInnerRect()
		c_mainc, c_combc, c_style, c_width := Screen.screen.GetContent(x, y)

		l.WithFields(logrus.Fields{
			"x": x, "y": y,
			"width": c_width,
			"mainc": c_mainc,
			"combc": c_combc,
			"style": c_style,
		}).Info("Mouse Event")
		if !r.InRect(x, y) {
			return false, nil
		}

		if action == tview.MouseLeftClick {
			setFocus(r)
			index := y - rectY
			if index >= 0 && index < len(r.options) && r.currentOption != index && ((index + 1) <= len(get_sexpect_sessions())) {
				r.currentOption = index
				consumed = true
			}
			cur_sessions := get_sexpect_sessions()
			cs := cur_sessions[r.currentOption]
			if consumed {
				var do_tail_file = func(log_file_path string, log_window *tuix.Window) {
					text_view := log_window.GetClient().(*tview.TextView)
					l.WithFields(logrus.Fields{
						"pid":   os.Getpid(),
						"title": log_window.GetTitle(),
						"log":   log_file_path,
					}).Info("Tailing Log Files")
					t, err := tail.TailFile(log_file_path, tail.Config{
						Follow: false,
					})
					if err != nil {
						panic(err)
					}
					for line := range t.Lines {
						fmt.Fprintf(text_view, "<%d> @%s #%d/%d Offset|%d Whence> %s\n", cs.PID, line.Time, line.Num, line.SeekInfo.Offset, line.SeekInfo.Whence, line.Text)
					}
				}
				if false {
					go do_tail_file(cs.StdoutLog, stdout_window)
					go do_tail_file(cs.StderrLog, stderr_window)
				}
			}
			l.WithFields(logrus.Fields{
				"qty":         len(r.options),
				"cur":         r.currentOption,
				"consumed":    consumed,
				"session_pid": cs.PID,
				"x":           x, "y": y,
				"pid": os.Getpid(),
			}).Info("Selected Session")
		}
		return
	})
}
