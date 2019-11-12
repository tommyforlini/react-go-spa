import React from 'react';
import { BrowserRouter as Router, Route, Link, Switch } from 'react-router-dom';

import logo from './logo.svg';
import './App.css';

import Landing from './Landing';
import Public from './Public';
import Private from './Private';

const Home = () => {

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        {/* <p>
          Edit <code>src/Home.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a> */}
      </header>

      <Router>
          <Link to="/">Home</Link> 
          &nbsp;&nbsp; | &nbsp;&nbsp;
          <Link to="/public">Public Content</Link> 
          &nbsp;&nbsp; | &nbsp;&nbsp;
          <Link to="/private">Private Content</Link>

          <div>
            <Switch>
              <Route path="/public"><Public /></Route>
              <Route path="/private"><Private /></Route>
              <Route path="/"><Landing /></Route>
            </Switch>
          </div>
        </Router>

    </div>
  );

}

export default Home;
