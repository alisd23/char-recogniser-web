import React from 'react';

export default ({ active, children }) => (
  <div
    style={{ display: active ? 'block' : 'none' }}
    className="page-wrapper"
  >
    {children}
  </div>
);
