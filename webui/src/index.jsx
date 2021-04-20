import * as React from 'react'
import * as ReactDOM from "react-dom";
import * as Server from 'react-dom/server'

let Greet = () => <h1>Hello, world!</h1>

console.log(Server.renderToString('React app started.'));

ReactDOM.render(
    <Greet />,
    document.getElementById('root')
);