import React, { Component } from 'react';
import classnames from 'classnames';
import Predict from '../Predict';
import Model from '../Model';
import Navbar from '../Navbar';
import './App.scss';

const PAGES = {
  PREDICT: 1,
  MODEL: 2
};

class App extends Component {
  state = {
    page: PAGES.PREDICT,
    filters: []
  }

  getPageLinkFragment = (page, text) => {
    const classes = classnames({
      'page-link': true,
      active: page === this.state.page
    });

    return (
      <div
        className={classes}
        onClick={() => this.setState({ page })}
      >
        {text}
      </div>
    )
  }

  getPageFragment = () => {
    switch (this.state.page) {
      case PAGES.PREDICT: {
        return (
          <Predict />
        )
      }
      case PAGES.MODEL: {
        return (
          <Model />
        )
      }
      default:
        return null;
    }
  }

  render() {
    return (
      <div className="app">
        <Navbar>
          {this.getPageLinkFragment(PAGES.PREDICT, 'Predict')}
          {this.getPageLinkFragment(PAGES.MODEL, 'Model')}
        </Navbar>
        <div className="app-main">
          {this.getPageFragment()}
        </div>
      </div>
    );
  }
}

export default App;
