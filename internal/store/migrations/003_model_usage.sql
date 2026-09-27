-- 003_model_usage.sql
CREATE TABLE IF NOT EXISTS model_usage (
  model_id TEXT PRIMARY KEY,
  last_used_at TEXT NOT NULL,
  FOREIGN KEY (model_id) REFERENCES models(id) ON DELETE CASCADE
);
