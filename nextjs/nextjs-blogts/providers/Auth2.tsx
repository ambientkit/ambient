import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { useRouter } from 'next/router'

type AuthContext = {
    isAuthenticated: boolean;
    setAuthenticated: React.Dispatch<React.SetStateAction<boolean>>;
}

const AuthContext = React.createContext({})

function AuthProvider({ children }: { children: ReactNode }) {
    const { pathname, events } = useRouter()
    const [user, setUser] = useState(false)

    async function getUser() {
        setUser(true)
        // try {
        //     const response = await fetch('/api/me')
        //     const profile = await response.json()
        //     if (profile.error) {
        //         setUser(false)
        //     } else {
        //         setUser(profile)
        //     }
        // } catch (err) {
        //     console.error(err)
        // }
    }

    useEffect(() => {
        getUser()
    }, [pathname])

    useEffect(() => {
        // Check that a new route is OK
        const handleRouteChange = (url: string) => {
            if (url !== '/' && !user) {
                window.location.href = '/'
            }
        }

        // Check that initial route is OK
        if (pathname !== '/' && user === null) {
            window.location.href = '/'
        }

        // Monitor routes
        events.on('routeChangeStart', handleRouteChange)
        return () => {
            events.off('routeChangeStart', handleRouteChange)
        }
    }, [user])

    return (
        <AuthContext.Provider value={{ user }}>
            {children}
        </AuthContext.Provider>
    )
}

const useAuth = () => useContext(AuthContext)

export { AuthProvider, useAuth }