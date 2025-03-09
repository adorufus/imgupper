CREATE TABLE
    IF NOT EXISTS files (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        filename VARCHAR(255) NOT NULL,
        filesize BIGINT NOT NULL,
        mime_type VARCHAR(100) NOT NULL,
        bucket_url TEXT NOT NULL,
        created_at TIMESTAMP
        WITH
            TIME ZONE NOT NULL,
            updated_at TIMESTAMP
        WITH
            TIME ZONE NOT NULL,
            CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
    );

-- Create an index on the user_id column to improve query performance
CREATE INDEX IF NOT EXISTS idx_files_user_id ON files (user_id);