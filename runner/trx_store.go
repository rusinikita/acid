package runner

import (
	"database/sql"
	"fmt"
	"github.com/rusinikita/acid/call"
	"maps"
	"slices"
	"sync"
)

type TrxStore struct {
	db      *sql.DB
	trxMap  map[call.TrxID]*sql.Tx
	running map[call.TrxID]struct{}
	mutex   sync.Mutex
}

func NewTrxStore(db *sql.DB) *TrxStore {
	return &TrxStore{
		db:      db,
		trxMap:  make(map[call.TrxID]*sql.Tx),
		running: make(map[call.TrxID]struct{}),
	}
}

func (t *TrxStore) Do(id call.TrxID, command call.TrxCommandType) error {
	tx := t.trxMap[id]

	switch command {
	case call.TrxBegin:
		trx, err := t.db.Begin()
		if err != nil {
			return fmt.Errorf("begin trx %s: %v", id, err)
		}

		t.trxMap[id] = trx
	case call.TrxCommit:
		if tx == nil {
			return fmt.Errorf("trx %s not found to commit", id)
		}
		err := tx.Commit()
		if err != nil {
			return fmt.Errorf("commit trx %s: %v", id, err)
		}
	case call.TrxRollback:
		if tx == nil {
			return fmt.Errorf("trx %s not found to rollback", id)
		}
		err := tx.Rollback()
		if err != nil {
			return fmt.Errorf("rollback trx %s: %v", id, err)
		}
	default:
		return nil
	}

	return nil
}

func (t *TrxStore) Get(id call.TrxID) (call.DBExec, error) {
	var exec DBExec = t.db
	if id != "" {
		tx := t.trxMap[id]
		if tx == nil {
			return nil, fmt.Errorf("trx %s is not started with Begin", id)
		}

		exec = tx
	}

	return &execWrapper{
		id:    id,
		exec:  exec,
		store: t,
	}, nil
}

func (t *TrxStore) start(id call.TrxID) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.running[id] = struct{}{}
}

func (t *TrxStore) finish(id call.TrxID) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	delete(t.running, id)
}

func (t *TrxStore) Running() []call.TrxID {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return slices.Collect(maps.Keys(t.running))
}

type execWrapper struct {
	id    call.TrxID
	exec  DBExec
	store *TrxStore
}

func (e *execWrapper) Query(query string, args ...any) (*sql.Rows, error) {
	e.store.start(e.id)
	defer e.store.finish(e.id)

	return e.exec.Query(query, args...)
}

func (e *execWrapper) Exec(query string, args ...any) (sql.Result, error) {
	e.store.start(e.id)
	defer e.store.finish(e.id)

	return e.exec.Exec(query, args...)
}
