import * as React from "react";

interface Props {
  title: string;
  description: string;
  children: JSX.Element[] | JSX.Element;
}

function Page(props: Props): JSX.Element {
  return (
    <div>
      <section className="section">
        <div className="container">
          <h1 className="title">{props.title}</h1>
          <h2 className="subtitle">{props.description}</h2>
          {props.children}
        </div>
      </section>
    </div>
  );
}

export default Page;
