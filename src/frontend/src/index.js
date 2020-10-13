import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import {GeneralInfoPage, ProcessInfoPage, LogsInfoPage, SiteInfoPage, SearchBookPage, BookInfoPage} from './App';
import * as serviceWorker from './serviceWorker';
import { BrowserRouter as Router,Route, Switch } from 'react-router-dom';

function Routes(){
  return (
    <Router>
      <div className="App">
        <Switch>
          <Route path="/" exact component={GeneralInfoPage} />
          <Route path="/site/" component={SiteInfoPage} />
          <Route path="/process/" component={ProcessInfoPage} />
          <Route path="/logs/" component={LogsInfoPage} />
          <Route path="/:name/" exact component={SiteInfoPage} />
          <Route path="/:name/search" component={SearchBookPage} />
          <Route path="/:name/book/:num" exact component={BookInfoPage} />
        </Switch>
      </div>
    </Router>
  );
}

ReactDOM.render(
  <Routes/>,
  document.getElementById('root')
);



// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
