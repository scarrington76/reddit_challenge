# reddit_challenge

## Environment setup

You need to have [Go](https://go.dev/),
installed on your computer.

Verify the required tools by running the following commands:

```sh
go version
```

**A Reddit API token is necessary to run the application and must
be loaded into the env variables under the var name ACCESS_TOKEN**

By default, the r/AskReddit thread will be tracked

## Start in development mode

Navigate to the `main` folder and start the back end:

```sh
cd main
go run main.go
```
The back end will serve on http://localhost:8080.

## Unit Testing

Unit tests can also be ran in the application from /main

**A Reddit API token is necessary to run unit tests. The token
string must be populated into the tok variable in main/server/test/test.go**

```sh
cd main
go test .///
```

## Postman

A postman collection is provided in /postman to hit the endpoints to obtain
user/post statistics. Simply import the collection into postman as a collection.