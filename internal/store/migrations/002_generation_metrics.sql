CREATE TABLE IF NOT EXISTS generation_metrics (
    id TEXT PRIMARY KEY,
    conversation_id TEXT,
    conversation_title TEXT,
    message_id TEXT,
    profile_id TEXT,
    profile_name TEXT,
    model_id TEXT NOT NULL DEFAULT '',
    runtime_id TEXT NOT NULL DEFAULT 'llamacpp',
    prompt_tokens INTEGER NOT NULL DEFAULT 0,
    completion_tokens INTEGER NOT NULL DEFAULT 0,
    total_tokens INTEGER NOT NULL DEFAULT 0,
    ttft_ms REAL NOT NULL DEFAULT 0,
    prompt_ms REAL NOT NULL DEFAULT 0,
    eval_ms REAL NOT NULL DEFAULT 0,
    total_ms REAL NOT NULL DEFAULT 0,
    prompt_tok_per_sec REAL NOT NULL DEFAULT 0,
    eval_tok_per_sec REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_generation_metrics_created ON generation_metrics(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_generation_metrics_model ON generation_metrics(model_id);
CREATE INDEX IF NOT EXISTS idx_generation_metrics_conversation ON generation_metrics(conversation_id);
