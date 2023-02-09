import { HashRouter } from 'react-router-dom';
import Router from './routers';

import Header from './components/Header';

import './App.css';

function App() {
  return (
    <>
      <Header />
      <HashRouter>
        <Router />
      </HashRouter>
    </>
  );
}

export default App;
