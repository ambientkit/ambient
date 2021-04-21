import * as React from "react";

interface Props {
  label?: string;
  name: string;
  type?: string;
  required?: boolean;
  value: string;
  onChange: (e: string) => void;
}

function View(props: Props): JSX.Element {
  return (
    <div className="field">
      <label className="label">{props.label}</label>
      <div className="control">
        <input
          name={props.name}
          type={props.type || "text"}
          className="input"
          data-cy={props.name}
          required={props.required}
          value={props.value}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            props.onChange(e.currentTarget.value);
          }}
        ></input>
      </div>
    </div>
  );
}

export default View;
