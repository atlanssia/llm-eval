-- Evaluations table
CREATE TABLE IF NOT EXISTS evaluations (
    id TEXT PRIMARY KEY,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    status TEXT NOT NULL,
    config TEXT NOT NULL, -- JSON
    total_cases INTEGER NOT NULL DEFAULT 0,
    completed_cases INTEGER NOT NULL DEFAULT 0,
    error TEXT
);

-- Model results table
CREATE TABLE IF NOT EXISTS model_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    evaluation_id TEXT NOT NULL,
    model_name TEXT NOT NULL,
    predictions TEXT NOT NULL, -- JSON array
    "references" TEXT NOT NULL, -- JSON array
    latencies TEXT NOT NULL, -- JSON array
    tokens_per_sec TEXT, -- JSON array
    metrics TEXT NOT NULL, -- JSON object
    error_count INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (evaluation_id) REFERENCES evaluations(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_model_results_evaluation_id ON model_results(evaluation_id);
CREATE INDEX IF NOT EXISTS idx_evaluations_status ON evaluations(status);
CREATE INDEX IF NOT EXISTS idx_evaluations_created_at ON evaluations(created_at DESC);
