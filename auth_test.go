package blia

import (
	"net/http"
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

func TestJWTAuthParserWithAuthenticator(t *testing.T) {
	type UserModel struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	jwt := NewJWTAuthParser[*UserModel]("secert", time.Hour)
	_ = NewAuthenticator[*UserModel](jwt) // TODO: add test

	var user = &UserModel{ID: 1, Name: "hello"}
	token, err := jwt.JWTSign(user)
	if err != nil {
		t.Error(err, token, user)
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add(AuthorizationKey, AuthorizationBearer+" "+token)

	user2, err := jwt.ParseFromRequest(r)
	if err != nil || user2.ID != user.ID || user2.Name != user.Name {
		t.Error(err, token, user2)
	}
}
