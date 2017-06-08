import React, { Component } from 'react';
import DigitCanvas from '../DigitCanvas';
import Results from './Results';
import './Predict.scss';

export default class Predict extends Component {
  canvasComponent = null;

  state = {
    results: null,
    error: false
  }

  onClear = () => {
    this.canvasComponent && this.canvasComponent.clear();
  }

  onSubmit = () => {
    this.setState({
      error: false
    })
    fetch('/api/predict', {
      method: 'POST',
      mode: 'cors',
      body: JSON.stringify({
        image: this.canvasComponent.canvas.toDataURL("image/png")
      })
    })
      .then(res => res.json())
      .then(res => {
        if (res.error) {
          this.setState({
            error: true
          })
        }
        this.setState({
          results: {
            predictions: res.predictions,
            image: res.image
          }
        })
      });
  }

  getResultsFragment = () => {
    const { error, results } = this.state;

    if (error) {
      return (
        <div className='results-error'>
          An error occurred <i class="material-icons md-32">sentiment_very_dissatisfied</i>
        </div>
      )
    }

    if (results) {
      return (
        <Results
          image={this.state.results.image}
          predictions={this.state.results.predictions}
        />
      )
    } else {
      return (
        <div className="no-results">
          <i className="material-icons md-48 gesture">gesture</i>
          <span>
            Draw a character (left) and press <strong>submit</strong> to see results
          </span>
        </div>
      )
    }
  }

  render() {
    return (
      <div className="predict-page">
        <div className="canvas-sidebar">
          <h4 className="sidebar-title">Predict</h4>
          <button
            className="btn-red sidebar-row"
            onClick={this.onClear}
          >
            Clear
          </button>
          <div className="sidebar-separator" />
          <button
            className="btn-green sidebar-row"
            onClick={this.onSubmit}
          >
            Submit
          </button>
        </div>
        <div className="canvas-wrapper">
          <DigitCanvas
            ref={el => (this.canvasComponent = el)}
            penColour="black"
            penRadius={8}
            size={400}
          />
        </div>
        <div className="results-panel">
          {this.getResultsFragment()}
        </div>
      </div>
    );
  }
}
