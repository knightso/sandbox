package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/knightso/base/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	aelog "google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

func init() {
	errors.ShowStackTraceOnError = true
}

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
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Book2 is sample model.
type Book2 struct {
	Number        int        `datastore:",noindex"`
	Title         string     `datastore:",noindex"`
	TitleIndex    []string   `datastore:",noindex"`
	TitlePrefix   []string   `datastore:",noindex"`
	Price         int        `datastore:",noindex"`
	PriceRange    string     `datastore:",noindex"`
	Category      string     `datastore:",noindex"`
	Status        BookStatus `datastore:",noindex"`
	StatusORIndex []int      `datastore:",noindex"`
	IsPublished   bool       `datastore:",noindex"`
	IsHobby       bool       `datastore:",noindex"`
	Indexes       []string   `datastore:",noindex"` // for XIAN
	Shard256      string     `datastore:",noindex"` // 分散処理用shard番号(256分割)
	Shard16       string     `datastore:",noindex"` // 分散処理用shard番号(16分割)
	CreatedAt     time.Time  `datastore:",noindex"`
	UpdatedAt     time.Time  `datastore:",noindex"`
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
	m := md5.New()
	io.WriteString(m, s)
	h := hex.EncodeToString(m.Sum(nil))
	return fmt.Sprintf("%s-%s", h[:3], s)
}

func call(c context.Context, start, count int) error {
	task := taskqueue.NewPOSTTask("/put-testbooks/create-data", url.Values{
		"start": {strconv.Itoa(start)},
		"count": {strconv.Itoa(count)},
	})

	_, err := taskqueue.Add(c, task, "create-test-data")

	return err
}

func putTestBooksTask(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil {
		msg := fmt.Sprintf("invalid start: %v", err.Error())
		putlogf(ctx, true, msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		msg := fmt.Sprintf("invalid count: %v", err.Error())
		putlogf(ctx, true, msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	end, err := strconv.Atoi(r.FormValue("end"))
	if err != nil {
		msg := fmt.Sprintf("invalid end: %v", err.Error())
		putlogf(ctx, true, msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	for i := start; i <= end; i += count {
		count := count
		if i+count > end {
			count = end - i + 1
		}
		if err = call(ctx, i, count); err != nil {
			putlogf(ctx, true, "error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	putlogf(ctx, false, "done")
	fmt.Fprintf(w, "done")
}

func putTestBooks(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	putlogf(ctx, false, "putTestBooks")

	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil {
		msg := fmt.Sprintf("invalid start: %v", err.Error())
		putlogf(ctx, true, msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		msg := fmt.Sprintf("invalid count: %v", err.Error())
		putlogf(ctx, true, msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	putlogf(ctx, false, "start=%d, count=%d", start, count)

	batchSize := 500

	keys := make([]*datastore.Key, 0, batchSize)
	books := make([]*Book2, 0, batchSize)

	titles := []string{"aPPle", "piNEApple", "banANA", "foobar", "hogefugapiyo"}

	for i := start; i < start+count; i++ {
		bookID := addHashPrefix(fmt.Sprintf("%d", i))
		key := datastore.NewKey(ctx, "Book2", bookID, 0, nil)

		book := &Book2{
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
			return
		}

		// shard番号設定
		if book.Shard256 == "" && book.Shard16 == "" {
			book.Shard256 = bookID[:2]
			book.Shard16 = bookID[:1]
		}

		now := time.Now()
		book.CreatedAt = now
		book.UpdatedAt = now

		keys = append(keys, key)
		books = append(books, book)

		if len(keys) == batchSize {
			if err := func() error {
				ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
				defer cancel()

				if _, err := datastore.PutMulti(ctx, keys, books); err != nil {
					return errors.WrapOr(err)
				}
				keys = keys[0:0]
				books = books[0:0]
				return nil
			}(); err != nil {
				putlogf(ctx, true, "error: %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if len(keys) > 0 {
		if err := func() error {
			ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
			defer cancel()
			if _, err := datastore.PutMulti(ctx, keys, books); err != nil {
				return errors.WrapOr(err)
			}
			return nil
		}(); err != nil {
			putlogf(ctx, true, "error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	putlogf(ctx, false, "done")
	fmt.Fprintf(w, "done")
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

func putlogf(c context.Context, iserr bool, msg string, args ...interface{}) {
	l := fmt.Sprintf(msg, args...)
	// log.Println(l)
	if iserr {
		aelog.Errorf(c, l)
	} else {
		aelog.Infof(c, l)
	}
}
