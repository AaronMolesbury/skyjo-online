package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"nhooyr.io/websocket"
)

type card struct {
	value  int
	faceUp bool
}

type player struct {
	conn  *websocket.Conn
	hand  [][]*card
	ready bool
}

type socketDataCard struct {
	Value *int `json:"value"`
}

type socketDataPlayer struct {
	Hand [][]*socketDataCard `json:"hand"`
}

type socketData struct {
	DiscardCard      *socketDataCard    `json:"lastDiscardedCard"`
	CardInHand       *socketDataCard    `json:"cardInHand"`
	GameState        string             `json:"gameState"`
	PlayerID         int                `json:"playerId"`
	Players          []socketDataPlayer `json:"players"`
	CardClickEnabled bool               `json:"cardClickEnabled"`
}

const (
	MIN_PLAYER_COUNT   int = 2
	MAX_PLAYER_COUNT   int = 8
	DECK_SIZE          int = 150
	HAND_SIZE          int = 12
	MAX_CARD_FREQUENCY int = 15
)

var (
	players           []*player
	deck              []card
	lastDiscardedCard card
	currentCard       card
	currentGameState  string
	canClick          bool
	cardsFlipped      int
)

func main() {
	httpPort := flag.Int("http-port", 8080, "Sets http server to port X, defaults to :8080")
	flag.Parse()
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	player := &player{conn: conn}
	players = append(players, player)

	defer conn.Close(websocket.StatusNormalClosure, "")

	for {
		msgType, msg, err := conn.Read(context.Background())
		if err != nil || msgType != websocket.MessageText {
			return
		}

		m := strings.Split(string(msg), ";")
		msgCmd := m[0]

		switch msgCmd {
		case "ready":
			ready(player)
		case "take-from-deck":
			currentCard = takeFromDeck()
			currentGameState = "swap-discard"
			canClick = true
		case "take-from-discard":
			currentCard = lastDiscardedCard
			currentGameState = "force-swap"
			canClick = true
		case "discard":
			lastDiscardedCard = currentCard
			currentGameState = "flip"
		case "card-clicked":
			col, err := strconv.Atoi(m[1])
			if err != nil {
				return
			}
			row, err := strconv.Atoi(m[2])
			if err != nil {
				return
			}
			cardClicked(player.hand[row][col])
		}

		broadcastState()
	}
}

func broadcastState() {
	var socketPlayers []socketDataPlayer
	for _, player := range players {
		socketPlayer := socketDataPlayer{}
		for _, row := range player.hand {
			var socketRow []*socketDataCard
			for _, card := range row {
				if card == nil {
					socketRow = append(socketRow, nil)
				} else if card.faceUp {
					socketRow = append(socketRow, &socketDataCard{Value: &card.value})
				} else {
					socketRow = append(socketRow, &socketDataCard{})
				}
			}
			socketPlayer.Hand = append(socketPlayer.Hand, socketRow)
		}
		socketPlayers = append(socketPlayers, socketPlayer)
	}
	for i, player := range players {
		b, _ := json.Marshal(socketData{
			DiscardCard:      &socketDataCard{Value: &lastDiscardedCard.value},
			CardInHand:       &socketDataCard{Value: &currentCard.value},
			GameState:        currentGameState,
			PlayerID:         i,
			Players:          socketPlayers,
			CardClickEnabled: canClick,
		})
		player.conn.Write(context.Background(), websocket.MessageText, b)
	}
}

func ready(player *player) {
	player.ready = true
	for _, p := range players {
		if !p.ready {
			return
		}
	}
	beginGame()
}

func beginGame() {
	resetDeck()
	lastDiscardedCard = takeFromDeck()
	for _, player := range players {
		for row := 0; row < 3; row++ {
			player.hand = append(player.hand, make([]*card, 4))
			for col := 0; col < 4; col++ {
				newCard := takeFromDeck()
				player.hand[row][col] = &newCard
			}
		}
	}

	currentGameState = "flip-two"
	canClick = true
}

func resetDeck() {
	deck = deck[:0]
	currentCard = card{}

	for i := 0; i < 5; i++ {
		deck = append(deck, card{value: -2})
	}

	for i := 0; i < 10; i++ {
		for j := -1; j < 13; j++ {
			if j == 0 {
				continue
			}
			deck = append(deck, card{value: j})
		}
	}

	for i := 0; i < 15; i++ {
		deck = append(deck, card{value: 0})
	}
}

func takeFromDeck() card {
	randI := rand.Intn(len(deck))
	card := deck[randI]
	// TODO
	deck = append(deck[:randI], deck[randI+1:]...)
	return card
}

func cardClicked(cardClicked *card) {
	if currentGameState == "flip-two" {
		if cardsFlipped < 1 {
			cardClicked.faceUp = true
			cardsFlipped += 1
		} else {
			cardClicked.faceUp = true
			cardsFlipped += 1
			currentGameState = "take-from"
			canClick = false
		}
		return
	} else if currentGameState == "flip" {
		cardClicked.faceUp = true
		currentGameState = "take-from"
		canClick = false
		return
	} else { // Swap
		lastDiscardedCard = card{value: cardClicked.value}
		cardClicked.value = currentCard.value
		cardClicked.faceUp = true
		canClick = false
		currentGameState = "take-from"
	}
}
