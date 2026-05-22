import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { getAdminDashboard } from '../api/client'

export default function AdminDashboard() {
  const [data, setData] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    getAdminDashboard()
      .then(setData)
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <p style={styles.muted}>Downloading...</p>
  if (error)   return <p style={{ color: '#dc2626' }}>{error}</p>
  if (!data)   return null

  const s = data.summary

  return (
    <div>
      <div style={styles.header}>
        <h1 style={styles.h1}>Admin pane</h1>
        <Link to="/admin/users" style={styles.usersBtn}>Users management →</Link>
      </div>

      {}
      <h2 style={styles.h2}>Users</h2>
      <div style={styles.grid}>
        <StatCard label="Total"        value={s.total_users}    color="#6366f1" />
        <StatCard label="Clients"     value={s.total_clients}  color="#0ea5e9" />
        <StatCard label="Providers"  value={s.total_providers} color="#10b981" />
        <StatCard label="Admins"      value={s.total_admins}   color="#f59e0b" />
      </div>

      {/* Услуги */}
      <h2 style={styles.h2}>Services</h2>
      <div style={styles.grid}>
        <StatCard label="Total services"   value={s.total_services}  color="#8b5cf6" />
        <StatCard label="Active"      value={s.active_services} color="#10b981" />
      </div>

      {}
      <h2 style={styles.h2}>Bookings</h2>
      <div style={styles.grid}>
        <StatCard label="Total"        value={s.total_appointments}     color="#64748b" />
        <StatCard label="Waiting"      value={s.pending_appointments}   color="#f59e0b" />
        <StatCard label="Approved" value={s.confirmed_appointments} color="#10b981" />
        <StatCard label="Cancelled"     value={s.cancelled_appointments} color="#ef4444" />
        <StatCard label="Finished"    value={s.completed_appointments} color="#6366f1" />
      </div>

      {}
      {data.recent_users?.length > 0 && (
        <>
          <h2 style={styles.h2}>Last users</h2>
          <div style={styles.table}>
            <div style={styles.tableHeader}>
              <span>Name</span><span>Email</span><span>Role</span>
            </div>
            {data.recent_users.map(u => (
              <div key={u.id} style={styles.tableRow}>
                <span style={{ fontWeight: 500 }}>{u.full_name}</span>
                <span style={{ color: '#666' }}>{u.email}</span>
                <RoleBadge role={u.role} />
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  )
}

function StatCard({ label, value, color }) {
  return (
    <div style={{ ...styles.card, borderTop: `3px solid ${color}` }}>
      <div style={{ fontSize: 28, fontWeight: 700, color }}>{value ?? 0}</div>
      <div style={{ fontSize: 13, color: '#888', marginTop: 4 }}>{label}</div>
    </div>
  )
}

function RoleBadge({ role }) {
  const map = {
    admin:    { label: 'Admin',     bg: '#fef3c7', color: '#92400e' },
    provider: { label: 'Provider', bg: '#d1fae5', color: '#065f46' },
    client:   { label: 'Client',    bg: '#dbeafe', color: '#1e40af' },
  }
  const s = map[role] || map.client
  return (
    <span style={{ background: s.bg, color: s.color, padding: '2px 10px', borderRadius: 20, fontSize: 12, fontWeight: 500 }}>
      {s.label}
    </span>
  )
}

const styles = {
  header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 },
  h1: { fontSize: 24, fontWeight: 700, color: '#111', margin: 0 },
  h2: { fontSize: 16, fontWeight: 600, color: '#444', margin: '24px 0 12px' },
  usersBtn: {
    padding: '8px 16px', background: '#1a1a2e', color: '#fff',
    borderRadius: 8, textDecoration: 'none', fontSize: 14,
  },
  grid: { display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(140px, 1fr))', gap: 12 },
  card: {
    background: '#fff', border: '1px solid #e5e7eb',
    borderRadius: 10, padding: '16px',
  },
  table: { background: '#fff', border: '1px solid #e5e7eb', borderRadius: 10, overflow: 'hidden' },
  tableHeader: {
    display: 'grid', gridTemplateColumns: '1fr 1fr 100px',
    padding: '10px 16px', background: '#f9fafb',
    fontSize: 12, fontWeight: 600, color: '#888',
    borderBottom: '1px solid #e5e7eb',
  },
  tableRow: {
    display: 'grid', gridTemplateColumns: '1fr 1fr 100px',
    padding: '10px 16px', fontSize: 13,
    borderBottom: '1px solid #f3f4f6', alignItems: 'center',
  },
  muted: { color: '#888', fontSize: 14 },
}