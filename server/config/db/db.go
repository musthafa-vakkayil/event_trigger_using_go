package config

const (
	CREATE_TABLE = `
			CREATE TABLE IF NOT EXISTS users(
  				id VARCHAR(50) PRIMARY KEY,
  				username VARCHAR(128) NOT NULL,
  				email VARCHAR(128) NOT NULL,
  				password_hash VARCHAR(128) NOT NULL,
  				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  			);

			CREATE TABLE IF NOT EXISTS triggers(
  				id VARCHAR(50) PRIMARY KEY,
  				name VARCHAR(128) NOT NULL,
  				type VARCHAR(20) CHECK (type IN ('SCHEDULED', 'API')),
  				schedule_time TIMESTAMP DEFAULT NULL,
  				interval_seconds INT DEFAULT NULL,
  				api_endpoint VARCHAR(255) DEFAULT NULL,
  				api_payload JSON DEFAULT NULL,
  				is_recurring BOOLEAN DEFAULT FALSE,
  				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  			);
  
  			CREATE TABLE IF NOT EXISTS events(
    			id VARCHAR(50) PRIMARY KEY,
    			trigger_id VARCHAR(50) NOT NULL,
    			event_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    			status VARCHAR(20) CHECK (status IN ('ACTIVE', 'ARCHIVE')),
    			archived_time TIMESTAMP DEFAULT NULL,
    			is_manual BOOLEAN DEFAULT FALSE,
    			FOREIGN KEY(trigger_id) REFERENCES triggers(id)
    	);
	`
)
