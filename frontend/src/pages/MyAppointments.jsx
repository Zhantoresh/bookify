import { useEffect, useState } from 'react'
import { getMyAppointments, cancelAppointment, confirmAppointment, completeAppointment } from '../api/client'
import { useAuth } from '../context/AuthContext'

const STATUS_LABELS = {
  pending:   { label: 'Waiting',     color: '#d97706', bg: '#fffbeb' },
  confirmed: { label: 'Approved', color: '#059669', bg: '#ecfdf5' },
  cancelled: { label: 'Cancelled',     color: '#dc2626', bg: '#fef2f2' },
  completed: { label: 'Finished',    color: '#6b7280', bg: '#f9fafb' },
}

export default function MyAppointments() {
  const { user } = useAuth()
  const [appointments, setAppointments] = useState([])
  const [loading, setLoading] = useState(true)
  const [actionLoading, setActionLoading] = useState(null) 

  function fetchAll() {
    setLoading(true)
    getMyAppointments(1, 50)
      .then(data => setAppointments(data.appointments || data.data || []))
      .catch(() => setAppointments([]))
      .finally(() => setLoading(false))
  }

  useEffect(() => { fetchAll() }, [])

  async function handleAction(fn, id) {
    setActionLoading(id)
    try {
      await fn(id)
      fetchAll()
    } catch (err) {
      alert(err.message)
    } finally {
      setActionLoading(null)
    }
  }

  async function handleCancel(id) {
    const reason = prompt('Reason of cancellation (Not necessary):') ?? ''
    setActionLoading(id)
    try {
      await cancelAppointment(id, reason)
      fetchAll()
    } catch (err) {
      alert(err.message)
    } finally {
      setActionLoading(null)
    }
  }

  if (loading) return <p style={styles.muted}>Downloading...</p>

  return (
    <div>
      <h1 style={styles.h1}>My appointmens</h1>

      {appointments.length === 0 && (
        <p style={styles.muted}>You do not have appointments yet.</p>
      )}

      {appointments.map(appt => {
        const s = STATUS_LABELS[appt.status] || STATUS_LABELS.pending
        const start = new Date(appt.start_time).toLocaleString('ru-RU', {
          day: 'numeric', month: 'long', year: 'numeric',
          hour: '2-digit', minute: '2-digit',
        })
        const end = new Date(appt.end_time).toLocaleTimeString('ru-RU', {
          hour: '2-digit', minute: '2-digit',
        })
        const busy = actionLoading === appt.id

        return (
          <div key={appt.id} style={styles.card}>
            <div style={styles.cardHeader}>
              <div style={styles.serviceName}>{appt.service_name}</div>
              <span style={{ ...styles.badge, color: s.color, background: s.bg }}>
                {s.label}
              </span>
            </div>

            <div style={styles.info}>
              <span>🕐 {start} — {end}</span>
            </div>
            <div style={styles.info}>
              {user?.role === 'client'
                ? <span>👤 Provider: {appt.provider_name}</span>
                : <span>👤 Client: {appt.client_name} ({appt.client_email})</span>
              }
            </div>
            {appt.notes && (
              <div style={styles.notes}>📝 {appt.notes}</div>
            )}
            {appt.cancellation_reason && (
              <div style={{ ...styles.notes, color: '#dc2626' }}>
                Reason of cancellation: {appt.cancellation_reason}
              </div>
            )}

            {}
            <div style={styles.actions}>
              {}
              {user?.role === 'client' && appt.status === 'pending' && (
                <button
                  style={styles.cancelBtn}
                  onClick={() => handleCancel(appt.id)}
                  disabled={busy}
                >
                  {busy ? '...' : 'Cancel'}
                </button>
              )}

              {}
              {user?.role === 'provider' && appt.status === 'pending' && (
                <button
                  style={styles.confirmBtn}
                  onClick={() => handleAction(confirmAppointment, appt.id)}
                  disabled={busy}
                >
                  {busy ? '...' : 'Approve'}
                </button>
              )}

              {}
              {user?.role === 'provider' && appt.status === 'confirmed' && (
                <>
                  <button
                    style={styles.confirmBtn}
                    onClick={() => handleAction(completeAppointment, appt.id)}
                    disabled={busy}
                  >
                    {busy ? '...' : 'Finish'}
                  </button>
                  <button
                    style={styles.cancelBtn}
                    onClick={() => handleCancel(appt.id)}
                    disabled={busy}
                  >
                    {busy ? '...' : 'Cancel'}
                  </button>
                </>
              )}
            </div>
          </div>
        )
      })}
    </div>
  )
}

const styles = {
  h1: { fontSize: 24, fontWeight: 700, color: '#111', marginBottom: 20 },
  card: {
    background: '#fff',
    border: '1px solid #e5e7eb',
    borderRadius: 12,
    padding: '16px 20px',
    marginBottom: 12,
  },
  cardHeader: {
    display: 'flex', justifyContent: 'space-between',
    alignItems: 'center', marginBottom: 8,
  },
  serviceName: { fontWeight: 600, fontSize: 15, color: '#111' },
  badge: {
    padding: '3px 12px', borderRadius: 20,
    fontSize: 12, fontWeight: 500,
  },
  info: { fontSize: 13, color: '#666', marginBottom: 4 },
  notes: { fontSize: 12, color: '#888', marginTop: 4, fontStyle: 'italic' },
  actions: { display: 'flex', gap: 10, marginTop: 12 },
  confirmBtn: {
    padding: '7px 16px', background: '#059669',
    color: '#fff', border: 'none', borderRadius: 8,
    fontSize: 13, cursor: 'pointer', fontWeight: 500,
  },
  cancelBtn: {
    padding: '7px 16px', background: 'transparent',
    color: '#dc2626', border: '1px solid #dc2626',
    borderRadius: 8, fontSize: 13, cursor: 'pointer', fontWeight: 500,
  },
  muted: { color: '#888', fontSize: 14 },
}