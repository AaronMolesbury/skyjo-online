# Skyjo
## Overview
In-browser online multiplayer application for the card game 'Skyjo' built in Vanilla React (TS) with a Go based WebSocket server.

## Project Status
Complete & playable but still has some Quality of Life/UX TODOs (see [TODOS](#TODOs))

## Screenshots



## Installation & Setup
I used Vite to run the client locally and standard Go for the server:

Clone this repository, since it's a React/Go project you will require `node` (see [here](https://nodejs.org/en/learn/getting-started/how-to-install-nodejs)) and `npm` for running the client and `go` for the server (see [here](https://go.dev/doc/install)).

### Install Dependencies

`npm install`

### Run Client
Ensure you're in the `client` directory

`cd ./client`

and run

`npm run dev`

### Run Server

Ensure you're in the `server` directory

`cd ./server`

and run

`go run main.go`

## <a name="TODOs"></a> TODOs 
There are a few Quality of Life and UX changes I'd like to make when I have the time:

- Preview other players hands when it's there go so all users understand what's going on
- Improve error handling both in UI & Server
- Animations?
