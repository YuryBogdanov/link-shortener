package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
		{
			"Negative case #3 (some gibberish instead of a link)",
			"some_link_really",
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

			fn := handler.HandleNewLinkRegistration()
			fn.ServeHTTP(w, request)

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
	t.Setenv("SHORTENER_ENVIRONMENT", "test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.id) != 0 && len(tt.want.location) != 0 {
				storage.Links = make(map[string]string)
				storage.Links[tt.id] = tt.want.location
			}
			request := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			w := httptest.NewRecorder()
			fn := handler.HandleExistingLinkRequest()
			fn.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}

func Test_handleShortenRequest(t *testing.T) {
	type want struct {
		code                 int
		jsonResponseAsString string
	}
	tests := []struct {
		name                string
		requestBodyAsString string
		headerKey           string
		headerValue         string
		want                want
	}{
		{
			"Successful case",
			`{"url":"https://practicum.yandex.ru"}`,
			"Content-Type",
			"application/json",
			want{
				201,
				`{"result":"http://localhost:8080/6bdb5b0e"}`,
			},
		},
		{
			"Content negotiation",
			`{"url":"https://practicum.yandex.ru"}`,
			"Content-Type",
			"text/html; charset=utf-8",
			want{
				400,
				"",
			},
		},
	}
	t.Setenv("SHORTENER_ENVIRONMENT", "test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := strings.NewReader(tt.requestBodyAsString)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bodyReader)
			request.Header.Add(tt.headerKey, tt.headerValue)
			w := httptest.NewRecorder()
			fn := handler.HandleShortenRequest()
			fn.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			resultBody, err := io.ReadAll(result.Body)
			assert.Nil(t, err)

			resultBodyAsString := string(resultBody)
			assert.Equal(t, tt.want.jsonResponseAsString, resultBodyAsString)
		})
	}
}

func Test_CompressedPayloadHandling(t *testing.T) {
	type want struct {
		code                 int
		jsonResponseAsString string
	}
	tests := []struct {
		name                   string
		requestBodyAsString    string
		requestLink            string
		compressionHeaderKey   string
		compressionHeaderValue string
		want                   want
	}{
		{
			"Successful case #1",
			`{"url":"https://practicum.yandex.ru"}`,
			"https://practicum.yandex.ru",
			"Content-Encoding",
			"gzip",
			want{
				201,
				`{"result":"http://localhost:8080/6bdb5b0e"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			zw, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
			if err != nil {
				assert.Fail(t, "failed to create gzip writer")
			}
			_, _ = zw.Write([]byte(tt.requestBodyAsString))
			_ = zw.Close()
			reader := bytes.NewReader(buf.Bytes())
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", reader)
			request.Header.Add(tt.compressionHeaderKey, tt.compressionHeaderValue)
			w := httptest.NewRecorder()
			fn := handler.HandleShortenRequest()

			fn.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()
			resultBody, err := io.ReadAll(result.Body)
			assert.Nil(t, err)

			assert.Equal(t, tt.want.code, result.StatusCode)

			gzipReader, gunzipErr := gzip.NewReader(bytes.NewBuffer(resultBody))
			if gunzipErr != nil {
				assert.Fail(t, "failed to gunzip a response")
			}

			var gunzipedBuffer bytes.Buffer
			_, gErr := gunzipedBuffer.ReadFrom(gzipReader)
			if gErr != nil {
				assert.Fail(t, "couldn't read response")
			}

			resData := gunzipedBuffer.Bytes()
			assert.Equal(t, tt.want.jsonResponseAsString, string(resData))
		})
	}
}
