# gator

CLI for RSS feed management: user auth, feeds, follow/unfollow, aggregation.

## Requirements

- **Go** 1.26+
- **PostgreSQL**

## Installation

Install the binary with `go install`:

```bash
go install github.com/realdnchka/gator-go@latest
```

Ensure `$GOPATH/bin` (or `$HOME/go/bin`) is in your `PATH` so the `gator` binary is available.

## Configuration

gator reads config from `~/.gatorconfig.json`. It must contain a valid PostgreSQL connection string (`db_url`) and optionally the current user (`current_user_name`).

Example:

```json
{
  "db_url": "postgres://user:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Create the database and run migrations (see `sql/schema/`) before using the CLI.

## Commands

| Command    | Description                    |
|------------|--------------------------------|
| `register` | Create a new user              |
| `login`    | Log in                         |
| `reset`    | Reset users db                 |
| `users`    | List users                     |
| `agg`      | Aggregate feeds (requires auth)|
| `addfeed`  | Add a feed (requires auth)     |
| `feeds`    | List all feeds                 |
| `follow`   | Follow a feed (requires auth)  |
| `following`| List followed feeds (requires auth) |
| `unfollow` | Unfollow a feed (requires auth)|
| `browse`   | Browse posts (requires auth)   |

Usage:

```bash
gator <command> [args...]
```
