# Gator

Gator is a command-line RSS/blog aggregator written in Go. It lets users register, add RSS feeds, follow feeds, fetch posts from those feeds, and browse recent posts from a PostgreSQL database.

## Requirements

To run Gator, you need:

- Go installed
- PostgreSQL installed and running
- A PostgreSQL database for Gator
- `goose` installed if you want to run the included database migrations

## Install

Install the Gator CLI with `go install`:

```bash
go install github.com/kayrabulbul/gator@latest
```

Make sure your Go binary directory is on your `PATH`. This is usually one of:

```bash
$HOME/go/bin
```

or:

```bash
$GOPATH/bin
```

After installation, the `gator` command should be available:

```bash
gator <command> [arguments]
```

## Database Setup

Create a PostgreSQL database for Gator. For example:

```bash
createdb gator
```

Then run the migrations from a clone of this repository:

```bash
goose -dir sql/schema postgres "postgres://user:password@localhost:5432/gator?sslmode=disable" up
```

Replace `user`, `password`, host, port, and database name with your PostgreSQL settings.

## Config Setup

Gator reads its config from a file in your home directory:

```text
~/.gatorconfig.json
```

Create that file with the following shape:

```json
{
  "connection_string": "postgres://user:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Update `connection_string` to match your PostgreSQL database. The `current_user_name` field can start empty. It will be updated when you register or log in.

## Running Gator

Commands follow this format:

```bash
gator <command> [arguments]
```

Example workflow:

```bash
gator register alice
gator addfeed "Hacker News" "https://hnrss.org/frontpage"
gator agg 1m
gator browse 5
```

The `agg` command runs continuously and fetches feeds on the interval you provide. Stop it with `Ctrl+C`.

## Commands

| Command     | Usage                        | Description                                                     |
| ----------- | ---------------------------- | --------------------------------------------------------------- |
| `register`  | `gator register <username>`  | Create a new user and set them as the current user.             |
| `login`     | `gator login <username>`     | Set an existing user as the current user.                       |
| `users`     | `gator users`                | List all users. The current user is marked.                     |
| `reset`     | `gator reset`                | Delete all users from the database.                             |
| `addfeed`   | `gator addfeed <name> <url>` | Add a new feed and automatically follow it as the current user. |
| `feeds`     | `gator feeds`                | List all feeds.                                                 |
| `follow`    | `gator follow <url>`         | Follow an existing feed as the current user.                    |
| `following` | `gator following`            | List feeds followed by the current user.                        |
| `unfollow`  | `gator unfollow <url>`       | Unfollow a feed as the current user.                            |
| `agg`       | `gator agg <duration>`       | Fetch feeds repeatedly, for example `30s`, `1m`, or `1h`.       |
| `browse`    | `gator browse [limit]`       | Show recent posts. The default limit is `2`.                    |

## Notes

- `addfeed`, `follow`, and `unfollow` require a current user.
- Use `register` to create your first user.
- Feed URLs and post URLs are unique in the database, so duplicates are not stored.
- `reset` deletes users from the database.
