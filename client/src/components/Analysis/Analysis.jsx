import React from 'react';
import Predictions from '../Predictions';
import './Analysis.scss';

export default ({ activations, predictions, image }) => {
  return (
    <div className="analysis-page">
      <div className="activations-section">
        <h3 className="section-header">Activations for Convolutional layer 1</h3>
        <div className="activations">
          {
            activations && activations.map((a, i) => (
              <img
                key={i}
                alt={`Activation ${i+1}`}
                src={`data:image/png;base64,${a}`}
                className="activation"
              />
            ))
          }
        </div>
      </div>
      <div className="predictions-section">
        <h3 className="section-header">Predictions</h3>
        <div className="predictions-wrapper">
          {
            predictions && image && (
              <Predictions
                key="results"
                image={image}
                predictions={predictions}
              />
            )
          }
        </div>
      </div>
    </div>
  )
}
