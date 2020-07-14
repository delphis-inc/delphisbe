package util

import (
	"errors"
	"strings"

	"github.com/delphis-inc/delphisbe/graph/model"

	"github.com/sirupsen/logrus"
)

func ReturnParsedEntityID(entityID string) (*model.ParsedEntityID, error) {
	s := strings.Split(entityID, ":")
	if len(s) != 2 {
		err := errors.New("entity string has more than one colon")
		logrus.WithError(err)
		return nil, err
	}

	return &model.ParsedEntityID{
		ID:   s[1],
		Type: s[0],
	}, nil
}
