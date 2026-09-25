package repository

import "gorm.io/gorm/clause"

// lockedClause 返回 SELECT ... FOR UPDATE 子句，用于并发保护。
func lockedClause() clause.Locking {
	return clause.Locking{Strength: "UPDATE"}
}
