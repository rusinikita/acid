package runner

import (
	"database/sql"
	"fmt"
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

func (r *Runner) Run(sequence call.Sequence, next <-chan struct{}) <-chan ui.Event {
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
	fmt.Println("step", step.Trx, step.Code)

	events <- ui.Call(step)

	go func() {
		result := step.Exec(r.store, nil)

		events <- ui.Result(step, result)
		//if result.Error != nil {
		//	fmt.Println("error:", result.Error)
		//	return
		//}
		//
		//if result.Rows == nil {
		//	fmt.Println("rows affected:", result.RowsAffected)
		//	return
		//}
		//
		//fmt.Println("columns:", result.Rows.Columns)
		//for i, row := range result.Rows.Rows {
		//	fmt.Println("row:", i, row)
		//}
	}()

	time.Sleep(200 * time.Millisecond)
}
