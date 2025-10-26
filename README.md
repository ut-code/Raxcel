# Raxcel

## Prerequisites

- Wails
- Bun >= 1.2
- Vercel CLI

## Setup

```sh
bun install
cd desktop
bun install
cp .env.sample .env
cd ../server
touch .env # Then paste the secrets shared on Discord
```

## Development
In one terminal,
```sh
cd frontend
wails dev
wails dev -tags webkit2_41 # develop on linux
wails build
wails build -tags webkit2_41 # build on linux
./build/bin/Raxcel # execute binary
```
In another terminal,
```sh
cd server
go run main.go
```

## Deployment
```sh
cd server
vc --prod
```
