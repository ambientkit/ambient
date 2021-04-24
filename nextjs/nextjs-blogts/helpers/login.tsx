import { useCookies } from "react-cookie";

export function isLoggedIn(): boolean {
    const [cookie] = useCookies(["auth"]);

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
}