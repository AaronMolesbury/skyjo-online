package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

var (
	currentPlayerId int
	players         map[int]*websocket.Conn
)

func main() {
	httpPort := flag.Int("http-port", 8080, "Sets http server to port X, defaults to :8080")
	flag.Parse()
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}

	defer conn.Close(websocket.StatusNormalClosure, "")

	if r.URL.Path == "/" {
		currentPlayerId += 1
		players[currentPlayerId] = conn

		// TODO - Handle player commands
	}
}

/* GAME LOGIC */

// Used for random weighting
var baseDeck = []int{
	-2, -2, -2, -2, -2,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
	11, 11, 11, 11, 11, 11, 11, 11, 11, 11,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
}

var inPlayDeck []int

/* TODO - Game Functions*/

func initInPlayDeck() {
	inPlayDeck = baseDeck
	deal()
}

// Resets inplay deck to be the base deck
func resetDeck() {

}

// Deals every player hand from inplay deck
func deal() {

}

// Gives player random card from deck
// - must be followed by swap or discard
func takeFromDeck() {

}

// PLAYER CALLED
// Swaps card taken from deck with card in hand
// - must be followed by discard (of card originally in hand)
// - must end turn
func swapCard() {

}

// PLAYERCALLED
// Player discards card either taken from deck or swapped from hand by deck card
// - if card from deck, must be followed by flipcard
func discardCard() {

}

// PLAYERCALLED
// Player flips selected card
// - must end turn
func flipCard() {

}
