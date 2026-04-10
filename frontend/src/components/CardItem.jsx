import { useState } from 'react'
import { deleteCard } from '../api/client.js'
import CardModal from './CardModal.jsx'
import MoveCardMenu from './MoveCardMenu.jsx'

export default function CardItem({ card, allTabs, onMutate }) {
  const [modalOpen, setModalOpen] = useState(false)

  async function handleDelete(e) {
    e.stopPropagation()
    await deleteCard(card.id)
    onMutate()
  }

  return (
    <>
      <div
        className="bg-white rounded-lg px-3 py-2 shadow-sm cursor-pointer hover:shadow-md transition-shadow group flex items-start justify-between gap-2"
        onClick={() => setModalOpen(true)}
      >
        <div className="flex-1 min-w-0">
          <p className="text-sm font-medium text-gray-800 truncate">{card.title}</p>
          {card.description && (
            <p className="text-xs text-gray-500 truncate mt-0.5">{card.description}</p>
          )}
        </div>
        <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
          <MoveCardMenu card={card} allTabs={allTabs} onMutate={onMutate} />
          <button
            onClick={handleDelete}
            className="text-gray-300 hover:text-red-500 text-base leading-none"
            title="Delete card"
          >
            ×
          </button>
        </div>
      </div>
      {modalOpen && (
        <CardModal card={card} onClose={() => setModalOpen(false)} onMutate={onMutate} />
      )}
    </>
  )
}
