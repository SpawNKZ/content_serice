package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (ch *ContentHistory) ToRepositoryModel() *ContentHistoryModel {
	return &ContentHistoryModel{
		ID:            primitive.NewObjectID(),
		ContentId:     ch.ContentId,
		UserId:        ch.UserId,
		Action:        ch.Action,
		PreviousValue: ch.PreviousValue,
		NewValue:      ch.NewValue,
		CreatedAt:     time.Now(),
	}
}

func (ch *ContentHistory) FromRepositoryModel(repoCh *ContentHistoryModel) {
	ch.ID = repoCh.ID.Hex()
	ch.ContentId = repoCh.ContentId
	ch.UserId = repoCh.UserId
	ch.Action = repoCh.Action
	ch.PreviousValue = repoCh.PreviousValue
	ch.NewValue = repoCh.NewValue
	ch.CreatedAt = repoCh.CreatedAt
}
