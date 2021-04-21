import * as React from "react";
import { withKnobs, text } from "@storybook/addon-knobs";
import { withA11y } from "@storybook/addon-a11y";
import SimplePage from "@/component/simple-page";
import "~/style/main.scss";

export default {
  title: "Component/Simple Page",
  component: SimplePage,
  decorators: [withKnobs, withA11y],
};

export const View = function (): JSX.Element {
  return (
    <SimplePage
      title={text("Title", "This is the Title")}
      description={text("Description", "This is a subtitle or description.")}
    >
      <div>{text("Content", "This is the content.")}</div>
    </SimplePage>
  );
};

View.story = {
  name: "Simple Page",
};
