import React, { Component } from 'react';
import DigitCanvas from '../DigitCanvas';
import './Predict.scss';

export default class Predict extends Component {
  canvasComponent = null;

  onClear = () => {
    this.canvasComponent && this.canvasComponent.clear();
  }

  onSubmit = () => {
    fetch('/api/predict', {
      method: 'POST',
      mode: 'cors',
      body: JSON.stringify({
        image: this.canvasComponent.canvas.toDataURL("image/png")
      })
    })
      .then(res => res.json())
      .then(res => {
        console.log(res);
      });
  }

  render () {
    return (
      <div className="canvas-page predict-page">
        <div className="canvas-sidebar">
          <h4 className="sidebar-title">Predict</h4>
          <button
            className="btn-red sidebar-row"
            onClick={this.onClear}
          >
            Clear Canvas
          </button>
          <div className="sidebar-separator" />
          <button
            className="btn-green sidebar-row"
            onClick={this.onSubmit}
          >
            Submit
          </button>
        </div>
        <DigitCanvas
          ref={el => (this.canvasComponent = el)}
          penColour="black"
          penRadius={10}
          size={400}
        />
      </div>
    );
  }
}
