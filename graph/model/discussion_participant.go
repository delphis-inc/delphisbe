package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DiscussionParticipant struct {
	DiscussionID  string `json:"discussionID"`
	ParticipantID int    `json:"participantID"`
}

func (d DiscussionParticipant) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	av.S = aws.String(fmt.Sprintf("%s.%d", d.DiscussionID, d.ParticipantID))
	return nil
}

func (d *DiscussionParticipant) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if av.S == nil {
		return fmt.Errorf("Invalid nil input to unmarshal")
	}
	parts := strings.Split(*av.S, ".")
	if len(parts) != 2 {
		return fmt.Errorf("Incorrectly marshaled object: %s", *av.S)
	}
	d.DiscussionID = parts[0]
	participantID, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("Could not convert participant ID to int: %s", parts[1])
	}
	d.ParticipantID = participantID
	return nil
}
