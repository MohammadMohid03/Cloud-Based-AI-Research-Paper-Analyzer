-- =============================================================================
-- AI-Powered Research Paper Analyzer — Database Schema
-- =============================================================================
-- PostgreSQL 15+
-- This file can be used as an alternative to GORM auto-migration.
-- Usage:
--   psql -U postgres -d research_paper_analyzer -f database/schema.sql
-- =============================================================================

-- ── Extensions ──────────────────────────────────────────────────────────────

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";   -- trigram index for full-text search

-- ── Users ───────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    first_name      VARCHAR(100) NOT NULL DEFAULT '',
    last_name       VARCHAR(100) NOT NULL DEFAULT '',
    avatar_url      TEXT         NOT NULL DEFAULT '',
    role            VARCHAR(20)  NOT NULL DEFAULT 'user'
                        CHECK (role IN ('user', 'admin')),
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users (email);

-- ── Papers ──────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS papers (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           VARCHAR(500) NOT NULL DEFAULT 'Untitled Paper',
    authors         TEXT         NOT NULL DEFAULT '',        -- comma-separated
    abstract        TEXT         NOT NULL DEFAULT '',
    journal         VARCHAR(300) NOT NULL DEFAULT '',
    publication_year INTEGER,
    doi             VARCHAR(255),
    file_path       TEXT         NOT NULL DEFAULT '',        -- local path or S3 key
    file_name       VARCHAR(500) NOT NULL DEFAULT '',
    file_size       BIGINT       NOT NULL DEFAULT 0,         -- bytes
    file_type       VARCHAR(50)  NOT NULL DEFAULT 'application/pdf',
    status          VARCHAR(30)  NOT NULL DEFAULT 'uploaded'
                        CHECK (status IN ('uploaded', 'processing', 'analyzed', 'failed')),
    analysis_error  TEXT         NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_papers_user_id   ON papers (user_id);
CREATE INDEX idx_papers_status    ON papers (status);
CREATE INDEX idx_papers_title_trgm ON papers USING gin (title gin_trgm_ops);

-- ── Analyses ────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS analyses (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    paper_id        UUID         NOT NULL REFERENCES papers(id) ON DELETE CASCADE,
    summary         TEXT         NOT NULL DEFAULT '',
    key_findings    TEXT         NOT NULL DEFAULT '',        -- JSON array stored as text
    methodology     TEXT         NOT NULL DEFAULT '',
    limitations     TEXT         NOT NULL DEFAULT '',
    future_work     TEXT         NOT NULL DEFAULT '',
    keywords        TEXT         NOT NULL DEFAULT '',        -- JSON array stored as text
    confidence_score NUMERIC(4,2) DEFAULT 0.00
                        CHECK (confidence_score >= 0 AND confidence_score <= 1),
    model_id        VARCHAR(200) NOT NULL DEFAULT '',        -- AI model used
    processing_time_ms INTEGER   NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_analyses_paper_id ON analyses (paper_id);

-- ── Chat Sessions ───────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS chat_sessions (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    paper_id        UUID         NOT NULL REFERENCES papers(id) ON DELETE CASCADE,
    user_id         UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           VARCHAR(300) NOT NULL DEFAULT 'New Chat',
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_chat_sessions_paper_id ON chat_sessions (paper_id);
CREATE INDEX idx_chat_sessions_user_id  ON chat_sessions (user_id);

-- ── Chat Messages ───────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS chat_messages (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id      UUID         NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    role            VARCHAR(20)  NOT NULL
                        CHECK (role IN ('user', 'assistant', 'system')),
    content         TEXT         NOT NULL DEFAULT '',
    tokens_used     INTEGER      NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_chat_messages_session_id  ON chat_messages (session_id);
CREATE INDEX idx_chat_messages_created_at  ON chat_messages (created_at);

-- ── Bookmarks / Favorites ───────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS bookmarks (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    paper_id        UUID         NOT NULL REFERENCES papers(id) ON DELETE CASCADE,
    note            TEXT         NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, paper_id)
);

CREATE INDEX idx_bookmarks_user_id  ON bookmarks (user_id);
CREATE INDEX idx_bookmarks_paper_id ON bookmarks (paper_id);

-- ── Activity Log (Audit Trail) ──────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS activity_logs (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID         REFERENCES users(id) ON DELETE SET NULL,
    action          VARCHAR(100) NOT NULL,
    entity_type     VARCHAR(50)  NOT NULL DEFAULT '',
    entity_id       UUID,
    metadata        JSONB        DEFAULT '{}',
    ip_address      VARCHAR(45)  NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_user_id    ON activity_logs (user_id);
CREATE INDEX idx_activity_logs_action     ON activity_logs (action);
CREATE INDEX idx_activity_logs_created_at ON activity_logs (created_at);

-- ── Updated-At Trigger ──────────────────────────────────────────────────────

CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to all tables with an updated_at column
DO $$
DECLARE
    tbl TEXT;
BEGIN
    FOR tbl IN
        SELECT table_name
        FROM information_schema.columns
        WHERE column_name = 'updated_at'
          AND table_schema = 'public'
    LOOP
        EXECUTE format(
            'CREATE TRIGGER set_updated_at BEFORE UPDATE ON %I
             FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();',
            tbl
        );
    END LOOP;
END;
$$;

-- =============================================================================
-- SEED DATA — Sample records for development / demo
-- =============================================================================

-- Default admin user (password: "admin123" — bcrypt hash)
INSERT INTO users (id, email, password_hash, first_name, last_name, role)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'admin@example.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'Admin',
    'User',
    'admin'
) ON CONFLICT (email) DO NOTHING;

-- Default demo user (password: "demo1234" — bcrypt hash)
INSERT INTO users (id, email, password_hash, first_name, last_name, role)
VALUES (
    'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'demo@example.com',
    '$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36PQm2Pro7DNzBEOxHXJQe6',
    'Demo',
    'User',
    'user'
) ON CONFLICT (email) DO NOTHING;

-- Sample paper for the demo user
INSERT INTO papers (id, user_id, title, authors, abstract, journal, publication_year, status)
VALUES (
    'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
    'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'Attention Is All You Need',
    'Ashish Vaswani, Noam Shazeer, Niki Parmar, Jakob Uszkoreit, Llion Jones, Aidan N. Gomez, Lukasz Kaiser, Illia Polosukhin',
    'The dominant sequence transduction models are based on complex recurrent or convolutional neural networks that include an encoder and a decoder. The best performing models also connect the encoder and decoder through an attention mechanism. We propose a new simple network architecture, the Transformer, based solely on attention mechanisms, dispensing with recurrence and convolutions entirely.',
    'Advances in Neural Information Processing Systems',
    2017,
    'analyzed'
) ON CONFLICT DO NOTHING;

-- Sample analysis for the demo paper
INSERT INTO analyses (id, paper_id, summary, key_findings, methodology, keywords, confidence_score, model_id)
VALUES (
    'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44',
    'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
    'This paper introduces the Transformer architecture, which relies entirely on self-attention mechanisms for sequence transduction tasks. The model achieves state-of-the-art results on machine translation benchmarks while being more parallelizable and requiring significantly less time to train.',
    '["Transformer architecture eliminates recurrence and convolutions","Self-attention mechanism captures global dependencies","Multi-head attention allows attending to different representation subspaces","Achieves 28.4 BLEU on WMT 2014 English-to-German translation"]',
    'The authors propose a novel neural network architecture based on self-attention. They evaluate the model on WMT 2014 English-to-German and English-to-French translation tasks, comparing against existing RNN and CNN-based approaches.',
    '["transformer","self-attention","neural machine translation","sequence-to-sequence","multi-head attention"]',
    0.92,
    'mock-analyzer'
) ON CONFLICT DO NOTHING;

-- =============================================================================
-- END OF SCHEMA
-- =============================================================================
