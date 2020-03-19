import React from "react";
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";
import './App.scss';
import UserQueries from './queries/user';
import Discussions from './components/Discussions';
import Discussion from './components/Discussion';
import Layout from './components/Layout';
import { useQuery } from "react-apollo";

function DelphisApp() {
  const { loading, error, data } = useQuery(UserQueries.me);

  const loadingDiv = <div>Loading...</div>;
  const loginDiv = (
    <div>
      <a href="https://staging.delphishq.com/twitter/login">
        <img src="/twitter-signin.png" alt="sign in with twitter"/>
      </a>
    </div>
  );
  
  const router = (
    <Router>
      <Switch>
        <Route exact path="/" component={Discussions}/>
        <Route path="/d/:id" component={Discussion}/>
      </Switch>
    </Router>
  )
  let renderedDiv = router;
  if (loading) {
    renderedDiv = loadingDiv;
  }
  if (error) {
    // TODO: Ensure this error actually is an auth error.
    // In this case the user needs to login.
    renderedDiv = loginDiv;
  }
  return (
    <div className="delphis-app">
      <Layout>
        {renderedDiv}
      </Layout>
    </div>
  );
}

export default DelphisApp