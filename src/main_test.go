package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPost(t *testing.T) {
	data = &dataInMem{
		URLKey: make(map[string]string),
		KeyURL: make(map[string]string),
	}

	uniq := make(map[string]struct{})

	tt := []struct {
		name       string
		input      string
		len        int
		statusCode int
		statusOK   bool
	}{
		{
			name:       "youtube",
			input:      "www.youtube.com",
			len:        10,
			statusCode: http.StatusOK,
			statusOK:   false,
		},
		{

			name:       "google",
			input:      "www.google.com",
			len:        10,
			statusCode: http.StatusOK,
			statusOK:   false,
		},
		{

			name:       "ozon",
			input:      "www.ozon.ru",
			len:        10,
			statusCode: http.StatusOK,
			statusOK:   false,
		},
		{

			name:       "ozon again",
			input:      "www.ozon.ru",
			len:        10,
			statusCode: http.StatusOK,
			statusOK:   true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/"+tc.input, nil)
			w := httptest.NewRecorder()

			handleSwitch(w, req)

			if w.Code != tc.statusCode {
				t.Errorf("want status '%d', got '%d'", tc.statusCode, w.Code)
			}

			res := w.Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			get := getURL{}
			json.Unmarshal(data, &get)

			if len(get.URL) != tc.len {
				t.Errorf("expected link length of 10 instead of %v", len(get.URL))
			}

			if _, ok := uniq[get.URL]; ok != tc.statusOK {
				if tc.name != "ozon again" {
					t.Errorf("expected uniq link for %v, short link: %v", tc.name, get.URL)
				} else {
					t.Errorf("The shortened link must already exist, short link: %v", get.URL)
				}
			}
			uniq[get.URL] = struct{}{}
		})
	}
}

func initData() error {
	const lenURL = 10
	tt := []struct {
		name  string
		input string
	}{
		{
			name:  "youtube",
			input: "www.youtube.com",
		},
		{

			name:  "google",
			input: "www.google.com",
		},
		{

			name:  "ozon",
			input: "www.ozon.ru",
		},
	}

	for _, tc := range tt {
		req := httptest.NewRequest(http.MethodPost, "/"+tc.input, nil)
		w := httptest.NewRecorder()

		handleSwitch(w, req)

		if w.Code != http.StatusOK {
			return fmt.Errorf("want status '%d', got '%d'", http.StatusOK, w.Code)
		}

		res := w.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		get := getURL{}
		json.Unmarshal(data, &get)

		if len(get.URL) != lenURL {
			return fmt.Errorf("expected link length of 10 instead of %v", len(get.URL))
		}

	}
	return nil
}

func TestGet(t *testing.T) {
	data = &dataInMem{
		URLKey: make(map[string]string),
		KeyURL: make(map[string]string),
	}

	if err := initData(); err != nil {
		t.Errorf("init error: %v", err)
	}

	tt := []struct {
		name       string
		input      string
		want       string
		statusCode int
	}{
		{
			name:       "youtube",
			input:      data.(*dataInMem).URLKey["www.youtube.com"],
			want:       "www.youtube.com",
			statusCode: http.StatusOK,
		},
		{

			name:       "google",
			input:      data.(*dataInMem).URLKey["www.google.com"],
			want:       "www.google.com",
			statusCode: http.StatusOK,
		},
		{

			name:       "ozon",
			input:      data.(*dataInMem).URLKey["www.ozon.ru"],
			want:       "www.ozon.ru",
			statusCode: http.StatusOK,
		},
		{

			name:       "key not found",
			want:       "",
			input:      "www.bad_key.ru",
			statusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/"+tc.input, nil)
			w := httptest.NewRecorder()

			handleSwitch(w, req)

			if w.Code != tc.statusCode {
				t.Errorf("want status '%d', got '%d'", tc.statusCode, w.Code)
			}

			res := w.Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			get := getURL{}
			json.Unmarshal(data, &get)

			if get.URL != tc.want {
				t.Errorf("get.URL != tc.want, get.URL: %v, tc.wand: %v", get.URL, tc.want)
			}
		})
	}
}
