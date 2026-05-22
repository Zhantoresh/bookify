import { useEffect, useState } from 'react'
import { getAdminUsers, updateUserRole, deleteUser } from '../api/client'
import { useAuth } from '../context/AuthContext'

const ROLES = ['client', 'provider']

const ROLE_LABELS = {
  admin:    { label: 'Admin',     bg: '#fef3c7', color: '#92400e' },
  provider: { label: 'Provider', bg: '#d1fae5', color: '#065f46' },
  client:   { label: 'Client',    bg: '#dbeafe', color: '#1e40af' },
}

export default function AdminUsers() {
  const { user: me } = useAuth()
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [filterRole, setFilterRole] = useState('')
  const [actionId, setActionId] = useState(null)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  function fetchUsers() {
    setLoading(true)
    getAdminUsers(filterRole)
      .then(data => setUsers(data.data || []))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => { fetchUsers() }, [filterRole])

  async function onRoleChange(id, newRole) {
    setActionId(id)
    setError('')
    setSuccess('')
    try {
      await updateUserRole(id, newRole)
      setSuccess('Role is updated')
      fetchUsers()
    } catch (err) {
      setError(err.message)
    } finally {
      setActionId(null)
    }
  }

  async function onDelete(id, name) {
    if (!window.confirm(`Delete user "${name}"? This action is irreversible.`)) return
    setActionId(id)
    setError('')
    setSuccess('')
    try {
      await deleteUser(id)
      setSuccess(`User "${name}" is deleted`)
      fetchUsers()
    } catch (err) {
      setError(err.message)
    } finally {
      setActionId(null)
    }
  }

  return (
    <div>
      <div style={styles.header}>
        <h1 style={styles.h1}>Users</h1>
        <span style={styles.count}>{users.length} people.</span>
      </div>

      {success && <div style={styles.successBox}>{success}</div>}
      {error   && <div style={styles.errorBox}>{error}</div>}

      {}
      <div style={styles.filters}>
        <button
          style={{ ...styles.filterBtn, ...(filterRole === '' ? styles.filterActive : {}) }}
          onClick={() => setFilterRole('')}
        >All</button>
        {ROLES.map(r => (
          <button
            key={r}
            style={{ ...styles.filterBtn, ...(filterRole === r ? styles.filterActive : {}) }}
            onClick={() => setFilterRole(r)}
          >
            {ROLE_LABELS[r].label}
          </button>
        ))}
      </div>

      {loading && <p style={styles.muted}>Downloading...</p>}

      {!loading && users.length === 0 && (
        <p style={styles.muted}>Users not found.</p>
      )}

      {}
      <div style={styles.list}>
        {users.map(u => {
          const rl = ROLE_LABELS[u.role] || ROLE_LABELS.client
          const isMe = u.id === me?.id
          const busy = actionId === u.id

          return (
            <div key={u.id} style={styles.card}>
              {}
              <div style={styles.cardLeft}>
                <div style={styles.avatar}>
                  {u.full_name?.[0]?.toUpperCase() || '?'}
                </div>
                <div>
                  <div style={styles.name}>
                    {u.full_name}
                    {isMe && <span style={styles.meBadge}>It is you</span>}
                  </div>
                  <div style={styles.email}>{u.email}</div>
                  {u.phone && <div style={styles.phone}>{u.phone}</div>}
                </div>
              </div>

              {}
              <div style={styles.cardRight}>
                {}
                <span style={{ ...styles.roleBadge, background: rl.bg, color: rl.color }}>
                  {rl.label}
                </span>

                {}
                {!isMe && (
                  <select
                    style={styles.select}
                    value={u.role}
                    disabled={busy}
                    onChange={e => onRoleChange(u.id, e.target.value)}
                  >
                    {ROLES.map(r => (
                      <option key={r} value={r}>{ROLE_LABELS[r].label}</option>
                    ))}
                  </select>
                )}

                {}
                {!isMe && (
                  <button
                    style={styles.deleteBtn}
                    onClick={() => onDelete(u.id, u.full_name)}
                    disabled={busy}
                  >
                    {busy ? '...' : 'Delete'}
                  </button>
                )}
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}

const styles = {
  header: { display: 'flex', alignItems: 'center', gap: 12, marginBottom: 20 },
  h1: { fontSize: 24, fontWeight: 700, color: '#111', margin: 0 },
  count: { fontSize: 13, color: '#888', background: '#f3f4f6', padding: '3px 10px', borderRadius: 20 },
  filters: { display: 'flex', gap: 8, marginBottom: 20, flexWrap: 'wrap' },
  filterBtn: {
    padding: '6px 14px', background: '#fff',
    border: '1px solid #d1d5db', borderRadius: 20,
    fontSize: 13, cursor: 'pointer', color: '#555',
  },
  filterActive: { background: '#1a1a2e', color: '#fff', borderColor: '#1a1a2e' },
  list: { display: 'flex', flexDirection: 'column', gap: 10 },
  card: {
    background: '#fff', border: '1px solid #e5e7eb',
    borderRadius: 12, padding: '14px 18px',
    display: 'flex', justifyContent: 'space-between',
    alignItems: 'center', gap: 12, flexWrap: 'wrap',
  },
  cardLeft: { display: 'flex', alignItems: 'center', gap: 12 },
  avatar: {
    width: 40, height: 40, borderRadius: '50%',
    background: '#e0e7ff', color: '#4338ca',
    display: 'flex', alignItems: 'center', justifyContent: 'center',
    fontWeight: 700, fontSize: 16, flexShrink: 0,
  },
  name: { fontWeight: 600, fontSize: 14, color: '#111', display: 'flex', alignItems: 'center', gap: 8 },
  email: { fontSize: 12, color: '#888', marginTop: 2 },
  phone: { fontSize: 12, color: '#aaa' },
  meBadge: {
    fontSize: 11, background: '#eff6ff', color: '#1d4ed8',
    padding: '1px 8px', borderRadius: 20, fontWeight: 400,
  },
  cardRight: { display: 'flex', alignItems: 'center', gap: 10, flexWrap: 'wrap' },
  roleBadge: { padding: '3px 12px', borderRadius: 20, fontSize: 12, fontWeight: 500 },
  select: {
    padding: '5px 8px', border: '1px solid #d1d5db',
    borderRadius: 8, fontSize: 13, cursor: 'pointer', outline: 'none',
  },
  deleteBtn: {
    padding: '6px 14px', background: 'transparent',
    color: '#dc2626', border: '1px solid #dc2626',
    borderRadius: 8, fontSize: 13, cursor: 'pointer', fontWeight: 500,
  },
  successBox: {
    background: '#ecfdf5', border: '1px solid #6ee7b7',
    color: '#065f46', borderRadius: 8, padding: '10px 14px', fontSize: 13, marginBottom: 16,
  },
  errorBox: {
    background: '#fef2f2', border: '1px solid #fecaca',
    color: '#b91c1c', borderRadius: 8, padding: '10px 14px', fontSize: 13, marginBottom: 16,
  },
  muted: { color: '#888', fontSize: 14 },
}