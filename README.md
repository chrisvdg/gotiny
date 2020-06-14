# Gotiny server

A Tiny URL (API) server in Go

## Prerequisites

Have [Go > 1.13](https://golang.org/dl/) installed

## Installation

```sh
go get github.com/chrisvdg/gotiny
cd %GOPATH/src/github.com/chrisvdg/gotiny
go build -mod vendor
```

## Docker

```sh
# Build images
docker build -t gotiny .

# Run basic http server
docker run -p 80:80 gotiny

# Run with custom arguments
docker run --rm gotiny --help
```

## Usage

A client library is available in the [gotiny_client repository.](https://github.com/chrisvdg/gotiny_client)  
The following examples use `CURL` instead.

Launch a gotiny server on the default port (":8080"),  
without any authentication and
pretty print the json output.

```sh
./gotiny -j


# Create a new entry using CURL (in another terminal window)
curl -d "id=google&url=google.com" -X POST http://localhost:8080/api/tiny
{
	"id": "google",
	"url": "http://google.com",
	"created": 1590330837
}
```

In a webbrowser surf to [http://localhost:8080/api/tiny/google](http://localhost:8080/api/tiny/google)  
where it should directed to the Google home page.

```sh
# Create a new entry with a generated ID
curl -d "url=chris.gent" -X POST http://localhost:8080/api/tiny
{
	"id": "b2p76",
	"url": "http://chris.gent",
	"created": 1590331741
}
```

Now in a webbrowser go to [http://localhost:8080/api/tiny/b2p76](http://localhost:8080/api/tiny/b2p76)  
where it should directed to my home page.


```sh
# List entries (you can also enter this URL into a browser)
curl http://localhost:8080/api/tiny
[
	{
		"id": "google",
		"url": "http://google.com",
		"created": 1590330837
	},
	{
		"id": "b2p76",
		"url": "http://chris.gent",
		"created": 1590331741
	}
]


# Get details of an entry (you can also enter this URL into a browser)
curl http://localhost:8080/api/tiny/google/expand
{
	"id": "google",
	"url": "http://google.com",
	"created": 1590330837
}
```
