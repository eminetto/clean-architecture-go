# Clean Architecture in Go

**This old repository doesn't represent what I'm using nowadays. In 2023, I am using and recommending what my colleagues and I have described in [this post](https://medium.com/inside-picpay/organizing-projects-and-defining-names-in-go-7f0eab45375d)**

[![Build Status](https://travis-ci.org/eminetto/clean-architecture-go.svg?branch=master)](https://travis-ci.org/eminetto/clean-architecture-go)

Clean Architecture sample

## Post

[https://medium.com/@eminetto/clean-architecture-using-golang-b63587aa5e3f](https://medium.com/@eminetto/clean-architecture-using-golang-b63587aa5e3f)


## Build

  make

## Run tests

  make test

## API requests 

### Add a bookmark

```
curl -X "POST" "http://localhost:8080/v1/bookmark" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d $'{
  "tags": [
    "git",
    "social"
  ],
  "name": "Github",
  "description": "Github site",
  "link": "http://github.com"
}'
```
### Search a bookmark

```
curl "http://localhost:8080/v1/bookmark?name=github" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json'
```

### Show all bookmarks

```
curl "http://localhost:8080/v1/bookmark" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json'
```

## CMD 

### Search for a bookmark

```
./bin/search github
```
