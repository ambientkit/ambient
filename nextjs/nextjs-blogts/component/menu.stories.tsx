import * as React from "react";
import { withKnobs } from "@storybook/addon-knobs";
import { withA11y } from "@storybook/addon-a11y";
import Menu from "./menu";
import "~/style/main.scss";

export default {
  title: "Component/Menu",
  component: Menu,
  decorators: [withKnobs, withA11y],
};

export const View = function (): JSX.Element {
  return <Menu />;
};

View.story = {
  name: "Menu",
};
