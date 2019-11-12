import React, { useState, useEffect } from 'react';
import axios from 'axios';

import { API_URL } from './constants';

import Home from './Home';
import Loading from './Loading';

const App = () => {

  const [isAuthenticated, setIsAuthenticated] = useState(null);

  useEffect(() => {

    async function fetchData() {
      const result = await axios(`/auth/validate`);
      console.debug(result);

      let authenticated = result.data.authenticated;
      console.log(`Is session user authenticated ? ${authenticated}`)

      setIsAuthenticated(authenticated);
      if (!authenticated) {
        window.location.href = `${API_URL}/auth`
      }
    }

    fetchData();

  }, []);

  if (isAuthenticated === null) {
    return <Loading />
  }
  return isAuthenticated ? <Home /> : <Loading />;
}

export default App;
