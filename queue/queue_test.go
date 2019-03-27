package queue

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dsociative/text-to-speech/tts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type QueueTestSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	store   *Store
	ttsMock *tts.MockTTS
}

func (s *QueueTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.store = NewStore()
	s.ttsMock = tts.NewMockTTS(s.ctrl)
}

func (s *QueueTestSuite) TestNewRequest() {
	ctx, _ := context.WithTimeout(context.Background(), 200*time.Millisecond)
	r := tts.NewRequest(ctx, "text", "ru-RU")

	ttsChan := make(chan tts.Request, 1)
	q := NewQueue(s.store, nil, ttsChan, nil)
	q.Request(r)
	s.Equal(r, <-ttsChan)
	s.EqualValues([]tts.Request{r}, q.wait["textru-RU"])
}

func (s *QueueTestSuite) TestCachedRequest() {
	r := tts.NewRequest(nil, "cache", "ru-RU")
	s.store.Set(r.Key(), []byte("cache"))
	q := NewQueue(s.store, nil, nil, nil)
	q.Request(r)
	s.Nil(q.wait["textru-RU"])
	s.Equal(tts.Response{Data: []byte("cache")}, r.Wait())
}

func (s *QueueTestSuite) TestDone() {
	r := tts.NewRequest(context.Background(), "text", "ru-RU")
	ttsChan := make(chan tts.Request, 1)
	q := NewQueue(s.store, nil, ttsChan, nil)

	q.Request(r)
	s.Equal(r, <-ttsChan)

	resp := tts.NewResponse(r, []byte("ok"), nil)
	q.done(resp)
	s.Equal([]tts.Request(nil), q.wait["textru-RU"])
	s.Equal([]byte("ok"), s.store.Get(r.Key()))
	s.Equal(tts.Response{Data: []byte("ok")}, r.Wait())
}

func (s *QueueTestSuite) TestDoneError() {
	r := tts.NewRequest(context.Background(), "text", "ru-RU")
	ttsChan := make(chan tts.Request, 1)
	q := NewQueue(s.store, nil, ttsChan, nil)

	q.Request(r)
	s.Equal(r, <-ttsChan)

	q.done(tts.NewResponse(r, []byte("error"), errors.New("error")))
	s.Equal([]tts.Request(nil), q.wait["textru-RU"])
	s.Equal([]byte(nil), s.store.Get(r.Key()))
}

func TestQueueSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}
