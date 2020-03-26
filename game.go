package main

import "sync"

//import "fmt"

// ActionBase base type for actions
type ActionBase struct {
	ActionType string `json:"type"`
}

// ActionTypePlay action type for play action
const ActionTypePlay = "play"

// ActionPlay action data for play action
type ActionPlay struct {
	ActionBase
	Index int `json:"index"`
}

// ActionTypeHint action type for hint action
const ActionTypeHint = "hint"

// ActionHint action data for hint action
type ActionHint struct {
	ActionBase
}

// ActionTypeDiscard action type for discard action
const ActionTypeDiscard = "discard"

// ActionDiscard action data for discard action
type ActionDiscard struct {
	ActionBase
	Index int `json:"index"`
}

// ActionTypeTilt action type for tilt card action
const ActionTypeTilt = "tilt"

// ActionTilt action data for tilt card action
type ActionTilt struct {
	ActionBase
	Index int  `json:"index"`
	Tilt  bool `json:"tilt"`
}

// ActionTypeMove action type for move card action
const ActionTypeMove = "move"

// ActionMove action data for move card action
type ActionMove struct {
	ActionBase
	From int `json:"from"`
	To   int `json:"to"`
}

// EventBase base type for events
type EventBase struct {
	EventType string `json:"type"`
}

// EventTypeSetup type for setup event
var EventTypeSetup = EventBase{EventType: "setup"}

// EventSetup data for setup event
type EventSetup struct {
	EventBase
	NumPlayers int `json:"n_players"`
	YourIndex  int `json:"your_player"`
}

// EventTypeSetNames type for set names event
var EventTypeSetNames = EventBase{EventType: "set_names"}

// EventSetNames data for set names event
type EventSetNames struct {
	EventBase
	PlayerNames []string `json:"player_names"`
}

// EventTypeTurnChange type for turn change event
var EventTypeTurnChange = EventBase{EventType: "turn_change"}

// EventTurnChange data for turn change event
type EventTurnChange struct {
	EventBase
	NextPlayer int `json:"next_player"`
}

// EventTypeHint type for hint event
var EventTypeHint = EventBase{EventType: "hint"}

// EventHint data for hint event
type EventHint struct {
	EventBase
	RemainingHints int `json:"remaining_hints"`
}

// EventTypePlay type for play event
var EventTypePlay = EventBase{EventType: "play"}

// EventPlay data for play event
type EventPlay struct {
	EventBase
	Show []int `json:"show"`
}

// EventTypeFail type for fail event
var EventTypeFail = EventBase{EventType: "fail"}

// EventFail data for fail event
type EventFail struct {
	EventBase
	Discard  *Card `json:"discard"`
	NumFails int   `json:"n_fails"`
}

// EventTypeHandDiscard type for hand discard event
var EventTypeHandDiscard = EventBase{EventType: "hand_discard"}

// EventHandDiscard event data for hand discard event
type EventHandDiscard struct {
	EventBase
	Player int    `json:"player"`
	Index  int    `json:"index"`
	Hand   []Card `json:"hand"`
}

// EventTypeBlindHandDiscard type fro blind hand discard event
var EventTypeBlindHandDiscard = EventBase{EventType: "blind_hand_discard"}

// EventBlindHandDiscard data for blind hand discard event
type EventBlindHandDiscard struct {
	EventBase
	Player   int `json:"player"`
	Index    int `json:"index"`
	HandSize int `json:"hand_size"`
}

// EventTypeHandDraw type for hand draw event
var EventTypeHandDraw = EventBase{EventType: "hand_draw"}

// EventHandDraw data for hand draw event
type EventHandDraw struct {
	EventBase
	Player int    `json:"player"`
	Hand   []Card `json:"hand"`
}

// EventTypeBlindHandDraw type for blind hand draw event
var EventTypeBlindHandDraw = EventBase{EventType: "blind_hand_draw"}

// EventBlindHandDraw data for blind hand draw event
type EventBlindHandDraw struct {
	EventBase
	Player int    `json:"player"`
	Hand   []bool `json:"hand"`
}

// EventTypeDeckDraw type for deck draw event
var EventTypeDeckDraw = EventBase{EventType: "deck_draw"}

