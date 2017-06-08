import React from 'react';

export default ({ predictions, image }) => {

  return (
    <div className="results">
      <img
        className="example"
        alt="Compressed character"
        src={`data:image/png;base64,${image}`}
      />
      {
        predictions.map(({ charcode, confidence }, i) => (
          <div key={i}>
            <h5>Character: {String.fromCharCode(charcode)}</h5>
            <h6>Confidence: {Number(confidence * 100).toFixed(2)}%</h6>
          </div>
        ))
      }
    </div>
  )
}
