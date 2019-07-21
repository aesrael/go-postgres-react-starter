import React from "react";
import { BrowserRouter as Router, Route } from "react-router-dom";
import Register from "./components/Register";
import Login from "./components/Login";
import Session from "./components/Session";
import "./App.css"

function App() {
  return (
      <Router>
           <Route exact path="/" component={Login} />
          <Route path="/register" component={Register} />
          <Route path="/login" component={Login} />
          <Route path="/session" component={Session} />
      </Router>
  );
}

export default App;
