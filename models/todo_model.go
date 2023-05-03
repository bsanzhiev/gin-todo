package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID      int    `json:"id" bson:"id"`
	Text    string `json:"text" validate:"required"`
	Checked bool   `json:"checked" validate:"required"`
}
