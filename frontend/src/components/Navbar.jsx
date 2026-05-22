import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function Navbar() {
  const { user, handleLogout } = useAuth()
  const navigate = useNavigate()

  function onLogout() {
    handleLogout()
    navigate('/login')
  }

  return (
    <nav style={styles.nav}>
      <div style={styles.inner}>
        <Link to="/dashboard" style={styles.logo}>Bookify</Link>

        <div style={styles.links}>
          {user ? (
            <>
              <Link to="/services" style={styles.link}>Услуги</Link>
                {user.role === 'provider' && (
              <Link to="/my-services" style={styles.link}>Мои услуги</Link>
            )}
                {user.role === 'admin' && (
            <>
              <Link to="/admin" style={styles.link}>Дашборд</Link>
              <Link to="/admin/users" style={styles.link}>Пользователи</Link>
            </>
            )}
              <Link to="/my-appointments" style={styles.link}>Мои записи</Link>
              <span style={styles.email}>{user.full_name}</span>
              <button onClick={onLogout} style={styles.btn}>Выйти</button>
            </>
          ) : (
            <>
              <Link to="/login"    style={styles.link}>Войти</Link>
              <Link to="/register" style={styles.link}>Регистрация</Link>
            </>
          )}
        </div>
      </div>
    </nav>
  )
}

const styles = {
  nav: {
    background: '#1a1a2e',
    padding: '0 16px',
    position: 'sticky',
    top: 0,
    zIndex: 100,
  },
  inner: {
    maxWidth: 900,
    margin: '0 auto',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    height: 56,
  },
  logo: {
    color: '#e94560',
    fontWeight: 700,
    fontSize: 22,
    textDecoration: 'none',
    letterSpacing: 1,
  },
  links: {
    display: 'flex',
    alignItems: 'center',
    gap: 20,
  },
  link: {
    color: '#ccc',
    textDecoration: 'none',
    fontSize: 14,
  },
  email: {
    color: '#888',
    fontSize: 13,
  },
  btn: {
    background: 'transparent',
    border: '1px solid #e94560',
    color: '#e94560',
    padding: '5px 14px',
    borderRadius: 6,
    cursor: 'pointer',
    fontSize: 13,
  },
}