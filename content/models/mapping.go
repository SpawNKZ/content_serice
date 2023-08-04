package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (ct *Content) ToRepositoryModel() *ContentModel {
	return &ContentModel{
		ID:           primitive.NewObjectID(),
		Locale:       ct.Locale,
		Body:         ct.Body,
		Description:  ct.Description,
		Resources:    ct.Resources,
		SubjectId:    ct.SubjectId,
		MicrotopicId: ct.MicrotopicId,
		StatusId:     ct.StatusId,
		AuthorId:     ct.AuthorId,
		Difficulty:   ct.Difficulty,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (ct *Content) FromRepositoryModel(repoCt *ContentModel) {
	ct.ID = repoCt.ID.Hex()
	ct.Locale = repoCt.Locale
	ct.Body = repoCt.Body
	ct.Description = repoCt.Description
	ct.Resources = repoCt.Resources
	ct.SubjectId = repoCt.SubjectId
	ct.MicrotopicId = repoCt.MicrotopicId
	ct.StatusId = repoCt.StatusId
	ct.AuthorId = repoCt.AuthorId
	ct.Difficulty = repoCt.Difficulty
	ct.CreatedAt = repoCt.CreatedAt
	ct.UpdatedAt = repoCt.UpdatedAt
}

func (ct *CreateRequest) ToContent() *Content {
	return &Content{
		Locale:       ct.Locale,
		Body:         ct.Body,
		Description:  ct.Description,
		Resources:    ct.Resources,
		SubjectId:    ct.SubjectId,
		MicrotopicId: ct.MicrotopicId,
		StatusId:     ct.StatusId,
		AuthorId:     ct.AuthorId,
		Difficulty:   ct.Difficulty,
	}
}

func (ur *UpdateRequest) ToContent() map[string]interface{} {
	args := make(map[string]interface{})

	if ur.Description != "" {
		args["description"] = ur.Description
	}

	if ur.Difficulty != 0 {
		args["difficulty"] = ur.Difficulty
	}

	if ur.Body != "" {
		args["body"] = ur.Body
	}

	args["updated_at"] = time.Now()
	return args
}

func (ur *AssignAuthorRequest) ToContent() map[string]interface{} {
	args := make(map[string]interface{})

	if ur.AuthorId != "" {
		args["author_id"] = ur.AuthorId
	}

	args["updated_at"] = time.Now()
	return args
}

func (ur *ChangeStatusRequest) ToContent() map[string]interface{} {
	args := make(map[string]interface{})

	if ur.StatusId != "" {
		args["status_id"] = ur.StatusId
	}

	args["updated_at"] = time.Now()
	return args
}
