package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
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

func (r *RadioButtons) Draw(screen tcell.Screen) {
	radio_mutex.Lock()
	defer radio_mutex.Unlock()
	r.Box.DrawForSubclass(screen, r)
	x, y, width, height := r.GetInnerRect()
	for index, option := range r.options {
		if index >= height {
			break
		}
		radioButton := " \u25ef" // Unchecked.
		if index == r.currentOption {
			radioButton = " \u25c9" // Checked.
		}
		line := fmt.Sprintf(`%s[white]  %s`, radioButton, option)
		tview.Print(screen, line, x, y+index, width, tview.AlignLeft, tcell.ColorYellow)
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
		if !r.InRect(x, y) {
			return false, nil
		}

		if action == tview.MouseLeftClick {
			setFocus(r)
			index := y - rectY
			if index >= 0 && index < len(r.options) {
				r.currentOption = index
				consumed = true
			}
			l.WithFields(logrus.Fields{
				"qty": len(r.options),
				"cur": r.currentOption,
				"x":   x, "y": y,
			}).Info("RadioButtons LeftClick")
		}
		return
	})
}
