CREATE TABLE IF NOT EXISTS cpu_usage (
    ts    TIMESTAMPTZ,
    host  TEXT,
    usage DOUBLE PRECISION
);
