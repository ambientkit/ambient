import * as React from "react";
import {
  withKnobs,
  text,
  button,
  select,
  number,
  boolean,
} from "@storybook/addon-knobs";
import { withA11y } from "@storybook/addon-a11y";
import Flash, { messageType, showFlash } from "@/component/flash";
import "~/style/main.scss";

export default {
  title: "Component/Flash",
  component: Flash,
  decorators: [withKnobs, withA11y],
};

interface Controls {
  timeout: number;
  prepend: boolean;
}

const Knobs = function (style: messageType): Controls {
  const timeout = number("Timeout (milliseconds)", 4000);
  const message = text("Text", "This is a flash message.");
  const prepend = boolean("Prepend", false);

  button("Show", function (): boolean {
    showFlash(message, style);
    return false; // False prevents re-rendering of the story.
  });

  return { timeout, prepend };
};

export const Success = function (): JSX.Element {
  const knobs = Knobs(messageType.success);
  return <Flash timeout={knobs.timeout} prepend={knobs.prepend} />;
};

export const Failed = function (): JSX.Element {
  const knobs = Knobs(messageType.failed);
  return <Flash timeout={knobs.timeout} prepend={knobs.prepend} />;
};

export const Warning = function (): JSX.Element {
  const knobs = Knobs(messageType.warning);
  return <Flash timeout={knobs.timeout} prepend={knobs.prepend} />;
};

export const Primary = function (): JSX.Element {
  const knobs = Knobs(messageType.primary);
  return <Flash timeout={knobs.timeout} prepend={knobs.prepend} />;
};

export const Link = function (): JSX.Element {
  const knobs = Knobs(messageType.link);
  return <Flash timeout={knobs.timeout} prepend={knobs.prepend} />;
};

export const Info = function (): JSX.Element {
  const knobs = Knobs(messageType.info);
  return <Flash timeout={knobs.timeout} prepend={knobs.prepend} />;
};

export const Dark = function (): JSX.Element {
  const knobs = Knobs(messageType.dark);
  return <Flash timeout={knobs.timeout} prepend={knobs.prepend} />;
};

export const Action = function (): JSX.Element {
  const s = select("Message Type", messageType, messageType.success);
  const knobs = Knobs(s);
  return <Flash timeout={knobs.timeout} prepend={knobs.prepend} />;
};
