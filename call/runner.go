package call

import (
	"database/sql"
)

type TrxStore interface {
	Do(id TrxID, command TrxCommandType) error
	GetWithLock(id TrxID) (DBExec, error)
	Unlock(id TrxID)
	Locked() []TrxID
}

type DBExec interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}
