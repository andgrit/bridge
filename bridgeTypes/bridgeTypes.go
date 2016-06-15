package bridgeTypes

import "math/rand"
import "gopkg.in/mgo.v2/bson"

const CARDS_IN_DECK = 52

type Card int // 0..52, 0 no card, 1..13 - Spade, 14..26 heart, 27..39 diamon, 40..52 club, 1-Ace Spade, 2-2 Spade, 52 K club
type Deck []Card // indexed by 0..51 - deck of cards
type Position int // 0..3, 0-North, 1-East, 2-South, 3-West, 0,2 against 1,3
type Suit int // 0..3 0 spades, 1 hearts, 2 diamonds, 3 clubs
type Level int // 1 7, 2 8, 3 9, 4 10, 5 11, 6 12, 7 slam, 0 - no bid level
type Bid struct {
	Position Position
	Suit     Suit
	Level    Level // 0 is the zero value
	Players  []Player // indexed by position
}
type Trick struct{
	Cards [4]Card // indexed by Position, value of 0 means not played
	Players	[]Player
}
type Deal struct {
	Deck    Deck // deck[0] == 0 is the zero value
	Bidding []Bid
	Play    []Trick
}

type Match struct {
	Deals []Deal
}

/// Table table
type Table struct {
	OId    bson.ObjectId `bson:"_id,omitempty"`
	Version int
	Match Match
	Players []bson.ObjectId // collection of 4 players
}

/// player table
type Player struct {
	OId    bson.ObjectId `bson:"_id,omitempty"`
	Username string
}

// dealt returns true of the cards have been dealt
func (deal *Deal) dealt() bool {
	return deal.Deck[0] != Card(0)
}

// randomDeal returns a random deal of the cards
func randomDeal() Deck {
	var deck Deck
	// create the deck
	for i := Card(0); i < CARDS_IN_DECK; i++ {
		deck[i] = i
	}
	// shuffle the deck
	for i := CARDS_IN_DECK - 1; i >= 0; i-- {
		cardI := rand.Intn(CARDS_IN_DECK - i) // choose a random card to swap with the last card
		deck[i], deck[cardI] = deck[cardI], deck[i]
	}
	return deck
}
