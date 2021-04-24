import cookie from "cookie"
import { IncomingMessage } from 'http';

interface Request {
    headers: Header;
}

interface Header {
    cookie: string;
}

export function parseCookies(req: IncomingMessage | undefined) {
    return cookie.parse(req ? req.headers.cookie || "" : document.cookie)
}
