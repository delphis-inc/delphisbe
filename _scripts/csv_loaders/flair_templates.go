package csv_loaders

import (
    "io"
    "encoding/csv"

    "github.com/nedrocks/delphisbe/internal/backend"
    "github.com/sirupsen/logrus"
)

const (
    FieldsPerRecord = 3
)

func create_from_csv(db backend.DelphisBackend, reader *csv.Reader) {
    expectedHeader := [3]string{"displayName", "imageURL", "source"}
    reader.FieldsPerRecord = FieldsPerRecord

    // Get and validate header row
	var header [3]string
    row, _ := reader.Read()
    copy(header[:], row)
    if header != expectedHeader {
		logrus.WithFields(logrus.Fields{
		  "header": header,
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

		logrus.Debugf("Creating flair template: %s\n", data)
		displayName := data[0]
		imageURL    := data[1]
		source      := data[2]
		flairTemplate, err := db.CreateFlairTemplate(
            nil, &displayName, &imageURL, source)
		if err != nil || flairTemplate == nil {
			logrus.WithError(err).Error("Failed to create flair template")
		} else {
			created += 1
		}
	}
	logrus.Infof("Created %d flair templates", created)
}
