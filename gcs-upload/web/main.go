package main

import (
	"io"
	"math"
	"net/http"

	"github.com/knightso/sandbox/newgetall/util"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const chunkSize = 1024 * 1024 * 20

func main() {
	http.HandleFunc("/upload/small", uploadSmall)
	http.HandleFunc("/upload/medium", uploadMedium)
	http.HandleFunc("/upload/large", uploadLarge)
	http.HandleFunc("/upload/veryLarge", uploadVeryLarge)

	appengine.Main()
}

func uploadSmall(w http.ResponseWriter, r *http.Request) {
	uploadFile("small.txt", 1024*1024*11, w, r)
}

func uploadMedium(w http.ResponseWriter, r *http.Request) {
	uploadFile("medium.txt", 1024*1024*100, w, r)
}

func uploadLarge(w http.ResponseWriter, r *http.Request) {
	uploadFile("large.txt", 1024*1024*1024, w, r)
}

func uploadVeryLarge(w http.ResponseWriter, r *http.Request) {
	uploadFile("veryLarge.txt", 1024*1024*1024*5, w, r)
}

func uploadFile(fileName string, size int, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		uploadedSize := 0
		for {
			dataSize := int(math.Min(float64(size-uploadedSize), chunkSize))
			data := newData(dataSize)

			if _, err := pw.Write(data); err != nil {
				log.Errorf(c, "error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			uploadedSize += dataSize
			if uploadedSize >= size {
				break
			}
		}
	}()

	_, err := util.UploadStreamToGCS(c, pr, fileName, "text/plain", "")

	if err != nil {
		log.Errorf(c, "error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func newData(size int) []byte {
	b := make([]byte, size)
	for i := 0; i < size; i++ {
		b[i] = 'a'
	}
	return b
}
