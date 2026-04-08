package client

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/rusinikita/acid/protocol"
	"github.com/rusinikita/acid/runner"
	"github.com/rusinikita/acid/sequence"
)

// Run connects to the server at serverAddr, executes seq against db,
// and streams events to the server for visualization.
func Run(db *sql.DB, seq sequence.Sequence, serverAddr string) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("connect to server %s: %w", serverAddr, err)
	}
	defer conn.Close()

	r := runner.New(db)
	iter := runner.NewIterator(r)
	iter.Run(seq)

	for {
		e, ok := iter.Next()
		if !ok {
			break
		}
		if err := protocol.WriteEvent(conn, e); err != nil {
			return fmt.Errorf("send event: %w", err)
		}
	}

	return nil
}
