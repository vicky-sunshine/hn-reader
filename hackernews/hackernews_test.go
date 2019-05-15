package hackernews

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mock hacker news
type HNServer struct{}

func (srv *HNServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check method and path
	fmt.Println(r.Method)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "/v0/topstories.json" {
		w.Header().Add("Content-Type", "application/json")
		dat, _ := ioutil.ReadFile("testdata/topstories.json")
		w.Write([]byte(dat))
		return
	}

	if r.URL.Path == "/v0/item/19893682.json" {
		w.Header().Add("Content-Type", "application/json")
		dat, _ := ioutil.ReadFile("testdata/19893682.json")
		w.Write([]byte(dat))
		return
	}

	// nothing match
	w.WriteHeader(http.StatusNotFound)
	return
}

func TestTopArticles(t *testing.T) {
	srv := HNServer{}
	hksrv := httptest.NewServer(&srv)
	defer hksrv.Close()

	apiBase := fmt.Sprintf("%v/v0", hksrv.URL)
	hnr := NewHNReader(apiBase)

	ids, _ := hnr.TopArticles(5)
	target := []string{"19893682", "19895335", "19894798", "19895218", "19895766"}
	if len(ids) != len(target) {
		t.Errorf("got: %#v; expect: %#v", ids, target)
	}

	for i, v := range target {
		if v != target[i] {
			t.Errorf("got: %#v; expect: %#v", v, target[i])
		}
	}
}

func TestGetArticle(t *testing.T) {
	srv := HNServer{}
	hksrv := httptest.NewServer(&srv)
	defer hksrv.Close()

	apiBase := fmt.Sprintf("%v/v0", hksrv.URL)
	hnr := NewHNReader(apiBase)

	article, _ := hnr.GetArticle("19893682")
	if id := article.GetID(); id != "19893682" {
		t.Errorf("got: %#v; expect: %#v", id, "19893682")
	}
	if title := article.GetTitle(); title != "Facebook sues analytics firm Rankwave over data misuse" {
		t.Errorf("got: %#v; expect: %#v", title, "Facebook sues analytics firm Rankwave over data misuse")
	}
	if auth := article.GetAuthor(); auth != "JumpCrisscross" {
		t.Errorf("got: %#v; expect: %#v", auth, "JumpCrisscross")
	}
	if ts := article.GetTimestamp(); ts != 1557686813 {
		t.Errorf("got: %#v; expect: %#v", ts, 1557686813)
	}
}
