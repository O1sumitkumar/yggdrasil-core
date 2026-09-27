-- Multi-node Team performance: per-role steps + cross-machine flag.
ALTER TABLE generation_metrics ADD COLUMN role_steps_json TEXT NOT NULL DEFAULT '[]';
ALTER TABLE generation_metrics ADD COLUMN cross_machine INTEGER NOT NULL DEFAULT 0;
ALTER TABLE generation_metrics ADD COLUMN node_count INTEGER NOT NULL DEFAULT 0;
