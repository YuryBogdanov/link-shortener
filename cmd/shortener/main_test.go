package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YuryBogdanov/link-shortener/internal/handler"
	"github.com/YuryBogdanov/link-shortener/internal/storage"
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
				"http://localhost:8080/6bdb5b0e",
			},
		},
		{
			"Positive case #2 (short link)",
			"https://go.dev",
			want{
				201,
				"http://localhost:8080/1dd1701d",
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

			handler.HandleNewLinkRegistration(w, request)

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

func Test_handleExistingLinkRequest(t *testing.T) {
	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name string
		id   string
		want want
	}{
		{
			"Positive case #1",
			"6bdb5b0e",
			want{
				307,
				"https://practicum.yandex.ru",
			},
		},
		{
			"Negative case #1",
			"invalid",
			want{
				400,
				"",
			},
		},
		{
			"Negative case #2",
			"",
			want{
				400,
				"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.id) != 0 && len(tt.want.location) != 0 {
				storage.Links = make(map[string]string)
				storage.Links[tt.id] = tt.want.location
			}
			request := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			w := httptest.NewRecorder()
			handler.HandleExistingLinkRequest(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}
