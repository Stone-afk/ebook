package ioc

import "net/http"

func InitDefaultHttpClient() *http.Client {
	return http.DefaultClient
}
