package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rusinikita/acid/protocol"
	"github.com/rusinikita/acid/runner"
	"github.com/rusinikita/acid/sequence"
)

// Run connects to the acid server at serverAddr, executes seq against db,
// and POSTs each event to the server for visualization.
func Run(db *sql.DB, seq sequence.Sequence, serverAddr string) error {
	base := "http://" + serverAddr

	r := runner.New(db)
	iter := runner.NewIterator(r)
	iter.Run(seq)

	for {
		e, ok := iter.Next()
		if !ok {
			break
		}
		if err := postEvent(base, protocol.Marshal(e)); err != nil {
			return err
		}
	}

	resp, err := http.Post(base+"/done", "application/json", http.NoBody)
	if err != nil {
		return fmt.Errorf("send done: %w", err)
	}
	resp.Body.Close()
	return nil
}

func postEvent(base string, msg protocol.EventMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	resp, err := http.Post(base+"/event", "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("send event: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}
	return nil
}
