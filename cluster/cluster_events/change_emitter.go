package cluster_events

import (
	"sync"
)

type CancelFunc func()

type SetOfValueChannels struct {
	m map[chan string]struct{}
}

func NewSetOfValueChannels() *SetOfValueChannels {
	return &SetOfValueChannels{
		m: make(map[chan string]struct{}),
	}
}

func (s *SetOfValueChannels) Add(ch chan string) {
	s.m[ch] = struct{}{}
}

func (s *SetOfValueChannels) Delete(ch chan string) {
	delete(s.m, ch)
}

func (s *SetOfValueChannels) Contains(ch chan string) bool {
	_, ok := s.m[ch]
	return ok
}

func (s *SetOfValueChannels) Values() []chan string {
	values := make([]chan string, 0, len(s.m))
	for ch := range s.m {
		values = append(values, ch)
	}
	return values
}

func (s *SetOfValueChannels) Empty() bool {
	return len(s.m) == 0
}

type Subscription struct {
	Ch     chan string
	Cancel CancelFunc
}

type ChangeEmitter struct {
	mu            sync.RWMutex
	Subscriptions map[string]*SetOfValueChannels
}

func NewChangeEmitter() *ChangeEmitter {
	return &ChangeEmitter{
		Subscriptions: make(map[string]*SetOfValueChannels),
	}
}

func (w *ChangeEmitter) Subscribe(key string) *Subscription {
	w.mu.Lock()
	defer w.mu.Unlock()

	ch := make(chan string)
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
		w.Subscriptions[key] = NewSetOfValueChannels()
	}
	w.Subscriptions[key].Add(ch)

	return sub
}

func (w *ChangeEmitter) Publish(key string, value string) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if _, ok := w.Subscriptions[key]; !ok {
		return
	}

	for _, ch := range w.Subscriptions[key].Values() {
		ch <- value
	}
}
