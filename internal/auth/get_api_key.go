package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authorizationToken := headers.Get("Authorization")
	if authorizationToken == "" {
		return "", errors.New("no authorization header")
	}

	apiKey, ok := strings.CutPrefix(authorizationToken, "ApiKey ")
	if !ok {
		return "", errors.New("invalid authorization header")
	}

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return "", errors.New("empty key")
	}

	return apiKey, nil
}
