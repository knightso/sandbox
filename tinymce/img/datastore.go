package img

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

const (
	thumbnailsLongestSide int = 100
)

type Image struct {
	ServingURL   string    `datastore:",noindex" json:"url"`
	FileName     string    `datastore:",noindex" json:"filename"`
	GCSPath      string    `datastore:",noindex" json:"-"` // gcs.go:Store()の戻り値absFilename
	Date         time.Time `json:"date"`
	ThumbnailURL string    `datastore:"-" json:"thumbnail"`
}

// Image構造体のThumbnailURLを設定します。
// size = 0で元の大きさで表示されるようだ。1600以上でも同様のようです。
func (i *Image) setThumbnailURL(size int, isCrop bool) error {
	u := i.ServingURL
	u += fmt.Sprintf("=s%d", size)
	if isCrop {
		u += "-c"
	}
	pu, err := url.Parse(u)
	if err != nil {
		return err
	}
	i.ThumbnailURL = pu.String()
	return nil
}

// アップロードされた画像のメタデータをDSに保存します。
func StoreImage(c appengine.Context, servingURL, fileName, gcsPath string) (err error) {
	e := &Image{
		ServingURL: servingURL,
		FileName:   fileName,
		GCSPath:    gcsPath,
		Date:       time.Now(),
	}

	key := datastore.NewIncompleteKey(c, "Image", nil)
	_, err = datastore.Put(c, key, e)
	if err != nil {
		c.Errorf("datastore.go:55: %s", err.Error())
		return err
	}

	return nil
}

// 画像のメタデータ一覧をDSから取得します。
// TODO: 表示する画像数を絞る必要がないなら、Cursor必要ないかも。
func GetImages(c appengine.Context, cursorStr string) ([]Image, string, error) {
	q := datastore.NewQuery("Image").Order("-Date")

	if len(cursorStr) != 0 {
		cursor, err := datastore.DecodeCursor(cursorStr)
		if err != nil {
			return []Image{}, "", err
		}

		q = q.Start(cursor)
	}

	images := []Image{}
	iter := q.Run(c)
	isNext := true
	for {
		var img Image
		_, err := iter.Next(&img)
		if err == datastore.Done {
			isNext = false
			break
		}
		if err != nil {
			c.Errorf("fetching next Person: %s", err.Error())
			break
		}

		err = img.setThumbnailURL(thumbnailsLongestSide, false)
		if err != nil {
			c.Errorf("%s", err.Error())
			break
		}
		images = append(images, img)
	}

	if isNext {
		next_cursor, err := iter.Cursor()
		if err != nil {
			c.Errorf("%s", err.Error())
			return []Image{}, "", err
		}
		return images, next_cursor.String(), nil
	} else {
		return images, "", nil
	}
}

type ResponseData struct {
	Files  []Image `json:"files"`
	Cursor string  `json:"cursor"`
}

// ファイルの一覧をクライアントに返すため、ImageのスライスをJSONシリアライズします。
func CreateJSONResponse(c appengine.Context, images []Image, cursor string) ([]byte, error) {
	fl := &ResponseData{Files: images, Cursor: cursor}
	res, err := json.Marshal(fl)
	if err != nil {
		c.Errorf("%s", err.Error())
		return []byte{}, err
	}

	return res, nil
}
