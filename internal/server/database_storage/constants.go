package dbstorage

var maximumNumberOfRetries = []int{1, 3, 6}

const schema = `
	CREATE TABLE IF NOT EXISTS metrics (
	    "_id" SERIAL,
		"name" TEXT NOT NULL,
		"mtype" VARCHAR(12) NOT NULL DEFAULT 'gauge',
		"delta" BIGINT NOT NULL DEFAULT 0,
		"value" DOUBLE PRECISION NOT NULL DEFAULT 0.0,
		CONSTRAINT unique_id_mtype UNIQUE (name, mtype),
		PRIMARY KEY (_id)
	)
`

const getGaugeMetricSQLRequest = `SELECT value FROM metrics WHERE name = :name AND mtype = 'gauge'`
const getCounterMetricSQLRequest = `SELECT delta FROM metrics WHERE name = :name AND mtype = 'counter'`

const setOrUpdateMetricSQLRequest = `
			INSERT INTO metrics (name, mtype, delta, value) 
				VALUES (:name, :mtype, :delta, :value)
			ON CONFLICT (name, mtype) DO 
			    UPDATE SET delta = metrics.delta + excluded.delta, value = excluded.value
		`
