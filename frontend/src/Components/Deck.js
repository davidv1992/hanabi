import React from 'react';
import './Deck.css';

class Deck extends React.Component {
	render() {
		return (
			<div className="Deck">
				<div>{this.props.count}</div>
			</div>
		)
	}
}

export default Deck;
