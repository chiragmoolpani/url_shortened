package services

import (
	"assignment/url/helper"
	"assignment/url/models"
	"assignment/url/restapi/operations/url_shortened_api"
	"assignment/url/structs"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

func GetShortenerURLService(params url_shortened_api.GetShortenerURLParams) middleware.Responder {
	fmt.Println("GetShortenedURLService API Invoked")

	payload := models.GetShortenedURLResp{}

	// Searching in file
	mappings, fetchErr := fetchFromFile(helper.DataFile)
	if fetchErr != nil {
		details := "Error while fetching data from file : " + fetchErr.Error()
		status := int32(url_shortened_api.GetShortenerURLInternalServerErrorCode)
		fmt.Println(details)
		return url_shortened_api.NewGetShortenerURLInternalServerError().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
	}
	shortURL, isThere := mappings[params.URL]

	if !isThere {
		fmt.Println("------------------------- Generating URL --------------------------")
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
		req.Header.Set("Authorization", "Bearer "+helper.AccessToken)
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
			details := "Unauthorised :: Error in validating token"
			status := int32(url_shortened_api.GetShortenerURLUnauthorizedCode)
			fmt.Println(details)
			return url_shortened_api.NewGetShortenerURLUnauthorized().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
		} else if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				details := "Error :: http Response parsing failed : " + err2.Error()
				status := int32(url_shortened_api.GetShortenerURLInternalServerErrorCode)
				fmt.Println()
				return url_shortened_api.NewGetShortenerURLInternalServerError().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
			}

			fmt.Println("Response received from Bitly : ", string(bodyBytes))
			if string(bodyBytes) != "" {
				result := structs.BitlyResponse{}
				json.Unmarshal([]byte(bodyBytes), &result)

				// storing in file as map[string]string
				mappings[params.URL] = result.Link
				err = saveMapsInFile(helper.DataFile, mappings)
				if err != nil {
					details := "Error while saving mappings in file : " + err.Error()
					status := int32(url_shortened_api.GetShortenerURLInternalServerErrorCode)
					fmt.Println(details)
					return url_shortened_api.NewGetShortenerURLInternalServerError().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
				}
				fmt.Println("----- Short Link : ", result.Link)
				payload.ShortURL = result.Link
				return url_shortened_api.NewGetShortenerURLOK().WithPayload(&payload)
			} else {
				details := "Error :: Bitly Get request failed"
				status := int32(url_shortened_api.GetShortenerURLInternalServerErrorCode)
				fmt.Println(details)
				return url_shortened_api.NewGetShortenerURLInternalServerError().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
			}
		} else {
			statusCode := strconv.Itoa(resp.StatusCode)
			details := "Error :: failed due to bitly is down, Please try again after sometime... status :: " + statusCode
			status := int32(url_shortened_api.GetShortenerURLInternalServerErrorCode)
			fmt.Println(details)
			return url_shortened_api.NewGetShortenerURLInternalServerError().WithPayload(&models.ErrorInformation{Detail: &details, Status: &status})
		}

	} else {
		fmt.Println("------------------------- Returning from file --------------------------")
		payload.ShortURL = shortURL
		return url_shortened_api.NewGetShortenerURLOK().WithPayload(&payload)
	}
}

// fetchFromFile() : function will fetch URLs from a file
func fetchFromFile(filename string) (map[string]string, error) {
	mappings := make(map[string]string)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 2 {
			mappings[parts[0]] = parts[1]
		}
	}
	return mappings, nil
}

// saveMapsInFile() : function will store URLs insto a file
func saveMapsInFile(filename string, mappings map[string]string) error {
	var lines []string
	for longURL, shortURL := range mappings {
		lines = append(lines, fmt.Sprintf("%s,%s", longURL, shortURL))
	}
	data := []byte(strings.Join(lines, "\n"))

	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
