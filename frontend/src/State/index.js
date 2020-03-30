import hands from "./handstate";
import players from "./player";
import table from "./table";

import { combineReducers } from 'redux';


export default combineReducers({
    hands,
    players,
    table,
});