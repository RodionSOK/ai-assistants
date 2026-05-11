CREATE TABLE assistants (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id         UUID NOT NULL REFERENCES categories(id),
    name                VARCHAR(255) NOT NULL,
    description         TEXT NOT NULL,
    model               VARCHAR(255) NOT NULL,
    system_prompt       TEXT NOT NULL,
    example_user_prompt TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_assistants_category_id ON assistants(category_id);

CREATE INDEX idx_assistants_name ON assistants(name);
CREATE INDEX idx_assistants_description ON assistants(description);

CREATE INDEX idx_assistants_is_active ON assistants(is_active);