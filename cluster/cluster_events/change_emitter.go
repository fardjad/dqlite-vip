package cluster_events

import (
	"sync"
)

type CancelFunc func()

type SetOfChangeChannels struct {
	m map[chan Change]struct{}
}

func NewSetOfChangeChannels() *SetOfChangeChannels {
	return &SetOfChangeChannels{
		m: make(map[chan Change]struct{}),
	}
}

func (s *SetOfChangeChannels) Add(ch chan Change) {
	s.m[ch] = struct{}{}
}

func (s *SetOfChangeChannels) Delete(ch chan Change) {
	delete(s.m, ch)
}

func (s *SetOfChangeChannels) Contains(ch chan Change) bool {
	_, ok := s.m[ch]
	return ok
}

func (s *SetOfChangeChannels) Values() []chan Change {
	values := make([]chan Change, 0, len(s.m))
	for ch := range s.m {
		values = append(values, ch)
	}
	return values
}

func (s *SetOfChangeChannels) Empty() bool {
	return len(s.m) == 0
}

type Change struct {
	Previous string
	Current  string
}

type Subscription struct {
	Ch     chan Change
	Cancel CancelFunc
}

type ChangeEmitter struct {
	mu            sync.RWMutex
	Subscriptions map[string]*SetOfChangeChannels
}

func NewChangeEmitter() *ChangeEmitter {
	return &ChangeEmitter{
		Subscriptions: make(map[string]*SetOfChangeChannels),
	}
}

func (w *ChangeEmitter) Subscribe(key string) *Subscription {
	w.mu.Lock()
	defer w.mu.Unlock()

	ch := make(chan Change)
	sub := &Subscription{
		Ch: ch,
		Cancel: func() {
			w.mu.Lock()
			defer w.mu.Unlock()

			chSet := w.Subscriptions[key]
			chSet.Delete(ch)
			close(ch)

			if chSet.Empty() {
				delete(w.Subscriptions, key)
			}
		},
	}

	if _, ok := w.Subscriptions[key]; !ok {
		w.Subscriptions[key] = NewSetOfChangeChannels()
	}
	w.Subscriptions[key].Add(ch)

	return sub
}

func (w *ChangeEmitter) Publish(key string, change Change) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if _, ok := w.Subscriptions[key]; !ok {
		return
	}

	for _, ch := range w.Subscriptions[key].Values() {
		ch <- change
	}
}
