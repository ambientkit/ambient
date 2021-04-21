import * as React from "react";
import { useState } from "react";
import { withKnobs, select, text } from "@storybook/addon-knobs";
import { withA11y } from "@storybook/addon-a11y";
import Input from "./input";
import "~/style/main.scss";

export default {
  title: "Component/Input",
  component: Input,
  decorators: [withKnobs, withA11y],
};

export const input = function (): JSX.Element {
  const label = text("Label", "First Name");
  const type = select(
    "Type",
    {
      text: "text",
      color: "color",
      date: "date",
      "datetime-local": "datetime-local",
      email: "email",
      hidden: "hidden",
      month: "month",
      number: "number",
      password: "password",
      range: "range",
      search: "search",
      time: "time",
      week: "week",
    },
    "text"
  );

  // Set the state.
  const [state, setState] = useState<string>("");

  return (
    <Input
      label={label}
      name=""
      type={type}
      required={true}
      value={state}
      onChange={(e: string) => setState(e)}
    />
  );
};