// EventDeckDraw event data for deck draw event
type EventDeckDraw struct {
	EventBase
	DeckRemaining int `json:"deck_remaining"`
}

// EventTypeTilt type for tilt event
var EventTypeTilt = EventBase{EventType: "tilt"}

// EventTilt event data for card tilt event
type EventTilt struct {
	EventBase
	Player int  `json:"player"`
	Index  int  `json:"index"`
	Tilt   bool `json:"tilt"`
}

// EventTypeMove type for move event
var EventTypeMove = EventBase{EventType: "move"}

// EventMove event data for card move event
type EventMove struct {
	EventBase
	Player int `json:"player"`
	From   int `json:"from"`
	To     int `json:"to"`
}

// GameStateListener listen channel for events related to game.
type GameStateListener struct {
	l           chan<- interface{}
	playerIndex int
}

// GameState state for game
type GameState struct {
	sync.Mutex

	nPlayers    int
	curPlayer   int
	lastPlayer  int
	playerNames []string

	hints int
	fails int

	deck        []Card
	discard     *Card
	playerHands [][]Card

	fieldStatus [5]int

	listeners []*GameStateListener
}

// NewGame create new game
func NewGame(nPlayers int) *GameState {
	result := &GameState{
		nPlayers:    nPlayers,
		curPlayer:   0,
		lastPlayer:  -1,
		playerNames: make([]string, nPlayers),

		hints: 8,
		fails: 0,

		deck:        randomDeck(),
		discard:     nil,
		playerHands: make([][]Card, nPlayers),

		fieldStatus: [5]int{0, 0, 0, 0, 0},
	}

	for i := range result.playerHands {
		result.playerHands[i] = []Card{
			result.drawCard(),
			result.drawCard(),
			result.drawCard(),
			result.drawCard(),
		}
		if nPlayers <= 3 {
			result.playerHands[i] = append(result.playerHands[i], result.drawCard())
		}
	}

	return result
}

func (s *GameState) drawCard() Card {
	res := s.deck[len(s.deck)-1]
	s.deck = s.deck[:len(s.deck)-1]
	return res
}

func (s *GameState) nextPlayer() {
	if s.curPlayer == s.lastPlayer {
		s.curPlayer = -1
	} else {
		s.curPlayer = s.curPlayer + 1
		if s.curPlayer == s.nPlayers {
			s.curPlayer = 0
		}
	}
}

// AddListener add listener to game event stream
func (s *GameState) AddListener(l *GameStateListener) {
	s.Lock()
	defer s.Unlock()

	l.l <- EventSetup{EventBase: EventTypeSetup, NumPlayers: s.nPlayers, YourIndex: l.playerIndex}
	namesCopy := make([]string, len(s.playerNames))
	copy(namesCopy, s.playerNames)
	l.l <- EventSetNames{EventBase: EventTypeSetNames, PlayerNames: namesCopy}
	l.l <- EventDeckDraw{EventBase: EventTypeDeckDraw, DeckRemaining: len(s.deck)}

	for i, hand := range s.playerHands {
		if i == l.playerIndex {
			handCopy := make([]bool, len(hand))
			for i := range hand {
				handCopy[i] = hand[i].Tilt
			}
			l.l <- EventBlindHandDraw{EventBase: EventTypeBlindHandDraw, Player: i, Hand: handCopy}
		} else {
			handCopy := make([]Card, len(hand))
			copy(handCopy, hand)
			l.l <- EventHandDraw{EventBase: EventTypeHandDraw, Player: i, Hand: handCopy}
		}
	}

	fieldCopy := make([]int, 5)
	copy(fieldCopy, s.fieldStatus[:])
	l.l <- EventPlay{EventBase: EventTypePlay, Show: fieldCopy}
	if s.discard != nil {
		discardPlay := *s.discard
		l.l <- EventFail{EventBase: EventTypeFail, Discard: &discardPlay, NumFails: s.fails}
	} else {
		l.l <- EventFail{EventBase: EventTypeFail, Discard: nil, NumFails: s.fails}
	}
	l.l <- EventHint{EventBase: EventTypeHint, RemainingHints: s.hints}
	l.l <- EventTurnChange{EventBase: EventTypeTurnChange, NextPlayer: s.curPlayer}

	s.listeners = append(s.listeners, l)
}

