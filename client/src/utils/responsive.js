import React, { Component } from 'react';

export default (Comp) =>
  class Responsive extends Component {
    state = {
      width: window.innerWidth,
      height: window.innerHeight
    }

    onResize = (e) => this.setState({
      width: window.innerWidth,
      height: window.innerHeight
    })

    componentDidMount() {
      window.addEventListener('resize', this.onResize);
    }
    componentWillUnmount() {
      window.removeEventListener('resize', this.onResize);
    }

    render() {
      return (
        <Comp
          {...this.props}
          width={this.state.width}
          height={this.state.height}
        />
      );
    }
  }
