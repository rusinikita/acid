package sequence

import "github.com/rusinikita/acid/call"

var (
	tx1 = call.TrxID("first")
	tx2 = call.TrxID("second")
	tx3 = call.TrxID("third")
)

var Sequences = append(TrxSpeech, Common...)

var TrxSpeech = []Sequence{
	// UPDATE
	{
		Name:        "Lost Update example",
		Description: "There is a potential bug",
		Calls: []call.Step{
			call.Setup("drop table if exists speaker_slots"),
			call.Setup(`CREATE TABLE speaker_slots (
    meetup_id INTEGER NOT NULL, 
    speaker_id INTEGER DEFAULT NULL,
    start_time INTEGER NOT NULL DEFAULT 1, 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY(meetup_id, start_time)
)`),
			call.Setup("insert into speaker_slots (meetup_id, start_time) VALUES (1, 1)"),
			call.Setup("insert into speaker_slots (meetup_id, start_time) VALUES (1, 2)"),
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call("select meetup_id, start_time, speaker_id from speaker_slots where meetup_id = 1", tx1),
			call.Call("select meetup_id, start_time, speaker_id from speaker_slots where meetup_id = 1", tx2),
			call.Call("update speaker_slots set speaker_id = 1 where meetup_id = 1 and start_time = 1", tx1),
			call.Call("update speaker_slots set speaker_id = 1 where meetup_id = 1 and start_time = 1", tx2),
			call.Commit(tx1),
			call.Commit(tx2),
			call.Call("select meetup_id, start_time, speaker_id from speaker_slots where meetup_id = 1"),
		},
	},
	{
		Name:        "Lost Update fix",
		Description: "Fix lost update",
		Calls: []call.Step{
			call.Setup("drop table if exists speaker_slots"),
			call.Setup(`CREATE TABLE speaker_slots (
    meetup_id INTEGER NOT NULL, 
    speaker_id INTEGER DEFAULT NULL,
    start_time INTEGER NOT NULL DEFAULT 1, 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY(meetup_id, start_time)
)`),
			call.Setup("insert into speaker_slots (meetup_id, start_time) VALUES (1, 1)"),
			call.Setup("insert into speaker_slots (meetup_id, start_time) VALUES (1, 2)"),
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call("select meetup_id, start_time, speaker_id from speaker_slots where meetup_id = 1", tx1),
			call.Call("select meetup_id, start_time, speaker_id from speaker_slots where meetup_id = 1", tx2),
			call.Call("update speaker_slots set speaker_id = 1 where meetup_id = 1 and start_time = 1 and speaker_id is null", tx1),
			call.Call("update speaker_slots set speaker_id = 1 where meetup_id = 1 and start_time = 1 and speaker_id is null", tx2),
			call.Commit(tx1),
			call.Commit(tx2),
			call.Call("select meetup_id, start_time, speaker_id from speaker_slots where meetup_id = 1"),
		},
	},
	// insert + update
	{
		Name:        "INSERT constraint err example",
		Description: "Shows how error appear during transaction",
		Calls: []call.Step{
			call.Setup("drop table if exists meetups"),
			call.Setup("drop table if exists visitors"),
			call.Setup(`CREATE TABLE meetups (
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    creator_id INTEGER NOT NULL,
    max_seats INTEGER NOT NULL DEFAULT 1,
    booked_seats INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    CHECK(booked_seats <= max_seats)
)`),
			call.Setup(`CREATE TABLE visitors (
    meetup_id INTEGER NOT NULL,
    visitor_id INTEGER NOT NULL,
    seats INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (meetup_id, visitor_id)
)`),
			call.Setup("insert into meetups (creator_id, name, max_seats) VALUES (123, 'Meet with Biba&Boba', 10)"),
			call.Setup("insert into visitors (meetup_id, visitor_id, seats) VALUES (1, 1, 8)"),
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call("SELECT sum(seats) FROM visitors WHERE meetup_id = 1", tx1),
			call.Call("insert into visitors (meetup_id, visitor_id, seats) VALUES (1, 2, 2)", tx1),
			call.Call("SELECT sum(seats) FROM visitors WHERE meetup_id = 1", tx2),
			call.Call("insert into visitors (meetup_id, visitor_id, seats) VALUES (1, 2, 2)", tx2),
			call.Commit(tx1),
			call.Commit(tx2),
			call.Call("SELECT sum(seats) FROM visitors WHERE meetup_id = 1"),
		},
	},
	{
		Name:        "INSERT check constraint err example",
		Description: "Shows how error appear during transaction",
		Calls: []call.Step{
			call.Setup("drop table if exists meetups"),
			call.Setup("drop table if exists visitors"),
			call.Setup(`CREATE TABLE meetups (
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    creator_id INTEGER NOT NULL,
    max_seats INTEGER NOT NULL DEFAULT 1,
    booked_seats INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    CHECK ( booked_seats <= max_seats )
)`),
			call.Setup(`CREATE TABLE visitors (
    meetup_id INTEGER NOT NULL,
    visitor_id INTEGER NOT NULL,
    seats INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (meetup_id, visitor_id)
)`),
			call.Setup("insert into visitors (meetup_id, visitor_id, seats) VALUES (1, 1, 8)"),
			call.Setup("insert into meetups (creator_id, name, max_seats, booked_seats) VALUES (123, 'Meet with Biba&Boba', 10, 8)"),
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call("insert into visitors (meetup_id, visitor_id, seats) VALUES (1, 2, 2)", tx1),
			call.Call("update meetups set booked_seats = booked_seats + 2 where id = 1", tx1),
			call.Call("insert into visitors (meetup_id, visitor_id, seats) VALUES (1, 3, 2)", tx2),
			call.Call("update meetups set booked_seats = booked_seats + 2 where id = 1", tx2),
			call.Commit(tx1),
			call.Commit(tx2),
			call.Call("SELECT sum(seats) FROM visitors WHERE meetup_id = 1"),
		},
	},
}

