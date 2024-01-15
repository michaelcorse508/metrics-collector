package postgresdb

// CreateTableQuery is the template for creating table metrics in Postgres DB
const CreateTableQuery = `create table metrics
(
    id      varchar(30) not null,
    mtype varchar(10) not null,
    value double precision,
    delta int,
    PRIMARY KEY (id, mtype)
);

`

// InsertUpdateRowQuery is the template to insert line in the table if such row
// doesn't exist and to update if it is already exist.
const InsertUpdateRowQuery = `INSERT INTO metrics (id, mtype, value, delta)
    VALUES ($1, $2, $3, $4)
        ON CONFLICT (id, mtype) DO UPDATE 
            SET value = excluded.value, delta = CASE
                WHEN excluded.delta IS NOT NULL OR (SELECT delta FROM metrics m WHERE m.id = excluded.id) IS NOT NULL
                    THEN COALESCE(excluded.delta, 0) + COALESCE((SELECT delta FROM metrics m WHERE m.id = excluded.id), 0)
                                   ELSE (SELECT delta FROM metrics m WHERE m.id = excluded.id)
            END;`

// SelectRowQuery is the template to select one row from the table metric with given ID and MType.
const SelectRowQuery = `SELECT id, mtype, delta, value FROM metrics m
                               WHERE m.id = $1 AND m.mtype = $2;`

// SelectAllRowsQuery is the template to select all rows from table metrics.
const SelectAllRowsQuery = `SELECT * FROM metrics;`
