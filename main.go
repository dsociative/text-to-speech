package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/dsociative/text-to-speech/handler"
	"github.com/dsociative/text-to-speech/queue"
	"github.com/dsociative/text-to-speech/tts"
)

var (
	yandexFolder = flag.String("yandex_folder", "", "yandex folder")
	yandexIAM    = flag.String("yandex_iam", "", "yandex iam token")
	bind         = flag.String("bind", ":8080", "bind")
	timeout      = flag.Duration("timeout", 500*time.Millisecond, "request timeout")
)

func main() {
	flag.Parse()
	if *yandexFolder == "" || *yandexIAM == "" {
		log.Fatal("enter yandex credentials")
	}
	yandex := tts.NewYandex(*yandexFolder, *yandexIAM)

	queueChan := make(chan tts.Request)
	ttsChan := make(chan tts.Request)
	resultChan := make(chan tts.Response)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for r := range ttsChan {
				b, err := yandex.TTS(r.Ctx, r.Text, r.Lang)
				resultChan <- tts.NewResponse(r, b, err)
			}
		}()
	}

	q := queue.NewQueue(
		queue.NewStore(),
		queueChan,
		ttsChan,
		resultChan,
	)

	go q.Pool()

	h := handler.NewHandler(queueChan, *timeout)
	http.Handle("/tts", h)
	log.Fatal(http.ListenAndServe(*bind, nil))
}
