import React, { Component } from 'react';
import uuid from 'uuid';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';
import DigitCanvas from '../DigitCanvas';
import Results from './Results';
import './Predict.scss';

export default class Predict extends Component {
  canvasComponent = null;

  state = {
    results: null,
    error: false,
    loading: false,
    requestID: null
  }

  onClear = () => {
    this.canvasComponent && this.canvasComponent.clear();
  }

  onSubmit = () => {
    const requestID = uuid.v4();
    this.setState({
      error: false,
      loading: true,
      requestID
    });
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
            error: true,
            loading: false
          })
        }
        // Only show results if this is the most recent request
        if (requestID === this.state.requestID) {
          console.log(res);
          this.setState({
            loading: false,
            results: {
              predictions: res.predictions,
              image: res.image
            }
          })
        }
      })
      .catch(err => this.setState({ error: true }));
  }

  getResultsFragment = () => {
    const { error, results, loading } = this.state;

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

    if (results) {
      return (
        <Results
          key="results"
          image={this.state.results.image}
          predictions={this.state.results.predictions}
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
                penColour="black"
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
