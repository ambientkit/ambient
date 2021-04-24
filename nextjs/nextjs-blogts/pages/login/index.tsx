import Head from 'next/head'
import Layout, { siteTitle } from '~/components/layout'
import { useState } from "react";
import Submit from "~/module/submit";
import Input from "~/component/input";
import { messageType, showFlash } from "~/component/flash";
import { useRouter } from 'next/router';
import { useAuth } from '~/providers/Auth';
import withoutAuth from '~/hocs/withoutAuth';

interface defaultProps {
    email?: string;
    password?: string;
}

interface User {
    email: string;
    password: string;
}

export default withoutAuth(function Page(props: defaultProps) {
    const router = useRouter();

    const { setAuthenticated } = useAuth();

    const data = {
        title: "Login",
        subtitle: "Enter your login information below.",
    };

    const clear = () => {
        setUser({ email: "", password: "" });
    };

    const [user, setUser] = useState<User>({
        email: props.email || "",
        password: props.password || "",
    });

    function toRegister(e: { preventDefault: () => void }) {
        e.preventDefault();
        router.push("/register");
    }

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

                                    fetch("/api/v1/login", {
                                        method: "POST",
                                        headers: { "Content-Type": "application/json" },
                                        body: JSON.stringify(user),
                                    })
                                        .then((response) => {
                                            const auth = {
                                                accessToken: "",
                                                loggedIn: false,
                                            };

                                            if (response.status === 200) {
                                                response.json().then(function (data) {
                                                    Submit.finish();

                                                    auth.loggedIn = true;
                                                    auth.accessToken = data.token;
                                                    setAuthenticated(true);

                                                    showFlash("Login successful.", messageType.success);

                                                    router.push("/");
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
                                            Submit
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

                                    <p className="control">
                                        <a href="#" className="button is-light" onClick={toRegister}>
                                            Register
                                        </a>
                                    </p>
                                </div>
                            </form>
                        </div>
                    </section>
                </div>
            </section>
        </Layout>
    );
})