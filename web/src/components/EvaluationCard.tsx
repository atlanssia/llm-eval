import type { Evaluation } from '../lib/types'

interface EvaluationCardProps {
  evaluation: Evaluation
}

export function EvaluationCard({ evaluation }: EvaluationCardProps) {
  const progress = (evaluation.completed_cases / evaluation.total_cases) * 100

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <div className="flex justify-between items-center mb-2">
        <h3 className="text-lg font-semibold">{evaluation.id}</h3>
        <span className={`px-2 py-1 rounded text-sm ${
          evaluation.status === 'running' ? 'bg-blue-100 text-blue-800' :
          evaluation.status === 'completed' ? 'bg-green-100 text-green-800' :
          evaluation.status === 'failed' ? 'bg-red-100 text-red-800' :
          'bg-gray-100 text-gray-800'
        }`}>
          {evaluation.status.charAt(0).toUpperCase() + evaluation.status.slice(1)}
        </span>
      </div>

      {evaluation.status === 'running' && (
        <div className="mb-2">
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div
              className="bg-blue-600 h-2 rounded-full transition-all"
              style={{ width: `${progress}%` }}
            />
          </div>
          <p className="text-sm text-gray-600 mt-1">{progress.toFixed(0)}%</p>
        </div>
      )}

      <div className="text-sm text-gray-600">
        <p>Models: {evaluation.config.models.join(', ')}</p>
        <p>Datasets: {evaluation.config.datasets.join(', ')}</p>
      </div>
    </div>
  )
}
