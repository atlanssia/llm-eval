import { useQuery } from '@tanstack/react-query'
import { getEvaluations } from '../lib/api'
import { EvaluationCard } from '../components/EvaluationCard'

export function Dashboard() {
  const { data: evaluations, isLoading, error } = useQuery({
    queryKey: ['evaluations'],
    queryFn: getEvaluations,
    refetchInterval: 5000,
  })

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600">Failed to load evaluations</p>
      </div>
    )
  }

  const runningEvaluations = evaluations?.filter(e => e.status === 'running') || []
  const recentEvaluations = evaluations?.filter(e => e.status !== 'running') || []

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">Dashboard</h2>
      </div>

      {runningEvaluations.length > 0 && (
        <section className="mb-8">
          <h3 className="text-xl font-semibold mb-4">Active Evaluations</h3>
          <div className="grid gap-4 md:grid-cols-2">
            {runningEvaluations.map(evaluation => (
              <EvaluationCard key={evaluation.id} evaluation={evaluation} />
            ))}
          </div>
        </section>
      )}

      <section>
        <h3 className="text-xl font-semibold mb-4">Recent Evaluations</h3>
        {recentEvaluations.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            No evaluations yet. Create one to get started.
          </div>
        ) : (
          <div className="grid gap-4 md:grid-cols-2">
            {recentEvaluations.map(evaluation => (
              <EvaluationCard key={evaluation.id} evaluation={evaluation} />
            ))}
          </div>
        )}
      </section>
    </div>
  )
}
