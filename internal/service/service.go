package service

import "errors"

var slug_map = map[string]string{
	"abc": "https://google.com",
}

func GetRedirectLink(slug string) (string, error) {
	link := slug_map[slug]
	if link == "" {
		return "", errors.New("slug not found")
	}
	return link, nil
}
