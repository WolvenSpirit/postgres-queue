CREATE TABLE IF NOT EXISTS notify_events 
(
    id          serial PRIMARY KEY,
    publisher   VARCHAR(255)    NOT NULL,
    channel     VARCHAR(255)    NOT NULL,
    consumer    VARCHAR(255)    DEFAULT NULL,
    ack         BOOLEAN         DEFAULT NULL,
    data        TEXT            NOT NULL,
    created     TIMESTAMP       DEFAULT NOW(),
    duration    NUMERIC(9,2)    DEFAULT NULL,
    completed   TIMESTAMP       DEFAULT NULL    
);