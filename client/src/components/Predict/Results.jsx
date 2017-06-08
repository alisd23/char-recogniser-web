import React from 'react';

export default ({ predictions, image }) => {

  return (
    <div className="results">
      <div className="results-header">
        <img
          className="example"
          alt="Compressed character"
          src={`data:image/png;base64,${image}`}
        />
        <h2>Top 3 predictions</h2>
      </div>
      <div className="predictions">
        {
          predictions.map(({ charcode, confidence }, i) => (
            <div
              className="prediction"
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
        }
      </div>
    </div>
  )
}
