package blia

import (
	"testing"
	"time"
)

func TestJWTAuthParser(t *testing.T) {
	type UserModel struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var user = &UserModel{ID: 1, Name: "hello"}

	jwt := NewJWTAuthParser[*UserModel]("secert", time.Hour)
	token, err := jwt.JWTSign(user)
	if err != nil {
		t.Error(err, token, user)
	}

	user2, err := jwt.JWTParse(token)
	if err != nil || user2.ID != user.ID || user2.Name != user.Name {
		t.Error(err, token, user2)
	}
}
