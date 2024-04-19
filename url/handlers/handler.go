package handlers

import (
	"assignment/url/restapi/operations/url_shortened_api"
	"assignment/url/services"

	"github.com/go-openapi/runtime/middleware"
)

func GetShortURL(params url_shortened_api.GetShortenerURLParams) middleware.Responder {
	return services.GetShortenerURLService(params)
}
