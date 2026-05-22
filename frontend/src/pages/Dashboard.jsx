import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { getMyAppointments } from '../api/client'

const STATUS_LABELS = {
  pending:   { label: 'Waiting',   color: '#d97706', bg: '#fffbeb' },
  confirmed: { label: 'Approved', color: '#059669', bg: '#ecfdf5' },
  cancelled: { label: 'Cancelled',   color: '#dc2626', bg: '#fef2f2' },
  completed: { label: 'Finished',  color: '#6b7280', bg: '#f9fafb' },
}

export default function Dashboard() {
  const { user } = useAuth()
  const [appointments, setAppointments] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getMyAppointments(1, 5)
      .then(data => setAppointments(data.appointments || data.data || []))
      .catch(() => setAppointments([]))
      .finally(() => setLoading(false))
  }, [])

  const roleLabel = user?.role === 'provider' ? 'Provider' : user?.role === 'admin' ? 'Admin' : 'Client'

  return (
    <div>
      {}
      <div style={styles.welcome}>
        <div>
          <h1 style={styles.h1}>Hello, {user?.full_name} 👋</h1>
          <span style={styles.roleBadge}>{roleLabel}</span>
        </div>
      </div>

      {}
      <div style={styles.grid}>
        <Link to="/services" style={styles.actionCard}>
          <div style={styles.actionIcon}>🔍</div>
          <div style={styles.actionTitle}>Find a service</div>
          <div style={styles.actionSub}>View all available services</div>
        </Link>

        <Link to="/my-appointments" style={styles.actionCard}>
          <div style={styles.actionIcon}>📅</div>
          <div style={styles.actionTitle}>My appoinments</div>
          <div style={styles.actionSub}>All my bookings</div>
        </Link>

        {user?.role === 'provider' && (
          <Link to="/services" style={styles.actionCard}>
            <div style={styles.actionIcon}>➕</div>
            <div style={styles.actionTitle}>My services</div>
            <div style={styles.actionSub}>Service management</div>
          </Link>
        )}
      </div>

      {}
      <h2 style={styles.h2}>Last bookings</h2>

      {loading && <p style={styles.muted}>Downloading...</p>}

      {!loading && appointments.length === 0 && (
        <div style={styles.empty}>
          <p>You do not have bookings yet.</p>
          <Link to="/services" style={styles.link}>Find service →</Link>
        </div>
      )}

      {!loading && appointments.map(appt => {
        const s = STATUS_LABELS[appt.status] || STATUS_LABELS.pending
        const date = new Date(appt.start_time).toLocaleString('ru-RU', {
          day: 'numeric', month: 'long', hour: '2-digit', minute: '2-digit',
        })
        return (
          <div key={appt.id} style={styles.apptCard}>
            <div style={styles.apptLeft}>
              <div style={styles.apptService}>{appt.service_name}</div>
              <div style={styles.apptMeta}>
                {user?.role === 'client' ? `Provider: ${appt.provider_name}` : `Client: ${appt.client_name}`}
              </div>
              <div style={styles.apptMeta}>{date}</div>
            </div>
            <span style={{ ...styles.statusBadge, color: s.color, background: s.bg }}>
              {s.label}
            </span>
          </div>
        )
      })}

      {!loading && appointments.length > 0 && (
        <Link to="/my-appointments" style={styles.link}>All bookings →</Link>
      )}
    </div>
  )
}

const styles = {
  welcome: {
    marginBottom: 32,
  },
  h1: {
    fontSize: 26,
    fontWeight: 700,
    color: '#111',
    marginBottom: 8,
  },
  h2: {
    fontSize: 18,
    fontWeight: 600,
    color: '#111',
    margin: '32px 0 16px',
  },
  roleBadge: {
    background: '#eff6ff',
    color: '#1d4ed8',
    padding: '3px 12px',
    borderRadius: 20,
    fontSize: 12,
    fontWeight: 500,
  },
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
    gap: 16,
    marginBottom: 8,
  },
  actionCard: {
    background: '#fff',
    border: '1px solid #e5e7eb',
    borderRadius: 12,
    padding: '20px',
    textDecoration: 'none',
    transition: 'border-color 0.2s',
  },
  actionIcon: {
    fontSize: 28,
    marginBottom: 10,
  },
  actionTitle: {
    fontSize: 15,
    fontWeight: 600,
    color: '#111',
    marginBottom: 4,
  },
  actionSub: {
    fontSize: 13,
    color: '#888',
  },
  apptCard: {
    background: '#fff',
    border: '1px solid #e5e7eb',
    borderRadius: 10,
    padding: '14px 18px',
    marginBottom: 10,
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  apptLeft: {
    display: 'flex',
    flexDirection: 'column',
    gap: 3,
  },
  apptService: {
    fontWeight: 600,
    fontSize: 14,
    color: '#111',
  },
  apptMeta: {
    fontSize: 12,
    color: '#888',
  },
  statusBadge: {
    padding: '4px 12px',
    borderRadius: 20,
    fontSize: 12,
    fontWeight: 500,
    whiteSpace: 'nowrap',
  },
  muted: {
    color: '#888',
    fontSize: 14,
  },
  empty: {
    textAlign: 'center',
    padding: '32px 0',
    color: '#888',
    fontSize: 14,
  },
  link: {
    color: '#e94560',
    textDecoration: 'none',
    fontSize: 14,
    fontWeight: 500,
  },
}