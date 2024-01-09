package data

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Movies MovieModel
	Users  UserModel
	Token  TokenModel
}

func NewModels(movies, users, token *mongo.Collection) Models {
	return Models{
		Movies: MovieModel{Mongo: movies},
		Users:  UserModel{Mongo: users},
		Token:  TokenModel{Mongo: token},
	}
}
