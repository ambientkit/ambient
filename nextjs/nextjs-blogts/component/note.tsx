import * as React from "react";
import { useState } from "react";
import { useCookies } from "react-cookie";
import Debounce from "@/module/debounce";
import { showFlash, messageType } from "@/component/flash";

interface defaultProps {
  id: string;
  message?: string;
  onChange: (e: string) => void;
  removeNote: (e: string) => void;
}

function View(props: defaultProps): JSX.Element {
  const [cookie] = useCookies(["auth"]);
  const [saving, setSaving] = useState<string>("");
  const [message, setMessage] = useState<string>(props.message);

  const getToken = function (): string {
    let token = "";
    if (cookie.auth && cookie.auth.accessToken) {
      token = cookie.auth.accessToken;
    }
    return token;
  };

  const update = function (id: string, text: string): void {
    fetch("/api/v1/note/" + id, {
      method: "PUT",
      headers: {
        Authorization: "Bearer " + getToken(),
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ message: text }),
    })
      .then((response) => {
        if (response.status !== 200) {
          response.json().then(function (data) {
            showFlash(
              "Could not update note: " + data.message,
              messageType.warning
            );
          });
        }
      })
      .catch((err) => {
        console.log("Error needs to be handled!", err);
      });
  };

  const runDelete = function (id: string) {
    fetch("/api/v1/note/" + id, {
      method: "DELETE",
      headers: {
        Authorization: "Bearer " + getToken(),
      },
    })
      .then((response) => {
        if (response.status === 200) {
          response.json().then(function () {
            showFlash("Note deleted.", messageType.success);

            props.removeNote(id);
          });
        } else {
          response.json().then(function (data) {
            showFlash(
              "Could not update note: " + data.message,
              messageType.warning
            );
          });
        }
      })
      .catch((err) => {
        console.log("Error needs to be handled!", err);
      });
  };

  return (
    <li style={{ marginTop: "12px" } as React.CSSProperties}>
      <div className="box">
        <div className="content">
          <div className="editable">
            <input
              id={props.id}
              type="text"
              className="input individual-note"
              value={message}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                props.onChange(e.currentTarget.value);
                setMessage(e.currentTarget.value);
              }}
              onKeyUp={() => {
                Debounce.run(
                  props.id,
                  () => {
                    setSaving("Saving...");
                    update(props.id, message);
                    setTimeout(() => {
                      setSaving("");
                    }, 1000);
                  },
                  1000
                );
              }}
            />
          </div>
        </div>
        <nav className="level is-mobile">
          <div className="level-left">
            <a
              title="Delete note"
              className="level-item"
              onClick={() => {
                runDelete(props.id);
              }}
            >
              <span className="icon is-small has-text-danger">
                <i className="fas fa-trash" data-cy="delete-note-link"></i>
              </span>
            </a>
          </div>
          <div
            className="level-right"
            style={{ minHeight: "1.2rem" } as React.CSSProperties}
          >
            <span className="is-size-7 has-text-grey">{saving}</span>
          </div>
        </nav>
      </div>
    </li>
  );
}

export default View;
