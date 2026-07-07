package util

import (
	"fmt"
	"net/url"
	"strings"
)

func BuildURL(server string, path string) (string, error) {
	return BuildURLParams(server, path, nil)
}

func BuildURLParams(server string, path string, params map[string][]string) (string, error) {
	arr := strings.Split(server, "://")
	if len(arr) != 2 || (arr[0] != "http" && arr[0] != "https") {
		return "", fmt.Errorf("Invalid server %s", server)
	}

	u := url.URL{
		Scheme: arr[0],
		Host:   arr[1],
		Path:   path,
	}

	if params != nil {
		u.RawQuery = url.Values(params).Encode()
	}

	return u.String(), nil
}
