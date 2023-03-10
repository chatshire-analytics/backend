CREATE TABLE flipside_query_result (
                                       id SERIAL PRIMARY KEY,
                                       results JSONB NULL, -- formerly JSONB or TEXT which processes the result afterward by having batch/processing server
                                       query VARCHAR(255) NOT NULL,
                                       sentence VARCHAR(255) NOT NULL,
                                       gpt_id VARCHAR(255) NULL,
                                       token VARCHAR(255) NULL,
                                       address VARCHAR(255) NOT NULL,
                                       started_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                       ended_at TIMESTAMP NULL,
                                       status VARCHAR(255) NOT NULL DEFAULT 'QUERY_CREATED'
);