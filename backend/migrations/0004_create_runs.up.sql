CREATE TABLE runs (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assistant_id   UUID NOT NULL REFERENCES assistants(id),
    user_id        UUID NOT NULL REFERENCES users(id),
    model          VARCHAR(255) NOT NULL,
    user_prompt    TEXT NOT NULL,
    output         TEXT,
    status         VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'success', 'failed')),
    error          TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_runs_user_id ON runs(user_id);

CREATE INDEX idx_runs_assistant_id ON runs(assistant_id);

CREATE INDEX idx_runs_status ON runs(status);

CREATE INDEX idx_runs_created_at ON runs(created_at DESC);