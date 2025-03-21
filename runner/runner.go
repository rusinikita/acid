package runner

import (
	"database/sql"
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/event"
	"time"
)

type Runner struct {
	store *TrxStore
}

func New(db *sql.DB) *Runner {
	return &Runner{
		store: NewTrxStore(db),
	}
}

type DBExec interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

func (r *Runner) Run(sequence call.Sequence) <-chan event.Event {
	results := make(chan event.Event)

	go func() {
		for i, step := range sequence.Calls {
			r.RunSQL(step, results, i == len(sequence.Calls)-1)
		}
	}()

	return results
}

func (r *Runner) RunSQL(step call.Step, events chan event.Event, lastStep bool) {
	events <- event.Call(step, r.store.Running())

	go func() {
		result := step.Exec(r.store, nil)

		events <- event.Result(step, result, r.store.Running())

		if lastStep {
			close(events)
		}
	}()

	time.Sleep(200 * time.Millisecond)
}
