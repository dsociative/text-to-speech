package tts

import (
	context "context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const yandexURL = "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize"

type Yandex struct {
	folder string
	iam    string
}

func NewYandex(folder, iam string) *Yandex {
	return &Yandex{folder: folder, iam: iam}
}

func (y *Yandex) TTS(ctx context.Context, text, lang string) ([]byte, error) {
	data := url.Values{
		"text":     {text},
		"lang":     {lang},
		"folderId": {y.folder},
	}
	request, err := http.NewRequest(
		"POST",
		yandexURL,
		strings.NewReader(data.Encode()),
	)
	request.Header.Set("Authorization", "Bearer "+y.iam)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(request.WithContext(ctx))
	if err == nil {
		if resp.StatusCode != http.StatusOK {
			b, err := httputil.DumpResponse(resp, resp.Body != nil)
			if err != nil {
				return nil, fmt.Errorf("yandex api error can't dump reponse %s", err)
			}
			return nil, errors.New("yandex api error: " + string(b))
		}
		return ioutil.ReadAll(resp.Body)
	}
	return nil, err
}
