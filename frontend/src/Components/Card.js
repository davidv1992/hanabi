import React from 'react';
import './Card.css';

const colorMap = {
	"R": "#FF0000",
	"Y": "#FFFF00",
	"G": "#00FF00",
	"B": "#008FFF",
	"W": "#FFFFFF",
	"?": "#000000",
}

class Card extends React.Component {
	render() {
		var transform = "";
		transform = "translate(calc("+((2*this.props.position-this.props.total+1)*2.871)+"vw - 50%), 0)";
		if (this.props.tilt)
			transform += " rotate(10deg)"
		if (this.props.lrmargin)
		  transform = "";
		const className="Card"+(this.props.select?" Card-Selected":"")+(this.props.lrmargin?" Card-Margin":"");
		var cardImg = '' + this.props.color + this.props.number + '.png'
		if (this.props.color === '?')
			cardImg = 'Back.png';
		return (
			<div className={className} style={{transform, color: colorMap[this.props.color], backgroundColor: "#000000"}} onClick={this.props.onClick}>
				<div className="Color">
					{this.props.number}
				</div>
				<div className="Color">
					{this.props.color}
				</div>
			</div>
		);
	}
}

export default Card;
