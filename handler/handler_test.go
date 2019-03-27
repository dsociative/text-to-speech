package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dsociative/text-to-speech/tts"
	"github.com/stretchr/testify/assert"
)

func TestTimeoutWaitTTS(t *testing.T) {
	queueChan := make(chan tts.Request)
	handler := NewHandler(queueChan, 100*time.Millisecond)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	response := w.Result()
	assert.Equal(t, http.StatusRequestTimeout, response.StatusCode)
}

func TestMulitContextUsage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	<-time.After(110 * time.Millisecond)
	select {
	case <-ctx.Done():
		t.Log("context is done")
	default:
		t.Log("context is ok")
	}

	t.Log("wait")
	<-ctx.Done()
	t.Log("ctx done")
}

func TestParalllelContextUsage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	done := make(chan bool, 2)

	for i := 0; i < 2; i++ {
		go func() {
			<-ctx.Done()
			done <- true
			t.Log(ctx.Err())
		}()
	}

	<-time.After(120 * time.Millisecond)

	ok := 0
	for ok < 2 {
		select {
		case <-done:
			t.Log(ok)
			ok++
		default:
			t.Fatal("channel dont have 2 result")
		}
	}
	assert.Equal(t, 2, ok)
}
