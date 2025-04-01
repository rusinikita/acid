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
			call.Call("drop table if exists exec_test"),
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
			call.Call("drop table if exists exec_test"),
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
	{
		Name:        "Lost update",
		Description: "Demonstrates a lost update scenario where one transaction overwrites changes of another transaction across two related tables",
		Calls: []call.Step{
			call.Call("drop table if exists cars"),
			call.Call("drop table if exists invoices"),
			call.Call("CREATE TABLE cars (id INTEGER PRIMARY KEY, model TEXT, price INTEGER, buyer TEXT)"),
			call.Call("CREATE TABLE invoices (id INTEGER PRIMARY KEY, car_id INTEGER, amount INTEGER, buyer TEXT)"),
			call.Call("insert into cars (id, model, price, buyer) values (1, 'Tesla Model S', 80000, 'A')"),
			call.Call("insert into invoices (id, car_id, amount, buyer) values (1, 1, 80000, 'A')"),
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call("update cars set buyer = 'X' where id = 1", tx1),
			call.Call("update cars set buyer = 'Y' where id = 1", tx2),
			call.Call("update invoices set buyer = 'Y' where car_id = 1", tx2),
			call.Call("update invoices set buyer = 'X' where car_id = 1", tx1),
			call.Commit(tx1),
			call.Commit(tx2),
			call.Call("select * from cars"),
			call.Call("select * from invoices"),
		},
	},
	{
		Name:        "Count locking",
		Description: "Demonstrates a try to lock and check count before insert",
		Calls: []call.Step{
			call.Call("drop table if exists exec_test"),
			call.Call("CREATE TABLE exec_test (id SERIAL PRIMARY KEY, name TEXT)"),
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call("select count(id) from exec_test where name = 'biba' for update", tx1),
			call.Call("select count(id) from exec_test where name = 'biba' for update", tx2),
			call.Call("insert into exec_test (name) values ('biba')", tx1),
			call.Call("insert into exec_test (name) values ('biba')", tx2),
			call.Commit(tx1),
			call.Commit(tx2),
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call("select count(id) from exec_test where name = 'biba' for update", tx1),
			call.Call("select count(id) from exec_test where name = 'biba' for update", tx2),
			call.Call("insert into exec_test (name) values ('biba')", tx1),
			call.Call("insert into exec_test (name) values ('biba')", tx2),
			call.Commit(tx1),
			call.Commit(tx2),
		},
		LearningLinks: nil,
	},
}
