package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (pt *Post) ToRepositoryModel() *PostModel {
	return &PostModel{
		ID:          primitive.NewObjectID(),
		Category:    pt.Category,
		Description: pt.Description,
		Resources:   pt.Resources,
		ContentId:   pt.ContentId,
	}
}

func (pt *Post) FromRepositoryModel(repoPt *PostModel) {
	pt.ID = repoPt.ID.Hex()
	pt.Description = repoPt.Description
	pt.Resources = repoPt.Resources
	pt.Category = repoPt.Category
	pt.ContentId = repoPt.ContentId
}

func (pt *CreateRequest) ToPost() *Post {
	return &Post{
		Description: pt.Description,
		Resources:   pt.Resources,
		Category:    pt.Category,
		ContentId:   pt.ContentId,
	}
}

func (ur *UpdateRequest) ToPost() map[string]interface{} {
	args := make(map[string]interface{})

	if ur.Description != "" {
		args["description"] = ur.Description
	}

	if len(ur.Resources) > 0 {
		args["resources"] = ur.Resources
	}

	return args
}
