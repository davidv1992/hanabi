import React from "react";
import "./Position.css";

class PositionIndicator extends React.Component {
	render() {
		const transform = "translate(calc("+((2*this.props.position-this.props.total)*2.871)+"vw - 50%), 0)";
		return (
			<div className="Position-Arrow" style={{transform}} onClick={this.props.onClick}/>
		);
	}
}

export default PositionIndicator;
