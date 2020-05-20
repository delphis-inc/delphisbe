package backend

import (
	"context"
	"io/ioutil"
	"mime/multipart"

	"go.uber.org/multierr"

	"github.com/nedrocks/delphisbe/graph/model"

	"github.com/nedrocks/delphisbe/internal/util"

	"github.com/sirupsen/logrus"
)

func (d *delphisBackend) GetMediaRecord(ctx context.Context, mediaID string) (*model.Media, error) {
	mediaRecord, err := d.db.GetMediaRecordByID(ctx, mediaID)
	if err != nil {
		logrus.WithError(err).Error("failed to get MediaRecord from db")
		return nil, err
	}

	// Generate asset URL from media package
	assetLoc, err := d.mediadb.GetAssetLocation(ctx, mediaRecord.ID, *mediaRecord.MediaType)
	if err != nil {
		logrus.WithError(err).Error("failed to get asset location")
		return nil, err
	}

	mediaRecord.AssetLocation = &assetLoc
	return mediaRecord, nil
}

func (d *delphisBackend) UploadMedia(ctx context.Context, media multipart.File) (uuid string, mimeType string, err error) {
	uuid = util.UUIDv4()

	mediaBytes, err := ioutil.ReadAll(media)
	if err != nil {
		logrus.WithError(err).Error("failed to read all media bytes")
		return "", "", err
	}

	// Pass in size into s3
	mimeType, err = d.mediadb.UploadMedia(ctx, uuid, mediaBytes)
	if err != nil {
		logrus.WithError(err).Error("failed to upload media to s3")
		return "", "", err
	}

	mediaSize := getMediaSize(ctx, mimeType, mediaBytes)

	// Create record within Media table
	mediaObj := model.Media{
		ID:        uuid,
		MediaType: &mimeType,
		MediaSize: &mediaSize,
	}

	if err := d.writeMediaRecord(ctx, mediaObj); err != nil {
		logrus.WithError(err).Error("failed to put media record in db")
		return "", "", err
	}

	return uuid, mimeType, nil
}

func (d *delphisBackend) writeMediaRecord(ctx context.Context, mediaObj model.Media) error {
	// Begin tx
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return err
	}

	if err := d.db.PutMediaRecord(ctx, tx, mediaObj); err != nil {
		// Rollback on errors
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("failed to rollback tx")
			return multierr.Append(err, txErr)
		}
		return err
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return err
	}

	return nil
}

func getMediaSize(ctx context.Context, mimeType string, media []byte) model.MediaSize {
	fileSize := len(media)

	// TODO: Get dimensions of image
	//switch mediaType {
	//case mediadb.ImageMedia:
	//	//image, _, err := image2.DecodeConfig()
	//
	//}

	return model.MediaSize{
		Height: 0,
		Width:  0,
		SizeKb: float64(fileSize) / 1024.0,
	}

}
