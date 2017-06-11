import React, { Component } from 'react';
import classnames from 'classnames';
import uuid from 'uuid';
import Predict from '../Predict';
import Model from '../Model';
import Analysis from '../Analysis';
import Navbar from '../Navbar';
import Page from './Page';
import './App.scss';

const PAGES = {
  PREDICT: 1,
  ANALYSIS: 2,
  MODEL: 3
};

class App extends Component {
  state = {
    page: PAGES.PREDICT,
    data: null,
    error: false,
    loading: false,
    requestID: null
  }

  getPageLinkFragment = (page, text, active = true) => {
    const classes = classnames({
      'page-link': true,
      active: page === this.state.page,
      disabled: !active
    });

    return (
      <div
        className={classes}
        onClick={() => active && this.setState({ page })}
      >
        {text}
      </div>
    )
  }

  onSubmit = (imageURL) => {
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
        image: imageURL
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
          this.setState({
            loading: false,
            data: {
              predictions: res.predictions,
              activations: res.activations,
              image: res.image
            }
          })
        }
      })
      .catch(err => this.setState({ error: true }));
    }

  getPageFragment = () => {
    const { page, error, data, loading } = this.state;

    return (
      <div className="app-main">
        <Page active={page === PAGES.PREDICT}>
          <Predict
            onSubmit={this.onSubmit}
            error={error}
            data={data}
            loading={loading}
          />
        </Page>
        <Page active={page === PAGES.MODEL}>
          <Model />
        </Page>
        <Page active={page === PAGES.ANALYSIS}>
          <Analysis activations={data && data.activations} />
        </Page>
      </div>
    );
  }

  render() {
    return (
      <div className="app">
        <Navbar>
          {this.getPageLinkFragment(PAGES.PREDICT, 'Predict')}
          {this.getPageLinkFragment(PAGES.ANALYSIS, 'Analysis', this.state.data)}
          {this.getPageLinkFragment(PAGES.MODEL, 'Model')}
        </Navbar>
        {this.getPageFragment()}
      </div>
    );
  }
}

export default App;
