package cluster_events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ChangeEmitterSuite struct {
	suite.Suite
	emitter *ChangeEmitter
}

func (s *ChangeEmitterSuite) SetupTest() {
	s.emitter = NewChangeEmitter()
}

func (s *ChangeEmitterSuite) TestBasicSubscribeAndPublish() {
	key := "test-key"
	value := "test-value"

	sub := s.emitter.Subscribe(key)
	defer sub.Cancel()

	go s.emitter.Publish(key, value)

	// Should receive the value
	select {
	case received := <-sub.Ch:
		s.Equal(value, received)
	case <-time.After(time.Second):
		s.Fail("Timeout waiting for value")
	}
}

func (s *ChangeEmitterSuite) TestMultipleSubscribers() {
	key := "test-key"
	value := "test-value"

	sub1 := s.emitter.Subscribe(key)
	sub2 := s.emitter.Subscribe(key)
	defer sub1.Cancel()
	defer sub2.Cancel()

	go s.emitter.Publish(key, value)

	for _, ch := range []chan string{sub1.Ch, sub2.Ch} {
		select {
		case received := <-ch:
			s.Equal(value, received)
		case <-time.After(time.Second):
			s.Fail("Timeout waiting for value")
		}
	}
}

func (s *ChangeEmitterSuite) TestCancellation() {
	key := "test-key"
	sub := s.emitter.Subscribe(key)

	sub.Cancel()

	_, ok := <-sub.Ch
	s.False(ok, "Channel should be closed")

	_, exists := s.emitter.Subscriptions[key]
	s.False(exists, "Subscription should be removed after cancellation")
}

func (s *ChangeEmitterSuite) TestMultipleKeys() {
	key1 := "key1"
	key2 := "key2"
	value1 := "value1"
	value2 := "value2"

	sub1 := s.emitter.Subscribe(key1)
	sub2 := s.emitter.Subscribe(key2)
	defer sub1.Cancel()
	defer sub2.Cancel()

	go s.emitter.Publish(key1, value1)
	go s.emitter.Publish(key2, value2)

	select {
	case received := <-sub1.Ch:
		s.Equal(value1, received)
	case <-time.After(time.Second):
		s.Fail("Timeout waiting for value1")
	}

	select {
	case received := <-sub2.Ch:
		s.Equal(value2, received)
	case <-time.After(time.Second):
		s.Fail("Timeout waiting for value2")
	}
}

func (s *ChangeEmitterSuite) TestSetOfValueChannels() {
	set := NewSetOfValueChannels()
	ch := make(chan string)

	set.Add(ch)
	s.True(set.Contains(ch))

	values := set.Values()
	s.Len(values, 1)
	s.Equal(ch, values[0])

	set.Delete(ch)
	s.False(set.Contains(ch))
	s.True(set.Empty())
}

func TestChangeEmitterSuite(t *testing.T) {
	suite.Run(t, new(ChangeEmitterSuite))
}
