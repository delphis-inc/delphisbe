import React from "react";
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";
import Discussions from './components/Discussions';
import Discussion from './components/Discussion';

// This site has 3 pages, all of which are rendered
// dynamically in the browser (not server rendered).
//
// Although the page does not ever refresh, notice how
// React Router keeps the URL up to date as you navigate
// through the site. This preserves the browser history,
// making sure things like the back button and bookmarks
// work properly.

function BasicExample() {
  return (
    <Router>
        <Switch>
          <Route exact path="/" component={Discussions}/>
          <Route path="/d/:id" component={Discussion}/>
        </Switch>
    </Router>
  );
}

export default BasicExample