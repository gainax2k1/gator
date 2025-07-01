# gator project for boot.dev

Requires Postgres, Goose, and Go installed to run Gator.
To install Postgres (linux):
`apt install postgresql`
Then, to create your user account in Postgres, run:
`sudo -u postgres psql`
Inside Postgres,
`CREATE USER <username> WITH PASSWORD '<password>'`
Next, you'll need to create the database:
`CREATE DATABASE gator OWNER <username>;`
You can now quit out of Postgres with `quit`
To install Goose:
`go install github.com/pressly/goose/v3/cmd/goose@latest`

Install Gator CLI using `go install` while inside the root of the Gator repository

You'll need to use Goose to run the migrations for the initial database creation. From inside the /gator/sql/schema/ folder of the repository, run:
`goose postgres "postgres://<username>:<password>@localhost:5432/gator?sslmode=disable" up`

In order to run, Gator relies on a config file in home directory:

filename: "~/.gatorconfig.json"

with content:
`
{
"db_url": "postgres://<username>:<password>@localhost:5432/gator?sslmode=disable"
}
`

where <username> and <password> are the apropriate credentials for your postgresql

# Available commands:

`gator register <username>` creates user with <username> and sets to current active user

`gator login <username>` switches active user to <username>

`gator reset` resets database

`gator users` lists all usernames

`ator agg <time>` scrapes all followed feeds and stores posts in database at <time> interval. For example, every five minutes: `5m` 

`gator addfeed <"name"> <url>` adds feed to database

`gator feeds` lists feeds in database

`gator follow <url>` follows feed to current active username

`gator following` lists feeds current active user is following

`gator unfollow <url>` current active username unfollows <url>

`gator browse <optional number>` shows <number> of latest posts, if no number given, defaults to 2 
