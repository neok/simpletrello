import { useState, useRef, useEffect } from 'react'
import { updateCard } from '../api/client.js'

export default function MoveCardMenu({ card, allTabs, onMutate }) {
  const [open, setOpen] = useState(false)
  const ref = useRef(null)

  useEffect(() => {
    if (!open) return
    function handleClick(e) {
      if (ref.current && !ref.current.contains(e.target)) setOpen(false)
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [open])

  async function handleMove(e, targetTabId) {
    e.stopPropagation()
    setOpen(false)
    await updateCard(card.id, { tab_id: targetTabId })
    onMutate()
  }

  const otherTabs = allTabs.filter(t => t.id !== card.tab_id)

  if (otherTabs.length === 0) return null

  return (
    <div ref={ref} className="relative" onClick={e => e.stopPropagation()}>
      <button
        onClick={e => { e.stopPropagation(); setOpen(v => !v) }}
        className="text-gray-300 hover:text-blue-500 text-xs leading-none px-0.5"
        title="Move to tab"
      >
        ⇄
      </button>
      {open && (
        <div className="absolute right-0 top-5 bg-white border border-gray-200 rounded-lg shadow-lg z-10 min-w-32 py-1">
          {otherTabs.map(t => (
            <button
              key={t.id}
              onClick={e => handleMove(e, t.id)}
              className="block w-full text-left px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50 truncate"
            >
              {t.name}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
