package blia

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	AuthorizationKey    = "Authorization"
	AuthorizationBearer = "Bearer"
)

var authorizationBEARER = strings.ToUpper(AuthorizationBearer)

func ReadBearerToken(r *http.Request) string {
	bearer := r.Header.Get(AuthorizationKey)
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == authorizationBEARER {
		return bearer[7:]
	}
	return ""
}

var _ AuthParser[any] = new(JWTAuthParser[any])

type JWTAuthParser[T any] struct {
	secert      string
	expire      time.Duration
	readSecerts []string
}

func NewJWTAuthParser[T any](secert string, expire time.Duration, oldSecerts ...string) *JWTAuthParser[T] {
	return &JWTAuthParser[T]{
		secert:      secert,
		expire:      expire,
		readSecerts: append([]string{secert}, oldSecerts...),
	}
}

func (a *JWTAuthParser[T]) ParseFromRequest(r *http.Request) (T, error) {
	return a.JWTParse(ReadBearerToken(r))
}

type jwtDataClaims struct {
	*jwt.StandardClaims
	Data string
}

func (a *JWTAuthParser[T]) JWTSign(user T) (string, error) {
	data, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	now := time.Now()
	token.Claims = &jwtDataClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: now.Add(a.expire).Unix(),
			IssuedAt:  now.Unix(),
		},
		Data: string(data),
	}
	return token.SignedString([]byte(a.secert))
}

var (
	ErrJwtTokenParseFailed = errors.New("jwt token prase failed")
	ErrJwtTokenInvalidAlg  = errors.New("jwt token prase with invalid alg")
)

func (a *JWTAuthParser[T]) jwtParse(token string, secert string) (T, error) {
	var user T
	tk, err := jwt.ParseWithClaims(
		token,
		&jwtDataClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%w :%v", ErrJwtTokenInvalidAlg, token.Header["alg"])
			}
			return []byte(secert), nil
		})
	if err != nil {
		return user, err
	}

	if claims, ok := tk.Claims.(*jwtDataClaims); ok && tk.Valid {
		err = json.Unmarshal([]byte(claims.Data), &user)
		if err != nil {
			return user, err
		}
		return user, nil
	}
	return user, ErrJwtTokenParseFailed
}

func (a *JWTAuthParser[T]) JWTParse(token string) (T, error) {
	var user T
	var err error
	for _, secert := range a.readSecerts {
		user, err = a.jwtParse(token, secert)
		if err == nil {
			return user, nil
		}
	}
	return user, err
}
