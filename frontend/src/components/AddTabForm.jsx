import { useState } from 'react'
import { createTab } from '../api/client.js'

export default function AddTabForm({ onMutate }) {
  const [open, setOpen] = useState(false)
  const [name, setName] = useState('')

  async function handleSubmit(e) {
    e.preventDefault()
    if (!name.trim()) return
    await createTab(name.trim())
    setName('')
    setOpen(false)
    onMutate()
  }

  function handleCancel() {
    setName('')
    setOpen(false)
  }

  if (!open) {
    return (
      <button
        onClick={() => setOpen(true)}
        className="bg-white/20 hover:bg-white/30 text-white rounded-xl px-4 py-3 text-sm font-medium w-72 shrink-0 text-left transition-colors"
      >
        + Add another list
      </button>
    )
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="bg-gray-100 rounded-xl w-72 shrink-0 p-3 flex flex-col gap-2"
    >
      <input
        autoFocus
        className="border border-gray-300 rounded-lg px-3 py-2 text-sm outline-none focus:border-blue-400"
        placeholder="List name…"
        value={name}
        onChange={e => setName(e.target.value)}
        onKeyDown={e => e.key === 'Escape' && handleCancel()}
      />
      <div className="flex gap-2">
        <button
          type="submit"
          disabled={!name.trim()}
          className="px-3 py-1.5 text-sm rounded-lg bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50"
        >
          Add list
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
