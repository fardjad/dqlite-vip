package change_emitter

import (
	"sync"
)

type Value string
type Key string
type CancelFunc func()

type SetOfValueChannels struct {
	m map[chan Value]struct{}
}

func NewSetOfValueChannels() *SetOfValueChannels {
	return &SetOfValueChannels{
		m: make(map[chan Value]struct{}),
	}
}

func (s *SetOfValueChannels) Add(ch chan Value) {
	s.m[ch] = struct{}{}
}

func (s *SetOfValueChannels) Delete(ch chan Value) {
	delete(s.m, ch)
}

func (s *SetOfValueChannels) Contains(ch chan Value) bool {
	_, ok := s.m[ch]
	return ok
}

func (s *SetOfValueChannels) Values() []chan Value {
	values := make([]chan Value, 0, len(s.m))
	for ch := range s.m {
		values = append(values, ch)
	}
	return values
}

func (s *SetOfValueChannels) Empty() bool {
	return len(s.m) == 0
}

type Subscription struct {
	Ch     chan Value
	Cancel CancelFunc
}

type ChangeEmitter struct {
	mu            sync.RWMutex
	subscriptions map[Key]*SetOfValueChannels
}

func NewChangeEmitter() *ChangeEmitter {
	return &ChangeEmitter{
		subscriptions: make(map[Key]*SetOfValueChannels),
	}
}

func (w *ChangeEmitter) Subscribe(key Key) *Subscription {
	w.mu.Lock()
	defer w.mu.Unlock()

	ch := make(chan Value)
	sub := &Subscription{
		Ch: ch,
		Cancel: func() {
			w.mu.Lock()
			defer w.mu.Unlock()

			chSet := w.subscriptions[key]
			chSet.Delete(ch)
			close(ch)

			if chSet.Empty() {
				delete(w.subscriptions, key)
			}
		},
	}

	if _, ok := w.subscriptions[key]; !ok {
		w.subscriptions[key] = NewSetOfValueChannels()
	}
	w.subscriptions[key].Add(ch)

	return sub
}

func (w *ChangeEmitter) Publish(key Key, value Value) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if _, ok := w.subscriptions[key]; !ok {
		return
	}

	for _, ch := range w.subscriptions[key].Values() {
		ch <- value
	}
}
