package api

import (
	"errors"
	"regexp"
)

func (cfg *ApiConfig) SteralizeChirp(body string) (string, error) {
	// check length
	if len(body) > 140 {
		return "", errors.New("chirp is too long")
	}

	// handle dirty words
	dirtyWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, word := range dirtyWords {
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(word))
		body = re.ReplaceAllString(body, "****")
	}
	return body, nil
}
