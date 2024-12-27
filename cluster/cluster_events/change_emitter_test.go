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

	change := Change{Previous: "", Current: value}
	go s.emitter.Publish(key, change)

	select {
	case received := <-sub.Ch:
		s.Equal(change, received)
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

	change := Change{Previous: "", Current: value}
	go s.emitter.Publish(key, change)

	for _, ch := range []chan Change{sub1.Ch, sub2.Ch} {
		select {
		case received := <-ch:
			s.Equal(change, received)
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

	change1 := Change{Previous: "", Current: value1}
	go s.emitter.Publish(key1, change1)
	change2 := Change{Previous: "", Current: value2}
	go s.emitter.Publish(key2, change2)

	select {
	case received := <-sub1.Ch:
		s.Equal(change1, received)
	case <-time.After(time.Second):
		s.Fail("Timeout waiting for value1")
	}

	select {
	case received := <-sub2.Ch:
		s.Equal(change2, received)
	case <-time.After(time.Second):
		s.Fail("Timeout waiting for value2")
	}
}

func (s *ChangeEmitterSuite) TestSetOfValueChannels() {
	set := NewSetOfChangeChannels()
	ch := make(chan Change)

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
