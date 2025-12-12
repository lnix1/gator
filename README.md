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

You will then want to use the commands below to:

1. login to your user
2. add feeds
3. start the aggregator

If you add more than one user and want to follow a feed that was added by another user, use the "follow" command to subscribe.

## Commands to use:

- Register a user: ``` gator register <choose a user name> ```
- Login to a user: ``` gator login <user_name> ```
- Define a feed to track: ``` gator addfeed "Feed Name" "https://feedRssUrlExample.com" ```
- Subscribe to a feed: ``` gator follow "https://feedRssUrlExample.com" ```
- Unsubscribe to a feed: ``` gator unfollow "https://feedRssUrlExample.com" ```
- List users: ``` gator users ```
- List feeds: ``` gator feeds ```
- List feeds user is following: ``` gator following ```
- Browse feeds current user follows: ``` gator browse <optional number of feeds to list> ```
- Start continuously running aggregator: ``` gator agg 1d ```
	- This will start a continuously running aggregate which will poll a new feed every 1d. You can adjust the "1d" argument to whatever time interval you want (1s, 1m, etc...), but be careful about using to short of a timeframe as you will continually spam the feed URLs until you kill the process.
	- Kill this process with ctrl+c
