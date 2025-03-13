package call

import (
	db2 "github.com/rusinikita/acid/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCall_Exec(t *testing.T) {
	db := db2.Connect()

	_, err := db.Exec(`DROP TABLE IF EXISTS exec_test`)
	require.NoError(t, err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS exec_test (
    id SERIAL PRIMARY KEY,
    name TEXT
)`)
	require.NoError(t, err)

	pp := Params{
		"nameme": "biba",
	}

	insertCall := Step{
		Code: "insert into exec_test (name) values ('biba')",
	}

	result := insertCall.Exec(db, pp)
	assert.NoError(t, result.Error)
	assert.Nil(t, result.Rows)
	assert.Equal(t, int64(1), result.RowsAffected)

	selectCall := Step{
		Code: "select * from exec_test",
	}

	result = selectCall.Exec(db, pp)
	assert.NoError(t, result.Error)
	require.NotNil(t, result.Rows)
	assert.Len(t, result.Rows.Rows, 1)
	assert.Len(t, result.Rows.Columns, 2)
	assert.Equal(t, []string{"1", "biba"}, result.Rows.Rows[0])
}
