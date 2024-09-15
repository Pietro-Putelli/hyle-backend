-- create users table
CREATE TABLE books (
    id SERIAL PRIMARY KEY NOT NULL,
    guid UUID DEFAULT uuid_generate_v4 () NOT NULL UNIQUE,
    
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX books_user_idx ON books (user_id);

CREATE TABLE topics (
    id SERIAL PRIMARY KEY NOT NULL,

    user_id BIGINT NOT NULL,
    topic VARCHAR(255) NOT NULL,
    color VARCHAR(255) NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (user_id, topic)
);

CREATE INDEX topics_topic_color_idx ON topics (topic, color);

CREATE TABLE book_topics (
    id SERIAL PRIMARY KEY NOT NULL,

    book_id BIGINT NOT NULL,
    topic_id BIGINT NOT NULL,

    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE,
    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE
);

CREATE INDEX book_topics_book_idx ON book_topics (book_id);

CREATE TABLE book_picks (
    id SERIAL PRIMARY KEY NOT NULL,
    guid UUID DEFAULT uuid_generate_v4 () NOT NULL UNIQUE,

    user_id BIGINT NOT NULL,
    book_id BIGINT NOT NULL,
    title VARCHAR(255) NULL,
    content JSONB NOT NULL,
    content_text TEXT NOT NULL,
    "index" INT NOT NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX picks_book_idx ON book_picks (book_id);
CREATE INDEX picks_book_trgm_idx ON book_picks USING GIN (content_text gin_trgm_ops);

CREATE TABLE topic_colors (
    id SERIAL PRIMARY KEY NOT NULL,
    color VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE pick_search_keywords (
    id SERIAL PRIMARY KEY NOT NULL,

    user_id BIGINT NOT NULL,
    pick_id BIGINT NOT NULL,
    keyword VARCHAR(255) NOT NULL,

    FOREIGN KEY (pick_id) REFERENCES book_picks (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX pick_search_keywords_user_idx ON pick_search_keywords (user_id);
CREATE INDEX pick_search_keywords_trgm_idx ON pick_search_keywords USING GIN (keyword gin_trgm_ops);