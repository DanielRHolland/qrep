package models



type TrackedItemType struct {
	Name   string      `json:"name" bson:"name"`
	Issues []IssueType `json:"issues" bson:"issues"`
	Id     string      `json:"id" bson:"_id"`
}

type IssueType struct {
	Description string `json:"description" bson:"description"`
	Resolved    bool   `json:"resolved" bson:"resolved"`
	Id          string `json:"id" bson:"_id"`
}


