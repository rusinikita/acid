package sequence

import "github.com/rusinikita/acid/call"

var (
	tx1 = call.TrxID("first")
	tx2 = call.TrxID("second")
)

var Sequences = []Sequence{{
	Calls: []call.Step{
		call.Call("CREATE TABLE IF NOT EXISTS exec_test (id SERIAL PRIMARY KEY, name TEXT)"),
		call.Call("TRUNCATE TABLE exec_test"),
		call.Begin(tx1),
		call.Begin(tx2),
		call.Call("insert into exec_test (name) values ('biba')", tx1),
		call.Call("select * from exec_test", tx2),
		call.Commit(tx1),
		call.Call("select * from exec_test", tx2),
		call.Call("select * from exec_test"),
	},
}}
