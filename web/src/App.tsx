import { useState } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Dashboard } from './pages/Dashboard'
import { NewEvaluation } from './pages/NewEvaluation'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5000,
      retry: 1,
    },
  },
})

type Page = 'dashboard' | 'new'

function App() {
  const [currentPage, setCurrentPage] = useState<Page>('dashboard')

  return (
    <QueryClientProvider client={queryClient}>
      <div className="min-h-screen bg-gray-50">
        <header className="bg-white shadow">
          <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
            <h1 className="text-2xl font-bold text-gray-900">LLM Evaluation Tool</h1>
            <nav className="flex gap-4">
              <button
                onClick={() => setCurrentPage('dashboard')}
                className={`px-4 py-2 rounded-md transition-colors ${
                  currentPage === 'dashboard'
                    ? 'bg-blue-600 text-white'
                    : 'text-gray-600 hover:bg-gray-100'
                }`}
              >
                Dashboard
              </button>
              <button
                onClick={() => setCurrentPage('new')}
                className={`px-4 py-2 rounded-md transition-colors ${
                  currentPage === 'new'
                    ? 'bg-blue-600 text-white'
                    : 'text-gray-600 hover:bg-gray-100'
                }`}
              >
                New Evaluation
              </button>
            </nav>
          </div>
        </header>
        <main className="max-w-7xl mx-auto px-4 py-8">
          {currentPage === 'dashboard' && <Dashboard />}
          {currentPage === 'new' && <NewEvaluation />}
        </main>
      </div>
    </QueryClientProvider>
  )
}

export default App
