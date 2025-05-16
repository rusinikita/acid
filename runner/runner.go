package runner

import (
	"database/sql"
	"errors"
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/event"
	"github.com/rusinikita/acid/sequence"
	"sync"
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

func (r *Runner) Run(sequence sequence.Sequence) <-chan event.Event {
	results := make(chan event.Event)

	go func() {
		r.store.ResetAll()

		wg := sync.WaitGroup{}

		for i, step := range sequence.Calls {
			if i > 0 && step.TestSetup && !sequence.Calls[i-1].TestSetup {
				results <- event.Call(step, nil)
				results <- event.Result(step, call.ExecResult{
					Error: errors.New("setup steps are allowed only before call steps"),
				}, nil)

				break
			}

			if step.Trx == "" {
				r.RunSQL(step, results)

				continue
			}

			wg.Add(1)
			go func(s call.Step) {
				r.RunSQL(step, results)

				wg.Done()
			}(step)

			time.Sleep(200 * time.Millisecond)
		}

		wg.Wait()
		close(results)
	}()

	return results
}

func (r *Runner) RunSQL(step call.Step, events chan event.Event) {
	events <- event.Call(step, r.store.Running())

	result := step.Exec(r.store, nil)

	events <- event.Result(step, result, r.store.Running())
}
