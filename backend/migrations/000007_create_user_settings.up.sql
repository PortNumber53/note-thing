CREATE TABLE user_settings (
    user_id             UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    default_notebook_id UUID REFERENCES notebooks(id) ON DELETE SET NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
