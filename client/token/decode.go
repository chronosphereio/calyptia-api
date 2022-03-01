package token

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

var errInvalidToken = errors.New("invalid project token")

const (
	tokenPartsSeparator = "."
	// default string parts separated by tokenPartsSeparator to compose a valid token.
	tokenParts = 2
)

type projectTokenPayload struct {
	ProjectID string // no json tag
}

// Decode decodes a project token without verifying its signature
// and getting its inner project ID.
func Decode(token []byte) (string, error) {
	parts := bytes.Split(token, []byte(tokenPartsSeparator))
	if len(parts) != tokenParts {
		return "", errInvalidToken
	}

	encodedPayload := parts[0]

	payload := make([]byte, base64.RawURLEncoding.DecodedLen(len(encodedPayload)))
	n, err := base64.RawURLEncoding.Decode(payload, encodedPayload)
	if err != nil {
		return "", errInvalidToken
	}

	payload = payload[:n]

	var out projectTokenPayload
	err = json.Unmarshal(payload, &out)
	if err != nil {
		return "", fmt.Errorf("could not json parse project token payload: %w", err)
	}

	return out.ProjectID, nil
}
