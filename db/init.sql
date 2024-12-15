SELECT 'CREATE DATABASE eventdb'
WHERE NOT EXISTS (
    SELECT FROM pg_database WHERE datname = 'eventdb'
)\gexec

CREATE OR REPLACE FUNCTION notify_trigger() 
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('trigger_channel', NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;