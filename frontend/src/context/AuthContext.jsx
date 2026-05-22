import { createContext, useContext, useState, useEffect } from 'react'
import { getUser, isLoggedIn, logout } from '../api/client'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (isLoggedIn()) {
      setUser(getUser())
    }
    setLoading(false)
  }, [])

  function handleLogin(userData) {
    setUser(userData)
  }

  function handleLogout() {
    logout()
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, loading, handleLogin, handleLogout }}>
      {!loading && children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}