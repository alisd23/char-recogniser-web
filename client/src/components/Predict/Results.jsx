import React, { Component } from 'react';

export default class Results extends Component {
  state = {
    showMore: false
  }

  // Reset show more button when new image is predicted
  componentWillReceiveProps(nextProps) {
    if (nextProps.image !== this.props.image) {
      this.setState({
        showMore: false
      });
    }
  }

  onShowMore = () => this.setState({
    showMore: true
  })

  getTop3 = () => {
    return (
      this.props.predictions
        .slice(0, 3)
        .map(({ charcode, confidence }, i) => (
          <div
            className="prediction-top"
            key={i}
          >
            <span className="index">{i + 1})</span>
            <div className="character-wrapper">
              <span>Character (code: <strong>{charcode}</strong>)</span>
              <span><span>{String.fromCharCode(charcode)}</span></span>
            </div>
            <div className="confidence-wrapper">
              <span>Confidence</span>
              <span><span>{Number(confidence * 100).toFixed(2)}%</span></span>
            </div>
          </div>
        ))
    );
  }

  getRest = () => {
    return (
      this.props.predictions
        .slice(3)
        .map(({ charcode, confidence }, i) => (
          <div
            className="prediction-rest"
            key={i}
          >
            <span className="index">{i + 4})&nbsp;</span>
            <div className="character-wrapper">
              <div className="keys">
                <span>Character</span>
                <span>Character Code</span>
              </div>
              <div className="values">
                <span>{String.fromCharCode(charcode)}</span>
                <span>{charcode}</span>
              </div>
            </div>
            <div className="confidence-wrapper">
              Confidence: <strong>{Number(confidence * 100).toFixed(2)}%</strong>
            </div>
          </div>
        ))
    );
  }

  render() {
    const { image } = this.props;

    return (
      <div className="results">
        <div className="results-header">
          <img
            className="example"
            alt="Compressed character"
            src={`data:image/png;base64,${image}`}
          />
          <h2>Predictions</h2>
        </div>
        <div className="predictions">
          {this.getTop3()}
          {
            !this.state.showMore && (
              <div
                className="show-more"
                onClick={this.onShowMore}
              >
                Show more <i className="material-icons md-32">keyboard_arrow_down</i>
              </div>
            )
          }
          {this.state.showMore && this.getRest()}
        </div>
      </div>
    );
  }
}
