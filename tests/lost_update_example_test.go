package tests

import (
	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/event"
	"github.com/rusinikita/acid/runner"
	"github.com/rusinikita/acid/sequence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *LostUpdateSuite) TestLostUpdateExample() {
	tx1 := call.TrxID("first")
	tx2 := call.TrxID("second")

	// Arrange: inline copy of the "Lost Update example" sequence
	seq := sequence.Sequence{
		Name: "Lost Update example",
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
	}

	// Act: run via StepIterator and collect all events
	r := runner.New(s.db)
	iter := runner.NewIterator(r)
	iter.Run(seq)

	var events []event.Event
	for {
		e, ok := iter.Next()
		if !ok {
			break
		}
		events = append(events, e)
	}

	// Assert: 13 steps × 2 events (Call + Result) each
	assert.Len(s.T(), events, 26)

	var updateResults []event.Event
	var finalSelectResult *call.ExecResult

	for _, e := range events {
		res := e.Result()
		if res == nil {
			continue
		}

		require.NoError(s.T(), res.Error, "unexpected error in sequence step")

		// UPDATE result: has a TrxID, no rows, and affected at least one row
		if e.Trx() != "" && res.Rows == nil && res.RowsAffected > 0 {
			updateResults = append(updateResults, e)
		}

		// Final SELECT result: no TrxID, not a setup step, has rows
		if !e.TestSetup() && e.Trx() == "" && res.Rows != nil {
			finalSelectResult = res
		}
	}

	// Both transactions updated the same row — this is the lost update anomaly:
	// neither transaction checked whether the slot was already taken.
	assert.Len(s.T(), updateResults, 2, "both transactions should have updated the same row")
	for _, e := range updateResults {
		assert.Equal(s.T(), int64(1), e.Result().RowsAffected)
	}

	// Final state: all rows are present, slot start_time=1 has speaker_id set by the last writer
	require.NotNil(s.T(), finalSelectResult)
	assert.Len(s.T(), finalSelectResult.Rows.Rows, 2)
}
