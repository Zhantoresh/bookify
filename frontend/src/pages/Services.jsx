import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { getServices } from '../api/client'
import { useAuth } from '../context/AuthContext'
import BookModal from '../components/BookModal'

export default function Services() {
  const { user } = useAuth()
  const navigate = useNavigate()

  const [services, setServices] = useState([])
  const [pagination, setPagination] = useState({ page: 1, pages: 1 })
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)
  const [selectedService, setSelectedService] = useState(null) 

  function fetchServices(page = 1) {
    setLoading(true)
    getServices({ page, limit: 9, search })
      .then(data => {
        setServices(data.services || data.data || [])
        setPagination(data.pagination || { page: 1, pages: 1 })
      })
      .catch(() => setServices([]))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    fetchServices(1)
  }, [])

  function onSearch(e) {
    e.preventDefault()
    fetchServices(1)
  }

  return (
    <div>
      <h1 style={styles.h1}>Услуги</h1>

      {/* Поиск */}
      <form onSubmit={onSearch} style={styles.searchRow}>
        <input
          style={styles.searchInput}
          placeholder="Поиск по названию..."
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        <button style={styles.searchBtn} type="submit">Найти</button>
      </form>

      {loading && <p style={styles.muted}>Загрузка...</p>}

      {!loading && services.length === 0 && (
        <p style={styles.muted}>Услуги не найдены.</p>
      )}

      {/* Сетка карточек */}
      <div style={styles.grid}>
        {services.map(svc => (
          <div key={svc.id} style={styles.card}>
            <div style={styles.cardTop}>
              <div style={styles.cardName}>{svc.name}</div>
              <div style={styles.cardPrice}>{svc.price} ₸</div>
            </div>
            {svc.description && (
              <div style={styles.cardDesc}>{svc.description}</div>
            )}
            <div style={styles.cardMeta}>
              ⏱ {svc.duration_minutes} мин &nbsp;·&nbsp; 👤 {svc.provider_name || 'Провайдер'}
            </div>
            {user?.role === 'client' && (
              <button
                style={styles.bookBtn}
                onClick={() => setSelectedService(svc)}
              >
                Записаться
              </button>
            )}
          </div>
        ))}
      </div>

      {/* Пагинация */}
      {pagination.pages > 1 && (
        <div style={styles.pagination}>
          {Array.from({ length: pagination.pages }, (_, i) => i + 1).map(p => (
            <button
              key={p}
              style={{
                ...styles.pageBtn,
                ...(p === pagination.page ? styles.pageBtnActive : {}),
              }}
              onClick={() => fetchServices(p)}
            >
              {p}
            </button>
          ))}
        </div>
      )}

      {/* Модалка бронирования */}
      {selectedService && (
        <BookModal
          service={selectedService}
          onClose={() => setSelectedService(null)}
          onSuccess={() => {
            setSelectedService(null)
            navigate('/my-appointments')
          }}
        />
      )}
    </div>
  )
}

const styles = {
  h1: { fontSize: 24, fontWeight: 700, color: '#111', marginBottom: 20 },
  searchRow: { display: 'flex', gap: 10, marginBottom: 24 },
  searchInput: {
    flex: 1,
    padding: '10px 14px',
    border: '1px solid #d1d5db',
    borderRadius: 8,
    fontSize: 14,
    outline: 'none',
  },
  searchBtn: {
    padding: '10px 20px',
    background: '#e94560',
    color: '#fff',
    border: 'none',
    borderRadius: 8,
    fontSize: 14,
    cursor: 'pointer',
    fontWeight: 500,
  },
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(260px, 1fr))',
    gap: 16,
  },
  card: {
    background: '#fff',
    border: '1px solid #e5e7eb',
    borderRadius: 12,
    padding: '18px',
    display: 'flex',
    flexDirection: 'column',
    gap: 8,
  },
  cardTop: { display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' },
  cardName: { fontWeight: 600, fontSize: 15, color: '#111' },
  cardPrice: { fontWeight: 700, fontSize: 15, color: '#e94560' },
  cardDesc: { fontSize: 13, color: '#666', lineHeight: 1.5 },
  cardMeta: { fontSize: 12, color: '#888' },
  bookBtn: {
    marginTop: 8,
    padding: '8px 0',
    background: '#e94560',
    color: '#fff',
    border: 'none',
    borderRadius: 8,
    fontSize: 13,
    cursor: 'pointer',
    fontWeight: 500,
  },
  muted: { color: '#888', fontSize: 14 },
  pagination: { display: 'flex', gap: 8, marginTop: 24, justifyContent: 'center' },
  pageBtn: {
    padding: '6px 12px',
    border: '1px solid #d1d5db',
    borderRadius: 6,
    background: '#fff',
    cursor: 'pointer',
    fontSize: 13,
  },
  pageBtnActive: { background: '#e94560', color: '#fff', borderColor: '#e94560' },
}