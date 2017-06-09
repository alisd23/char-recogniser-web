import React, { Component } from 'react';
import './Analysis.scss';

export default class Analysis extends Component {
  state = {
    filters: null
  }

  componentDidMount() {
    fetch('/api/filters', {
      mode: 'cors'
    })
      .then(r => r.json())
      .then(data => this.setState({ filters: data.filters }));
  }

  render() {
    const { filters } = this.state;

    return (
      <div className="analysis-page">
        {
          !filters && (
            <div className="loading">
              <div className="spinner" />
            </div>
          )
        }
        {
          filters && (
            <div className="filters-sections">
              <h3 className="section-header">Convolutional layer 1 filters</h3>
              <div className="filters">
                {
                  filters.map((image, i) => (
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
          )
        }
      </div>
    )
  }
}