// RemoveListener remove listener from game event stream
func (s *GameState) RemoveListener(l *GameStateListener) {
	s.Lock()
	defer s.Unlock()

	for i, ll := range s.listeners {
		if ll == l {
			s.listeners[i] = s.listeners[len(s.listeners)-1]
			s.listeners = s.listeners[:len(s.listeners)-1]
			return
		}
	}
}

// DoHint play a hint action as player
func (s *GameState) DoHint(player int) {
	s.Lock()
	defer s.Unlock()

	if s.curPlayer != player {
		return
	}

	if s.hints == 0 {
		return
	}

	// Use hint
	s.hints = s.hints - 1
	// Advance turn
	s.nextPlayer()

	for _, l := range s.listeners {
		l.l <- EventHint{EventBase: EventTypeHint, RemainingHints: s.hints}
		l.l <- EventTurnChange{EventBase: EventTypeTurnChange, NextPlayer: s.curPlayer}
	}
}

// DoDiscard discard card with index as player
func (s *GameState) DoDiscard(player, index int) {
	s.Lock()
	defer s.Unlock()

	if s.curPlayer != player {
		return
	}

	if index < 0 || index >= len(s.playerHands[player]) {
		return
	}

	// Increase hints available if possible
	if s.hints < 8 {
		s.hints = s.hints + 1
	}
	// Discard card
	discard := s.playerHands[player][index]
	s.discard = &discard
	copy(s.playerHands[player][index:], s.playerHands[player][index+1:])
	s.playerHands[player] = s.playerHands[player][:len(s.playerHands[player])-1]
	for _, l := range s.listeners {
		l.l <- EventHint{EventBase: EventTypeHint, RemainingHints: s.hints}
		if l.playerIndex == player {
			l.l <- EventBlindHandDiscard{EventBase: EventTypeBlindHandDiscard, Player: player, Index: index, HandSize: len(s.playerHands[player])}
		} else {
			handCopy := make([]Card, len(s.playerHands[player]))
			copy(handCopy, s.playerHands[player])
			l.l <- EventHandDiscard{EventBase: EventTypeHandDiscard, Player: player, Index: index, Hand: handCopy}
		}
		discardCopy := discard
		l.l <- EventFail{EventBase: EventTypeFail, Discard: &discardCopy, NumFails: s.fails}
	}

	// Draw new
	if len(s.deck) != 0 {
		s.playerHands[player] = append(s.playerHands[player], s.drawCard())
		for _, l := range s.listeners {
			if l.playerIndex == player {
				handCopy := make([]bool, len(s.playerHands[player]))
				for i := range s.playerHands[player] {
					handCopy[i] = s.playerHands[player][i].Tilt
				}
				l.l <- EventBlindHandDraw{EventBase: EventTypeBlindHandDraw, Player: player, Hand: handCopy}
			} else {
				handCopy := make([]Card, len(s.playerHands[player]))
				copy(handCopy, s.playerHands[player])
				l.l <- EventHandDraw{EventBase: EventTypeHandDraw, Player: player, Hand: handCopy}
			}
			l.l <- EventDeckDraw{EventBase: EventTypeDeckDraw, DeckRemaining: len(s.deck)}
		}
	}

	// Advance player and detect game end
	s.nextPlayer()
	if len(s.deck) == 0 && s.lastPlayer == -1 {
		s.lastPlayer = player
	}
	for _, l := range s.listeners {
		l.l <- EventTurnChange{EventBase: EventTypeTurnChange, NextPlayer: s.curPlayer}
	}
}

