import Head from 'next/head'
import Layout, { siteTitle } from '~/components/layout'
import { useRouter } from "next/router";

export default function Page() {
    const title = "Welcome";
    const subtitle = "Login was successful";

    const router = useRouter();

    function onclick(e: { preventDefault: () => void }) {
        e.preventDefault();
        router.push("/notepad");
    }

    return (
        <Layout>
            <Head>
                <title>{siteTitle}</title>
            </Head>
            <>
                <section className="hero is-primary">
                    <div className="hero-body">
                        <div className="container">
                            <h1 className="title">{title}</h1>
                            <h2 className="subtitle">{subtitle}</h2>
                        </div>
                    </div>
                </section>
                <br />
                <section>
                    <div className="container">
                        <a href="#" onClick={onclick} data-cy="notepad-link">
                            Click here to access your Notepad.
                        </a>
                    </div>
                </section>
            </>
        </Layout>
    );
}
