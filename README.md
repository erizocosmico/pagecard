# pagecard [![Build Status](https://travis-ci.org/mvader/pagecard.svg)](https://travis-ci.org/mvader/pagecard) [![Coverage Status](https://coveralls.io/repos/mvader/pagecard/badge.svg?branch=master&service=github)](https://coveralls.io/github/mvader/pagecard?branch=master) [![GoDoc](https://godoc.org/github.com/mvader/pagecard?status.svg)](http://godoc.org/github.com/mvader/pagecard)

Golang library to retrieve some useful data to represent a webpage in other applications from metadata exposed by the webpage.
To do so, this library reads all the `meta` tags on the page and builds a data structure containing the information from OpenGraph, Twitter and some other useful tags.

## Install

```
go get github.com/mvader/pagecard
```

## Usage

```go
func main() {
  info, err := pagecard.Get("http://www.imdb.com/title/tt0094721/")
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(info.OpenGraph.Title) // Beetlejuice (1988)
}
```

## Future additions

* [ ] Retrieve color exposed with `<meta name="theme-color">`
* [ ] Retrieve iTunes app with `<meta name="apple-itunes-app">`
* [ ] Retrieve icon sets with the Apple Icons and Microsoft tiles.
