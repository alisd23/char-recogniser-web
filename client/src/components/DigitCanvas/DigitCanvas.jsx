import React, { Component, PropTypes } from 'react';

class DigitCanvas extends Component {
  drawing = false;
  canvas = null;
  context = null;
  lastMousePos = null;

  componentDidMount() {
    this.context = this.canvas.getContext('2d');

    this.clear();

    // Add all event listeners
    Object
      .keys(this.events)
      .forEach(eventName => {
        this.canvas.addEventListener(eventName, this.events[eventName])
      });
  }

  componentWillUnmount() {
    if (this.canvas) {
      // Remove all event listeners
      Object
        .keys(this.events)
        .forEach(eventName => {
          this.canvas.removeEventListener(eventName, this.events[eventName])
        });
    }
  }

  draw = ({ x, y }) => {
    this.context.beginPath();

    // Set line colour and width
    this.context.fillStyle = this.props.penColour;
    this.context.strokeStyle = this.props.penColour;
    this.context.lineWidth = this.props.penRadius * 2;

    // Draw line from last recorded mouse position to ensure no gaps in the line
    if (this.lastMousePos) {
      this.context.moveTo(this.lastMousePos.x, this.lastMousePos.y);
      this.context.lineTo(x, y);
      this.context.closePath();
      this.context.stroke();
    }
    // Draw circle at current mouse position to round off any lines drawn
    this.context.arc(x, y, this.props.penRadius, 0, 2 * Math.PI);
    this.context.fill();
    this.lastMousePos = { x, y };
  }

  getMousePos = (e) => {
    const rect = this.canvas.getBoundingClientRect();
    return {
      x: e.clientX - rect.left,
      y: e.clientY - rect.top
    };
  }

  // Specific TOUCH handlers - calls the draw methods below
  onTouchStart = (e) => {
    if (e.touches.length) {
      this.onDrawStart(e.touches[0]);
      // Stop scrolling
      document.querySelector('.app-main').style.overflow = 'hidden';
    }
  }
  onTouchMove = (e) => {
    if (e.touches.length) {
      this.onDraw(e.touches[0]);
    }
  }
  onTouchEnd = (e) => {
    this.onDrawStop();
    document.querySelector('.app-main').style.overflow = 'auto';
  }

  // DRAW handlers (start, move, stop)
  onDrawStart = (e) => {
    this.drawing = true;
    this.draw(this.getMousePos(e));
  }
  onDraw = (e) => {
    if (this.drawing) {
      this.draw(this.getMousePos(e));
    }
  }
  onDrawStop = () => {
    this.drawing = false;
    this.lastMousePos = null;
  }

  clear = () => {
    this.context.fillStyle = '#ffffff';
    this.context.rect(0, 0, this.props.size, this.props.size);
    this.context.fill();
  }

  events = {
    // Mouse events
    mousedown: this.onDrawStart,
    mouseout: this.onDrawStop,
    mouseup: this.onDrawStop,
    mouseleave: this.onDrawStop,
    mousemove: this.onDraw,
    // Touch events
    touchstart: this.onTouchStart,
    touchmove: this.onTouchMove,
    touchend: this.onTouchEnd,
    touchcancel: this.onTouchEnd,
  };

  render() {
    return (
      <canvas
        ref={el => (this.canvas = el)}
        width={this.props.size}
        height={this.props.size}
      />
    )
  }
}

DigitCanvas.propTypes = {
  penRadius: PropTypes.number.isRequired,
  size: PropTypes.number.isRequired,
  penColour: PropTypes.string.isRequired
}

export default DigitCanvas;
