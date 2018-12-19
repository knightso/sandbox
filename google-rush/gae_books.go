package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/taskqueue"
)

// BookStatus describes status of Book
type BookStatus int

// BookStatus constants
const (
	BookStatusUnpublished BookStatus = 1 << iota
	BookStatusPublished
	BookStatusDiscontinued
)

// Book is sample model.
type Book struct {
	Number        int
	Title         string
	TitleIndex    []string
	TitlePrefix   []string
	Price         int
	PriceRange    string
	Category      string
	Status        BookStatus
	StatusORIndex []int
	IsPublished   bool
	IsHobby       bool
	Indexes       []string // for XIAN
	Shard256      string   // 分散処理用shard番号(256分割)
	Shard16       string   // 分散処理用shard番号(16分割)
}

// BookStatuses is the list of all BookStatuses
var BookStatuses = []BookStatus{
	BookStatusUnpublished,
	BookStatusPublished,
	BookStatusDiscontinued,
}

// BookCategories is the list of all Book categories
var BookCategories = []string{
	"sports",
	"cooking",
	"education",
	"cartoons",
}

func addHashPrefix(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x-%s", h.Sum(nil)[:3], s)
}

func call(r *http.Request, start, count int) error {
	task := taskqueue.NewPOSTTask("/put-testbooks/create-data", url.Values{
		"start": {strconv.Itoa(start)},
		"count": {strconv.Itoa(count)},
	})

	ctx := appengine.NewContext(r)
	_, err := taskqueue.Add(ctx, task, "create-test-data")

	return err
}

func putTestBooksTask(w http.ResponseWriter, r *http.Request) {
	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil {
		log.Printf("invalid start: %v", err.Error())
		return
	}

	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		log.Printf("invalid count: %v", err.Error())
		return
	}

	end, err := strconv.Atoi(r.FormValue("end"))
	if err != nil {
		log.Printf("invalid end: %v", err.Error())
		return
	}

	for i := start; i < end; i += count {
		if err = call(r, i, count); err != nil {
			log.Printf("error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func putTestBooks(w http.ResponseWriter, r *http.Request) {
	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil {
		log.Printf("invalid start: %v", err.Error())
		return
	}

	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		log.Printf("invalid count: %v", err.Error())
		return
	}

	ctx := appengine.NewContext(r)

	keys := make([]*datastore.Key, 0, count)
	books := make([]*Book, 0, count)

	titles := []string{"aPPle", "piNEApple", "banANA", "foobar", "hogefugapiyo"}

	for i := start; i < start+count; i++ {
		bookID := addHashPrefix(fmt.Sprintf("%d", i))
		key := datastore.NewKey(ctx, "Book", bookID, 0, nil)

		book := &Book{
			Number:   i,
			Title:    titles[i%len(titles)],
			Price:    i * 100,
			Status:   BookStatuses[i%len(BookStatuses)],
			Category: BookCategories[i%len(BookCategories)],
		}

		// Book保存時に派生プロパティを補完
		book.TitleIndex = biunigrams(book.Title)
		book.TitlePrefix = prefixes(book.Title)
		book.IsPublished = book.Status == BookStatusPublished
		book.IsHobby = book.Category == "sports" || book.Category == "cooking"

		for j := 1; j < 1<<uint(len(BookStatuses))+1; j++ {
			if j&int(book.Status) != 0 {
				book.StatusORIndex = append(book.StatusORIndex, j)
			}
		}

		switch {
		case book.Price < 3000:
			book.PriceRange = "p<3000"
		case book.Price < 5000:
			book.PriceRange = "3000<=p<5000"
		case book.Price < 10000:
			book.PriceRange = "5000<=p<10000"
		default:
			book.PriceRange = "10000<=p"
		}

		if err := saveBookIndexes(book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// shard番号設定
		if book.Shard256 == "" && book.Shard16 == "" {
			book.Shard256 = bookID[:2]
			book.Shard16 = bookID[:1]
		}

		keys = append(keys, key)
		books = append(books, book)
	}

	if _, err := datastore.PutMulti(ctx, keys, books); err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func gaeNotEqual(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var books []Book
	_, err := datastore.NewQuery("Book").Filter("IsPublished =", false).GetAll(ctx, &books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(books); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func gaeIn(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var books []Book
	_, err := datastore.NewQuery("Book").Filter("IsHobby =", true).GetAll(ctx, &books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(books); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func gaeIn2(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var books []Book
	_, err := datastore.NewQuery("Book").Filter("StatusORIndex =", BookStatusUnpublished|BookStatusPublished).GetAll(ctx, &books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(books); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func gaeNumRange(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var books []Book
	_, err := datastore.NewQuery("Book").Filter("PriceRange =", "5000<=p<10000").GetAll(ctx, &books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(books); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func gaeLike(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	title := r.FormValue("title")

	q := datastore.NewQuery("Book")

	if runeLen := utf8.RuneCountInString(title); runeLen == 1 {
		// パラメータが1文字の場合はunigramで検索
		q = q.Filter("TitleIndex =", title)
	} else if runeLen > 1 {
		// パラメータが2文字以上の場合はbigramで検索
		for _, gram := range bigrams(title) {
			q = q.Filter("TitleIndex =", gram)
		}
	}

	var books []Book
	_, err := q.GetAll(ctx, &books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(books); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func gaePrefix(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	title := r.FormValue("title")

	var books []Book
	_, err := datastore.NewQuery("Book").Filter("TitlePrefix =", strings.ToLower(title)).GetAll(ctx, &books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(books); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func bigrams(s string) []string {
	tokens := make([]string, 0, 32)

	for bigram := range toBigrams(strings.ToLower(s)) {
		tokens = append(tokens, fmt.Sprintf("%c%c", bigram.a, bigram.b))
	}

	return tokens

}

func biunigrams(s string) []string {
	tokens := make([]string, 0, 32)

	for bigram := range toBigrams(strings.ToLower(s)) {
		tokens = append(tokens, fmt.Sprintf("%c%c", bigram.a, bigram.b))
	}
	for unigram := range toUnigrams(strings.ToLower(s)) {
		tokens = append(tokens, fmt.Sprintf("%c", unigram))
	}

	return tokens

}

func prefixes(s string) []string {
	prefixes := make(map[string]struct{})

	runes := make([]rune, 0, 64)

	for _, w := range strings.Split(strings.ToLower(s), " ") {
		if w == "" {
			continue
		}

		runes = runes[0:0]

		for _, c := range w {
			runes = append(runes, c)
			prefixes[string(runes)] = struct{}{}
		}
	}

	tokens := make([]string, 0, 32)

	for pref := range prefixes {
		tokens = append(tokens, pref)
	}

	return tokens
}
