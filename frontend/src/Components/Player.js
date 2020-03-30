import React from 'react';
import './Player.css';
import Card from './Card.js';
import Button from './Button.js';
import Position from './Position.js';
import {connect} from 'react-redux';

function mapStateToProps(state, props) {
	return {
		myturn: state.players.current == props.playerIndex,
		havehints: state.table.hints > 0,
		name: state.players.names[props.playerIndex] + ((props.playerIndex == state.players.current)?" (turn)":""),
		cards: state.hands.cards[props.playerIndex] ?? [],
	}
}

class Player extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			selected: -1,
		}
	}

	onClickCard(idx) {
		return () => {
			if (!this.props.localPlayer) return;
			if (this.state.selected == idx)
				this.setState({selected:-1});
			else
				this.setState({selected: idx});
		}
	}

	onMove(idx) {
		return () => {
			var target = idx
			if (target > this.state.selected)
				target--;
			this.props.connection.send(JSON.stringify({
				type:"move",
				from: this.state.selected,
				to: target,
			}));
			this.setState({selected:-1});
		}
	}

	render() {
		var cardList = []
		this.props.cards.forEach((a)=>{cardList.push(a)})
		const cards = cardList.sort((a,b)=>(a.ourIndex-b.ourIndex)).map(card => (
			<Card
				number={card.number}
				color={card.color}
				tilt={card.tilt}
				select={this.state.selected == card.serverIndex}
				key={card.ourIndex} position={card.serverIndex}
				total={this.props.cards.length} 
				onClick={this.onClickCard(card.serverIndex)}/>
		))
		
		var positions = [];
		
		var buttons = null;
		if (this.props.localPlayer) {
			var buttonList = []
			if (this.props.myturn && this.props.havehints) {
				buttonList.push(
					<Button onClick={()=>{this.props.connection.send(JSON.stringify({type:"hint"}));}}>
						Hint
					</Button>
				);
			}
			if (this.props.myturn && this.state.selected != -1) {
				buttonList.push(
					<Button onClick={()=>{
						this.props.connection.send(JSON.stringify({type:"discard", index: this.state.selected}));
						this.setState({selected:-1});
					}}>
						Discard
					</Button>
				);
				buttonList.push(
					<Button onClick={()=>{
						this.props.connection.send(JSON.stringify({type:"play", index: this.state.selected}));
						this.setState({selected:-1});
					}}>
						Play
					</Button>
				);
			}
			if (this.state.selected != -1) {
				buttonList.push(
					<Button onClick={()=>{
						this.props.connection.send(JSON.stringify({type:"tilt", index: this.state.selected, tilt: !this.props.cards[this.state.selected].tilt}));
					}}>
						Tilt
					</Button>
				);

				// Setup positions
				for (var i=0; i<=this.props.cards.length; i++) {
					if (i != this.state.selected && i != this.state.selected+1)
					positions.push(<Position key={"pos"+i} position={i} total={this.props.cards.length} onClick={this.onMove(i)}/>);
				}
			}
			buttons = (<div className="Button-Box">{buttonList}</div>);
		}
		
		return (
			<div className={"Player-Box "+this.props.position}>
				<div className="Name-Box">
					{this.props.name}
				</div>
				<div className="Card-Box">
					{cards}
					{positions}
				</div>
				{buttons}
			</div>
		);
	}
}

export default connect(mapStateToProps)(Player);
