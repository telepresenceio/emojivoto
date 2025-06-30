import React from 'react';
import ReactDOM from 'react-dom';
import { createRoot } from 'react-dom/client';
import Vote from './components/Vote.jsx';
import Leaderboard from './components/Leaderboard.jsx';
import gridStyles from './../css/grid.css';
import styles from './../css/styles.css';

import { BrowserRouter, Routes, Route } from 'react-router';

const appMain = createRoot(document.getElementById("main"))
appMain.render(
  <BrowserRouter>
      <div className="main-content">
        <Routes>
          <Route exact path="/" element = {<Vote/>} />
          <Route path="/leaderboard" element = {<Leaderboard/>} />
        </Routes>
      </div>
  </BrowserRouter>
);
