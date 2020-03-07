package model

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DiscussionViewerKey struct {
	DiscussionID string `json:"discussionID"`
	ViewerID     string `json:"participantID"`
}

type DiscussionViewerKeys struct {
	Keys []DiscussionViewerKey `json:"keys"`
}

func (d DiscussionViewerKeys) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	ss := make([]*string, 0)
	for _, k := range d.Keys {
		ss = append(ss, aws.String(fmt.Sprintf("%s.%s", k.DiscussionID, k.ViewerID)))
	}
	av.SS = ss
	return nil
}

func (d *DiscussionViewerKeys) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if av.SS == nil {
		return fmt.Errorf("Invalid nil input to unmarshal")
	}
	d.Keys = make([]DiscussionViewerKey, len(av.SS))
	for i, elem := range av.SS {
		key := DiscussionViewerKey{}
		parts := strings.Split(*elem, ".")
		if len(parts) != 2 {
			return fmt.Errorf("Incorrectly marshaled object: %s", *av.S)
		}
		key.DiscussionID = parts[0]
		key.ViewerID = parts[1]
		d.Keys[i] = key
	}
	return nil
}
