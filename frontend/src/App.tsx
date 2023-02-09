import { HashRouter } from "react-router-dom";

import CustomHead from "./header"

import Router from "./routers";

import "./App.css";

function App() {
  return (
    <HashRouter>
      <CustomHead />
      <Router />
    </HashRouter>
  );
}

export default App;
