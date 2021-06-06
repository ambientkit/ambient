import Head from 'next/head'
import Layout, { siteTitle } from '~/components/layout'

import { useState } from "react";
import Input from "~/component/input";
import Submit from "~/module/submit";
import { showFlash, messageType } from "~/component/flash";
import { useRouter } from 'next/router';

const data = {
  title: "Register",
  subtitle: "Enter your information below.",
};

interface User {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
}

interface defaultProps {
  firstName?: string;
  lastName?: string;
  email?: string;
  password?: string;
}

function Page(props: defaultProps) {
  const router = useRouter();

  const clear = () => {
    setUser({ firstName: "", lastName: "", email: "", password: "" });
  };

  const [user, setUser] = useState<User>({
    firstName: props.firstName || "",
    lastName: props.lastName || "",
    email: props.email || "",
    password: props.password || "",
  });

  return (
    <Layout>
      <Head>
        <title>{siteTitle}</title>
      </Head>
      <section>
        <div>
          <section className="section">
            <div className="container">
              <h1 className="title">{data.title}</h1>
              <h2 className="subtitle">{data.subtitle}</h2>
            </div>

            <div
              className="container"
              style={{ marginTop: "1em" } as React.CSSProperties}
            >
              <form
                name="login"
                onSubmit={(e) => {
                  e.preventDefault();
                  Submit.start(e);
                  fetch("/api/v1/register", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(user),
                  })
                    .then((response) => {
                      if (response.status === 201) {
                        response.json().then(function () {
                          clear();
                          Submit.finish();

                          showFlash("User registered.", messageType.success);

                          router.push("/login");
                        });
                      } else {
                        response.json().then(function (data) {
                          Submit.finish();
                          showFlash(data.message, messageType.warning);
                        });
                      }
                    })
                    .catch((err) => {
                      console.log("Error needs to be handled!", err);
                    });
                }}
              >
                <Input
                  label="First Name"
                  name="first_name"
                  type="text"
                  required={true}
                  onChange={(e: string) => {
                    const newUser = { ...user };
                    newUser.first_Name = e;
                    setUser(newUser);
                  }}
                  value={user.firstName}
                />

                <Input
                  label="Last Name"
                  name="last_name"
                  type="text"
                  required={true}
                  onChange={(e: string) => {
                    const newUser = { ...user };
                    newUser.lastName = e;
                    setUser(newUser);
                  }}
                  value={user.lastName}
                />

                <Input
                  label="Email"
                  name="email"
                  type="email"
                  required={true}
                  onChange={(e: string) => {
                    const newUser = { ...user };
                    newUser.email = e;
                    setUser(newUser);
                  }}
                  value={user.email}
                />

                <Input
                  label="Password"
                  name="password"
                  required={true}
                  onChange={(e: string) => {
                    const newUser = { ...user };
                    newUser.password = e;
                    setUser(newUser);
                  }}
                  value={user.password}
                  type="password"
                />

                <div className="field is-grouped">
                  <p className="control">
                    <button
                      id="submit"
                      type="submit"
                      data-cy="submit"
                      className="button is-primary"
                    >
                      Create Account
                  </button>
                  </p>

                  <p className="control">
                    <button
                      type="button"
                      className="button is-light"
                      onClick={() => {
                        clear();
                      }}
                    >
                      Clear
                  </button>
                  </p>
                </div>
              </form>
            </div>
          </section>
        </div>
      </section>
    </Layout>
  );
}

export default Page;
