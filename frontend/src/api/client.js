const BASE = '/api/v1'

async function request(method, path, body) {
  const res = await fetch(`${BASE}${path}`, {
    method,
    headers: body !== undefined ? { 'Content-Type': 'application/json' } : {},
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })
  const json = await res.json()
  if (!res.ok) throw new Error(json.error ?? 'request failed')
  return json.data
}

export const getTabs = () => request('GET', '/tabs')
export const createTab = (name) => request('POST', '/tabs', { name })
export const renameTab = (id, name) => request('PATCH', `/tabs/${id}`, { name })
export const deleteTab = (id) => request('DELETE', `/tabs/${id}`)

export const createCard = (tab_id, title, description = '') =>
  request('POST', '/cards', { tab_id, title, description })
export const updateCard = (id, fields) => request('PATCH', `/cards/${id}`, fields)
export const deleteCard = (id) => request('DELETE', `/cards/${id}`)
