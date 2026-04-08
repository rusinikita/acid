package runner

import (
	"github.com/rusinikita/acid/event"
	"github.com/rusinikita/acid/sequence"
)

// ChannelSource satisfies the ui/run runner interface for server mode.
// Events are provided externally via a channel populated by the network listener.
type ChannelSource struct {
	ch <-chan event.Event
}

func NewChannelSource(ch <-chan event.Event) *ChannelSource {
	return &ChannelSource{ch: ch}
}

// Run is a no-op; the server receives events from the network, not from a local DB.
func (cs *ChannelSource) Run(_ sequence.Sequence) {}

// Next blocks until an event is available or the channel is closed.
func (cs *ChannelSource) Next() (event.Event, bool) {
	e, ok := <-cs.ch
	return e, ok
}
