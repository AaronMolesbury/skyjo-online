package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"nhooyr.io/websocket"
)

// func main() {
// 	httpPort := flag.Int("http-port", 8080, "Sets http server to port X, defaults to :8080")
// 	flag.Parse()
// 	http.HandleFunc("/", handle)
// 	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil))
// }

// func handle(w http.ResponseWriter, r *http.Request) {
// 	conn, err := websocket.Accept(w, r, nil)
// 	if err != nil {
// 		return
// 	}

// 	defer conn.Close(websocket.StatusNormalClosure, "")

// 	if r.URL.Path == "/" {
// 		currentPlayerId += 1
// 		players[currentPlayerId] = conn

// 		// TODO - Handle player commands
// 	}
// }

/* GAME LOGIC */

type card struct {
	visible bool
	value   int
}

type player struct {
	id   int
	conn *websocket.Conn
	hand []card
}

var (
	deck        []card
	discardPile []card
	players     []player
)

const (
	MIN_PLAYER_COUNT   int = 2
	MAX_PLAYER_COUNT   int = 8
	DECK_SIZE          int = 150
	HAND_SIZE          int = 12
	MAX_CARD_FREQUENCY int = 15
)

/* TODO - Game Functions*/

// / Debug
func main() {
	beginGame()
}

///

func beginGame() {
	players = make([]player, 0, MAX_PLAYER_COUNT)
	players = append(players, player{1, nil, nil})
	players = append(players, player{2, nil, nil})
	resetDeck()
	deal()
	turn(players[0])
}

func resetDeck() {
	if deck == nil {
		deck = make([]card, 0, DECK_SIZE)
	} else {
		deck = deck[:0]
	}

	for i := 0; i < MAX_CARD_FREQUENCY; i++ {
		if i < 5 {
			deck = append(deck, card{visible: false, value: -2})
		}
		if i < 10 {
			deck = append(deck, card{visible: false, value: -1})
			deck = append(deck, card{visible: false, value: 1})
			deck = append(deck, card{visible: false, value: 2})
			deck = append(deck, card{visible: false, value: 3})
			deck = append(deck, card{visible: false, value: 4})
			deck = append(deck, card{visible: false, value: 5})
			deck = append(deck, card{visible: false, value: 6})
			deck = append(deck, card{visible: false, value: 7})
			deck = append(deck, card{visible: false, value: 8})
			deck = append(deck, card{visible: false, value: 9})
			deck = append(deck, card{visible: false, value: 10})
			deck = append(deck, card{visible: false, value: 11})
			deck = append(deck, card{visible: false, value: 12})
		}

		deck = append(deck, card{visible: false, value: 0})
	}
}

// Deals every player hand from inplay deck + starts discard pile
func deal() {
	for i := range players {
		players[i].hand = generateHand()
	}

	discardPile = append(discardPile, takeFromDeck())
}

func generateHand() []card {
	hand := []card{}

	for i := 0; i < HAND_SIZE; i++ {
		hand = append(hand, takeFromDeck())
	}

	return hand
}

// Gives player random card from deck
// - must be followed by swap or discard
func takeFromDeck() card {
	if len(deck) == 0 {
		resetDeck()
	}
	randI := rand.Intn(len(deck))
	card := deck[randI]
	deck = append(deck[:randI], deck[randI+1:]...)
	return card
}

func turn(player player) {
	var newCard card

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Discard pile top card: %v\n", discardPile[len(discardPile)-1])
	fmt.Print("Take from discard pile or deck\n")
	input, _ := reader.ReadString('\n')

	switch input {
	case "x\n":
		newCard = discardPile[len(discardPile)-1]
		discardPile = discardPile[:len(discardPile)-1]

	default:
		newCard = takeFromDeck()
	}

	fmt.Printf("player %d draws a %d \n", player.id, newCard.value)

	fmt.Printf("Current hand: %v\n", player.hand)
	fmt.Print("Swap or Discard\n")
	input, _ = reader.ReadString('\n')
	switch input {
	case "s\n":
		swapCard(player, newCard)

	case "d\n":
		discardPile = append(discardPile, newCard)
		//flipCard(player)
	}
}

// // PLAYER CALLED
// // Swaps card taken from deck with card in hand
// // - must be followed by discard (of card originally in hand)
// // - must end turn
func swapCard(player player, newCard card) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("index to swap\n")
	input, _ := reader.ReadString('\n')

	i, err := strconv.Atoi(input[:len(input)-1])
	if err != nil {
		print("bad conversion")
		return
	}

	player.hand[i] = newCard
	fmt.Printf("New hand: %v\n", player.hand)
}

// // PLAYERCALLED
// // Player flips selected card
// // - must end turn
// func flipCard(player player) {
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("index to flip\n")

// 	input, _ := reader.ReadString('\n')

// 	i, err := strconv.Atoi(input[:len(input)-1])
// 	if err != nil {
// 		print("bad conversion")
// 		return
// 	}

// 	player.hand[i] = newCard
// 	fmt.Printf("New hand: %v\n", player.hand)
// }
