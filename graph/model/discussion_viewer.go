package model

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DiscussionViewer struct {
	DiscussionID string `json:"discussionID"`
	ViewerID     string `json:"participantID"`
}

func (d DiscussionViewer) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	av.S = aws.String(fmt.Sprintf("%s.%s", d.DiscussionID, d.ViewerID))
	return nil
}

func (d *DiscussionViewer) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if av.S == nil {
		return fmt.Errorf("Invalid nil input to unmarshal")
	}
	parts := strings.Split(*av.S, ".")
	if len(parts) != 2 {
		return fmt.Errorf("Incorrectly marshaled object: %s", *av.S)
	}
	d.DiscussionID = parts[0]
	d.ViewerID = parts[1]
	return nil
}
