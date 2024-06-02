package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"nhooyr.io/websocket"
)

/* Structs */

type lobby struct {
	started            bool
	lobbyCode          int
	players            []*player
	deck               []card
	lastDiscardedCard  *card
	currentCard        *card
	currentPlayerIndex int
	winnerID           int
}

type player struct {
	lobby        *lobby
	conn         *websocket.Conn
	hand         [][]*card
	isReady      bool
	currentScore int
	turnType     string
}

type card struct {
	value  int
	faceUp bool
}

func (c *card) equalValue(c2 *card) bool {
	if c == nil || c2 == nil || !c.faceUp || !c2.faceUp {
		return false
	}

	return c.value == c2.value
}

/* WS Structs */

type socketDataCard struct {
	Value *int `json:"value"`
}

type socketDataPlayer struct {
	Hand     [][]*socketDataCard `json:"hand"`
	TurnType string              `json:"turnType"`
	Score    int                 `json:"score"`
}

type socketData struct {
	Started         bool               `json:"started"`
	LobbyCode       int                `json:"lobbyCode"`
	DiscardCard     *socketDataCard    `json:"lastDiscardedCard"`
	CardInHand      *socketDataCard    `json:"cardInHand"`
	PlayerID        int                `json:"playerId"`
	Players         []socketDataPlayer `json:"players"`
	CurrentPlayerID int                `json:"currentPlayerId"`
	WinnerID        int                `json:"winnerId"`
}

const (
	MIN_RAND_LOBBY_VAL  = 100000
	MAX_RAND_LOBBY_VAL  = 999999
	MAX_MINUS_TWO_CARDS = 5
	MAX_ZERO_CARDS      = 15
	MAX_OTHER_CARDS     = 10
	MIN_PLAYERS         = 2
)

var lobbies = make(map[int]*lobby)

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

	if r.URL.Path == "/" {
		createLobby(conn)
	} else {
		lobbyCode, err := strconv.Atoi(r.URL.Path[1:])
		if err != nil {
			fmt.Print("Bad join request")
			return
		}
		l := lobbies[lobbyCode]
		if l != nil {
			l.joinLobby(conn)
		}
	}
}

/* Lobby Methods */

func createLobby(conn *websocket.Conn) {
	c := 0
	for i := 0; i < 5; i++ {
		c = MIN_RAND_LOBBY_VAL + rand.Intn(MAX_RAND_LOBBY_VAL-MIN_RAND_LOBBY_VAL)

		if lobbies[c] == nil {
			break
		}
	}

	l := &lobby{
		lobbyCode: c,
		winnerID:  -1,
	}

	lobbies[c] = l

	l.joinLobby(conn)
}

func (l *lobby) joinLobby(conn *websocket.Conn) {
	p := &player{lobby: l, conn: conn}
	l.players = append(l.players, p)
	pID := len(l.players) - 1

	l.broadcastState()

	defer conn.Close(websocket.StatusNormalClosure, "")
	defer func() {
		i := slices.Index(l.players, p)
		l.players = slices.Delete(l.players, i, i+1)
		if len(l.players) == 0 {
			delete(lobbies, l.lobbyCode)
		}
	}()

	for {
		msgType, msg, err := conn.Read(context.Background())
		if err != nil || msgType != websocket.MessageText {
			return
		}

		m := strings.Split(string(msg), ";")
		msgCmd := m[0]

		if msgCmd == "ready" {
			p.ready()
			continue
		}

		if msgCmd == "reset" {
			l.beginGame()
			continue
		}

		if pID != l.currentPlayerIndex {
			continue
		}

		switch msgCmd {
		case "take-from-deck":
			c := l.takeFromDeck()
			l.currentCard = &c
			p.turnType = "swap-discard"
		case "take-from-discard":
			l.currentCard = l.lastDiscardedCard
			l.lastDiscardedCard = nil
			p.turnType = "force-swap"
		case "discard":
			l.lastDiscardedCard = l.currentCard
			l.currentCard = nil
			p.turnType = "flip"
		case "card-clicked":
			col, err := strconv.Atoi(m[1])
			if err != nil {
				return
			}
			row, err := strconv.Atoi(m[2])
			if err != nil {
				return
			}
			p.cardClicked(p.hand[row][col])
		}
		p.endMove()
	}
}

