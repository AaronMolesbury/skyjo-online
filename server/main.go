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
	DiscardCard *socketDataCard    `json:"lastDiscardedCard"`
	CardInHand  *socketDataCard    `json:"cardInHand"`
	GameState   string             `json:"gameState"`
	PlayerID    int                `json:"playerId"`
	Players     []socketDataPlayer `json:"players"`
	Score       int                `json:"score"`
}

const (
	MIN_PLAYER_COUNT   int = 2
	MAX_PLAYER_COUNT   int = 8
	DECK_SIZE          int = 150
	HAND_SIZE          int = 12
	MAX_CARD_FREQUENCY int = 15
)

var (
	players            []*player
	currentPlayerIndex int
	currentPlayer      *player
	deck               []card
	lastDiscardedCard  card
	currentCard        card
	currentGameState   string
	cardsFlipped       int
	score              int
	resetCard          card = card{value: -3}
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

		if msgCmd == "ready" {
			ready(player)
			continue
		}

		if player != currentPlayer {
			continue
		}

		switch msgCmd {
		case "take-from-deck":
			currentCard = takeFromDeck()
			currentGameState = "swap-discard"
		case "take-from-discard":
			currentCard = lastDiscardedCard
			lastDiscardedCard = resetCard
			currentGameState = "force-swap"
		case "discard":
			lastDiscardedCard = currentCard
			currentCard = resetCard
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
		endMove(*player)
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
			DiscardCard: &socketDataCard{Value: &lastDiscardedCard.value},
			CardInHand:  &socketDataCard{Value: &currentCard.value},
			GameState:   currentGameState,
			PlayerID:    i,
			Players:     socketPlayers,
			Score:       score,
		})
		player.conn.Write(context.Background(), websocket.MessageText, b)
	}
}

func updateCurrentPlayer(override int) {
	if override != -1 {
		currentPlayerIndex = override
	} else {
		currentPlayerIndex += 1
		if currentPlayerIndex == len(players) {
			currentPlayerIndex = 0
		}
	}
	currentPlayer = players[currentPlayerIndex]
}

func ready(player *player) {
	player.ready = true
	for i, p := range players {
		fmt.Printf("Player %d is %v \n", i, p.ready)
		if !p.ready {
			fmt.Println("not everyone is ready")
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
	updateCurrentPlayer(0)
	currentGameState = "flip-two"
	broadcastState()
}

func resetDeck() {
	deck = deck[:0]
	currentCard = resetCard

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
		}
		return
	} else if currentGameState == "flip" {
		cardClicked.faceUp = true
		currentGameState = "take-from"
		updateCurrentPlayer(-1)
		return
	} else { // Swap
		lastDiscardedCard = card{value: cardClicked.value}
		cardClicked.value = currentCard.value
		cardClicked.faceUp = true
		currentGameState = "take-from"
		updateCurrentPlayer(-1)
	}
	currentCard = resetCard
}

/* Special rule; if all cards in a col are face up & share the same value, that col will be removed */
func checkColEquality(hand [][]*card) {
	rows, cols := len(hand), len(hand[0])
	// Check columns
	for j := 0; j < cols; j++ {
		var comparisonCard *card
		// Check for non nil
		for i := 0; i < rows; i++ {
			if hand[i][j] != nil && hand[i][j].faceUp {
				comparisonCard = hand[i][j]
				break
			}
		}
		// Skip if whole col is nil
		if comparisonCard == nil {
			continue
		}

		equal := true
		for i := 0; i < rows; i++ {
			if hand[i][j] != nil && hand[i][j].value != comparisonCard.value {
				equal = false
				break
			}
		}
		if equal {
			lastDiscardedCard = *comparisonCard
			for i := 0; i < rows; i++ {
				hand[i][j] = nil
			}
		}
	}
}

/* Check for if all cards in hand are face up */
func checkScoreAndCompleteHand(hand [][]*card) bool {
	complete := true
	score = 0
	for _, row := range hand {
		for _, card := range row {
			if card != nil && !card.faceUp {
				complete = false
			}
			if card != nil && card.faceUp {
				score += card.value
			}
		}
	}
	return complete
}

func endMove(player player) {
	checkColEquality(player.hand)
	if checkScoreAndCompleteHand(player.hand) {
		currentGameState = "complete"
	}
	broadcastState()
}
