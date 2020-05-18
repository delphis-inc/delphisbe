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
	mediaType MediaType
	key       string
	mimeType  string
}

func (m *mediaDB) UploadMedia(ctx context.Context, filename string, media []byte) (MediaType, error) {
	fileInfo, err := m.getFileInfo(filename, media[:512])
	if err != nil {
		logrus.WithError(err).Error("failed to get bucket and key")
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

	return fileInfo.mediaType, nil
}

func (m *mediaDB) getFileInfo(filename string, partialMedia []byte) (fileInfo, error) {
	var mediaPrefix string
	var mediaType MediaType

	// Function return MIME Types
	mimeType := http.DetectContentType(partialMedia)

	s := strings.Split(mimeType, "/")
	if s[0] == "image" {
		if s[1] == "gif" {
			mediaPrefix = m.s3BucketConfig.GifKeyPrefix
			mediaType = GifMedia
		}
		mediaPrefix = m.s3BucketConfig.ImageKeyPrefix
		mediaType = ImageMedia
	} else if s[0] == "video" {
		mediaPrefix = m.s3BucketConfig.VideoKeyPrefix
		mediaType = VideoMedia
	}

	return fileInfo{
		mediaType: mediaType,
		key:       strings.Join([]string{m.s3BucketConfig.BaseKey, mediaPrefix, filename}, "/"),
		mimeType:  mimeType,
	}, nil

}
