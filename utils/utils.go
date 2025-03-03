package utils

import (
	"regexp"
	"strings"
)

func Profane(content string) string {
	contentslice := strings.Split(content, " ")

	for i, word := range contentslice {
		wordl := strings.ToLower(word)
		if wordl == "fuck" || wordl == "shit" || wordl == "fornax" {
			contentslice[i] = "****"
		}
	}

	return strings.Join(contentslice, " ")
}

func Mentions(content string) ([]string, error) {

	words := strings.Split(content, " ")
	uniqueusers := make(map[string]bool)
	users := []string{}

	pattern := regexp.MustCompile(`^@[a-zA-Z_][a-zA-Z0-9._%+-]{0,8}$`)

	for _, k := range words {

		if pattern.MatchString(k) {
			username := k[1:]

			if !uniqueusers[username] {
				if _, ok := uniqueusers[k[1:]]; !ok {
					uniqueusers[username] = true
					users = append(users, username)
				}
			}
		}
	}

	return users, nil
}
