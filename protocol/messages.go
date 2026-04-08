package protocol

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"

	"github.com/rusinikita/acid/call"
	"github.com/rusinikita/acid/event"
)

// EventMessage is the JSON wire format for a single event.
type EventMessage struct {
	Kind      string         `json:"kind"` // "call" | "result"
	Trx       string         `json:"trx"`
	TestSetup bool           `json:"test_setup"`
	Waiting   []string       `json:"waiting,omitempty"`
	Step      *StepMessage   `json:"step,omitempty"`
	Result    *ResultMessage `json:"result,omitempty"`
}

type StepMessage struct {
	Code       string `json:"code"`
	TrxCommand int    `json:"trx_command"`
}

type ResultMessage struct {
	RowsAffected int64          `json:"rows_affected"`
	Error        string         `json:"error,omitempty"`
	Rows         *SelectMessage `json:"rows,omitempty"`
}

type SelectMessage struct {
	Columns []string   `json:"columns"`
	Rows    [][]string `json:"rows"`
}

func Marshal(e event.Event) EventMessage {
	msg := EventMessage{
		Trx:       string(e.Trx()),
		TestSetup: e.TestSetup(),
	}
	for _, w := range e.Waiting() {
		msg.Waiting = append(msg.Waiting, string(w))
	}

	if step := e.Step(); step != nil {
		msg.Kind = "call"
		msg.Step = &StepMessage{
			Code:       step.Code,
			TrxCommand: int(step.TrxCommand),
		}
	} else {
		msg.Kind = "result"
		r := e.Result()
		rm := &ResultMessage{RowsAffected: r.RowsAffected}
		if r.Error != nil {
			rm.Error = r.Error.Error()
		}
		if r.Rows != nil {
			rm.Rows = &SelectMessage{
				Columns: r.Rows.Columns,
				Rows:    r.Rows.Rows,
			}
		}
		msg.Result = rm
	}
	return msg
}

func Unmarshal(msg EventMessage) event.Event {
	waiting := make([]call.TrxID, len(msg.Waiting))
	for i, w := range msg.Waiting {
		waiting[i] = call.TrxID(w)
	}

	trx := call.TrxID(msg.Trx)

	if msg.Kind == "call" && msg.Step != nil {
		step := call.Step{
			Code:       msg.Step.Code,
			Trx:        trx,
			TrxCommand: call.TrxCommandType(msg.Step.TrxCommand),
			TestSetup:  msg.TestSetup,
		}
		return event.Call(step, waiting)
	}

	step := call.Step{Trx: trx, TestSetup: msg.TestSetup}
	result := call.ExecResult{}
	if msg.Result != nil {
		result.RowsAffected = msg.Result.RowsAffected
		if msg.Result.Error != "" {
			result.Error = errors.New(msg.Result.Error)
		}
		if msg.Result.Rows != nil {
			result.Rows = &call.SelectResult{
				Columns: msg.Result.Rows.Columns,
				Rows:    msg.Result.Rows.Rows,
			}
		}
	}
	return event.Result(step, result, waiting)
}

// WriteEvent serializes e as a JSON line and writes it to w.
func WriteEvent(w io.Writer, e event.Event) error {
	data, err := json.Marshal(Marshal(e))
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = w.Write(data)
	return err
}

// ReadEvent reads one JSON line from r and deserializes it into an event.Event.
func ReadEvent(r *bufio.Reader) (event.Event, error) {
	line, err := r.ReadBytes('\n')
	if err != nil {
		return event.Event{}, err
	}
	var msg EventMessage
	if err := json.Unmarshal(line, &msg); err != nil {
		return event.Event{}, err
	}
	return Unmarshal(msg), nil
}