// Generates & sends data to ui
func (l *lobby) broadcastState() {
	var socketPlayers []socketDataPlayer
	for _, p := range l.players {
		socketPlayer := socketDataPlayer{
			TurnType: p.turnType,
			Score:    p.currentScore,
		}
		socketPlayer.TurnType = p.turnType
		for _, row := range p.hand {
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
	for i, p := range l.players {
		sd := socketData{
			Started:         l.started,
			LobbyCode:       l.lobbyCode,
			PlayerID:        i,
			Players:         socketPlayers,
			CurrentPlayerID: l.currentPlayerIndex,
			WinnerID:        l.winnerID,
		}

		if l.currentCard != nil {
			sd.CardInHand = &socketDataCard{Value: &l.currentCard.value}
		}

		if l.lastDiscardedCard != nil {
			sd.DiscardCard = &socketDataCard{Value: &l.lastDiscardedCard.value}
		}

		b, _ := json.Marshal(sd)
		p.conn.Write(context.Background(), websocket.MessageText, b)
	}
}

func (l *lobby) nextPlayer() {
	l.currentPlayerIndex += 1
	if l.currentPlayerIndex == len(l.players) {
		l.currentPlayerIndex = 0
	}
}

func (l *lobby) beginGame() {
	l.started = true
	l.resetDeck()
	l.winnerID = -1
	c := l.takeFromDeck()
	l.lastDiscardedCard = &c
	for _, p := range l.players {
		p.turnType = "flip"
		for row := 0; row < 3; row++ {
			p.hand = append(p.hand, make([]*card, 4))
			for col := 0; col < 4; col++ {
				nc := l.takeFromDeck()
				p.hand[row][col] = &nc
			}
		}
	}
	l.currentPlayerIndex = 0
	l.broadcastState()
}

func (l *lobby) endGame() {
	l.winnerID = 0
	for i, p := range l.players {
		p.turnType = "gameOver"
		if p.currentScore < l.players[l.winnerID].currentScore {
			l.winnerID = i
		}
	}
	l.broadcastState()
}

func (l *lobby) resetDeck() {
	l.deck = l.deck[:0]
	l.currentCard = nil

	for i := 0; i < MAX_MINUS_TWO_CARDS; i++ {
		l.deck = append(l.deck, card{value: -2})
	}

	for i := 0; i < MAX_OTHER_CARDS; i++ {
		for j := -1; j < 13; j++ {
			if j == 0 {
				continue
			}
			l.deck = append(l.deck, card{value: j})
		}
	}

	for i := 0; i < MAX_ZERO_CARDS; i++ {
		l.deck = append(l.deck, card{value: 0})
	}
}

func (l *lobby) takeFromDeck() card {
	randI := rand.Intn(len(l.deck))
	card := l.deck[randI]
	l.deck = append(l.deck[:randI], l.deck[randI+1:]...)
	return card
}

/* PLAYER METHODS */

func (p *player) ready() {
	p.isReady = true
	if len(p.lobby.players) < MIN_PLAYERS {
		return
	}
	for _, player := range p.lobby.players {
		if !player.isReady {
			return
		}
	}
	p.lobby.beginGame()
}

func (p *player) cardClicked(c *card) {
	if p.turnType == "flip" {
		c.faceUp = true
		if p.countFaceUpCards() >= 2 {
			p.turnType = "take-from"
			p.lobby.nextPlayer()
		}
	} else { // Swap
		p.lobby.lastDiscardedCard = &card{value: c.value, faceUp: true}
		c.value = p.lobby.currentCard.value
		c.faceUp = true
		p.turnType = "take-from"
		p.lobby.nextPlayer()
	}
	p.lobby.currentCard = nil
}

func (p *player) countFaceUpCards() int {
	ct := 0
	for _, row := range p.hand {
		for _, c := range row {
			if c != nil && c.faceUp {
				ct += 1
			}
		}
	}
	return ct
}

/* Special rule - if all cards in a col are face up & share the same value, that col will be removed */
func (p *player) checkColEquality() {
	rows, cols := len(p.hand), len(p.hand[0])
	// Check columns
	for j := 0; j < cols; j++ {
		var comparisonCard *card
		// Check for non nil
		for i := 0; i < rows; i++ {
			if p.hand[i][j] != nil && p.hand[i][j].faceUp {
				comparisonCard = p.hand[i][j]
				break
			}
		}
		// Skip if whole col is nil
		if comparisonCard == nil {
			continue
		}

		equal := true
		for i := 0; i < rows; i++ {
			if !p.hand[i][j].equalValue(comparisonCard) {
				equal = false
				break
			}
		}
		if equal {
			p.lobby.lastDiscardedCard = comparisonCard
			for i := 0; i < rows; i++ {
				p.hand[i][j] = nil
			}
		}
	}
}

/* Check for if all cards in hand are face up */
func (p *player) checkScoreAndCompleteHand() bool {
	complete := true
	p.currentScore = 0
	for _, row := range p.hand {
		for _, c := range row {
			if c != nil && !c.faceUp {
				complete = false
			}
			if c != nil && c.faceUp {
				p.currentScore += c.value
			}
		}
	}
	return complete
}

func (p *player) endMove() {
	p.checkColEquality()
	if p.checkScoreAndCompleteHand() {
		p.turnType = "complete"
	} else if p.turnType == "take-from" {
		for _, player := range p.lobby.players {
			if player.turnType == "complete" {
				p.turnType = "complete"
				p.flipAllCards()
				break
			}
		}
	}

	allComplete := true
	for _, player := range p.lobby.players {
		if player.turnType != "complete" {
			allComplete = false
			break
		}
	}

	if allComplete {
		p.lobby.endGame()
	}

	p.lobby.broadcastState()
}

func (p *player) flipAllCards() {
	s := 0
	for _, row := range p.hand {
		for _, c := range row {
			c.faceUp = true
			s += c.value
		}
	}
	p.currentScore = s
}
