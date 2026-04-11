package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/rusinikita/acid/event"
	"github.com/rusinikita/acid/protocol"
)

// Server is an HTTP server that receives events from a remote client.
// GET  /health        — liveness probe
// POST /event         — deliver one event (JSON body: protocol.EventMessage)
// POST /done          — signal end of stream; closes the event channel
// POST /toggle-mode   — toggle result visibility on the server TUI
type Server struct {
	addr     string
	bound    string
	ch       chan event.Event
	toggleCh chan struct{}
	mu       sync.Mutex
	visible  bool // false = hidden (onlyStepsMode=true), true = visible
}

func New(addr string) *Server {
	return &Server{
		addr:     addr,
		ch:       make(chan event.Event),
		toggleCh: make(chan struct{}, 1),
	}
}

// ToggleCh returns the channel that fires when a /toggle-mode request is received.
func (s *Server) ToggleCh() <-chan struct{} {
	return s.toggleCh
}

// Channel returns the read-only channel that the TUI drains for incoming events.
func (s *Server) Channel() <-chan event.Event {
	return s.ch
}

// Addr returns the actual bound address (e.g. "127.0.0.1:7331").
// Only valid after ListenAndServe returns without error.
func (s *Server) Addr() string {
	return s.bound
}

// ListenAndServe binds and starts the HTTP server in a background goroutine.
func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.bound = ln.Addr().String()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("POST /event", func(w http.ResponseWriter, r *http.Request) {
		var msg protocol.EventMessage
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.ch <- protocol.Unmarshal(msg)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("POST /start", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		s.visible = false
		s.mu.Unlock()
		s.ch <- event.Start()
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("POST /done", func(w http.ResponseWriter, r *http.Request) {
		s.ch <- event.Done()
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("POST /toggle-mode", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		select {
		case s.toggleCh <- struct{}{}:
			s.visible = !s.visible
			v := s.visible
			s.mu.Unlock()
			from, to := "hidden", "visible"
			if !v {
				from, to = "visible", "hidden"
			}
			fmt.Fprintf(w, "toggled: %s -> %s", from, to)
		default:
			s.mu.Unlock()
			http.Error(w, "toggle already pending", http.StatusConflict)
		}
	})

	go func() { _ = http.Serve(ln, mux) }()

	return nil
}
