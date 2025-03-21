package runner

import (
	"database/sql"
	"fmt"
	"github.com/rusinikita/acid/call"
)

type TrxStore struct {
	db     *sql.DB
	trxMap map[call.TrxID]*sql.Tx
}

func NewTrxStore(db *sql.DB) *TrxStore {
	return &TrxStore{
		db:     db,
		trxMap: make(map[call.TrxID]*sql.Tx),
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

func (t *TrxStore) GetWithLock(id call.TrxID) (call.DBExec, error) {
	if id == "" {
		return t.db, nil
	}

	tx := t.trxMap[id]
	if tx == nil {
		return nil, fmt.Errorf("trx %s is not started with Begin", id)
	}

	return tx, nil
}

func (t *TrxStore) Unlock(id call.TrxID) {
	//TODO implement me
	panic("implement me")
}

func (t *TrxStore) Locked() []call.TrxID {
	//TODO implement me
	panic("implement me")
}
