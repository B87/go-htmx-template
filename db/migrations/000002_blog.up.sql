BEGIN;

CREATE TYPE status AS ENUM ('draft', 'published');

CREATE TABLE IF NOT EXISTS public.blog_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL CONSTRAINT title_unique UNIQUE,
    slug VARCHAR(255) NOT NULL CONSTRAINT slug_unique UNIQUE,
    description VARCHAR(400) NOT NULL,
    content TEXT NOT NULL,
    thumbnail_url VARCHAR(255) NOT NULL,
    status status,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
