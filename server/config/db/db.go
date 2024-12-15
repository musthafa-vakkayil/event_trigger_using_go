package config

const (
	CREATE_TABLE = `
		CREATE TABLE IF NOT EXISTS users (
			user_id VARCHAR(50) PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(100) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS chat_room(
			room_id VARCHAR(50) PRIMARY KEY,
			room_name VARCHAR(50) UNIQUE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS chat_messages(
			message_id VARCHAR(50) PRIMARY KEY,
			sender_id VARCHAR(50) NOT NULL,
			message_text TEXT NOT NULL,
			timestamp VARCHAR(50) NOT NULL,
			chat_room_id VARCHAR(50) NOT NULL
		);
	`

	PG_CRON = `CREATE EXTENSION IF NOT EXISTS pg_cron;`

	STATS = `CREATE EXTENSION IF NOT EXISTS pg_stats_statements CASCADE;`

	TRIGEGR = `
		CREATE OR REPLACE FUNCTION notify_trigger() RETURNS TRIGGER AS $trigger$
DECLARE
    rec RECORD;
    dat RECORD;
    payload TEXT;
BEGIN

    -- Set record row depending on operation
    CASE TG_OP
        WHEN 'UPDATE' THEN
            rec := NEW;
            dat := OLD;
        WHEN 'INSERT' THEN
            rec := NEW;
        WHEN 'DELETE' THEN
            rec := OLD;
        ELSE 
            RAISE EXCEPTION 'Unknown TG_OP: "%". Should not occur!', TG_OP;
    END CASE;

    -- Build the payload
    payload := json_build_object(
        'timestamp', CURRENT_TIMESTAMP,
        'action', LOWER(TG_OP),
        'schema', TG_TABLE_SCHEMA,
        'table', TG_TABLE_NAME,
        'data', row_to_json(rec)
    );

    -- Notify the channel
    PERFORM pg_notify('chat_events', payload);

    RETURN rec;
	END;
	$trigger$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS event_notify ON public.chat_messages;

	CREATE TRIGGER event_notify
	AFTER INSERT ON public.chat_messages
	FOR EACH ROW EXECUTE PROCEDURE notify_trigger();
	`
)
