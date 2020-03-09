package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DiscussionParticipantKey struct {
	DiscussionID  string `json:"discussionID"`
	ParticipantID int    `json:"participantID"`
}

func (d DiscussionParticipantKey) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	av.S = aws.String(fmt.Sprintf("%s.%d", d.DiscussionID, d.ParticipantID))
	return nil
}

type DiscussionParticipantKeys struct {
	Keys []DiscussionParticipantKey `json:"keys" dynamodbav:",omitempty"`
}

func (d DiscussionParticipantKeys) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	ss := make([]*string, 0)
	for _, k := range d.Keys {
		ss = append(ss, aws.String(fmt.Sprintf("%s.%d", k.DiscussionID, k.ParticipantID)))
	}
	av.SS = ss
	return nil
}

func (d *DiscussionParticipantKeys) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if av.SS == nil {
		// I guess this is a string slice that made its way here?
		return fmt.Errorf("Invalid nil input to unmarshal")
	}
	d.Keys = make([]DiscussionParticipantKey, len(av.SS))
	for i, elem := range av.SS {
		key := DiscussionParticipantKey{}
		parts := strings.Split(*elem, ".")
		if len(parts) != 2 {
			return fmt.Errorf("Incorrectly marshaled object: %s", *av.S)
		}
		key.DiscussionID = parts[0]
		participantID, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("Could not convert participant ID to int: %s", parts[1])
		}
		key.ParticipantID = participantID
		d.Keys[i] = key
	}
	return nil
}
