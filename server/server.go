package server

import (
	"bufio"
	"io"
	"log"
	"net"

	"github.com/rusinikita/acid/event"
	"github.com/rusinikita/acid/protocol"
)

// Server listens for a single TCP client connection and forwards events to a channel.
type Server struct {
	addr  string
	bound string
	ch    chan event.Event
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
		ch:   make(chan event.Event),
	}
}

// Channel returns the read-only channel that the TUI drains for incoming events.
func (s *Server) Channel() <-chan event.Event {
	return s.ch
}

// Addr returns the actual bound address (e.g. "127.0.0.1:14322").
// Only valid after ListenAndServe returns without error.
func (s *Server) Addr() string {
	return s.bound
}

// ListenAndServe binds to the address and accepts one client in a background goroutine.
// It returns as soon as the listener is bound.
func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.bound = ln.Addr().String()

	go func() {
		defer ln.Close()
		defer close(s.ch)

		conn, err := ln.Accept()
		if err != nil {
			log.Printf("server accept: %v", err)
			return
		}
		defer conn.Close()

		s.readEvents(conn)
	}()

	return nil
}

func (s *Server) readEvents(r io.Reader) {
	br := bufio.NewReader(r)
	for {
		e, err := protocol.ReadEvent(br)
		if err != nil {
			if err != io.EOF {
				log.Printf("server read: %v", err)
			}
			return
		}
		s.ch <- e
	}
}
