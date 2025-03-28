package sequence

import "github.com/rusinikita/acid/call"

var (
	tx1 = call.TrxID("first")
	tx2 = call.TrxID("second")
)

var Sequences = []Sequence{
	{
		Name:        "Insert isolation level",
		Description: "Shows default isolation level and how it affects received data",
		Calls: []call.Step{
			call.Call("drop table exec_test"),
			call.Call("CREATE TABLE exec_test (id SERIAL PRIMARY KEY, name TEXT)"),
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call("insert into exec_test (name) values ('biba')", tx1),
			call.Call("select * from exec_test", tx2),
			call.Commit(tx1),
			call.Call("select * from exec_test", tx2),
			call.Call("select * from exec_test"),
		},
	},
	{
		Name:        "Update + condition",
		Description: "Shows that update locks row after finding, that can cause bugs",
		Calls: []call.Step{
			call.Call("drop table exec_test"),
			call.Call("CREATE TABLE exec_test (id INTEGER PRIMARY KEY, name TEXT, counter INTEGER)"),
			call.Call("insert into exec_test (id, name, counter) values (1, 'biba', 0)"),
			call.Call("update exec_test set counter = counter + 1"),
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call("update exec_test set counter = counter + 1 where id = 1 and counter < 2", tx1),
			call.Call("update exec_test set counter = counter + 1 where id = 1 and counter < 2", tx2),
			call.Call("update exec_test set counter = counter + 1 where id = 1 and counter < 2", tx1),
			call.Commit(tx1),
			call.Commit(tx2),
			call.Call("select * from exec_test"),
		},
		LearningLinks: nil,
	},
}
