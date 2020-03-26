import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import state from './State';
import {createStore} from 'redux';
import {Provider} from 'react-redux';

var cookieString = document.cookie;

if (!cookieString.startsWith("token=")) {
    document.location.replace("/");
} else {
    var store = createStore(state);
    var conn = new WebSocket("ws:walkintrack.nl/connect");
    conn.onopen = ()=>{conn.send(cookieString.substr(6));};
    conn.onmessage = (message)=>{
        var action = JSON.parse(message.data);
        store.dispatch(action)
    };

    ReactDOM.render(<Provider store={store}><App connection={conn}/></Provider>, document.getElementById('root'));
}

