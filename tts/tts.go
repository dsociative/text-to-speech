package tts

import (
	"context"
	"errors"
)

var errTimeout = errors.New("request timeout")

type TTS interface {
	TTS(ctx context.Context, text, lang string) ([]byte, error)
}

type Request struct {
	Text         string
	Lang         string
	Ctx          context.Context
	responseChan chan Response
}

type Response struct {
	Data    []byte
	Err     error
	Request Request
}

func NewResponse(request Request, data []byte, err error) Response {
	return Response{Request: request, Data: data, Err: err}
}

func NewRequest(ctx context.Context, text, lang string) Request {
	return Request{
		Text:         text,
		Lang:         lang,
		Ctx:          ctx,
		responseChan: make(chan Response, 1),
	}
}

func (r Request) Key() string {
	return r.Text + r.Lang
}

func (r Request) Done(b []byte, err error) {
	r.responseChan <- Response{Data: b, Err: err}
}

func (r Request) Timeout() {
	r.responseChan <- Response{Err: errTimeout}
}

func (r Request) Wait() Response {
	return <-r.responseChan
}
