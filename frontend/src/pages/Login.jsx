import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { login } from '../api/client'
import { useAuth } from '../context/AuthContext'

export default function Login() {
  const navigate = useNavigate()
  const { handleLogin } = useAuth()

  const [form, setForm] = useState({ email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  function onChange(e) {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  async function onSubmit(e) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const data = await login(form)
      handleLogin(data.user)
      navigate('/dashboard')
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={styles.wrap}>
      <div style={styles.card}>
        <h2 style={styles.title}>Вход в Bookify</h2>

        {error && <div style={styles.error}>{error}</div>}

        <form onSubmit={onSubmit}>
          <label style={styles.label}>Email</label>
          <input
            style={styles.input}
            type="email"
            name="email"
            value={form.email}
            onChange={onChange}
            placeholder="you@example.com"
            required
          />

          <label style={styles.label}>Пароль</label>
          <input
            style={styles.input}
            type="password"
            name="password"
            value={form.password}
            onChange={onChange}
            placeholder="••••••••"
            required
          />

          <button style={styles.btn} type="submit" disabled={loading}>
            {loading ? 'Загрузка...' : 'Войти'}
          </button>
        </form>

        <p style={styles.footer}>
          Нет аккаунта?{' '}
          <Link to="/register" style={styles.linkText}>Зарегистрироваться</Link>
        </p>
      </div>
    </div>
  )
}

const styles = {
  wrap: {
    display: 'flex',
    justifyContent: 'center',
    marginTop: 60,
  },
  card: {
    background: '#fff',
    border: '1px solid #e5e7eb',
    borderRadius: 12,
    padding: '36px 32px',
    width: '100%',
    maxWidth: 400,
  },
  title: {
    marginBottom: 24,
    fontSize: 22,
    fontWeight: 600,
    color: '#111',
  },
  label: {
    display: 'block',
    marginBottom: 6,
    fontSize: 13,
    color: '#555',
  },
  input: {
    display: 'block',
    width: '100%',
    padding: '10px 12px',
    marginBottom: 16,
    border: '1px solid #d1d5db',
    borderRadius: 8,
    fontSize: 14,
    boxSizing: 'border-box',
    outline: 'none',
  },
  btn: {
    width: '100%',
    padding: '11px 0',
    background: '#e94560',
    color: '#fff',
    border: 'none',
    borderRadius: 8,
    fontSize: 15,
    cursor: 'pointer',
    fontWeight: 600,
    marginTop: 4,
  },
  error: {
    background: '#fef2f2',
    border: '1px solid #fecaca',
    color: '#b91c1c',
    borderRadius: 8,
    padding: '10px 14px',
    fontSize: 13,
    marginBottom: 16,
  },
  footer: {
    marginTop: 20,
    textAlign: 'center',
    fontSize: 13,
    color: '#666',
  },
  linkText: {
    color: '#e94560',
    textDecoration: 'none',
    fontWeight: 500,
  },
}