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

type lobby struct {
	started            bool
	joinCode           int
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
	MIN_PLAYER_COUNT   int = 2
	MAX_PLAYER_COUNT   int = 8
	DECK_SIZE          int = 150
	HAND_SIZE          int = 12
	MAX_CARD_FREQUENCY int = 15
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

func createLobby(conn *websocket.Conn) {
	code := 0
	for i := 0; i < 5; i++ {
		code = 100000 + rand.Intn(999999-100000)

		if lobbies[code] == nil {
			break
		}
	}

	l := &lobby{
		joinCode: code,
		winnerID: -1,
	}

	lobbies[code] = l

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
			delete(lobbies, l.joinCode)
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
			newCard := l.takeFromDeck()
			l.currentCard = &newCard
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

func (l *lobby) broadcastState() {
	var socketPlayers []socketDataPlayer
	for _, player := range l.players {
		socketPlayer := socketDataPlayer{
			TurnType: player.turnType,
			Score:    player.currentScore,
		}
		socketPlayer.TurnType = player.turnType
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
	for i, player := range l.players {
		sd := socketData{
			Started:         l.started,
			LobbyCode:       l.joinCode,
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
		player.conn.Write(context.Background(), websocket.MessageText, b)
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
	newDiscardedCard := l.takeFromDeck()
	l.lastDiscardedCard = &newDiscardedCard
	for _, player := range l.players {
		player.turnType = "flip"
		for row := 0; row < 3; row++ {
			player.hand = append(player.hand, make([]*card, 4))
			for col := 0; col < 4; col++ {
				newCard := l.takeFromDeck()
				player.hand[row][col] = &newCard
			}
		}
	}
	l.currentPlayerIndex = 0
	l.broadcastState()
}

func (l *lobby) endGame() {
	l.winnerID = 0
	for i, player := range l.players {
		player.turnType = "gameOver"
		if player.currentScore < l.players[l.winnerID].currentScore {
			l.winnerID = i
		}
	}
	l.broadcastState()
}

func (l *lobby) resetDeck() {
	l.deck = l.deck[:0]
	l.currentCard = nil

	for i := 0; i < 5; i++ {
		l.deck = append(l.deck, card{value: -2})
	}

	for i := 0; i < 10; i++ {
		for j := -1; j < 13; j++ {
			if j == 0 {
				continue
			}
			l.deck = append(l.deck, card{value: j})
		}
	}

	for i := 0; i < 15; i++ {
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
	if len(p.lobby.players) < 2 {
		return
	}
	for _, player := range p.lobby.players {
		if !player.isReady {
			return
		}
	}
	p.lobby.beginGame()
}

func (p *player) cardClicked(cardClicked *card) {
	if p.turnType == "flip" {
		cardClicked.faceUp = true
		if p.countFaceUpCards() >= 2 {
			p.turnType = "take-from"
			p.lobby.nextPlayer()
		}
	} else { // Swap
		p.lobby.lastDiscardedCard = &card{value: cardClicked.value, faceUp: true}
		cardClicked.value = p.lobby.currentCard.value
		cardClicked.faceUp = true
		p.turnType = "take-from"
		p.lobby.nextPlayer()
	}
	p.lobby.currentCard = nil
}

func (p *player) countFaceUpCards() int {
	count := 0
	for _, row := range p.hand {
		for _, card := range row {
			if card != nil && card.faceUp {
				count += 1
			}
		}
	}
	return count
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
		for _, card := range row {
			if card != nil && !card.faceUp {
				complete = false
			}
			if card != nil && card.faceUp {
				p.currentScore += card.value
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
	totalScore := 0
	for _, row := range p.hand {
		for _, card := range row {
			card.faceUp = true
			totalScore += card.value
		}
	}
	p.currentScore = totalScore
}
