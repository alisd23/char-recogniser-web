import React, { Component } from 'react';
import classnames from 'classnames';
import Predict from '../Predict';
import Navbar from '../Navbar';
import './App.scss';

const PAGES = {
  PREDICT: 1
};

class App extends Component {
  state = {
    page: PAGES.PREDICT
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
      default:
        return null;
    }
  }

  render() {
    return (
      <div className="app">
        <Navbar>
          {this.getPageLinkFragment(PAGES.PREDICT, 'Predict')}
        </Navbar>
        <div className="app-main">
          {this.getPageFragment()}
        </div>
      </div>
    );
  }
}

export default App;
