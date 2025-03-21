package runner

import (
	"database/sql"
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/ui"
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

func (r *Runner) Run(sequence call.Sequence) <-chan ui.Event {
	results := make(chan ui.Event)

	go func() {
		for _, step := range sequence.Calls {
			r.RunSQL(step, results)
		}

		close(results)
	}()

	return results
}

func (r *Runner) RunSQL(step call.Step, events chan ui.Event) {
	events <- ui.Call(step, r.store.Running())

	go func() {
		result := step.Exec(r.store, nil)

		events <- ui.Result(step, result, r.store.Running())
	}()

	time.Sleep(200 * time.Millisecond)
}
