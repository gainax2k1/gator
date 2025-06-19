gator project for boot.dev

- Requires Postgres and Go installed to run Gator.

- Install Gator CLI using "go install" command while inside the root of the Gator repository


- In order to run, Gator relies on config file in home directory with the following contents:

  "~/.gatorconfig.json"

  with content:
    {
      "db_url": "postgres://example"
    }

* Available commands: *

gator register <username> - creates user with <username> and sets to current active user

gator login <username> - switches active user to <username>

gator reset - resets database

gator users - lists all usernames

gator agg - scrapes all followed feeds and stores posts in database

gator addfeed <"name"> <url> - adds feed to database

gator feeds - lists feeds in database

gator follow <url> - follows feed to current active username

gator following - lists feeds current active user is following

gator unfollow <url> - current active username unfollows <url>

gator browse <optional number> - shows <number> of latest posts, if no number given, defaults to 2 

