import * as React from "react";
import { useState } from "react";
import { withKnobs, boolean, button } from "@storybook/addon-knobs";
import { withA11y } from "@storybook/addon-a11y";
import Note from "./note";
import Flash from "@/component/flash";
import { rest } from "msw";
import { worker } from "@/mock/browser";
import "~/style/main.scss";

export default {
  title: "Component/Note",
  component: Note,
  decorators: [withKnobs, withA11y],
};

export const note = (): JSX.Element => {
  const shouldFail = boolean("Fail", false);
  const [state, setState] = useState<string>("");
  const [visible, setVisible] = useState<boolean>(true);

  button("Restore Note", function () {
    setVisible(true);
  });

  const removeNote = function (id: string): void {
    console.log("Remove note with id:", id);
    setVisible(false);
  };

  worker.use(
    ...[
      rest.put("/api/v1/note/1", (req, res, ctx) => {
        if (shouldFail) {
          return res(
            ctx.status(400),
            ctx.json({
              message: "There was an error.",
            })
          );
        } else {
          return res(
            ctx.status(200),
            ctx.json({
              message: "ok",
            })
          );
        }
      }),
      rest.delete("/api/v1/note/1", (req, res, ctx) => {
        if (shouldFail) {
          return res(
            ctx.status(400),
            ctx.json({
              message: "There was an error.",
            })
          );
        } else {
          return res(
            ctx.status(200),
            ctx.json({
              message: "ok",
            })
          );
        }
      }),
    ]
  );

  return (
    <main>
      {visible ? (
        <ul>
          <Note
            id="1"
            message={state}
            onChange={(e: string) => setState(e)}
            removeNote={removeNote}
          />
        </ul>
      ) : null}

      <Flash />
    </main>
  );
};
