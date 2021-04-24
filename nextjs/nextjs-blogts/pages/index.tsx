import Head from 'next/head'
import Layout, { siteTitle } from '../components/layout'
import Link from 'next/link'
import { useCookies } from "react-cookie";
import { useRouter } from 'next/router'
import { useEffect } from "react"
import withAuth from '~/hocs/withAuth';

export default withAuth(function Page() {
  const title = "Welcome";
  const subtitle = "Login was successful";

  const router = useRouter();
  const [cookie] = useCookies(["auth"]);
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
    if (!isLoggedIn()) {
      router.push("/login");
    }
  }, [])

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
            <Link href="/notepad">
              <a data-cy="notepad-link">
                Click here to access your Notepad.
                </a>
            </Link>
          </div>
        </section>
      </>
    </Layout>
  )
})