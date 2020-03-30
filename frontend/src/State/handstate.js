import colormap from "./cardColorMap";

const initialState = {
    cards: [],
    playerIndex: [],
};

export default function(state = initialState, action) {
    switch(action.type) {
        case "setup":
            return Object.assign({}, state, {
                cards: Array(action.n_players).fill([]),
                playerIndex: Array(action.n_players).fill(0),
            });
        case "hand_draw": {
            var newHand = [];
            state.cards[action.player].forEach((card)=>{
                newHand.push(card);
            });
            var newIndex = state.playerIndex[action.player];
            for (var i = newHand.length; i<action.hand.length; i++) {
                newHand.push({
                    color: colormap[action.hand[i].color],
                    number: action.hand[i].number,
                    tilt: action.hand[i].tilt,
                    serverIndex: i,
                    ourIndex: newIndex,
                });
                newIndex++;
            }
            var newHands = []
            var newIndices = []
            for (var i = 0; i<state.cards.length; i++) {
                if (i == action.player) {
                    newHands.push(newHand);
                    newIndices.push(newIndex);
                } else {
                    newHands.push(state.cards[i]);
                    newIndices.push(state.playerIndex[i]);
                }
            }
            return Object.assign({}, state, {
                cards: newHands,
                playerIndex: newIndices,
            });
        }
        case "blind_hand_draw": {
            var newHand = [];
            state.cards[action.player].forEach((card)=>{
                newHand.push(card);
            });
            var newIndex = state.playerIndex[action.player];
            for (var i = newHand.length; i<action.hand.length; i++) {
                newHand.push({
                    color: "?",
                    tilt: action.hand[i],
                    serverIndex: i,
                    ourIndex: newIndex,
                });
                newIndex++;
            }
            var newHands = []
            var newIndices = []
            for (var i = 0; i<state.cards.length; i++) {
                if (i == action.player) {
                    newHands.push(newHand);
                    newIndices.push(newIndex);
                } else {
                    newHands.push(state.cards[i]);
                    newIndices.push(state.playerIndex[i]);
                }
            }
            return Object.assign({}, state, {
                cards: newHands,
                playerIndex: newIndices,
            });
        }
        case "hand_discard":
        case "blind_hand_discard":
            var newHand = []
            for (var i=0; i<state.cards[action.player].length; i++) {
                if (i != action.index) {
                    newHand.push(Object.assign({}, state.cards[action.player][i],{
                        serverIndex: newHand.length,
                    }));
                }
            }
            var newHands = []
            for (var i = 0; i<state.cards.length; i++) {
                if (i == action.player) {
                    newHands.push(newHand);
                } else {
                    newHands.push(state.cards[i]);
                }
            }
            return Object.assign({}, state, {
                cards: newHands,
            });
        case "tilt":
            var newHand = []
            for (var i=0; i<state.cards[action.player].length; i++) {
                if (i == action.index) {
                    newHand.push(Object.assign({}, state.cards[action.player][i], {
                        tilt: action.tilt,
                    }))
                } else {
                    newHand.push(state.cards[action.player][i])
                }
            }
            var newHands = []
            for (var i = 0; i<state.cards.length; i++) {
                if (i == action.player) {
                    newHands.push(newHand);
                } else {
                    newHands.push(state.cards[i]);
                }
            }
            return Object.assign({}, state, {
                cards: newHands,
            });
        case "move":
            var inter = []
            for (var i=0; i<state.cards[action.player].length; i++) {
                if (i != action.from)
                    inter.push(state.cards[action.player][i]);
            }
            var newHand = []
            var j = 0;
            for (var i=0; i<state.cards[action.player].length; i++) {
                if (i == action.to) {
                    newHand.push(Object.assign({}, state.cards[action.player][action.from],{
                        serverIndex: i,
                    }));
                } else {
                    newHand.push(Object.assign({}, inter[j],{
                        serverIndex: i,
                    }));
                    j++;
                }
            }
            var newHands = []
            for (var i = 0; i<state.cards.length; i++) {
                if (i == action.player) {
                    newHands.push(newHand);
                } else {
                    newHands.push(state.cards[i]);
                }
            }
            return Object.assign({}, state, {
                cards: newHands,
            });
        default:
            return state;
    }
};