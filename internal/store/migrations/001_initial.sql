CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS nodes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    os TEXT,
    arch TEXT,
    status TEXT NOT NULL DEFAULT 'unknown',
    last_seen_at TEXT,
    hardware_json TEXT,
    is_local INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS node_trust (
    node_id TEXT PRIMARY KEY REFERENCES nodes(id) ON DELETE CASCADE,
    fingerprint TEXT NOT NULL,
    paired_at TEXT NOT NULL DEFAULT (datetime('now')),
    revoked_at TEXT,
    cert_pem TEXT
);

CREATE TABLE IF NOT EXISTS runtimes (
    id TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    version TEXT,
    status TEXT NOT NULL DEFAULT 'unknown',
    install_path TEXT,
    meta_json TEXT,
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS models (
    id TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    family TEXT,
    variant TEXT,
    size_bytes INTEGER,
    memory_needed INTEGER,
    context_length INTEGER,
    capabilities_json TEXT,
    source_json TEXT,
    catalog_json TEXT
);

CREATE TABLE IF NOT EXISTS installed_models (
    model_id TEXT PRIMARY KEY REFERENCES models(id),
    path TEXT NOT NULL,
    sha256 TEXT,
    size_bytes INTEGER,
    status TEXT NOT NULL DEFAULT 'installed',
    installed_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS profiles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    purpose TEXT,
    orchestrator_id TEXT NOT NULL,
    node_policy_json TEXT,
    tools_json TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS profile_roles (
    id TEXT PRIMARY KEY,
    profile_id TEXT NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    role TEXT NOT NULL,
    model_id TEXT,
    node_id TEXT,
    required INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS conversations (
    id TEXT PRIMARY KEY,
    title TEXT,
    profile_id TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    meta_json TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    profile_id TEXT,
    conversation_id TEXT,
    prompt TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    result TEXT,
    error TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    completed_at TEXT
);

CREATE TABLE IF NOT EXISTS task_steps (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    step_index INTEGER NOT NULL,
    role TEXT,
    node_id TEXT,
    model_id TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    input TEXT,
    output TEXT,
    meta_json TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    completed_at TEXT
);

CREATE TABLE IF NOT EXISTS tool_permissions (
    id TEXT PRIMARY KEY,
    profile_id TEXT,
    tool_id TEXT NOT NULL,
    policy TEXT NOT NULL DEFAULT 'ask',
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS api_keys (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    key_prefix TEXT NOT NULL,
    key_hash TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    revoked_at TEXT,
    last_used_at TEXT
);

CREATE TABLE IF NOT EXISTS downloads (
    id TEXT PRIMARY KEY,
    model_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    bytes_downloaded INTEGER NOT NULL DEFAULT 0,
    bytes_total INTEGER NOT NULL DEFAULT 0,
    temp_path TEXT,
    error TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
