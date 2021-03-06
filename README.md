Manga reader tracker API
========================

Getting started
---------------

Fetch the projet and install the dependencies:

```go
go get github.com/l-lin/mr-tracker-api
go get github.com/tools/godep
cd $GOPATH/github/l-lin/mr-tracker-api
godep go install
```

You need to set the following environment variables in the `startup.sh`:

|Variable name          |Description                            |Example                                                |
|-----------------------|---------------------------------------|-------------------------------------------------------|
|PORT                   |Server port                            |3000                                                   |
|GOOGLE_CLIENT_ID       |The client ID of your google API       |xxxx.apps.googleusercontent.com                        |
|GOOGLE_CLIENT_SECRET   |The client secret of your Google API   |ABCDEFGHIJKLMOPQRSTUVWXYZ                              |
|GOOGLE_REDIRECT_URL    |The redirect URL after being connected |http://localhost:3000/oauth2callback                   |
|DATABASE_URL           |The Database URL                       |postgres://postgres@localhost:5432/mr?sslmode=disable  |

Those variables are available in your [Google developer console](https://console.developers.google.com/project).

After the configuration, you just need to execute `startup.sh` and access to `http://localhost:3000/`

Notes
=====

This project use the [MangaFeeder project](https://github.com/cheeaun/mangafeeder) to fetch the news.
