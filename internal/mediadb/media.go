package mediadb

import (
	"bytes"
	"context"
	"errors"
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

func (m *mediaDB) GetAssetLocation(ctx context.Context, mediaID, mimeType string) (string, error) {
	fileInfo, err := m.getFileInfo(mediaID, mimeType)
	if err != nil {
		logrus.WithError(err).Error("failed to get file info for cloudfront")
		return "", err
	}

	return strings.Join([]string{m.s3BucketConfig.CloudFrontURL, fileInfo.key}, "/"), nil
}

func (m *mediaDB) UploadMedia(ctx context.Context, filename string, media []byte) (string, error) {
	if len(media) < 512 {
		err := errors.New("media file size is too small to detect")
		logrus.WithError(err)
		return "", err
	}

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
		return "", err
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
