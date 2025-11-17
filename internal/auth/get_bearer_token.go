package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authorizationToken := headers.Get("Authorization")
	if authorizationToken == "" {
		return "", errors.New("no authorization header")
	}

	jwtToken, ok := strings.CutPrefix(authorizationToken, "Bearer ")
	if !ok {
		return "", errors.New("invalid authorization header")
	}

	jwtToken = strings.TrimSpace(jwtToken)
	if jwtToken == "" {
		return "", errors.New("empty token")
	}

	return jwtToken, nil
}
