const initialState = {
    n: 0,
    me: -1,
    current: -1,
    names: [],
};

export default function(state = initialState, action) {
    switch(action.type) {
        case "setup":
            return Object.assign({}, state,{
                n: action.n_players,
                me: action.your_player,
                names: Array(action.n_players).fill(""),
            });
        case "set_names":
            return Object.assign([], state, {
                names: action.player_names,
            });
        case "turn_change":
            return Object.assign([], state, {
                current: action.next_player,
            })
        default:
            return state;
    }
}