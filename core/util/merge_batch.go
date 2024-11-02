package util

import "github.com/jackc/pgx/v5"

func MergeBatch(from pgx.Batch, to *pgx.Batch) {
	for _, q := range from.QueuedQueries {
		to.Queue(q.SQL, q.Arguments...)
	}
}
