import { useEffect, useState } from 'react'
import { getMyServices, createService, deleteService } from '../api/client'

export default function MyServices() {
  const [services, setServices] = useState([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [deleteId, setDeleteId] = useState(null)
  const [actionLoading, setActionLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const [form, setForm] = useState({
    name: '',
    description: '',
    price: '',
    duration_minutes: '',
  })

  function fetchServices() {
    setLoading(true)
    getMyServices(1, 50)
      .then(data => setServices(data.data || []))
      .catch(() => setServices([]))
      .finally(() => setLoading(false))
  }

  useEffect(() => { fetchServices() }, [])

  function onChange(e) {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  async function onCreateSubmit(e) {
    e.preventDefault()
    setError('')
    setSuccess('')
    setActionLoading(true)
    try {
      await createService({
        name: form.name,
        description: form.description,
        price: parseFloat(form.price),
        duration_minutes: parseInt(form.duration_minutes),
      })
      setSuccess(`Service "${form.name}" is created!`)
      setForm({ name: '', description: '', price: '', duration_minutes: '' })
      setShowForm(false)
      fetchServices()
    } catch (err) {
      setError(err.message ?? 'Error when creating')
    } finally {
      setActionLoading(false)
    }
  }

  async function onDelete(id, name) {
    if (!window.confirm(`Delete service "${name}"?`)) return
    setDeleteId(id)
    try {
      await deleteService(id)
      setSuccess(`Service "${name}" is deleted`)
      fetchServices()
    } catch (err) {
      setError(err.message ?? 'Error when deleting')
    } finally {
      setDeleteId(null)
    }
  }

  return (
    <div>
      <div style={styles.header}>
        <h1 style={styles.h1}>My services</h1>
        <button
          style={styles.createBtn}
          onClick={() => { setShowForm(!showForm); setError(''); setSuccess('') }}
        >
          {showForm ? '✕ Cancel' : '+ Create service'}
        </button>
      </div>

      {success && <div style={styles.successBox}>{success}</div>}
      {error   && <div style={styles.errorBox}>{error}</div>}

      {}
      {showForm && (
        <div style={styles.formCard}>
          <h3 style={styles.formTitle}>New service</h3>
          <form onSubmit={onCreateSubmit}>
            <label style={styles.label}>Name *</label>
            <input
              style={styles.input}
              name="name"
              value={form.name}
              onChange={onChange}
              placeholder="For instance: Hair cutting, Consultation..."
              required
            />

            <label style={styles.label}>Description</label>
            <textarea
              style={{ ...styles.input, resize: 'vertical' }}
              name="description"
              value={form.description}
              onChange={onChange}
              placeholder="A brief description of service"
              rows={2}
            />

            <div style={styles.row}>
              <div style={{ flex: 1 }}>
                <label style={styles.label}>Price (₸) *</label>
                <input
                  style={styles.input}
                  name="price"
                  type="number"
                  min="0"
                  step="0.01"
                  value={form.price}
                  onChange={onChange}
                  placeholder="1000"
                  required
                />
              </div>
              <div style={{ flex: 1 }}>
                <label style={styles.label}>Duration (min) *</label>
                <input
                  style={styles.input}
                  name="duration_minutes"
                  type="number"
                  min="5"
                  step="5"
                  value={form.duration_minutes}
                  onChange={onChange}
                  placeholder="30"
                  required
                />
              </div>
            </div>

            <button style={styles.submitBtn} type="submit" disabled={actionLoading}>
              {actionLoading ? 'Creating...' : 'Create service'}
            </button>
          </form>
        </div>
      )}

      {}
      {loading && <p style={styles.muted}>Downloading...</p>}

      {!loading && services.length === 0 && !showForm && (
        <div style={styles.empty}>
          <p>You do not have services yet.</p>
          <button style={styles.createBtn} onClick={() => setShowForm(true)}>
            + Create first service
          </button>
        </div>
      )}

      <div style={styles.list}>
        {services.map(svc => (
          <div key={svc.id} style={styles.card}>
            <div style={styles.cardLeft}>
              <div style={styles.cardName}>{svc.name}</div>
              {svc.description && (
                <div style={styles.cardDesc}>{svc.description}</div>
              )}
              <div style={styles.cardMeta}>
                ⏱ {svc.duration_minutes} min &nbsp;·&nbsp; {svc.price} ₸
                &nbsp;·&nbsp;
                <span style={{ color: svc.is_active ? '#059669' : '#dc2626' }}>
                  {svc.is_active ? '● Active' : '● Inactive'}
                </span>
              </div>
            </div>
            <button
              style={styles.deleteBtn}
              onClick={() => onDelete(svc.id, svc.name)}
              disabled={deleteId === svc.id}
            >
              {deleteId === svc.id ? '...' : 'Delete'}
            </button>
          </div>
        ))}
      </div>
    </div>
  )
}

const styles = {
  header: {
    display: 'flex', justifyContent: 'space-between',
    alignItems: 'center', marginBottom: 20,
  },
  h1: { fontSize: 24, fontWeight: 700, color: '#111', margin: 0 },
  createBtn: {
    padding: '9px 18px', background: '#e94560',
    color: '#fff', border: 'none', borderRadius: 8,
    fontSize: 14, cursor: 'pointer', fontWeight: 500,
  },
  formCard: {
    background: '#fff', border: '1px solid #e5e7eb',
    borderRadius: 12, padding: '20px 24px', marginBottom: 24,
  },
  formTitle: { fontSize: 16, fontWeight: 600, color: '#111', marginBottom: 16, marginTop: 0 },
  label: { display: 'block', fontSize: 13, color: '#555', marginBottom: 6, marginTop: 14 },
  input: {
    display: 'block', width: '100%', padding: '9px 12px',
    border: '1px solid #d1d5db', borderRadius: 8, fontSize: 14,
    boxSizing: 'border-box', outline: 'none',
  },
  row: { display: 'flex', gap: 16 },
  submitBtn: {
    marginTop: 18, padding: '10px 24px', background: '#e94560',
    color: '#fff', border: 'none', borderRadius: 8,
    fontSize: 14, cursor: 'pointer', fontWeight: 600,
  },
  list: { display: 'flex', flexDirection: 'column', gap: 10 },
  card: {
    background: '#fff', border: '1px solid #e5e7eb',
    borderRadius: 12, padding: '14px 18px',
    display: 'flex', justifyContent: 'space-between', alignItems: 'center',
  },
  cardLeft: { display: 'flex', flexDirection: 'column', gap: 4 },
  cardName: { fontWeight: 600, fontSize: 15, color: '#111' },
  cardDesc: { fontSize: 13, color: '#666' },
  cardMeta: { fontSize: 12, color: '#888' },
  deleteBtn: {
    padding: '6px 14px', background: 'transparent',
    color: '#dc2626', border: '1px solid #dc2626',
    borderRadius: 8, fontSize: 13, cursor: 'pointer', fontWeight: 500,
    whiteSpace: 'nowrap',
  },
  successBox: {
    background: '#ecfdf5', border: '1px solid #6ee7b7',
    color: '#065f46', borderRadius: 8, padding: '10px 14px',
    fontSize: 13, marginBottom: 16,
  },
  errorBox: {
    background: '#fef2f2', border: '1px solid #fecaca',
    color: '#b91c1c', borderRadius: 8, padding: '10px 14px',
    fontSize: 13, marginBottom: 16,
  },
  muted: { color: '#888', fontSize: 14 },
  empty: { textAlign: 'center', padding: '40px 0', color: '#888' },
}