import { useState } from 'react'
import { renameTab } from '../api/client.js'

export default function TabHeader({ tab, onRename, onDelete }) {
  const [editing, setEditing] = useState(false)
  const [name, setName] = useState(tab.name)

  async function handleSubmit(e) {
    e.preventDefault()
    if (!name.trim()) return
    await renameTab(tab.id, name.trim())
    setEditing(false)
    onRename()
  }

  function handleKeyDown(e) {
    if (e.key === 'Escape') {
      setName(tab.name)
      setEditing(false)
    }
  }

  return (
    <div className="flex items-center justify-between px-3 pt-3 pb-1">
      {editing ? (
        <form onSubmit={handleSubmit} className="flex-1 mr-2">
          <input
            autoFocus
            className="w-full rounded px-2 py-1 text-sm font-semibold border border-blue-400 outline-none"
            value={name}
            onChange={e => setName(e.target.value)}
            onKeyDown={handleKeyDown}
            onBlur={handleSubmit}
          />
        </form>
      ) : (
        <span
          className="font-semibold text-gray-800 cursor-pointer flex-1 truncate"
          onDoubleClick={() => setEditing(true)}
          title="Double-click to rename"
        >
          {tab.name}
        </span>
      )}
      <button
        onClick={onDelete}
        className="text-gray-400 hover:text-red-500 text-lg leading-none ml-1 shrink-0"
        title="Delete tab"
      >
        ×
      </button>
    </div>
  )
}
