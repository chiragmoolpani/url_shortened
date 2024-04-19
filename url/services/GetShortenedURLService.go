package services

import (
	"assignment/url/restapi/operations/url_shortened_api"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
)

func GetShortenerURLService(params url_shortened_api.GetShortenerURLParams) middleware.Responder {
	fmt.Println("API not implemented yet")
	return nil
}
