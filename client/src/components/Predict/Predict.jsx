import React, { Component } from 'react';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';
import DigitCanvas from '../DigitCanvas';
import Results from './Results';
import './Predict.scss';

export default class Predict extends Component {
  canvasComponent = null;

  onClear = () => {
    this.canvasComponent && this.canvasComponent.clear();
  }

  onSubmit = () => {
    this.props.onSubmit(this.canvasComponent.canvas.toDataURL("image/png"));
  }

  getResultsFragment = () => {
    const { error, data, loading } = this.props;

    if (error) {
      return (
        <div
          key="error"
          className="results-error"
        >
          <p>An error occurred whilst fetching the predictions</p>
          <i className="material-icons md-36">sentiment_very_dissatisfied</i>
        </div>
      )
    }

    if (loading) {
      return (
        <div
          key="loading"
          className="results-loading"
        >
          <div className="spinner" />
          <span>Fetching predictions...</span>
        </div>
      )
    }

    if (data) {
      return (
        <Results
          key="results"
          image={data.image}
          predictions={data.predictions}
        />
      )
    } else {
      return (
        <div
          key="no-results"
          className="no-results"
        >
          <i className="material-icons md-48 gesture">gesture</i>
          <div>
            <p>Draw a character on the canvas and press <strong>submit</strong> to see results.</p>
            <p>Draw in the center of the canvas and not too big.</p>
          </div>
        </div>
      )
    }
  }

  render() {
    return (
      <div className="predict-page">
        <div className="predict-page-inner">
          <div className="canvas-wrapper">
            <div className="canvas-toolbar">
              <i
                className="clear material-icons"
                onClick={this.onClear}
              >
                clear
              </i>
              <button
                className="btn-green submit"
                onClick={this.onSubmit}
              >
                submit
              </button>
            </div>
            <div className="canvas">
              <DigitCanvas
                ref={el => (this.canvasComponent = el)}
                penColour="black "
                penRadius={6}
                size={400}
              />
            </div>
          </div>
          <ReactCSSTransitionGroup
            className="results-wrapper"
            transitionName="fade"
            transitionEnterTimeout={250}
            transitionLeaveTimeout={250}
          >
            {this.getResultsFragment()}
          </ReactCSSTransitionGroup>
        </div>
      </div>
    );
  }
}
