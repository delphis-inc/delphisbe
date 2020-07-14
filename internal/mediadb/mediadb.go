package mediadb

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/delphis-inc/delphisbe/internal/config"
)

type MediaType string

const (
	GifMedia   MediaType = "gif"
	ImageMedia MediaType = "image"
	VideoMedia MediaType = "video"
)

type MediaDB interface {
	GetAssetLocation(ctx context.Context, mediaID, mimeType string) (string, error)
	UploadMedia(ctx context.Context, filename string, media []byte) (string, error)
}

type mediaDB struct {
	uploader       *s3manager.Uploader
	downloader     *s3manager.Downloader
	s3BucketConfig config.S3BucketConfig
}

func NewMediaDB(config config.Config, awsSession *session.Session) MediaDB {
	mySession := awsSession
	return &mediaDB{
		s3BucketConfig: config.S3BucketConfig,
		uploader:       s3manager.NewUploader(mySession),
		downloader:     s3manager.NewDownloader(mySession),
	}
}
