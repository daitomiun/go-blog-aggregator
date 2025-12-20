
# Gator: Blog aggregator

A Blog aggregator from the CLI.

> surf, save, follow and unfollow your registered rss blog feeds!


## Requirements

You'll need to install [postgresql](https://docs.fedoraproject.org/en-US/quick-docs/postgresql/) and [go](https://go.dev/doc/install) to make it work


## installation

Via go install
```bash
go install https://github.com/daitomiun/go-blog-aggregator
```

or

```bash
cd go-blog-aggregator/ && go build .
```

## Create and configure the .gatorconfig.json

```json
{"db_url":"postgres://postgres:password@localhost:5432/gator?sslmode=disable","current_user_name":"mark"}
```

> Add and configure the password for the table and schema


## List of commands

```bash
go-blog-aggregator register juan
```
> Registers a new user to the database

```bash
go-blog-aggregator login juan
```
> logs in to a existing user

```bash
go-blog-aggregator reset
```
> Resets and deletes database users 


```bash
go-blog-aggregator users
```
> list of current and existing users

```bash
go-blog-aggregator addfeed "ava's blog" "https://avas.bearblog.dev/feed/?type=rss"
```
> Adds a new feed (needs name and url) 

```bash
go-blog-aggregator agg
```
> From the added feeds, it gators to every post and saves it to the database


```bash
go-blog-aggregator feeds
```
> list all feeds from all users

```bash
go-blog-aggregator follow "https://avas.bearblog.dev/feed/?type=rss"
```
> Follows a selected feed, saves it for current user


```bash
go-blog-aggregator following
```
> list of following feeds for current user

```bash
go-blog-aggregator unfollow
```
> Unfollows a selected feed, saves it for current user

```bash
go-blog-aggregator browse 5
```
> List posts from user feeds, if no argument is passed, it defaults to 2



