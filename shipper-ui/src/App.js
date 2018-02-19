import React, { Component } from "react";
import "./App.css";
import CreateConsignment from "./CreateConsignment";
import Authenticate from "./Authenticate";

class App extends Component {
  state = {
    err: null,
    authenticated: false,
  };

  onAuth = token => {
    this.setState({
      authenticated: true,
    });
  };

  renderLogin = () => {
    return <Authenticate onAuth={this.onAuth} />;
  };

  renderAuthenticated = () => {
    return <CreateConsignment />;
  };

  getToken = () => {
    return localStorage.getItem("token") || false;
  };

  isAuthenticated = () => {
    return this.state.authenticated || this.getToken() || false;
  };

  render() {
    const authenticated = this.isAuthenticated();

    return (
      <div className="App">
        <header className="App-header">
          <h2 className="App-title">Shipper</h2>
        </header>
        <div className="App-intro">
          {authenticated ? this.renderAuthenticated() : this.renderLogin()}
        </div>
      </div>
    );
  }
}

export default App;
