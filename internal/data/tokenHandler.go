package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TokenModel struct {
	Mongo *mongo.Collection
}

func (m TokenModel) New(userId primitive.ObjectID, ttl time.Duration, scope string) (*Token, error) {

	token, err := generateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}
	return token, m.Insert(token)
}

func (m TokenModel) Insert(token *Token) error {
	_, err := m.Mongo.InsertOne(context.TODO(), token)
	if err != nil {
		return err
	}
	return nil
}

func (m TokenModel) DeleteAllForUser(scope string, userID primitive.ObjectID) error {

	filter := bson.M{"scope": scope, "_id": userID}
	_, err := m.Mongo.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
