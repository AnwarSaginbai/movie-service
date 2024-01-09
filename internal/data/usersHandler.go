package data

import (
	"context"
	"crypto/sha256"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type UserModel struct {
	Mongo *mongo.Collection
}

func (m UserModel) Insert(user *User) error {
	// Set up the MongoDB insert document.

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Выполните вставку в коллекцию MongoDB.
	_, err := m.Mongo.InsertOne(ctx, user)
	if err != nil {
		// Если возникает ошибка вставки, проверьте наличие дубликата по email.
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	filter := bson.M{"email": email}
	var user User
	err := m.Mongo.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	filter := bson.M{"_id": user.ID}

	update := bson.M{
		"$set": bson.M{
			"name":      user.Name,
			"email":     user.Email,
			"password":  user.Password,
			"activated": user.Activated,
		},
		"$inc": bson.M{
			"version": 1,
		},
	}

	_, err := m.Mongo.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	user.Version++

	return nil
}

func (m TokenModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {

	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Set up the MongoDB pipeline.
	pipeline := mongo.Pipeline{
		// Match stage
		{{"$match", bson.D{
			{"hash", tokenHash[:]},
			{"scope", tokenScope},
			{"expiry", bson.D{{"$gt", time.Now()}}},
		}}},
		// Lookup stage to join with the users collection
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		// Unwind stage to destructure the array produced by the lookup
		{{"$unwind", "$user"}},
		// Project stage to shape the output
		{{"$project", bson.D{
			{"_id", "$user._id"},
			{"created_at", "$user.createdAt"},
			{"name", "$user.name"},
			{"email", "$user.email"},
			{"password", "$user.password"},
			{"activated", "$user.activated"},
			{"version", "$user.version"},
		}}},
	}

	var users []User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cursor, err := m.Mongo.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("MongoDB Aggregation Error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		log.Println("Cursor All Error:", err)
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrRecordNotFound
	}

	return &users[0], nil

}
