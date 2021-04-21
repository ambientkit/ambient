import * as React from "react";
import { useState, useEffect, useRef } from "react";
import { uid } from "react-uid";
import EventEmitter from "~/module/event";

// Components allows you to display flash messages on the screen that are styled
// with Bulma: http://bulma.io/documentation/components/message/.
// To use this component, you should use add it to your page:
/*
<Flash timeout={number} prepend={boolean} />
*/
// And then dispatch events to the component:
/*
showMessage("This is the message", messageType.success);
*/

// Types of styles available for the flash messages.
export enum messageType {
  success = "is-success",
  failed = "is-danger",
  warning = "is-warning",
  primary = "is-primary",
  link = "is-link",
  info = "is-info",
  dark = "is-dark",
}

// Types of events this component supports.
export enum flashEvent {
  showMessage = "Flash.showMessage",
}

// flashMessage is used by the component and by others calling the component.
interface flashMessage {
  message: string;
  style: messageType;
}

// defaultProps is the optional values that can be passed to the component.
interface defaultProps {
  timeout?: number;
  prepend?: boolean;
}

// Track if component is mounted or not.
// https://www.debuggr.io/react-update-unmounted-component/
function useIsMountedRef() {
  const isMountedRef = useRef(null);
  useEffect(() => {
    isMountedRef.current = true;
    return () => (isMountedRef.current = false);
  });
  return isMountedRef;
}

export const showFlash = function (message: string, style: messageType): void {
  EventEmitter.dispatch(flashEvent.showMessage, {
    message: message,
    style: style,
  });
};

const Flash = function (props: defaultProps): JSX.Element {
  // Handle optional props and provide defaults.
  const showtime = props.timeout || 4000;
  const prepend = props.prepend || false;

  // Use useState from React 16.8 so we can leverage functional components
  // with state instead of using a class. We want to keep track of messages
  // and timers when this component is used across pages.
  const [list, setList] = useState<flashMessage[]>([]);

  // Use a reference to access the current value in an async callback. This
  // prevents stale closures which are callbacks that only update a point in
  // variable instead of the current variable.
  // https://github.com/facebook/react/issues/14010
  const listRef = useRef(list);
  listRef.current = list;

  // Use a reference to whether component is mounted or not.
  const isMountedRef = useIsMountedRef();

  const showMessage = (msg: flashMessage): void => {
    // Don't show a message if zero.
    if (showtime === 0) {
      return;
    }

    // Check if the messages should stack in reverse order.
    const newList = [...list];
    if (prepend === true) {
      newList.unshift(msg);
    } else {
      newList.push(msg);
    }
    setList(newList);

    // Show forever if -1.
    if (showtime > 0) {
      setTimeout(function () {
        removeFlash(msg);
      }, showtime);
    }
  };

  const removeFlash = (i: flashMessage): void => {
    // If the component is not mounted, then don't change state to prevent
    // the error below.
    // "Warning: Can't perform a React state update on an unmounted component.
    // This is a no-op, but it indicates a memory leak in your application.
    // To fix, cancel all subscriptions and asynchronous tasks in a useEffect
    // cleanup function."
    // https://stackoverflow.com/a/8860210/13953226
    // https://www.debuggr.io/react-update-unmounted-component/
    if (!isMountedRef.current) {
      return;
    }

    // Prevent stale closure: must use ref of list instead of list.
    // https://github.com/facebook/react/issues/14010
    setList(
      listRef.current.filter((v: flashMessage) => {
        return v !== i;
      })
    );
  };

  // The Effect Hook is available in 16.8 to allow functional components to be
  // stateful without a class.
  // https://reactjs.org/docs/hooks-effect.html
  // Equivalent to: componentDidMount()
  useEffect(() => {
    // Subscribe so messages can be added from other components.
    EventEmitter.subscribe(flashEvent.showMessage, (msg: flashMessage) => {
      showMessage(msg);
    });

    // Perform cleanup - equivalent to: componentWillUnmount()
    return () => {
      EventEmitter.unsubscribe(flashEvent.showMessage);
    };
  });

  return (
    <div
      style={
        {
          marginTop: "1em",
          position: "fixed",
          bottom: "1.5rem",
          right: "1.5rem",
          zIndex: 100,
          margin: 0,
        } as React.CSSProperties
      }
    >
      {list.map((i: flashMessage) => (
        <div key={uid(i)} className={`notification ${i.style}`}>
          {i.message}
          <button
            className="delete"
            onClick={() => {
              removeFlash(i);
            }}
          ></button>
        </div>
      ))}
    </div>
  );
};

export default Flash;
