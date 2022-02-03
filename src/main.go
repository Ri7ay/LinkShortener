//go:build !solution
// +build !solution

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
var data urlData

type urlDB struct {
	Key string `db:"shorturl"`
	URL string `db:"url"`
}

type getURL struct {
	URL string `json:"url"`
}

type urlData interface {
	store(key string, url string)
	loadKey(url string) (key string, ok bool)
	loadURL(key string) (url string, ok bool)
}

type dataInMem struct {
	sync.RWMutex
	URLKey, KeyURL map[string]string
}

func (dim *dataInMem) store(key string, url string) {
	dim.Lock()
	defer dim.Unlock()

	dim.URLKey[url] = key
	dim.KeyURL[key] = url
}

func (dim *dataInMem) loadKey(url string) (key string, ok bool) {
	dim.RLock()
	defer dim.RUnlock()

	key, ok = dim.URLKey[url]
	return
}

func (dim *dataInMem) loadURL(key string) (url string, ok bool) {
	dim.RLock()
	defer dim.RUnlock()

	url, ok = dim.KeyURL[key]
	return
}

type dataInSQL struct {
	db *sqlx.DB
}

func (dis *dataInSQL) store(key string, url string) {
	dis.db.MustExec("INSERT INTO surls (shorturl, url) VALUES ($1, $2)", key, url)
}

func (dis *dataInSQL) loadKey(url string) (key string, ok bool) {
	el := urlDB{}
	dis.db.Get(&el, "SELECT * FROM surls WHERE url=$1", url)
	key, ok = el.Key, el.Key != ""
	return
}

func (dis *dataInSQL) loadURL(key string) (url string, ok bool) {
	el := urlDB{}
	dis.db.Get(&el, "SELECT * FROM surls WHERE shorturl=$1", key)
	url, ok = el.URL, el.URL != ""
	return
}

func main() {
	port := os.Getenv("PORT")
	db := connectDB()
	defer db.Close()
	flagDB := os.Getenv("DBFLAG")
	if flagDB == "true" {
		data = &dataInSQL{db: db}
	} else {
		data = &dataInMem{
			URLKey: make(map[string]string),
			KeyURL: make(map[string]string),
		}
	}
	http.HandleFunc("/", handleSwitch)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic("error!")
	}
}

func handleSwitch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		handlePost(w, r)
	case http.MethodGet:
		hadleGet(w, r)
	}
}

func connectDB() *sqlx.DB {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v?sslmode=%v",
		"postgres",
		"password",
		host,
		"5432",
		"disable")
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	url := r.URL.Path[1:]
	if url == "" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if key, ok := data.loadKey(url); ok {
		createResp(w, key)
		return
	}

	key, ok := "", true
	for ok {
		key = randKey()
		_, ok = data.loadURL(key)
	}
	data.store(key, url)
	createResp(w, key)
}

func randKey() string {
	var build strings.Builder
	build.Grow(10)
	for i := 0; i < 10; i++ {
		build.WriteRune(letters[rand.Intn(len(letters))])
	}
	return build.String()
}

func hadleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	url, ok := data.loadURL(key)
	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	createResp(w, url)
}

func createResp(w http.ResponseWriter, url string) {
	jsonResp := getURL{URL: url}
	resp, err := json.Marshal(jsonResp)
	if err != nil {
		http.Error(w, "JSON is invalid", 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
