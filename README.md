# blog aggreGATOR

## Description

Gator is a command-line tool which allows you and other users to define a set of RSS feeds to track and continuously poll these feeds at a definable interval, then collect and store new articles. There is additional functionality to register and "login" a user (in a sense), then subscribe to feeds the tool is currently tracking. Finally, you can browse the lastest article titles and their descriptions.

## Dependencies

- Golang
- Postgres

## Setup

You can install this tool with the following command:
```go install github.com/lnix1/gator@latest```

You will then need to ensure you have a postgres server activated with a database created.

Then, using the connection string to the database, create a config file in your home directory under the name:
``` ~/.gatorconfig.json ```
with the following structure:
```
{"db_url":"postgres://{your_postgres_user}:{your_postgres_pass}@localhost:5432/{your_database_name}?sslmode=disable","current_user_name":""}
```

Then run the following command:

``` gator register <choose a user name> ```

## Commands to use:

- Register a user: ``` gator register <choose a user name> ```
- Login to a user: ``` gator login <user_name> ```
- Define a feed to track: ``` gator addfeed "Feed Name" "https://feedRssUrlExample.com" ```
- Subscribe to a few: ``` gator follow "https://feedRssUrlExample.com" ```
- List users: ``` gator users ```
- List Feeds: ``` gator feeds ```
- Browse Feeds current user follows: ``` gator browse <optional number of feeds to list> ```
