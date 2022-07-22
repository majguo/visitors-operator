import React, { Component } from 'react';

import VisitorsTable from './VisitorsTable.js';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      title: "Visitors Dashboard",
    };
    this.loadTitle = this.loadTitle.bind(this)
  }

  loadTitle() {
    fetch('/title')
      .then(results => {
        return results.json();
      }).then(data => {
        this.setState({ title: data.title || "Visitors Dashboard" });
        console.log('title', this.state.title)
      })

    console.log('visitors loadTitle')
  }

  componentDidMount() {
    this.loadTitle()
}

  render() {

    let title = process.env.REACT_APP_TITLE || 'Visitors Dashboard'

    return (
      <div className="container">
        <div className="page-header">
          <h2>{this.state.title}</h2>
        </div>
        <div>
          <VisitorsTable />
        </div>
      </div>
    );
  }
}

export default App;
