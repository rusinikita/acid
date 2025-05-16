package call

import (
	"fmt"
	"strings"
)

type Params map[string]any

func (p Params) CallParams(call Step) (sql string, args []any, err error) {
	if len(call.ParamNames) == 0 {
		return call.Code, nil, nil
	}

	panic("params not supported")

	//pp := make([]any, len(call.ParamNames))
	//var notFound []string
	//
	//for i, n := range call.ParamNames {
	//	value, ok := p[n]
	//	if !ok {
	//		notFound = append(notFound, n)
	//		continue
	//	}
	//
	//	pp[i] = sql.Named(n, value)
	//}
	//
	//if len(notFound) > 0 {
	//	return "", nil, errors.New("args not found: " + strings.Join(notFound, ", "))
	//}
	//
	//return "", pp, nil
}

type Step struct {
	Code       string
	ParamNames []string
	Trx        TrxID
	TrxCommand TrxCommandType
	TestSetup  bool
}

type TrxID string

type TrxCommandType int

const (
	TrxNone TrxCommandType = iota
	TrxBegin
	TrxCommit
	TrxRollback
)

func (t TrxCommandType) String() string {
	switch t {
	case TrxBegin:
		return "BEGIN"
	case TrxCommit:
		return "COMMIT"
	case TrxRollback:
		return "ROLLBACK"
	default:
		return "no transaction command"
	}
}

type ExecResult struct {
	Rows         *SelectResult
	RowsAffected int64
	Error        error
}

type SelectResult struct {
	Columns []string
	Rows    [][]string
}

func Select(exec DBExec, sql string, args []any) (*SelectResult, error) {
	rows, err := exec.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	resultRows := make([][]string, 0)
	columnCount := len(columns)

	// Create slices once outside the loop for efficiency
	values := make([]*string, columnCount)
	scanArgs := make([]any, columnCount)
	for i := range scanArgs {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		// Create a new row for each iteration
		row := make([]string, columnCount)
		for i := range row {
			if values[i] == nil {
				row[i] = "NULL"
				continue
			}

			row[i] = *values[i]
		}

		resultRows = append(resultRows, row)
	}

	// Check for errors after iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &SelectResult{
		Columns: columns,
		Rows:    resultRows,
	}, nil
}

func (c Step) Exec(store TrxStore, pp Params) (result ExecResult) {
	sql, args, err := pp.CallParams(c)
	if err != nil {
		result.Error = err

		return result
	}

	err = store.Do(c.Trx, c.TrxCommand)
	if err != nil {
		result.Error = err

		return result
	}

	if c.TrxCommand != TrxNone {
		return result
	}

	if c.Code == "" {
		result.Error = fmt.Errorf("no code specified")

		return result
	}

	exec, err := store.Get(c.Trx)
	if err != nil {
		result.Error = err

		return result
	}

	if strings.HasPrefix(strings.ToLower(c.Code), "select") || strings.HasPrefix(strings.ToLower(c.Code), "show") {
		rows, err := Select(exec, sql, args)

		result.Rows = rows
		result.Error = err

		return result
	}

	r, err := exec.Exec(sql, args...)
	if err != nil {
		result.Error = err

		return result
	}

	result.RowsAffected, err = r.RowsAffected()
	result.Error = err

	return result
}
