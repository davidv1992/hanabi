package main

import "math/rand"

// Card card data
type Card struct {
	Color  int  `json:"color"`
	Number int  `json:"number"`
	Tilt   bool `json:"tilt"`
}

func randomDeck() []Card {
	// First build deck
	deck := []Card{}
	for i := 0; i < 5; i++ {
		for j := 1; j <= 5; j++ {
			deck = append(deck, Card{Color: i, Number: j})
			if j != 5 {
				deck = append(deck, Card{Color: i, Number: j})
			}
			if j == 1 {
				deck = append(deck, Card{Color: i, Number: j})
			}
		}
	}

	rand.Shuffle(len(deck), func(i, j int) { t := deck[i]; deck[i] = deck[j]; deck[j] = t })

	return deck
}
