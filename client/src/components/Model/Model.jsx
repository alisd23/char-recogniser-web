import React, { Component } from 'react';
import './Model.scss';
import archImage from '../../../assets/convnet-architecture.jpg';

export default class Model extends Component {
  state = {
    model: null
  }

  componentDidMount() {
    fetch('/api/model', {
      mode: 'cors'
    })
      .then(r => r.json())
      .then(data => this.setState({
        model: {
          filters: data.filters,
          top1: data.top1,
          top3: data.top3
        }
      }));
  }

  render() {
    const { model } = this.state;

    return (
      <div className="model-page">
        {
          !model && (
            <div className="loading">
              <div className="spinner" />
            </div>
          )
        }
        {
          model && (
            <div>
              <div className="accuracy-section">
                <h3 className="section-header">
                  Accuracy
                  <span>(On test set of size <strong>140,000</strong>)</span>
                </h3>
                <div className="accuracy-wrapper">
                  <div className="accuracy-panel top1">
                    <h4 className="accuracy-panel-title">Top 1</h4>
                    <span className="accuracy-panel-score">
                      {Number(model.top1 * 100).toFixed(2)}%
                    </span>
                  </div>
                  <div className="accuracy-panel top3">
                    <h4 className="accuracy-panel-title">Top 3</h4>
                    <span className="accuracy-panel-score">
                      {Number(model.top3 * 100).toFixed(2)}%
                    </span>
                  </div>
                </div>
              </div>
              <div className="architecture-section">
                <h3 className="section-header">Network Architecture</h3>
                <img
                  src={archImage}
                  alt="ConvNet Architecture"
                />
              </div>
              <div className="filters-section">
                <h3 className="section-header">Convolutional layer 1 filters</h3>
                <div className="filters">
                  {
                    model.filters.map((image, i) => (
                      <img
                        key={i}
                        className="filter"
                        alt={`Conv layer 1 filter ${i + 1}`}
                        src={`data:image/png;base64,${image}`}
                      />
                    ))
                  }
                </div>
              </div>
            </div>
          )
        }
      </div>
    )
  }
}
