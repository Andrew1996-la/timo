CREATE TABLE IF NOT EXISTS tasks  (
   id SERIAL PRIMARY KEY,
   title TEXT NOT NULL,
   created_at TIMESTAMP DEFAULT now(),
   deleted_at TIMESTAMP
);