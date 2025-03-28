package run

import (
	"github.com/rusinikita/acid/event"
	"github.com/rusinikita/acid/sequence"
)

type runner interface {
	Run(s sequence.Sequence)
	Next() (event.Event, bool)
}
