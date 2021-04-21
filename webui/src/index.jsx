import * as React from 'react'
import * as ReactDOM from "react-dom";
import * as Server from 'react-dom/server'

let Greet = () => <h1>Hello, world!</h1>

console.log(Server.renderToString('React app started.'));

if (typeof window !== "undefined") {
    ReactDOM.render(
        <Greet />,
        document.getElementById('root')
    );
} else {
    reactssr.render(Server.renderToString(
        <React.StrictMode>
            <Greet />
        </React.StrictMode>
    ));
}