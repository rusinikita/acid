package run

import "github.com/rusinikita/acid/event"

type runner interface {
	Run()
	Next() (event.Event, bool)
}
