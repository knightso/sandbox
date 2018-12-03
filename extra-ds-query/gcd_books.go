package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"

	"cloud.google.com/go/datastore"
)

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

// GCDBookStatus describles status of Book
type GCDBookStatus int

// GCDBookStatus contants
const (
	GCDBookStatusUnpublished GCDBookStatus = 1 << iota
	GCDBookStatusPublished
	GCDBookStatusDiscontinued
)

// GCDBook is sample model.
type GCDBook struct {
	Title         string
	TitleIndex    []string
	TitlePrefix   []string
	Price         int
	PriceRange    string
	Category      string
	Status        GCDBookStatus
	StatusORIndex []int
	IsPublished   bool
	IsHobby       bool
}

// GCDBookStatuses is the list of all GCDBookStatuses
var GCDBookStatuses = []GCDBookStatus{
	GCDBookStatusUnpublished,
	GCDBookStatusPublished,
	GCDBookStatusDiscontinued,
}

// GCDBookCategories is the list of all Book categories
var GCDBookCategories = []string{
	"sports",
	"cooking",
	"education",
	"cartoons",
}

func putTestGCDBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	titles := []string{"aPPle", "piNEApple", "banANA", "foobar", "hogefugapiyo"}

	keys := make([]*datastore.Key, 0, 100)
	books := make([]*GCDBook, 0, 100)

	for i := 0; i < 100; i++ {
		bookID := fmt.Sprintf("book%04d", i)
		key := datastore.NameKey("GCDBook", bookID, nil)

		book := &GCDBook{
			Title:    titles[i%len(titles)],
			Price:    i * 100,
			Status:   GCDBookStatuses[i%len(GCDBookStatuses)],
			Category: GCDBookCategories[i%len(GCDBookCategories)],
		}

		// Book保存時に派生プロパティを補完
		book.TitleIndex = gcdBiunigrams(book.Title)
		book.TitlePrefix = prefixes(book.Title)
		book.IsPublished = book.Status == GCDBookStatusPublished
		book.IsHobby = book.Category == "sports" || book.Category == "cooking"

		for j := 1; j < 1<<uint(len(GCDBookStatuses))+1; j++ {
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

		keys = append(keys, key)
		books = append(books, book)
	}

	if _, err := client.PutMulti(ctx, keys, books); err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func gcdNotEqual(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	var books []GCDBook
	q := datastore.NewQuery("GCDBook").Filter("IsPublished =", false)
	_, err = client.GetAll(ctx, q, &books)
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

func gcdIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	var books []GCDBook
	q := datastore.NewQuery("GCDBook").Filter("IsHobby =", true)
	_, err = client.GetAll(ctx, q, &books)
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

func gcdIn2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	var books []GCDBook
	q := datastore.NewQuery("GCDBook").Filter("StatusORIndex =", int(GCDBookStatusUnpublished|GCDBookStatusPublished))
	_, err = client.GetAll(ctx, q, &books)
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

func gcdNumRange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	var books []GCDBook
	q := datastore.NewQuery("GCDBook").Filter("PriceRange =", "5000<=p<10000")
	_, err = client.GetAll(ctx, q, &books)
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

func gcdLike(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	title := r.FormValue("title")

	q := datastore.NewQuery("GCDBook")

	if runeLen := utf8.RuneCountInString(title); runeLen == 1 {
		// パラメータが1文字の場合はunigramで検索
		q = q.Filter("TitleIndex =", title)
	} else if runeLen > 1 {
		// パラメータが2文字以上の場合はbigramで検索
		for _, gram := range gcdBigrams(title) {
			q = q.Filter("TitleIndex =", gram)
		}
	}

	var books []GCDBook
	_, err = client.GetAll(ctx, q, &books)
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

func gcdPrefix(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	title := r.FormValue("title")
	var books []Book
	q := datastore.NewQuery("GCDBook").Filter("TitlePrefix =", strings.ToLower(title))
	_, err = client.GetAll(ctx, q, &books)
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

func gcdBigrams(s string) []string {
	tokens := make([]string, 0, 32)

	for bigram := range toBigrams(strings.ToLower(s)) {
		tokens = append(tokens, fmt.Sprintf("%c%c", bigram.a, bigram.b))
	}

	return tokens

}

func gcdBiunigrams(s string) []string {
	tokens := make([]string, 0, 32)

	for bigram := range toBigrams(strings.ToLower(s)) {
		tokens = append(tokens, fmt.Sprintf("%c%c", bigram.a, bigram.b))
	}
	for unigram := range toUnigrams(strings.ToLower(s)) {
		tokens = append(tokens, fmt.Sprintf("%c", unigram))
	}

	return tokens

}

func gcdPrefixes(s string) []string {
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
