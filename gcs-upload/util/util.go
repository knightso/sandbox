package util

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/knightso/base/errors"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
)

// UploadStreamToGCS は io.Reader を入力として、GCSにファイルをアップロードする
func UploadStreamToGCS(c context.Context, r io.Reader, fileName, mimeType, bucketName string) (absFileName string, err error) {
	if bucketName == "" {
		var err error
		if bucketName, err = file.DefaultBucketName(c); err != nil {
			log.Errorf(c, "failed to get default GCS bucket name: %v", err)
			return "", errors.WrapOr(err)
		}
	}

	client, err := storage.NewClient(c)
	if err != nil {
		log.Errorf(c, "failed to create storage client: %v", err)
		return "", errors.WrapOr(err)
	}
	defer client.Close()

	wc := client.Bucket(bucketName).Object(fileName).NewWriter(c)
	wc.ChunkSize = 1024 * 1024 * 5
	wc.ContentType = mimeType

	buf := make([]byte, 1024*1024*5)
	if _, err := io.CopyBuffer(wc, r, buf); err != nil {
		log.Errorf(c, "upload file: unable to write data to bucket %q, file %q: %v", bucketName, fileName, err)
		return "", errors.WrapOr(err)
	}
	if err := wc.Close(); err != nil {
		log.Errorf(c, "upload file: unable to close bucket %q, file %q: %v", bucketName, fileName, err)
		return "", errors.WrapOr(err)
	}

	return getAbsFilename(bucketName, fileName), nil
}

func getAbsFilename(bucketName string, fileName string) string {
	return fmt.Sprintf("/gs/%s/%s", bucketName, fileName)
}
