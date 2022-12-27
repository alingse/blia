package blia

import (
	"net/http"
	"strings"
)

const (
	AuthorizationKey    = "Authorization"
	AuthorizationBEARER = "BEARER"
)

func ReadBearerToken(r *http.Request) string {
	bearer := r.Header.Get(AuthorizationKey)
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == AuthorizationBEARER {
		return bearer[7:]
	}
	return ""
}
