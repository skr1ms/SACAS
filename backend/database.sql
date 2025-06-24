CREATE TABLE IF NOT EXISTS code_submissions (
    id VARCHAR(64) PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(16) NOT NULL,
    content TEXT NOT NULL,
    grade INTEGER,
    feedback TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_code_submissions_created_at ON code_submissions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_code_submissions_file_type ON code_submissions(file_type);
CREATE INDEX IF NOT EXISTS idx_code_submissions_grade ON code_submissions(grade);
