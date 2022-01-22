create table if not exists register_events
(
    id      serial PRIMARY KEY,
    name    varchar(255) NOT NULL,
    created timestamp DEFAULT NOW(),
    skip    BOOLEAN DEFAULT FALSE    
)