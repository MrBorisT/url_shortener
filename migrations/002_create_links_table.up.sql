CREATE TABLE links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    original_url TEXT NOT NULL,
    short_code TEXT NOT NULL UNIQUE,

    click_count BIGINT NOT NULL DEFAULT 0,

    disabled_at TIMESTAMP NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_links_user_id ON links(user_id);