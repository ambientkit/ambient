import * as React from "react";
import { useEffect, useState } from "react";
import { useCookies } from "react-cookie";
import { useRouter } from 'next/router';
import Link from 'next/link';
import { useAuth } from '../providers/Auth';

const Page = () => {
  const [cookie, , removeCookie] = useCookies(["auth"]);
  const [shownNavClass, setShownNavClass] = useState<string>("");
  const [shownMobileNavClass, setShownMobileNavClass] = useState<string>("");

  const router = useRouter();

  const clear = function (): void {
    removeCookie("auth", { path: "/" });
  };

  const { isAuthenticated } = useAuth();
  console.log(isAuthenticated);

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
              router.push('/');
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
                {!isAuthenticated && (
                  <Link href="/login">
                    <a className="navbar-item">Login</a>
                  </Link>
                )}
                <Link href="https://petstore.swagger.io/">
                  <a className="navbar-item">Swagger</a>
                </Link>
                <Link href="/about">
                  <a className="navbar-item">About</a>
                </Link>
                <hr className="navbar-divider" />
                {isAuthenticated && (
                  <a
                    className="dropdown-item"
                    onClick={() => {
                      clear();
                      router.push('/login');
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


export default Page;