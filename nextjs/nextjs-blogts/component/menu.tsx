import * as React from "react";
import { useEffect, useState } from "react";
import { navigate } from "hookrouter";
import { useCookies } from "react-cookie";

function View(): JSX.Element {
  const [cookie, , removeCookie] = useCookies(["auth"]);
  const [shownNavClass, setShownNavClass] = useState<string>("");
  const [shownMobileNavClass, setShownMobileNavClass] = useState<string>("");

  const clear = function (): void {
    removeCookie("auth", { path: "/" });
  };

  const isLoggedIn = function (): boolean {
    try {
      const auth = cookie.auth;
      if (auth === undefined) {
        return false;
      }
      return true;
    } catch (err) {
      console.log(err);
    }

    return false;
  };

  useEffect(() => {
    // Close the nav menus when an item is clicked.
    const links = document.querySelectorAll(".navbar-item");
    links.forEach((link) => {
      link.addEventListener("click", function () {
        setShownNavClass("");
        setShownMobileNavClass("");
      });
    });
  }, []);

  return (
    <main>
      <nav
        className="navbar is-black"
        role="navigation"
        aria-label="main navigation"
      >
        <div className="navbar-brand">
          <a
            className="navbar-item"
            data-cy="home-link"
            onClick={() => {
              navigate("/");
            }}
          >
            <strong>goreactapp</strong>
          </a>

          <a
            id="mobile-navbar-top"
            role="button"
            className={"navbar-burger burger " + shownMobileNavClass}
            aria-label="menu"
            aria-expanded="false"
            data-target="navbar-top"
            onClick={() => {
              if (shownMobileNavClass == "is-active") {
                setShownMobileNavClass("");
              } else {
                setShownMobileNavClass("is-active");
              }
            }}
          >
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
          </a>
        </div>

        <div id="navbar-top" className={"navbar-menu " + shownMobileNavClass}>
          <div className="navbar-end">
            <div
              id="ddmenu"
              className={`navbar-item has-dropdown ` + shownNavClass}
              onMouseEnter={() => {
                setShownNavClass("is-active");
              }}
              onMouseLeave={() => {
                setShownNavClass("");
              }}
            >
              <a className="navbar-link">Menu</a>

              <div className="navbar-dropdown is-right">
                {!isLoggedIn() && (
                  <a
                    className="navbar-item"
                    onClick={() => {
                      navigate("/login");
                    }}
                  >
                    Login
                  </a>
                )}
                <a
                  className="navbar-item"
                  href={`https://petstore.swagger.io/?url=${location.origin}/static/swagger.json`}
                >
                  Swagger
                </a>

                <a
                  className="navbar-item"
                  onClick={() => {
                    navigate("/about");
                  }}
                >
                  About
                </a>
                <hr className="navbar-divider" />
                {isLoggedIn() && (
                  <a
                    className="dropdown-item"
                    onClick={() => {
                      clear();
                      navigate("/login");
                    }}
                  >
                    Logout
                  </a>
                )}
                <div className="navbar-item">v1.0.0</div>
              </div>
            </div>
          </div>
        </div>
      </nav>
    </main>
  );
}

export default View;
