import { deleteTab } from '../api/client.js'
import TabHeader from './TabHeader.jsx'
import CardItem from './CardItem.jsx'
import AddCardForm from './AddCardForm.jsx'

export default function TabColumn({ tab, allTabs, onMutate }) {
  async function handleDeleteTab() {
    if (!confirm(`Delete tab "${tab.name}" and all its cards?`)) return
    await deleteTab(tab.id)
    onMutate()
  }

  return (
    <div className="bg-gray-100 rounded-xl w-72 shrink-0 flex flex-col">
      <TabHeader tab={tab} onRename={onMutate} onDelete={handleDeleteTab} />
      <div className="flex flex-col gap-2 px-3 py-2">
        {tab.cards.map(card => (
          <CardItem key={card.id} card={card} allTabs={allTabs} onMutate={onMutate} />
        ))}
      </div>
      <AddCardForm tabId={tab.id} onMutate={onMutate} />
    </div>
  )
}
