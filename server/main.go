package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"

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
	PlayerID int                `json:"playerId"`
	Players  []socketDataPlayer `json:"players"`
}

const (
	MIN_PLAYER_COUNT   int = 2
	MAX_PLAYER_COUNT   int = 8
	DECK_SIZE          int = 150
	HAND_SIZE          int = 12
	MAX_CARD_FREQUENCY int = 15
)

var (
	players     []*player
	deck        []card
	discardPile []card
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
		fmt.Println(string(msg))
		switch string(msg) {
		case "ready":
			ready(player)
		}
	}
}

func broadcastState() {
	fmt.Println("broadcasting")
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
			PlayerID: i,
			Players:  socketPlayers,
		})
		player.conn.Write(context.Background(), websocket.MessageText, b)
	}
}

func ready(player *player) {
	fmt.Println("ready")
	player.ready = true
	for _, p := range players {
		if !p.ready {
			return
		}
	}
	beginGame()
}

func beginGame() {
	fmt.Println("begin")
	resetDeck()
	discardPile = append(discardPile, takeFromDeck())
	for _, player := range players {
		for row := 0; row < 3; row++ {
			player.hand = append(player.hand, make([]*card, 4))
			for col := 0; col < 4; col++ {
				newCard := takeFromDeck()
				player.hand[row][col] = &newCard
			}
		}
	}
	broadcastState()
}

func resetDeck() {
	fmt.Println("reset")

	deck = deck[:0]

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
	fmt.Println("take")

	randI := rand.Intn(len(deck))
	card := deck[randI]
	// TODO
	deck = append(deck[:randI], deck[randI+1:]...)
	return card
}
