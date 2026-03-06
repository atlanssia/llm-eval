import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { createEvaluation } from '../lib/api'

export function NewEvaluation() {
  const [selectedModels, setSelectedModels] = useState<string[]>([])
  const [selectedDatasets, setSelectedDatasets] = useState<string[]>([])
  const [sampleSize, setSampleSize] = useState(100)

  const createMutation = useMutation({
    mutationFn: createEvaluation,
    onSuccess: (data) => {
      alert(`Evaluation created! ID: ${data.id}`)
      // Reset form
      setSelectedModels([])
      setSelectedDatasets([])
      setSampleSize(100)
    },
  })

  // Mock data for demo
  const availableModels = ['gpt-4', 'claude-3', 'llama-2']
  const availableDatasets = ['mmlu_anatomy', 'mmlu_history', 'cmmlu']

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (selectedModels.length === 0 || selectedDatasets.length === 0) {
      alert('Please select at least one model and one dataset')
      return
    }

    createMutation.mutate({
      models: selectedModels,
      datasets: selectedDatasets,
      config: {
        sample_size: sampleSize,
        max_workers: 4,
        ephemeral: false,
      },
    })
  }

  return (
    <div className="max-w-2xl mx-auto">
      <h2 className="text-2xl font-bold mb-6">Create New Evaluation</h2>

      <form onSubmit={handleSubmit} className="space-y-6 bg-white rounded-lg shadow p-6">
        {/* Models Selection */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Select Models
          </label>
          <div className="space-y-2">
            {availableModels.map(model => (
              <label key={model} className="flex items-center">
                <input
                  type="checkbox"
                  checked={selectedModels.includes(model)}
                  onChange={(e) => {
                    if (e.target.checked) {
                      setSelectedModels([...selectedModels, model])
                    } else {
                      setSelectedModels(selectedModels.filter(m => m !== model))
                    }
                  }}
                  className="mr-2"
                />
                <span className="text-gray-900">{model}</span>
              </label>
            ))}
          </div>
          <p className="text-sm text-gray-500 mt-1">
            Selected: {selectedModels.length === 0 ? 'None' : selectedModels.join(', ')}
          </p>
        </div>

        {/* Datasets Selection */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Select Datasets
          </label>
          <div className="space-y-2">
            {availableDatasets.map(dataset => (
              <label key={dataset} className="flex items-center">
                <input
                  type="checkbox"
                  checked={selectedDatasets.includes(dataset)}
                  onChange={(e) => {
                    if (e.target.checked) {
                      setSelectedDatasets([...selectedDatasets, dataset])
                    } else {
                      setSelectedDatasets(selectedDatasets.filter(d => d !== dataset))
                    }
                  }}
                  className="mr-2"
                />
                <span className="text-gray-900">{dataset}</span>
              </label>
            ))}
          </div>
          <p className="text-sm text-gray-500 mt-1">
            Selected: {selectedDatasets.length === 0 ? 'None' : selectedDatasets.join(', ')}
          </p>
        </div>

        {/* Sample Size */}
        <div>
          <label htmlFor="sampleSize" className="block text-sm font-medium text-gray-700 mb-2">
            Sample Size
          </label>
          <input
            type="number"
            id="sampleSize"
            value={sampleSize}
            onChange={(e) => setSampleSize(parseInt(e.target.value))}
            min={1}
            max={10000}
            className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
          />
          <p className="text-sm text-gray-500 mt-1">
            Number of test cases to evaluate (1-10000)
          </p>
        </div>

        {/* Submit Button */}
        <div className="flex justify-end gap-3">
          <button
            type="button"
            onClick={() => {
              setSelectedModels([])
              setSelectedDatasets([])
              setSampleSize(100)
            }}
            className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
          >
            Reset
          </button>
          <button
            type="submit"
            disabled={createMutation.isPending}
            className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
          >
            {createMutation.isPending ? 'Creating...' : 'Create Evaluation'}
          </button>
        </div>
      </form>

      {/* Info Section */}
      <div className="mt-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
        <h3 className="text-sm font-semibold text-blue-900 mb-2">About Evaluations</h3>
        <p className="text-sm text-blue-800">
          This tool evaluates LLM models on various datasets. Select models and datasets above,
          then click "Create Evaluation" to start. The evaluation will run in the background
          and you can monitor progress in real-time.
        </p>
      </div>
    </div>
  )
}
