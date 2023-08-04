package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (cs *ContentStatus) ToRepositoryModel() *ContentStatusModel {
	return &ContentStatusModel{
		ID:          primitive.NewObjectID(),
		Name:        cs.Name,
		IsRemovable: cs.IsRemovable,
	}
}

func (cs *ContentStatus) FromRepositoryModel(repoCs *ContentStatusModel) {
	cs.ID = repoCs.ID.Hex()
	cs.Name = repoCs.Name
	cs.IsRemovable = repoCs.IsRemovable
}

func (cr *CreateRequest) ToContentStatus() *ContentStatus {
	return &ContentStatus{
		Name:        cr.Name,
		IsRemovable: cr.IsRemovable,
	}
}

func (ur *UpdateRequest) ToContentStatus() map[string]interface{} {
	args := make(map[string]interface{})

	if ur.Name != "" {
		args["name"] = ur.Name
	}

	return args
}
