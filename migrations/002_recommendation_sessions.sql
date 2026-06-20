-- Adds the recommendation session table for the Gemini Q&A flow.
-- Run: mysql -u root -p letter_square < migrations/002_recommendation_sessions.sql

USE letter_square;

CREATE TABLE IF NOT EXISTS recommendation_sessions (
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id    BIGINT       NOT NULL,
    status     VARCHAR(20)  NOT NULL DEFAULT 'in_progress', -- in_progress | completed
    step       INT          NOT NULL DEFAULT 0,
    answers    TEXT         NOT NULL DEFAULT '[]',          -- JSON array of answer strings
    result     TEXT         NOT NULL DEFAULT '',             -- final Gemini recommendation text
    created_at DATETIME     DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_rec_session_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
