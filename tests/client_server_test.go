package tests

import (
	"os"

	"github.com/rusinikita/acid/client"
	"github.com/rusinikita/acid/config"
	"github.com/rusinikita/acid/protocol"
	"github.com/rusinikita/acid/runner"
	"github.com/rusinikita/acid/server"
)

func (s *PostgresSuite) TestClientServerRoundTrip() {
	// Load sequence from a temp TOML file (same content as toml_load_test)
	f, err := os.CreateTemp("", "client_server_*.toml")
	s.Require().NoError(err)
	defer os.Remove(f.Name())
	_, err = f.WriteString(lostUpdateTOML)
	s.Require().NoError(err)
	s.Require().NoError(f.Close())

	seq, err := config.Load(f.Name())
	s.Require().NoError(err)

	// Start a server on a random port
	srv := server.New("127.0.0.1:0")
	s.Require().NoError(srv.ListenAndServe())

	// Run client in background: connects to server, executes seq, streams events
	clientErr := make(chan error, 1)
	go func() {
		clientErr <- client.Run(s.db, seq, srv.Addr())
	}()

	// Collect all events the server received and marshal them for comparison
	var received []protocol.EventMessage
	for e := range srv.Channel() {
		msg := protocol.Marshal(e)
		msg.Waiting = nil // waiting is timing-dependent; excluded from comparison
		received = append(received, msg)
	}
	s.Require().NoError(<-clientErr)

	// Build reference by running the same sequence locally
	r := runner.New(s.db)
	iter := runner.NewIterator(r)
	iter.Run(seq)

	var expected []protocol.EventMessage
	for {
		e, ok := iter.Next()
		if !ok {
			break
		}
		msg := protocol.Marshal(e)
		msg.Waiting = nil
		expected = append(expected, msg)
	}

	s.Equal(expected, received)
}
