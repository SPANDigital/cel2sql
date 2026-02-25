package testcases

import "github.com/spandigital/cel2sql/v3/dialect"

// TimestampTests returns test cases for timestamp and duration operations.
func TimestampTests() []ConvertTestCase {
	return []ConvertTestCase{
		{
			Name:     "duration_second",
			CELExpr:  `duration("10s")`,
			Category: CategoryTimestamp,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "INTERVAL 10 SECOND",
				dialect.MySQL:      "INTERVAL 10 SECOND",
				dialect.SQLite:     "'+10 seconds'",
				dialect.DuckDB:     "INTERVAL 10 SECOND",
				dialect.BigQuery:   "INTERVAL 10 SECOND",
			},
		},
		{
			Name:     "duration_minute",
			CELExpr:  `duration("1h1m")`,
			Category: CategoryTimestamp,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "INTERVAL 61 MINUTE",
				dialect.MySQL:      "INTERVAL 61 MINUTE",
				dialect.SQLite:     "'+61 minutes'",
				dialect.DuckDB:     "INTERVAL 61 MINUTE",
				dialect.BigQuery:   "INTERVAL 61 MINUTE",
			},
		},
		{
			Name:     "duration_hour",
			CELExpr:  `duration("60m")`,
			Category: CategoryTimestamp,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "INTERVAL 1 HOUR",
				dialect.MySQL:      "INTERVAL 1 HOUR",
				dialect.SQLite:     "'+1 hours'",
				dialect.DuckDB:     "INTERVAL 1 HOUR",
				dialect.BigQuery:   "INTERVAL 1 HOUR",
			},
		},
		{
			Name:     "timestamp_getSeconds",
			CELExpr:  `created_at.getSeconds()`,
			Category: CategoryTimestamp,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "EXTRACT(SECOND FROM created_at)",
				dialect.MySQL:      "EXTRACT(SECOND FROM created_at)",
				dialect.SQLite:     "CAST(strftime('%S', created_at) AS INTEGER)",
				dialect.DuckDB:     "EXTRACT(SECOND FROM created_at)",
				dialect.BigQuery:   "EXTRACT(SECOND FROM created_at)",
			},
		},
		{
			Name:     "timestamp_getHours_withTimezone",
			CELExpr:  `created_at.getHours("Asia/Tokyo")`,
			Category: CategoryTimestamp,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "EXTRACT(HOUR FROM created_at AT TIME ZONE 'Asia/Tokyo')",
				dialect.MySQL:      "EXTRACT(HOUR FROM created_at AT TIME ZONE 'Asia/Tokyo')",
				dialect.SQLite:     "CAST(strftime('%H', created_at) AS INTEGER)",
				dialect.DuckDB:     "EXTRACT(HOUR FROM created_at AT TIME ZONE 'Asia/Tokyo')",
				dialect.BigQuery:   "EXTRACT(HOUR FROM created_at AT TIME ZONE 'Asia/Tokyo')",
			},
		},
		{
			Name:     "timestamp_sub_duration",
			CELExpr:  `created_at - duration("60m") <= timestamp("2021-09-01T18:00:00Z")`,
			Category: CategoryTimestamp,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "created_at - INTERVAL 1 HOUR <= CAST('2021-09-01T18:00:00Z' AS TIMESTAMP WITH TIME ZONE)",
				dialect.MySQL:      "created_at - INTERVAL 1 HOUR <= CAST('2021-09-01T18:00:00Z' AS DATETIME)",
				dialect.SQLite:     "datetime(created_at, '-'||'+1 hours') <= datetime('2021-09-01T18:00:00Z')",
				dialect.DuckDB:     "created_at - INTERVAL 1 HOUR <= CAST('2021-09-01T18:00:00Z' AS TIMESTAMPTZ)",
				dialect.BigQuery:   "TIMESTAMP_SUB(created_at, INTERVAL 1 HOUR) <= CAST('2021-09-01T18:00:00Z' AS TIMESTAMP)",
			},
		},
		{
			Name:     "interval_month",
			CELExpr:  `interval(1, MONTH)`,
			Category: CategoryTimestamp,
			EnvSetup: EnvWithTimestamp,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "INTERVAL 1 MONTH",
				dialect.MySQL:      "INTERVAL 1 MONTH",
				dialect.SQLite:     "'+'||1||' months'",
				dialect.DuckDB:     "INTERVAL 1 MONTH",
				dialect.BigQuery:   "INTERVAL 1 MONTH",
			},
		},
		{
			Name:     "date_getFullYear",
			CELExpr:  `birthday.getFullYear()`,
			Category: CategoryTimestamp,
			EnvSetup: EnvWithTimestamp,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "EXTRACT(YEAR FROM birthday)",
				dialect.MySQL:      "EXTRACT(YEAR FROM birthday)",
				dialect.SQLite:     "CAST(strftime('%Y', birthday) AS INTEGER)",
				dialect.DuckDB:     "EXTRACT(YEAR FROM birthday)",
				dialect.BigQuery:   "EXTRACT(YEAR FROM birthday)",
			},
		},
		{
			Name:     "datetime_getMonth",
			CELExpr:  `scheduled_at.getMonth()`,
			Category: CategoryTimestamp,
			EnvSetup: EnvWithTimestamp,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "EXTRACT(MONTH FROM scheduled_at) - 1",
				dialect.MySQL:      "EXTRACT(MONTH FROM scheduled_at) - 1",
				dialect.SQLite:     "CAST(strftime('%m', scheduled_at) AS INTEGER) - 1",
				dialect.DuckDB:     "EXTRACT(MONTH FROM scheduled_at) - 1",
				dialect.BigQuery:   "EXTRACT(MONTH FROM scheduled_at) - 1",
			},
		},
	}
}
