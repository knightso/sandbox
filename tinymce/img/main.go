package img

import (
	"appengine"
	"fmt"
	"io/ioutil"
	"net/http"
)

func init() {
	http.HandleFunc("/upload", UploadHandler)
	http.HandleFunc("/gcs/files", FileListViewHandler)
}

func UploadHandler(rw http.ResponseWriter, req *http.Request) {
	var maxMemory int64 = 10 * 1024 * 1024
	var formKey string = "filename"
	var bucketName string = "knightso-test"

	c := appengine.NewContext(req)

	err := req.ParseMultipartForm(maxMemory)
	if err != nil {
		if err.Error() == "permission denied" {
			fmt.Fprint(rw, "アップロード可能な容量を超えています。\n")
		} else {
			fmt.Fprintf(rw, "%s", err.Error())
		}
		return
	}

	file, fileHeader, err := req.FormFile(formKey)
	if err != nil {
		c.Errorf("main.go:30: %s", err.Error())
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.Errorf("main.go:37: %s", err.Error())
		return
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if len(mimeType) == 0 {
		c.Errorf("main.go:43: couldn't get mime-type of file.")
		return
	}

	absFilename, err := Store(c, data, fileHeader.Filename, mimeType, bucketName)
	if err != nil {
		return
	}

	servingURL, err := GetServingURL(c, absFilename)
	if err != nil {
		return
	}

	err = StoreImage(c, servingURL.String(), fileHeader.Filename, absFilename)
	if err != nil {
		return
	}

	fmt.Fprint(rw, "Success.")
}

func FileListViewHandler(rw http.ResponseWriter, req *http.Request) {
	var formKey string = "cursor"

	c := appengine.NewContext(req)

	cursorStr := req.FormValue(formKey)
	images, next_cursor, err := GetImages(c, cursorStr)
	if err != nil {
		return
	}

	res, err := CreateJSONResponse(c, images, next_cursor)
	if err != nil {
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Write(res)
}
