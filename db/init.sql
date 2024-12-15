SELECT 'CREATE DATABASE eventdb'
WHERE NOT EXISTS (
    SELECT FROM pg_database WHERE datname = 'eventdb'
)\gexec