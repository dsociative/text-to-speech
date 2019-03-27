mock:
	${GOPATH}/bin/mockgen -package=tts -destination tts/tts_mock.go github.com/dsociative/text-to-speech/tts TTS
