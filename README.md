# music-box-game

## Backend

Installation:

```sh
go install github.com/air-verse/air@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
go mod download
```

Start it up:

```sh
docker compose up -d
air
```

## Mobile

```sh
npm install
```

and

```sh
npx expo start
```

See [mobile/README.md](mobile/README.md) for details.
