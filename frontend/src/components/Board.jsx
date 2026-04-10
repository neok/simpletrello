import { useState, useCallback } from 'react'
import { getTabs } from '../api/client.js'
import TabColumn from './TabColumn.jsx'
import AddTabForm from './AddTabForm.jsx'

function getInitialTabs() {
  if (typeof window !== 'undefined' && window.__INITIAL_DATA__) {
    return window.__INITIAL_DATA__
  }
  return []
}

export default function Board() {
  const [tabs, setTabs] = useState(getInitialTabs)
  const [error, setError] = useState(null)

  const reload = useCallback(async () => {
    try {
      const data = await getTabs()
      setTabs(data)
      setError(null)
    } catch (e) {
      setError(e.message)
    }
  }, [])

  if (error) {
    return (
      <div className="flex items-center justify-center h-screen bg-blue-600">
        <span className="text-red-200 text-xl">{error}</span>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-blue-600 p-6">
      <h1 className="text-white text-2xl font-bold mb-6">SimpleTrello</h1>
      <div className="flex gap-4 items-start overflow-x-auto pb-4">
        {tabs.map(tab => (
          <TabColumn key={tab.id} tab={tab} allTabs={tabs} onMutate={reload} />
        ))}
        <AddTabForm onMutate={reload} />
      </div>
    </div>
  )
}
