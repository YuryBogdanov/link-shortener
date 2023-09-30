package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_handleNewLinkRegistration(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name string
		link string
		want want
	}{
		{
			"Positive case #1 (regular link)",
			"https://practicum.yandex.ru",
			want{
				201,
				"http://example.com/6bdb5b0e",
			},
		},
		{
			"Positive case #2 (short link)",
			"https://go.dev",
			want{
				201,
				"http://example.com/1dd1701d",
			},
		},
		{
			"Negative case #1 (empty body)",
			"",
			want{
				400,
				"",
			},
		},
		{
			"Negative case #2 (URL Scheme (https://, etc) only)",
			"https://",
			want{
				400,
				"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(tt.link)))
			w := httptest.NewRecorder()

			handleNewLinkRegistration(w, request)

			result := w.Result()
			assert.Equal(t, tt.want.code, result.StatusCode)

			defer result.Body.Close()
			userResult, err := io.ReadAll(result.Body)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.want.response, string(userResult))
		})
	}
}
