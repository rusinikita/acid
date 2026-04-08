package tests

import (
	"os"
	"testing"

	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const lostUpdateTOML = `
name        = "Lost Update example"
description = "There is a potential bug"

[[steps]]
sql    = "drop table if exists speaker_slots"
setup  = true

[[steps]]
sql   = """CREATE TABLE speaker_slots (
    meetup_id INTEGER NOT NULL,
    speaker_id INTEGER DEFAULT NULL,
    start_time INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY(meetup_id, start_time)
)"""
setup = true

[[steps]]
sql   = "insert into speaker_slots (meetup_id, start_time) VALUES (1, 1)"
setup = true

[[steps]]
sql   = "insert into speaker_slots (meetup_id, start_time) VALUES (1, 2)"
setup = true

[[steps]]
cmd = "begin"
trx = "first"

[[steps]]
cmd = "begin"
trx = "second"

[[steps]]
sql = "select meetup_id, start_time, speaker_id from speaker_slots where meetup_id = 1"
trx = "first"

[[steps]]
sql = "select meetup_id, start_time, speaker_id from speaker_slots where meetup_id = 1"
trx = "second"

[[steps]]
sql = "update speaker_slots set speaker_id = 1 where meetup_id = 1 and start_time = 1"
trx = "first"

[[steps]]
sql = "update speaker_slots set speaker_id = 1 where meetup_id = 1 and start_time = 1"
trx = "second"

[[steps]]
cmd = "commit"
trx = "first"

[[steps]]
cmd = "commit"
trx = "second"

[[steps]]
sql = "select meetup_id, start_time, speaker_id from speaker_slots where meetup_id = 1"
`

func TestLoadLostUpdateFromTOML(t *testing.T) {
	f, err := os.CreateTemp("", "lost_update_*.toml")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	_, err = f.WriteString(lostUpdateTOML)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	seq, err := config.Load(f.Name())
	require.NoError(t, err)

	assert.Equal(t, "Lost Update example", seq.Name)
	assert.Equal(t, "There is a potential bug", seq.Description)
	assert.Len(t, seq.Calls, 13)

	// First four steps are setup
	for i := 0; i < 4; i++ {
		assert.True(t, seq.Calls[i].TestSetup, "step %d should be setup", i)
	}

	// Steps 4 and 5: BEGIN for tx1 and tx2
	assert.Equal(t, call.TrxBegin, seq.Calls[4].TrxCommand)
	assert.Equal(t, call.TrxID("first"), seq.Calls[4].Trx)
	assert.Equal(t, call.TrxBegin, seq.Calls[5].TrxCommand)
	assert.Equal(t, call.TrxID("second"), seq.Calls[5].Trx)

	// Steps 8 and 9: UPDATE in tx1 and tx2
	assert.Equal(t, call.TrxID("first"), seq.Calls[8].Trx)
	assert.Equal(t, call.TrxID("second"), seq.Calls[9].Trx)

	// Steps 10 and 11: COMMIT
	assert.Equal(t, call.TrxCommit, seq.Calls[10].TrxCommand)
	assert.Equal(t, call.TrxCommit, seq.Calls[11].TrxCommand)

	// Last step: plain SELECT with no trx
	assert.Equal(t, call.TrxID(""), seq.Calls[12].Trx)
	assert.False(t, seq.Calls[12].TestSetup)
}
