import React from 'react';
import {connect} from 'react-redux';
import Player from "./Components/Player.js";
import Table from "./Components/Table.js";

function mapStateToProps(state) {
	const n = state.players.n;
	const me = state.players.me;

	const playersBeforeMe = Math.min(n-1, 3);

	return {
		n: n,
		me: me,
		tli: (n>1?((n+me-playersBeforeMe)%n):undefined),
		tri: (n>2?((n+me-playersBeforeMe+1)%n):undefined),
		bri: (n>3?((n+me-playersBeforeMe+2)%n):undefined),
		bli: (n>4?((me+1)%n):undefined),
	}
}

class App extends React.Component {
  render() {
	  	if (this.props.n == 0) return <div/>; // Don't render while loading to much
		return (
		<div>
			{(this.props.tli != undefined)?<Player playerIndex={this.props.tli} position="Player-Top-Left" />:null}
			{(this.props.tri != undefined)?<Player playerIndex={this.props.tri} position="Player-Top-Right" />:null}
			{(this.props.bli != undefined)?<Player playerIndex={this.props.bli} position="Player-Bottom-Left" />:null}
			{(this.props.bri != undefined)?<Player playerIndex={this.props.bri} position="Player-Bottom-Right" />:null}
			<Player playerIndex={this.props.me} position="Player-Center" localPlayer={true} connection={this.props.connection}/>
			<Table />
		</div>
		);
	}
}

export default connect(mapStateToProps)(App);
