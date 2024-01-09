package data

import (
	"github.com/AnwarSaginbai/netflix-service/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Movie struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt time.Time          `json:"-" bson:"created_at"`
	Title     string             `json:"title,omitempty" bson:"title"`
	Year      int32              `json:"year,omitempty" bson:"year"`
	Runtime   Runtime            `json:"runtime,omitempty" bson:"runtime"`
	Genres    []string           `json:"genres,omitempty" bson:"genres"`
	Version   int32              `json:"version" bson:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
