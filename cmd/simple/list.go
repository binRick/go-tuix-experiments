package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	ac "github.com/binRick/abduco-dev/go/abducoctl"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
	logrus "github.com/sirupsen/logrus"
)

var _list_mt sync.Mutex
var _list = tview.NewList().AddItem("Quit", "Press to exit", 'q', func() {
	app.Stop()
})

func init() {
	_list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		l.Info(fmt.Sprintf(`SELECTED ITEM %s`, mainText))
		l.WithFields(logrus.Fields{
			"session": mainText,
		}).Info(fmt.Sprintf(`Tail Log File`))
		go tail_log(mainText)
		if false {
			fmt.Fprintf(os.Stderr, "%s\n", pp.Sprintf(`%s`, app))
		}
	})
	_list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		//	pp.Println(`CHANGED BUT NOT SELECTED ITEM`, index, mainText, secondaryText, shortcut)
	})
}

func update_items(sessions []ac.AbducoSession) {
	_list_mt.Lock()
	defer _list_mt.Unlock()
	qty := _list.GetItemCount()
	session_names := []string{}
	rogue_indexes := []int{}
	for _, s := range sessions {
		session_names = append(session_names, s.Session)
	}
	for sindex, s := range sessions {
		index := -1
		if qty > 0 {
			on := 0
			for {
				_qty := _list.GetItemCount()
				if on >= (_qty) {
					break
				}
				mainText, secondaryText := _list.GetItemText(on)
				if mainText == `Quit` {
					on++
					continue
				}
				mm := pp.Sprintf("%s|%s\n%s\n", mainText, secondaryText, s)
				if false {
					fmt.Fprintf(os.Stderr, "%s\n", mm)
				}
				if strings.ToLower(mainText) == strings.ToLower(s.Session) {
					index = on
				} else {
					is_rogue := true
					for _, n := range session_names {
						if strings.ToLower(n) == strings.ToLower(mainText) {
							is_rogue = false
						}
					}
					if is_rogue {
						rogue_indexes = append(rogue_indexes, on)
						l.WithFields(logrus.Fields{
							"sessions_qty": len(sessions),
							"items_qty":    qty,
							"index":        on,
						}).Info(fmt.Sprintf("Removing Rogue Session %s", mainText))
					}
				}
				if on >= (qty - 1) {
					break
				}
				on++
			}
		}
		for _, index := range rogue_indexes {
			if index < 0 {
				continue
			}
			_qty := _list.GetItemCount()
			if index >= (_qty) {
				break
			}
			if index > 0 {
				_list.RemoveItem(index)
			}
		}
		if index == -1 {
			l.WithFields(logrus.Fields{
				"qty": len(sessions),
			}).Info(fmt.Sprintf("Need to add session %s", s.Session))
			ru := '1'
			if sindex == 1 {
				ru = '2'
			} else if sindex == 2 {
				ru = '3'
			} else if sindex == 3 {
				ru = '4'
			}
			_list.AddItem(fmt.Sprintf(`%s`, s.Session), fmt.Sprintf(`⛔ %d- %d Processes, %d Threads`, s.PID, len(s.PIDs), s.Threads), ru, nil)
			/*func() {
				lp := get_session_log_path(s.Session)
				l.WithFields(logrus.Fields{
					"sindex": sindex,
					"index":  index,
				}).Info(fmt.Sprintf(`ITEM...................%s`, lp))
				go tail_log(lp)
			})
			*/
		}
	}
}
