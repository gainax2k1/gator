gator project for boot.dev

- Requires Postgres, Goose, and Go installed to run Gator.
-  To install Postgres (linux):
  <tt>apt install postgresql</tt>
  - Then, to create your user account in Postgres, run:
  <tt>sudo -u postgres psql<tt>
  - Inside Postgres,
  <tt>CREATE USER <username> WITH PASSWORD '<password>';
  - Next, you'll need to create the database:
  <tt> CREATE DATABASE gator OWNER <username>;</tt>
  - You can now quit out of Postgres with <tt>quit</tt>
  
 - To install Goose:
   <tt>go install github.com/pressly/goose/v3/cmd/goose@latest</tt>

- Install Gator CLI using "go install" command while inside the root of the Gator repository

- You'll need to use Goose to run the migrations for the initial database creation. From inside the /gator/sql/schema/ folder of the repository, run:
  <tt>goose postgres "postgres://<username>:<password>@localhost:5432/gator?sslmode=disable" up </tt>

- In order to run, Gator relies on config file in home directory with the following contents:

  "~/.gatorconfig.json"

  with content:
<tt>
    {
      "db_url": "postgres://<username>:<password>@localhost:5432/gator?sslmode=disable"
    }
</tt>
  where <username> and <password> are the apropriate credentials for your postgresql

* Available commands: *

gator register <username> - creates user with <username> and sets to current active user

gator login <username> - switches active user to <username>

gator reset - resets database

gator users - lists all usernames

gator agg <time> - scrapes all followed feeds and stores posts in database at <time> interval. For example, every five minutes: <tt>5m</tt> 

gator addfeed <"name"> <url> - adds feed to database

gator feeds - lists feeds in database

gator follow <url> - follows feed to current active username

gator following - lists feeds current active user is following

gator unfollow <url> - current active username unfollows <url>

gator browse <optional number> - shows <number> of latest posts, if no number given, defaults to 2 