// DoPlay try to play card with index as player
func (s *GameState) DoPlay(player, index int) {
	s.Lock()
	defer s.Unlock()

	if s.curPlayer != player {
		return
	}

	if index < 0 || index >= len(s.playerHands[player]) {
		return
	}

	// Take card out of hand
	toPlay := s.playerHands[player][index]
	//fmt.Println("before: ", s.playerHands[player])
	copy(s.playerHands[player][index:], s.playerHands[player][index+1:])
	//fmt.Println("inbetween: ", s.playerHands[player])
	s.playerHands[player] = s.playerHands[player][:len(s.playerHands[player])-1]
	for _, l := range s.listeners {
		if l.playerIndex == player {
			l.l <- EventBlindHandDiscard{EventBase: EventTypeBlindHandDiscard, Player: player, Index: index, HandSize: len(s.playerHands[player])}
		} else {
			handCopy := make([]Card, len(s.playerHands[player]))
			copy(handCopy, s.playerHands[player])
			l.l <- EventHandDiscard{EventBase: EventTypeHandDiscard, Player: player, Index: index, Hand: handCopy}
		}
	}

	// And draw new one
	if len(s.deck) != 0 {
		s.playerHands[player] = append(s.playerHands[player], s.drawCard())
		for _, l := range s.listeners {
			if l.playerIndex == player {
				handCopy := make([]bool, len(s.playerHands[player]))
				for i := range s.playerHands[player] {
					handCopy[i] = s.playerHands[player][i].Tilt
				}
				l.l <- EventBlindHandDraw{EventBase: EventTypeBlindHandDraw, Player: player, Hand: handCopy}
			} else {
				handCopy := make([]Card, len(s.playerHands[player]))
				copy(handCopy, s.playerHands[player])
				l.l <- EventHandDraw{EventBase: EventTypeHandDraw, Player: player, Hand: handCopy}
			}
			l.l <- EventDeckDraw{EventBase: EventTypeDeckDraw, DeckRemaining: len(s.deck)}
		}
	}

	// Check whether it works
	if toPlay.Number == s.fieldStatus[toPlay.Color]+1 {
		s.fieldStatus[toPlay.Color] = s.fieldStatus[toPlay.Color] + 1
		for _, l := range s.listeners {
			fieldCopy := make([]int, 5)
			copy(fieldCopy, s.fieldStatus[:])
			l.l <- EventPlay{EventBase: EventTypePlay, Show: fieldCopy}
		}
	} else {
		s.fails++
		s.discard = &toPlay
		for _, l := range s.listeners {
			discardCopy := toPlay
			l.l <- EventFail{EventBase: EventTypeFail, Discard: &discardCopy, NumFails: s.fails}
		}
	}

	// And update turn
	s.nextPlayer()
	if s.fails == 3 {
		s.curPlayer = -1
	}
	if len(s.deck) == 0 && s.lastPlayer == -1 {
		s.lastPlayer = player
	}
	for _, l := range s.listeners {
		l.l <- EventTurnChange{EventBase: EventTypeTurnChange, NextPlayer: s.curPlayer}
	}
}

// DoTilt handle player tilting card
func (s *GameState) DoTilt(player, index int, tilt bool) {
	s.Lock()
	defer s.Unlock()

	if player < 0 || player >= s.nPlayers {
		return
	}

	if index < 0 || index >= len(s.playerHands[player]) {
		return
	}

	s.playerHands[player][index].Tilt = tilt

	for _, l := range s.listeners {
		l.l <- EventTilt{EventBase: EventTypeTilt, Player: player, Index: index, Tilt: tilt}
	}
}

// DoMove handle player moving a card
func (s *GameState) DoMove(player, from, to int) {
	s.Lock()
	defer s.Unlock()

	if player < 0 || player > s.nPlayers {
		return
	}

	if from < 0 || from >= len(s.playerHands[player]) || to < 0 || to >= len(s.playerHands[player]) {
		return
	}

	inter := []Card{}
	inter = append(inter, s.playerHands[player][:from]...)
	inter = append(inter, s.playerHands[player][from+1:]...)

	newHand := []Card{}
	newHand = append(newHand, inter[:to]...)
	newHand = append(newHand, s.playerHands[player][from])
	newHand = append(newHand, inter[to:]...)

	s.playerHands[player] = newHand

	for _, l := range s.listeners {
		l.l <- EventMove{EventBase: EventTypeMove, Player: player, From: from, To: to}
	}
}

// SetName set player name
func (s *GameState) SetName(player int, name string) {
	s.Lock()
	defer s.Unlock()

	if player < 0 || player >= s.nPlayers {
		return
	}

	s.playerNames[player] = name

	for _, l := range s.listeners {
		namesCopy := make([]string, len(s.playerNames))
		copy(namesCopy, s.playerNames)
		l.l <- EventSetNames{EventBase: EventTypeSetNames, PlayerNames: namesCopy}
	}
}
