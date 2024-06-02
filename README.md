# Skyjo
## Overview
In-browser online multiplayer application for the card game 'Skyjo' built in Vanilla React (TS) with a Go based WebSocket server.

## Project Status
Complete & playable but still has some Quality of Life/UX TODOs (see [TODOS](#TODOs))

## Screenshots

![image](https://github.com/AaronMolesbury/skyjo-online/assets/55638411/f3aa3e4d-e483-47ff-9cb3-66291f427de9)
![image](https://github.com/AaronMolesbury/skyjo-online/assets/55638411/7f8d34ac-a836-4a67-af02-ec55566997eb)
![image](https://github.com/AaronMolesbury/skyjo-online/assets/55638411/e02123fe-f83c-4db5-8d89-fb9a05bf3ca4)
![image](https://github.com/AaronMolesbury/skyjo-online/assets/55638411/f1533e2a-fb9c-4739-987c-5aeb1f69468b)


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
- Deploy it
- Preview other players hands when it's there go so all users understand what's going on
- Improve error handling both in UI & Server
- Animations?
