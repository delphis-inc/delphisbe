import React from "react";
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";
import './App.scss';
import Discussions from './components/Discussions';
import Discussion from './components/Discussion';
import Layout from './components/Layout';

function DelphisApp() {
  return (
    <div className="delphis-app">
      <Layout>
        <Router>
            <Switch>
              <Route exact path="/" component={Discussions}/>
              <Route path="/d/:id" component={Discussion}/>
            </Switch>
        </Router>
      </Layout>
    </div>
  );
}

export default DelphisApp