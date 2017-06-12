import React, { Component } from 'react';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group';
import DigitCanvas from '../DigitCanvas';
import Predictions from '../Predictions';
import responsive from '../../utils/responsive';
import './Predict.scss';

class Predict extends Component {
  canvasComponent = null;

  onClear = () => {
    this.canvasComponent && this.canvasComponent.clear();
  }

  onSubmit = () => {
    this.props.onSubmit(this.canvasComponent.canvas.toDataURL("image/png"));
  }

  componentDidMount() {
    // Force component to render - to use container width
    this.setState({});
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
        <Predictions
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
    const canvasSize = Math.min(
      400,
      (this.wrapper && this.wrapper.offsetWidth) || Infinity
    );

    return (
      <div
        className="predict-page"
        ref={el => (this.wrapper = el)}
      >
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
            <div
              className="canvas"
              style={{
                // Canvas border width = 2
                width: canvasSize + 4,
                height: canvasSize + 4
              }}
            >
              <DigitCanvas
                ref={el => (this.canvasComponent = el)}
                penColour="black"
                penRadius={6}
                size={canvasSize}
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

export default responsive(Predict)
