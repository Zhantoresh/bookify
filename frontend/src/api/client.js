const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'



function getToken() {
  return localStorage.getItem('token')
}

function saveToken(token) {
  localStorage.setItem('token', token)
}

function removeToken() {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
}

function saveUser(user) {
  localStorage.setItem('user', JSON.stringify(user))
}

function getUser() {
  const raw = localStorage.getItem('user')
  return raw ? JSON.parse(raw) : null
}

async function request(method, path, body = null, auth = true) {
  const headers = { 'Content-Type': 'application/json' }

  if (auth) {
    const token = getToken()
    if (token) headers['Authorization'] = `Bearer ${token}`
  }

  const options = { method, headers }
  if (body) options.body = JSON.stringify(body)

  const res = await fetch(`${BASE_URL}${path}`, options)

  if (res.status === 401) {
    removeToken()
    window.location.href = '/login'
    return
  }

  const data = await res.json().catch(() => ({}))

  if (!res.ok) {
    const message = data?.error?.message || data?.error || data?.message || 'Something went wrong'
    throw new Error(message)
  }

  return data
}



export async function register({ email, password, full_name, phone, role }) {
  const data = await request('POST', '/api/v1/auth/register', {
    email,
    password,
    full_name,
    phone,
    role, 
  }, false)
  return data
}

export async function login({ email, password }) {
  const data = await request('POST', '/api/v1/auth/login', { email, password }, false)
  saveToken(data.token)
  saveUser(data.user)
  return data
}

export function logout() {
  removeToken()
  window.location.href = '/login'
}

export function isLoggedIn() {
  return !!getToken()
}

export { getUser }


export async function getServices(params = {}) {
  const query = new URLSearchParams()
  if (params.page)      query.set('page', params.page)
  if (params.limit)     query.set('limit', params.limit)
  if (params.search)    query.set('search', params.search)
  if (params.min_price) query.set('min_price', params.min_price)
  if (params.max_price) query.set('max_price', params.max_price)

  const qs = query.toString() ? `?${query.toString()}` : ''
  return request('GET', `/api/v1/services${qs}`)
}


export async function getServiceById(id) {
  return request('GET', `/api/v1/services/${id}`)
}

export async function getMyServices(page = 1, limit = 10) {
  return request('GET', `/api/v1/services/my?page=${page}&limit=${limit}`)
}


export async function createService({ name, description, price, duration_minutes }) {
  return request('POST', '/api/v1/services', {
    name,
    description,
    price,
    duration_minutes,
  })
}


export async function deleteService(id) {
  return request('DELETE', `/api/v1/services/${id}`)
}


export async function getAvailableSlots(serviceId, date) {
  return request(
    'GET',
    `/api/v1/appointments/available-slots?service_id=${serviceId}&date=${date}`,
    null,
    false 
  )
}


export async function createAppointment({ service_id, start_time, notes }) {
  return request('POST', '/api/v1/appointments', {
    service_id,
    start_time,
    notes,
  })
}


export async function getMyAppointments(page = 1, limit = 10) {
  return request('GET', `/api/v1/appointments/my?page=${page}&limit=${limit}`)
}


export async function cancelAppointment(id, reason = '') {
  return request('PATCH', `/api/v1/appointments/${id}/cancel`, { reason })
}

export async function confirmAppointment(id) {
  return request('PATCH', `/api/v1/appointments/${id}/confirm`)
}


export async function completeAppointment(id) {
  return request('PATCH', `/api/v1/appointments/${id}/complete`)
}


export async function getMe() {
  return request('GET', '/api/v1/users/me')
}


export async function getAdminDashboard() {
  return request('GET', '/api/v1/admin/dashboard')
}

export async function getAdminUsers(role = '') {
  const qs = role ? `?role=${role}` : ''
  return request('GET', `/api/v1/admin/users${qs}`)
}

export async function updateUserRole(id, role) {
  return request('PATCH', `/api/v1/admin/users/${id}`, { role })
}

export async function deleteUser(id) {
  return request('DELETE', `/api/v1/admin/users/${id}`)
}