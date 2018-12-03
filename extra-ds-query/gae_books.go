package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// BookStatus describles status of Book
type BookStatus int

// BookStatus contants
const (
	BookStatusUnpublished BookStatus = 1 << iota
	BookStatusPublished
	BookStatusDiscontinued
)

// Book is sample model.
type Book struct {
	Title         string
	TitleIndex    []string
	Price         int
	PriceRange    string
	Category      string
	Status        BookStatus
	StatusORIndex []int
	IsPublished   bool
	IsHobby       bool
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

func putTestBooks(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	titles := []string{"aPPle", "piNEApple", "banANA", "foobar", "hogefugapiyo"}

	keys := make([]*datastore.Key, 0, 100)
	books := make([]*Book, 0, 100)

	for i := 0; i < 100; i++ {
		bookID := fmt.Sprintf("book%04d", i)
		key := datastore.NewKey(ctx, "Book", bookID, 0, nil)

		book := &Book{
			Title:    titles[i%len(titles)],
			Price:    i * 100,
			Status:   BookStatuses[i%len(BookStatuses)],
			Category: BookCategories[i%len(BookCategories)],
		}

		// Book保存時に派生プロパティを補完
		book.TitleIndex = biunigrams(book.Title)
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

		keys = append(keys, key)
		books = append(books, book)
	}

	if _, err := datastore.PutMulti(ctx, keys, books); err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
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
