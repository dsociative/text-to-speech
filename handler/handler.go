package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dsociative/text-to-speech/tts"
)

type handler struct {
	queueChan chan tts.Request
	timeout   time.Duration
}

func NewHandler(queueChan chan tts.Request, timeout time.Duration) handler {
	return handler{
		queueChan: queueChan,
		timeout:   timeout,
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	request := tts.NewRequest(ctx, r.Form.Get("text"), r.Form.Get("lang"))
	select {
	case h.queueChan <- request:
		ttsResponse := request.Wait()
		if ttsResponse.Err == nil {
			w.Header().Set("Content-Type", "audio/ogg;codecs=opus")
			w.Write(ttsResponse.Data)
		} else {
			log.Println("tts handler err:", ttsResponse.Err)
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}
	case <-ctx.Done():
		http.Error(
			w,
			http.StatusText(http.StatusRequestTimeout),
			http.StatusRequestTimeout,
		)
	}
}
