import React from 'react';
import './Card.css';

class Card extends React.Component {
	render() {
		var transform = "";
		transform = "translate(calc("+((2*this.props.position-this.props.total+1)*2.871)+"vw - 50%), 0)";
		if (this.props.tilt)
			transform += " rotate(10deg)"
		if (this.props.lrmargin)
		  transform = "";
		const className="Card"+(this.props.select?" Card-Selected":"")+(this.props.lrmargin?" Card-Margin":"");
		return (
			<div className={className} style={{transform}} onClick={this.props.onClick}>
				<div className="Number-Top-Left">{this.props.number}</div>
				<div className="Number-Top-Right">{this.props.number}</div>
				<div className="Number-Bottom-Left">{this.props.number}</div>
				<div className="Number-Bottom-Right">{this.props.number}</div>
				<div className="Color">{this.props.color}</div>
				<div className="Color">{this.props.color}</div>
			</div>
		);
	}
}

export default Card;
