package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"github.com/AnwarSaginbai/netflix-service/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string             `json:"token" bson:"token"`
	Hash      []byte             `json:"-" bson:"hash"`
	UserID    primitive.ObjectID `json:"-" bson:"user_id"`
	Expiry    time.Time          `json:"expiry" bson:"expiry"`
	Scope     string             `json:"-" bson:"scope"`
}

func generateToken(userID primitive.ObjectID, ttl time.Duration, scope string) (*Token, error) {

	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}
