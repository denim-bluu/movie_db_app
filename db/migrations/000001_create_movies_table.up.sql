CREATE TABLE IF NOT EXISTS movies (
    id BIGSERIAL PRIMARY KEY,  
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    title TEXT NOT NULL,
    year INTEGER NOT NULL,
    runtime INTEGER NOT NULL,
    genres TEXT[] NOT NULL,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS movies_title_idx ON movies USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS movies_genres_idx ON movies USING GIN (genres);

-- Insert statements for each movie
INSERT INTO movies (title, year, runtime, genres)
VALUES
    ('Black Panther', 2018, 134, '{"action","adventure"}'),
    ('Deadpool', 2016, 108, '{"action","comedy"}'),
    ('The Breakfast Club', 1986, 96, '{"drama"}');