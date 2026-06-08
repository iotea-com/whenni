package sqlc

import (
	sqldb "github.com/iotea-com/iotea/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connection pool to the database
var Pool *pgxpool.Pool

// sqlc query client
var Queries *sqldb.Queries
