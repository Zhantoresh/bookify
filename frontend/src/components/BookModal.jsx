import { useState, useEffect } from 'react'
import { getAvailableSlots, createAppointment } from '../api/client'

export default function BookModal({ service, onClose, onSuccess }) {
  const today = new Date().toISOString().split('T')[0]
  const [date, setDate] = useState(today)
  const [slots, setSlots] = useState([])
  const [selectedSlot, setSelectedSlot] = useState(null)
  const [notes, setNotes] = useState('')
  const [loading, setLoading] = useState(false)
  const [slotsLoading, setSlotsLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!date) return
    setSlotsLoading(true)
    getAvailableSlots(service.id, date)
      .then(data => {
        const raw = data.available_slots ?? data.slots ?? data.data ?? data ?? []
        setSlots(Array.isArray(raw) ? raw : [])
})
      .catch(() => setSlots([]))
      .finally(() => setSlotsLoading(false))
  }, [date, service.id])

  async function onBook() {
    if (!selectedSlot) return
    setError('')
    setLoading(true)
    try {
      await createAppointment({
        service_id: service.id,
        start_time: `${date}T${selectedSlot.start_time}:00Z`,
        notes,
      })
      onSuccess()
    } catch (err) {
      const msg = typeof err === 'string' ? err
        : err?.message ?? 'Something went wrong'

      if (msg.toLowerCase().includes('already booked') || msg.toLowerCase().includes('overlap') || msg.toLowerCase().includes('conflict')) {
        setError('That slot is already taken. Please select another slot..')
      } else if (msg.toLowerCase().includes('past')) {
        setError('You cannot use the past tense..')
      } else {
        setError(msg)
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={styles.overlay} onClick={e => e.target === e.currentTarget && onClose()}>
      <div style={styles.modal}>
        <div style={styles.header}>
          <h3 style={styles.title}>Booking to: {service.name}</h3>
          <button style={styles.close} onClick={onClose}>✕</button>
        </div>

        <div style={styles.meta}>
          ⏱ {service.duration_minutes} min &nbsp;·&nbsp; {service.price} ₸
        </div>

        <label style={styles.label}>Choose a date</label>
        <input
          style={styles.input}
          type="date"
          value={date}
          min={today}
          onChange={e => { setDate(e.target.value); setSelectedSlot(null) }}
        />

        <label style={styles.label}>Available slots</label>

        {slotsLoading && <p style={styles.muted}>Downloading slots...</p>}

        {!slotsLoading && slots.length === 0 && (
          <p style={styles.muted}>There are no available slots for this date..</p>
        )}

        <div style={styles.slotsGrid}>
          {slots.map((slot, i) => {
            const time = slot.start_time
            const isSelected = selectedSlot?.start_time === slot.start_time
            return (
              <button
                key={i}
                style={{ ...styles.slotBtn, ...(isSelected ? styles.slotActive : {}) }}
                onClick={() => setSelectedSlot(slot)}
              >
                {time}
              </button>
            )
          })}
        </div>

        <label style={styles.label}>Notes (Not necessary)</label>
        <textarea
          style={styles.textarea}
          value={notes}
          onChange={e => setNotes(e.target.value)}
          placeholder="Any wish..."
          rows={2}
        />

        {error && <div style={styles.error}>{error}</div>}

        <button
          style={{ ...styles.bookBtn, opacity: (!selectedSlot || loading) ? 0.6 : 1 }}
          onClick={onBook}
          disabled={!selectedSlot || loading}
        >
          {loading ? 'Booking...' : 'Approve booking'}
        </button>
      </div>
    </div>
  )
}

const styles = {
  overlay: {
    position: 'fixed', inset: 0,
    background: 'rgba(0,0,0,0.5)',
    display: 'flex', alignItems: 'center', justifyContent: 'center',
    zIndex: 200, padding: 16,
  },
  modal: {
    background: '#fff',
    borderRadius: 14,
    padding: '24px',
    width: '100%',
    maxWidth: 440,
    maxHeight: '90vh',
    overflowY: 'auto',
  },
  header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 },
  title: { fontSize: 17, fontWeight: 600, color: '#111', margin: 0 },
  close: { background: 'none', border: 'none', fontSize: 18, cursor: 'pointer', color: '#888' },
  meta: { fontSize: 13, color: '#888', marginBottom: 18 },
  label: { display: 'block', fontSize: 13, color: '#555', marginBottom: 6, marginTop: 14 },
  input: {
    display: 'block', width: '100%', padding: '9px 12px',
    border: '1px solid #d1d5db', borderRadius: 8, fontSize: 14,
    boxSizing: 'border-box', outline: 'none',
  },
  slotsGrid: { display: 'flex', flexWrap: 'wrap', gap: 8, marginBottom: 4 },
  slotBtn: {
    padding: '7px 14px', border: '1px solid #d1d5db',
    borderRadius: 8, fontSize: 13, cursor: 'pointer', background: '#fff',
  },
  slotActive: { background: '#e94560', color: '#fff', borderColor: '#e94560' },
  textarea: {
    display: 'block', width: '100%', padding: '9px 12px',
    border: '1px solid #d1d5db', borderRadius: 8, fontSize: 14,
    boxSizing: 'border-box', outline: 'none', resize: 'vertical',
  },
  bookBtn: {
    width: '100%', marginTop: 16, padding: '11px 0',
    background: '#e94560', color: '#fff', border: 'none',
    borderRadius: 8, fontSize: 15, cursor: 'pointer', fontWeight: 600,
  },
  error: {
    background: '#fef2f2', border: '1px solid #fecaca',
    color: '#b91c1c', borderRadius: 8, padding: '10px 14px', fontSize: 13, marginTop: 12,
  },
  muted: { color: '#888', fontSize: 13 },
}