package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MovieModel struct {
	Mongo *mongo.Collection
}

func (m MovieModel) Insert(movie *Movie) error {
	res, err := m.Mongo.InsertOne(context.TODO(), movie)
	if err != nil {
		return err
	}
	log.Printf("INSERTED ID: %v", res.InsertedID)
	return nil
}

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) {

	filter := bson.M{}

	if title != "" {
		filter["title"] = title
	}

	if len(genres) > 0 {
		filter["genres"] = bson.M{"$in": genres}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	findOptions := options.Find().SetLimit(int64(filters.limit())).SetSkip(int64(filters.offset()))

	cursor, err := m.Mongo.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	movies := []*Movie{}

	for cursor.Next(ctx) {

		var movie Movie

		if err := cursor.Decode(&movie); err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return movies, nil

}

func (m MovieModel) Get(id string) (*Movie, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var movie Movie
	err = m.Mongo.FindOne(context.TODO(), filter).Decode(&movie)
	if err != nil {

		return nil, err
	}

	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {

	filter := bson.M{"_id": movie.ID}

	update := bson.M{
		"$set": bson.M{
			"title":   movie.Title,
			"year":    movie.Year,
			"runtime": movie.Runtime,
			"genres":  movie.Genres,
		},
		"$inc": bson.M{
			"version": 1,
		},
	}

	_, err := m.Mongo.UpdateOne(context.TODO(), filter, update)
	if err != nil {

		return err
	}

	movie.Version++

	return nil
}

func (m MovieModel) Delete(id string) error {
	Obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": Obj}
	_, err = m.Mongo.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
