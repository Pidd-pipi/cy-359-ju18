package repository

import "gorm.io/gorm"

// Repositories 聚合所有仓储。
type Repositories struct {
	User         *UserRepository
	Activity     *ActivityRepository
	Checkpoint   *CheckpointRepository
	Team         *TeamRepository
	Registration *RegistrationRepository
	Checkin      *CheckinRepository
	Product      *ProductRepository
	Redemption   *RedemptionRepository
	Favorite     *FavoriteRepository
	Audit        *AuditRepository
}

// New 构造全部仓储。
func New(db *gorm.DB) *Repositories {
	return &Repositories{
		User:         NewUserRepository(db),
		Activity:     NewActivityRepository(db),
		Checkpoint:   NewCheckpointRepository(db),
		Team:         NewTeamRepository(db),
		Registration: NewRegistrationRepository(db),
		Checkin:      NewCheckinRepository(db),
		Product:      NewProductRepository(db),
		Redemption:   NewRedemptionRepository(db),
		Favorite:     NewFavoriteRepository(db),
		Audit:        NewAuditRepository(db),
	}
}
