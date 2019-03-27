package queue

import (
	"log"

	"github.com/dsociative/text-to-speech/tts"
)

type Queue struct {
	ttsChan    chan tts.Request
	queueChan  chan tts.Request
	resultChan chan tts.Response
	store      *Store
	wait       map[string][]tts.Request
}

func NewQueue(
	store *Store,
	queueChan chan tts.Request,
	ttsChan chan tts.Request,
	resultChan chan tts.Response,
) *Queue {
	return &Queue{
		store:      store,
		ttsChan:    ttsChan,
		queueChan:  queueChan,
		resultChan: resultChan,
		wait:       map[string][]tts.Request{},
	}
}

func (q *Queue) isWaitingSimilarKey(key string) bool {
	return q.wait[key] != nil
}

func (q *Queue) waitAppend(key string, request tts.Request) {
	q.wait[key] = append(q.wait[key], request)
}

func (q *Queue) Request(request tts.Request) {
	key := request.Key()
	data := q.store.Get(key)
	if data == nil {
		if q.isWaitingSimilarKey(key) {
			q.waitAppend(key, request)
		} else {
			select {
			case q.ttsChan <- request:
				q.waitAppend(key, request)
			case <-request.Ctx.Done():
				request.Timeout()
			}
		}
	} else {
		request.Done(data, nil)
	}
}

func (q *Queue) done(r tts.Response) {
	key := r.Request.Key()
	if waiting, ok := q.wait[key]; ok {
		for _, request := range waiting {
			request.Done(r.Data, r.Err)
		}
		delete(q.wait, key)
	} else {
		log.Println("finished the request with an empty wait list")
	}
	if r.Err == nil {
		q.store.Set(key, r.Data)
	}
}

func (q *Queue) Pool() {
	for {
		select {
		case request := <-q.queueChan:
			q.Request(request)
		case response := <-q.resultChan:
			q.done(response)
		}
	}
}
