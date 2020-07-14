package csv_loaders

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/delphis-inc/delphisbe/internal/backend"
	"github.com/sirupsen/logrus"
)

func CreateFlair(db backend.DelphisBackend, reader *csv.Reader) {
	expectedHeader := [2]string{"user_id", "template_id"}
	reader.FieldsPerRecord = 2

	// Get and validate header row
	var header [2]string
	row, err := reader.Read()
	if err != nil {
		logrus.Fatal(err)
	}
	copy(header[:], row)
	if header != expectedHeader {
		logrus.WithFields(logrus.Fields{
			"header":         header,
			"expectedHeader": expectedHeader,
		}).Fatal("Invalid header row")
	}

	created := 0
	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Debugf("Creating flair: %v\n", strings.Join(data, ","))
		userID := data[0]
		templateID := data[1]
		flair, err := db.CreateFlair(nil, userID, templateID)
		if err != nil || flair == nil {
			logrus.WithError(err).Error("Failed to create flair")
		} else {
			created += 1
		}
	}
	logrus.Infof("Created %d flair", created)
}
