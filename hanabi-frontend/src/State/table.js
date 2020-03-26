import colorMap from './cardColorMap';

const initialState = {
    hints: 0,
    fails: 0,
    deck: 0,
    discard: null,
    show: [0,0,0,0,0],
};

export default function(state = initialState, action) {
    switch(action.type) {
        case "hint":
            return Object.assign({}, state, {
                hints: action.remaining_hints,
            });
        case "fail": {
            var discard = null;
            if (action.discard != null) {
                discard = {
                    color: colorMap[action.discard.color],
                    number: action.discard.number,
                };
            }
            return Object.assign({}, state, {
                fails: action.n_fails,
                discard: discard,
            });
        }
        case "deck_draw":
            if (state.deck != 0 && action.deck_remaining == 0)
            	window.alert("Laatste ronde!");
            return Object.assign({}, state, {
                deck: action.deck_remaining,
            });
        case "play":
            return Object.assign({}, state, {
                show: action.show,
            })
        default:
            return state;
    }
}
