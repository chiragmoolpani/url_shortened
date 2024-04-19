package services

import (
	"assignment/url/models"
	"assignment/url/restapi/operations/url_shortened_api"
	"assignment/url/structs"
	"fmt"
	"io/ioutil"

	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

func GetShortenerURLService(params url_shortened_api.GetShortenerURLParams) middleware.Responder {
	fmt.Println("Inside GetShortenedURLService")

	// check file in db ---------

	const accessToken = "48849069ea9b2f5510e252d032f733409f1c9165"

	requestBody, err := json.Marshal(map[string]string{
		"long_url": params.URL,
	})

	if err != nil {
		details := "Error while Unmarshalling : " + err.Error()
		status := int32(url_shortened_api.GetShortenerURLInternalServerErrorCode)
		fmt.Println(details)
		return url_shortened_api.NewGetShortenerURLInternalServerError().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
	}

	apiURL := "https://api-ssl.bitly.com/v4/shorten"
	req, reqErr := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if reqErr != nil {
		details := "Error :: http.Request conversion got failed : " + reqErr.Error()
		status := int32(url_shortened_api.GetShortenerURLInternalServerErrorCode)
		fmt.Println(details)
		return url_shortened_api.NewGetShortenerURLInternalServerError().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	fmt.Println("Request sending to Bitly : ", req)

	client := &http.Client{}
	resp, err1 := client.Do(req)
	if err1 != nil {
		details := "Error while get request : " + err1.Error()
		status := int32(url_shortened_api.GetShortenerURLInternalServerErrorCode)
		fmt.Println(details)
		return url_shortened_api.NewGetShortenerURLInternalServerError().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
	}
	defer resp.Body.Close()
	fmt.Println("Response Code Received From Bitly : ", resp.StatusCode)

	if resp.StatusCode == url_shortened_api.GetShortenerURLUnauthorizedCode {
		details := "Unauthorised :: Error in token validation"
		status := int32(url_shortened_api.GetShortenerURLUnauthorizedCode)
		fmt.Println(details)
		return url_shortened_api.NewGetShortenerURLUnauthorized().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
	} else if resp.StatusCode != http.StatusOK {
		details := "Error :: failed due to bitly is down, Please try again after sometime..."
		status := int32(url_shortened_api.GetShortenerURLInternalServerErrorCode)
		fmt.Println(details)
		return url_shortened_api.NewGetShortenerURLInternalServerError().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
	}

	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		details := "Error:: http Response parsing failed : " + err2.Error()
		status := int32(url_shortened_api.GetShortenerURLInternalServerErrorCode)
		fmt.Println()
		return url_shortened_api.NewGetShortenerURLInternalServerError().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
	}

	fmt.Println("Response received from Bitly : ", string(bodyBytes))
	if string(bodyBytes) != "" {
		payload := models.GetShortenedURLResp{}
		result := structs.BitlyResponse{}
		json.Unmarshal([]byte(bodyBytes), &result)

		// storing in file ----------

		payload.ShortURL = result.Link
		return url_shortened_api.NewGetShortenerURLOK().WithPayload(&payload)
	} else {
		details := "Error :: Bitly Get request failed"
		status := int32(url_shortened_api.GetShortenerURLInternalServerErrorCode)
		fmt.Println(details)
		return url_shortened_api.NewGetShortenerURLInternalServerError().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
	}
}
