import React from 'react';
import "./Navbar.scss";

export default ({ children }) => (
  <div className="navbar">
    <span>Character Recognition</span>
    <div className="page-links">
      {children}
    </div>
  </div>
);