var Common = []Sequence{
	// other
	{
		Name:        "PG Lock",
		Description: "Lock pg",
		Calls: []call.Step{
			call.Setup("drop table if exists for_test"),
			call.Setup(`CREATE TABLE for_test (
    id INTEGER NOT NULL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW()
)`),
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call(`SELECT * from for_test where id = 1 FOR UPDATE`, tx1),
			call.Call(`SELECT * from for_test where id = 1 FOR UPDATE`, tx2),
			call.Call(`INSERT INTO for_test (id) VALUES (1)`, tx2),
			call.Commit(tx1),
			call.Commit(tx2),
		},
	},
	{
		Name:        "Insert isolation level",
		Description: "Shows default isolation level and how it affects received data",
		Calls: []call.Step{
			call.Setup("drop table if exists exec_test"),
			call.Setup("CREATE TABLE exec_test (id SERIAL PRIMARY KEY, name TEXT)"),
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
		Name:        "Insert isolation level 2",
		Description: "Shows default isolation level and how it affects received data",
		Calls: []call.Step{
			call.Call("drop table if exists exec_test"),
			call.Call("CREATE TABLE exec_test (id INTEGER, name TEXT)"),
			call.Call("INSERT INTO exec_test (id, name) VALUES (1, 'biba')"),
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call("update exec_test set name = 'biba1' where id = 1", tx1),
			call.Call("update exec_test set name = 'biba2' where id = 1", tx2),
			call.Call("select * from exec_test where id = 1 for update"),
			call.Commit(tx1),
			call.Commit(tx2),
			call.Call("select * from exec_test where id = 1 for update"),
		},
	},
	{
		Name:        "Advisory lock MySQL",
		Description: "Shows default isolation level and how it affects received data",
		Calls: []call.Step{
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call("SELECT IS_FREE_LOCK('lock')", tx1),
			call.Call("SELECT GET_LOCK('lock',10)", tx1),
			call.Call("SELECT GET_LOCK('lock',10) skip locked", tx2),
			call.Call("SELECT IS_FREE_LOCK('lock')", tx2),
			call.Commit(tx1),
			call.Commit(tx2),
			call.Call("SELECT IS_FREE_LOCK('lock')"),
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
	{
		Name:        "Deadlock",
		Description: "Demonstrates deadlock scenario where two transactions wait for each other to release locks",
		Calls: []call.Step{
			call.Call("drop table if exists resources"),
			call.Call("CREATE TABLE resources (id INTEGER PRIMARY KEY, name TEXT, value INTEGER)"),
			call.Call("insert into resources (id, name, value) values (1, 'Resource1', 100), (2, 'Resource2', 200)"),
			call.Begin(tx1),
			call.Begin(tx2),
			call.Call("update resources set value = value - 10 where id = 1", tx1), // tx1 locks Resource1
			call.Call("update resources set value = value + 10 where id = 2", tx2), // tx2 locks Resource2
			call.Call("update resources set value = value - 10 where id = 2", tx1), // tx1 waits for tx2
			call.Call("update resources set value = value + 10 where id = 1", tx2), // tx2 waits for tx1
			call.Commit(tx1), // Deadlock resolution
			call.Commit(tx2),
			call.Call("select * from resources"),
		},
		LearningLinks: []string{
			"https://en.wikipedia.org/wiki/Deadlock",
		},
	},
	{
		Name:        "Dirty Read",
		Description: "Demonstrates dirty read scenario where one transaction reads uncommitted changes made by another transaction",
		Calls: []call.Step{
			call.Call("drop table if exists users"),
			call.Call("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER)"),
			call.Call("insert into users (id, name, age) values (1, 'Alice', 20)"),
			call.Begin(tx1),
			call.Call("select age from users where id = 1", tx1),
			call.Begin(tx2),
			call.Call("update users set age = 21 where id = 1", tx2),
			call.Call("select age from users where id = 1", tx1),
			call.Commit(tx1),
			call.Rollback(tx2),
			call.Call("select * from users"),
		},
		LearningLinks: []string{
			"https://en.wikipedia.org/wiki/Isolation_(database_systems)#Dirty_reads",
		},
	},
	{
		Name:        "Phantom Reads",
		Description: "Demonstrates phantom reads where transaction sees different results for the same query due to another transaction inserting new rows",
		Calls: []call.Step{
			call.Call("drop table if exists users"),
			call.Call("CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT, age INTEGER)"),
			call.Call("insert into users (name, age) values ('Alice', 18)"),
			call.Call("insert into users (name, age) values ('Bob', 20)"),
			call.Begin(tx1),
			call.Call("select name from users where age > 17", tx1),
			call.Begin(tx2),
			call.Call("insert into users (name, age) values ('Carol', 26)", tx2),
			call.Commit(tx2),
			call.Call("select name from users where age > 17", tx1),
			call.Commit(tx1),
		},
		LearningLinks: []string{
			"https://en.wikipedia.org/wiki/Isolation_(database_systems)#Phantom_reads",
		},
	},
}
