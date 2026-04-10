import { useState } from 'react'
import { createCard } from '../api/client.js'

export default function AddCardForm({ tabId, onMutate }) {
  const [open, setOpen] = useState(false)
  const [title, setTitle] = useState('')

  async function handleSubmit(e) {
    e.preventDefault()
    if (!title.trim()) return
    await createCard(tabId, title.trim())
    setTitle('')
    setOpen(false)
    onMutate()
  }

  function handleCancel() {
    setTitle('')
    setOpen(false)
  }

  if (!open) {
    return (
      <button
        onClick={() => setOpen(true)}
        className="text-gray-500 hover:text-gray-800 hover:bg-gray-200 rounded-lg text-sm px-3 py-2 text-left transition-colors m-2"
      >
        + Add a card
      </button>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-2 p-2">
      <input
        autoFocus
        className="border border-gray-300 rounded-lg px-3 py-2 text-sm outline-none focus:border-blue-400"
        placeholder="Card title…"
        value={title}
        onChange={e => setTitle(e.target.value)}
        onKeyDown={e => e.key === 'Escape' && handleCancel()}
      />
      <div className="flex gap-2">
        <button
          type="submit"
          disabled={!title.trim()}
          className="px-3 py-1.5 text-sm rounded-lg bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50"
        >
          Add
        </button>
        <button
          type="button"
          onClick={handleCancel}
          className="px-3 py-1.5 text-sm rounded-lg text-gray-500 hover:bg-gray-200"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
