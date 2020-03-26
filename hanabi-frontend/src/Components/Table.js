import React from "react";
import {connect} from 'react-redux';
import "./Table.css"
import Card from "./Card.js";
import EmptyPile from "./EmptyPile.js";
import Deck from "./Deck.js";
import colorMap from '../State/cardColorMap';

function mapStateToProps(state) {
	var showdata = []
	for (var i=0; i<5; i++) {
		showdata.push({
			color: colorMap[i],
			number: state.table.show[i],
		});
	}

	return {
		hints: state.table.hints,
		fails: state.table.fails,
		deck: state.table.deck,
		show: showdata,
		discard: state.table.discard,
	}
}

class Table extends React.Component {
	render() {
		return (
			<div className="Table">
				<div className="Table-Row">
					{this.props.discard?<Card color={this.props.discard.color} number={this.props.discard.number} lrmargin={true}/>:<EmptyPile/>}
					<Deck count={this.props.deck}/>
					<div className="Table-Column">
						<div>Hints: {this.props.hints}</div>
						<div>Fails: {this.props.fails}</div>
					</div>
				</div>
				<div className="Table-Row">
					{this.props.show.map(showpart=>{
						if (showpart.number == 0) {
							return (<EmptyPile />);
						} else {
							return (<Card color={showpart.color} number={showpart.number} lrmargin={true} />);
						}
					})}
				</div>
			</div>
		);
	}
}

export default connect(mapStateToProps)(Table);
