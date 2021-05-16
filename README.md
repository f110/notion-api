notion-api
---

API client for Notion written by Go.

Currently under active development. All APIs will be changed possibly.

# Usage

First, import from your code.

```go
import "go.f110.dev/notion-api"
```

and you also need `golang.org/x/oauth2` module for *http.Client.

```go
import "golang.org/x/oauth2"
```

After import, you can create the client with *http.Client.

```go
ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
tc := oauth2.NewClient(context.Background(), ts)
client, err := notion.New(tc, notion.BaseURL)
```

And example code exists under [example directory](./example)

# Supported methods

* [x] [Retrieve a database](https://developers.notion.com/reference/get-database)
* [x] [Query a database](https://developers.notion.com/reference/post-database-query)
* [x] [List databases](https://developers.notion.com/reference/get-databases)
* [x] [Retrieve a user](https://developers.notion.com/reference/get-user)
* [x] [List all users](https://developers.notion.com/reference/get-users)
* [x] [Retrieve a page](https://developers.notion.com/reference/get-page)
* [x] [Create a page](https://developers.notion.com/reference/post-page)
* [ ] [Update a page properties](https://developers.notion.com/reference/patch-page)
* [x] [Retrieve block children](https://developers.notion.com/reference/get-block-children)
* [ ] [Append block children](https://developers.notion.com/reference/patch-block-children)

# Author

Fumihiro Ito