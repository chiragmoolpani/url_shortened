This project is a simple URL shortener service that accepts a URL as an argument over a REST API and returns a shortened URL as a result. The purpose of this assignment is to create an API-only version of a URL shortener similar to services like Bitly.

Features
URL Shortening: Accepts a URL via a REST API and returns a shortened URL.
In-File Storage: Stores the original URL and its shortened version in txt file present in url/data/url.txt.
Redirection: When the same original URL is requested again, it returns the previously generated shortened URL.

Usage
  Shorten a URL
  To shorten a URL, send a POST request to the following endpoint:

  -------------------------------------------------------------------------------------------------------------------------------
  Curl Request : curl -X 'GET' 'http://<Serving_IP>:<Serving_Port>/url/getShortUrl/<LONG_URL>'   -H 'accept: application/json'
  -------------------------------------------------------------------------------------------------------------------------------


  Replace Serving_IP, Serving_Port and LONG_URL
    In LONG_URL replace "/" with "%2F", ":" with "%3A", "?" with "%3F" and "=" with "%3D"
    ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------
    Example : If you want to convert LONG_URL : https://www.youtube.com/watch?v=l3-gpiMGOOg
    Example request : curl -X 'GET' 'http://<Serving_IP>:<Serving_Port>/url/getShortUrl/https%3A%2F%2Fwww.youtube.com%2Fwatch%3Fv%3Dl3-gpiMGOOg'   -H 'accept: application/json'
    ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------

  Response

  {
    "shortened_url": "https://short.url/abc123"
  }


Setup
  Clone this repository in your GOPATH/assignment/.
  Install the required dependencies.
  Start the server.

  Running Locally
    bash

    Redirect to yoru GOPATH
    # Create the directory
    mkdir assignment

    # Redirect to that directory
    cd assignment/
    
    Copy code
    # Clone the repository
    git clone https://github.com/chiragmoolpani/url_shortened.git
    cd url-shortener

    # Go inside url/helper/GlobalData.go and replace your DataFile and AccessToken fields
    cd url/helper/
    vi GlobalData.go

    DataFile is located in url/data/url.txt replace the path in DataFile
    AccessToken : you can replace your access token here.

    # Run the server
    cd /url/cmd/url-server
    go run main.go

    
  Once the server is up and running it will provide you IP and PORT where it is serving the server

  In Above Curl request Replace that IP and PORT and Hit the curl request.

  You will get the shotURL. and also it will be stored in map format in /url/data/url.txt.

  thanks...
