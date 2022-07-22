import React, { Component } from 'react';


class VisitorsTable extends Component {
    constructor(props) {
        super(props);
        this.state = {
            visitors: [],
        };
        this.loadChanges = this.loadChanges.bind(this)
    }

    generateUrl() {
        return '/api/visitors/';
    }

    loadChanges() {
        var url = this.generateUrl()

        fetch(url)
        .then(results => {
            return results.json();
        }).then(data => {
            let visitors = data.map((e) => {
                return(
                    <tr key={e.id}>
                        <td>{e.service_ip}</td>
                        <td>{e.client_ip}</td>
                        <td>{new Date(e.timestamp).toLocaleString()}</td>
                    </tr>
                )
            })
            this.setState({visitors: visitors})
            console.log('visitors', this.state.visitors)
        })

        console.log('visitors loadChanges')
    }

    componentDidMount() {

        // Send request to log visitor
        var url = this.generateUrl()
        fetch(url, {method: 'POST',})

        // Initial load of the visitors list
        this.loadChanges()

        // Repeat the visitors load
        setInterval(this.loadChanges, 1000)
    }

    render() {
        return(
            <table className="table table-striped table-bordered table-hover table-dark">
                <thead>
                    <tr>
                        <th scope="col">Service IP</th>
                        <th scope="col">Client IP</th>
                        <th scope="col">Timestamp</th>
                    </tr>
                </thead>
                <tbody>
                    {this.state.visitors}
                </tbody>
            </table>
        )
    }
}

export default VisitorsTable;