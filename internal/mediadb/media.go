package mediadb

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type fileInfo struct {
	key      string
	mimeType string
}

func (m *mediaDB) UploadMedia(ctx context.Context, filename string, media []byte) (string, error) {
	mimeType := http.DetectContentType(media[:512])

	fileInfo, err := m.getFileInfo(filename, mimeType)
	if err != nil {
		logrus.WithError(err).Error("failed to get file info for s3")
		return "", err
	}
	logrus.Debugf("FileInfo: %v\n", fileInfo)

	if _, err := m.uploader.Upload(&s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Body:        bytes.NewBuffer(media),
		Bucket:      aws.String(m.s3BucketConfig.MediaBucket),
		Key:         aws.String(fileInfo.key),
		ContentType: aws.String(fileInfo.mimeType),
	}); err != nil {
		logrus.WithError(err).Error("failed to upload image to s3")
	}

	return fileInfo.mimeType, nil
}

func (m *mediaDB) getFileInfo(filename string, mimeType string) (fileInfo, error) {
	var mediaPrefix string

	// Place gifs, images, and videos in different buckets
	s := strings.Split(mimeType, "/")
	if s[0] == "image" {
		if s[1] == "gif" {
			mediaPrefix = m.s3BucketConfig.GifKeyPrefix
		}
		mediaPrefix = m.s3BucketConfig.ImageKeyPrefix
	} else if s[0] == "video" {
		mediaPrefix = m.s3BucketConfig.VideoKeyPrefix
	}

	return fileInfo{
		key:      strings.Join([]string{m.s3BucketConfig.BaseKey, mediaPrefix, filename}, "/"),
		mimeType: mimeType,
	}, nil

}
