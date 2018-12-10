package main

import (
	"encoding/json"
	"net/http"

	"github.com/knightso/xian"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var bookIndexesConfig = xian.MustValidateConfig(&xian.Config{
	IgnoreCase: true,
})

const (
	// BookQueryLabelTitleIndex is a label to search Books with Title index.
	BookQueryLabelTitleIndex = "ti"
	// BookQueryLabelTitlePrefix is a label to search Books with Title prefix
	BookQueryLabelTitlePrefix = "tp"
	// BookQueryLabelIsPublished is a label to search Books which are published.
	BookQueryLabelIsPublished = "p"
	// BookQueryLabelIsHobby is a label to search Books which are hoby.
	BookQueryLabelIsHobby = "h"
	// BookQueryLabelStatusIN is a label to search Books which are IN some statuses.
	BookQueryLabelStatusIN = "s"
	// BookQueryLabelPriceRange is a label to search Books which prices are in a range.
	BookQueryLabelPriceRange = "pr"
)

func saveBookIndexes(book *Book) {
	idxs := xian.NewIndexes(bookIndexesConfig)
	idxs.AddBiunigrams(BookQueryLabelTitleIndex, book.Title)
	idxs.AddPrefixes(BookQueryLabelTitlePrefix, book.Title)
	idxs.AddSomething(BookQueryLabelIsPublished, book.Status == BookStatusPublished)
	idxs.AddSomething(BookQueryLabelIsHobby, book.Category == "sports" || book.Category == "cooking")

	for i := 1; i < 1<<uint(len(BookStatuses))+1; i++ {
		if i&int(book.Status) != 0 {
			idxs.AddSomething(BookQueryLabelStatusIN, i)
		}
	}

	switch {
	case book.Price < 3000:
		idxs.Add(BookQueryLabelPriceRange, "p<3000")
	case book.Price < 5000:
		idxs.Add(BookQueryLabelPriceRange, "3000<=p<5000")
	case book.Price < 10000:
		idxs.Add(BookQueryLabelPriceRange, "5000<=p<10000")
	default:
		idxs.Add(BookQueryLabelPriceRange, "10000<=p")
	}

	book.Indexes = idxs.MustBuild()
}

func gaeXianNotEqual(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	filters := xian.NewFilters(bookIndexesConfig).AddSomething(BookQueryLabelIsPublished, false)

	q := datastore.NewQuery("Book")

	for _, f := range filters.MustBuild() {
		q = q.Filter("Indexes =", f)
	}

	var books []Book
	if _, err := q.GetAll(ctx, &books); err != nil {
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

func gaeXianIn(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	filters := xian.NewFilters(bookIndexesConfig).AddSomething(BookQueryLabelIsHobby, true)

	q := datastore.NewQuery("Book")

	for _, f := range filters.MustBuild() {
		q = q.Filter("Indexes =", f)
	}

	var books []Book
	if _, err := q.GetAll(ctx, &books); err != nil {
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

func gaeXianIn2(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	filters := xian.NewFilters(bookIndexesConfig).AddSomething(BookQueryLabelStatusIN, BookStatusUnpublished|BookStatusPublished)

	q := datastore.NewQuery("Book")

	for _, f := range filters.MustBuild() {
		q = q.Filter("Indexes =", f)
	}

	var books []Book
	if _, err := q.GetAll(ctx, &books); err != nil {
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

func gaeXianNumRange(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	filters := xian.NewFilters(bookIndexesConfig).Add(BookQueryLabelPriceRange, "5000<=p<10000")

	q := datastore.NewQuery("Book")

	for _, f := range filters.MustBuild() {
		q = q.Filter("Indexes =", f)
	}

	var books []Book
	if _, err := q.GetAll(ctx, &books); err != nil {
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

func gaeXianLike(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	title := r.FormValue("title")

	filters := xian.NewFilters(bookIndexesConfig).AddBiunigrams(BookQueryLabelTitleIndex, title)

	q := datastore.NewQuery("Book")

	for _, f := range filters.MustBuild() {
		q = q.Filter("Indexes =", f)
	}

	var books []Book
	if _, err := q.GetAll(ctx, &books); err != nil {
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

func gaeXianPrefix(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	title := r.FormValue("title")

	filters := xian.NewFilters(bookIndexesConfig).AddPrefix(BookQueryLabelTitlePrefix, title)

	q := datastore.NewQuery("Book")

	for _, f := range filters.MustBuild() {
		q = q.Filter("Indexes =", f)
	}

	var books []Book
	if _, err := q.GetAll(ctx, &books); err != nil {
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
