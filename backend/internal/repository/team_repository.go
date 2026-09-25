package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/orienteering/platform/internal/constants"
	"github.com/orienteering/platform/internal/model"
)

// TeamRepository 团队仓储。
type TeamRepository struct {
	db *gorm.DB
}

// NewTeamRepository 构造团队仓储。
func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(tx *gorm.DB, team *model.Team) error {
	if err := tx.Create(team).Error; err != nil {
		return fmt.Errorf("create team: %w", err)
	}
	return nil
}

func (r *TeamRepository) GetByID(id int64) (*model.Team, error) {
	var team model.Team
	if err := r.db.Preload("Members").First(&team, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get team by id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get team by id: %w", err)
	}
	return &team, nil
}

func (r *TeamRepository) GetByIDForUpdate(tx *gorm.DB, id int64) (*model.Team, error) {
	var team model.Team
	if err := tx.Clauses(lockedClause()).First(&team, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get team for update: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get team for update: %w", err)
	}
	return &team, nil
}

func (r *TeamRepository) ExistsByName(name string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Team{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, fmt.Errorf("exists team by name: %w", err)
	}
	return count > 0, nil
}

func (r *TeamRepository) ListByUser(userID int64) ([]model.Team, error) {
	var teams []model.Team
	if err := r.db.Preload("Members").
		Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ?", userID).
		Find(&teams).Error; err != nil {
		return nil, fmt.Errorf("list teams by user: %w", err)
	}
	return teams, nil
}

func (r *TeamRepository) ListAll(page, pageSize, offset int) ([]model.Team, int64, error) {
	var total int64
	if err := r.db.Model(&model.Team{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count teams: %w", err)
	}
	var teams []model.Team
	if err := r.db.Preload("Members").Order("id DESC").Limit(pageSize).Offset(offset).Find(&teams).Error; err != nil {
		return nil, 0, fmt.Errorf("list teams: %w", err)
	}
	return teams, total, nil
}

func (r *TeamRepository) AddMember(tx *gorm.DB, member *model.TeamMember) error {
	if err := tx.Create(member).Error; err != nil {
		return fmt.Errorf("add team member: %w", err)
	}
	return nil
}

func (r *TeamRepository) CountMembers(teamID int64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.TeamMember{}).Where("team_id = ?", teamID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count team members: %w", err)
	}
	return count, nil
}

func (r *TeamRepository) IsMember(teamID, userID int64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.TeamMember{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check team member: %w", err)
	}
	return count > 0, nil
}

func (r *TeamRepository) Leave(teamID, userID int64) error {
	if err := r.db.Where("team_id = ? AND user_id = ?", teamID, userID).
		Delete(&model.TeamMember{}).Error; err != nil {
		return fmt.Errorf("leave team: %w", err)
	}
	return nil
}

// RegistrationRepository 报名仓储。
type RegistrationRepository struct {
	db *gorm.DB
}

// NewRegistrationRepository 构造报名仓储。
func NewRegistrationRepository(db *gorm.DB) *RegistrationRepository {
	return &RegistrationRepository{db: db}
}

func (r *RegistrationRepository) Create(tx *gorm.DB, reg *model.Registration) error {
	if err := tx.Create(reg).Error; err != nil {
		return fmt.Errorf("create registration: %w", err)
	}
	return nil
}

func (r *RegistrationRepository) GetByID(id int64) (*model.Registration, error) {
	var reg model.Registration
	if err := r.db.First(&reg, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get registration by id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get registration by id: %w", err)
	}
	return &reg, nil
}

func (r *RegistrationRepository) GetByTeamAndActivity(teamID, activityID int64) (*model.Registration, error) {
	var reg model.Registration
	if err := r.db.Where("team_id = ? AND activity_id = ?", teamID, activityID).First(&reg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get registration by team/activity: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get registration by team/activity: %w", err)
	}
	return &reg, nil
}

// GetByTeamAndActivityTx 事务内查询某团队在某活动的报名记录（不存在返回 ErrNotFound）。
func (r *RegistrationRepository) GetByTeamAndActivityTx(tx *gorm.DB, teamID, activityID int64) (*model.Registration, error) {
	var reg model.Registration
	if err := tx.Where("team_id = ? AND activity_id = ?", teamID, activityID).First(&reg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get registration by team/activity tx: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get registration by team/activity tx: %w", err)
	}
	return &reg, nil
}

func (r *RegistrationRepository) UpdateStatus(tx *gorm.DB, id int64, status string) error {
	if err := tx.Model(&model.Registration{}).Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("update registration status: %w", err)
	}
	return nil
}

// UpdateStatusIf 仅当报名记录当前状态为 expect 时才更新为 next（CAS 条件更新，防止两个管理员并发操作）。
// 返回受影响行数：0 表示状态已被其他事务改动。
func (r *RegistrationRepository) UpdateStatusIf(tx *gorm.DB, id int64, expect, next string) (int64, error) {
	res := tx.Model(&model.Registration{}).
		Where("id = ? AND status = ?", id, expect).
		Update("status", next)
	if res.Error != nil {
		return 0, fmt.Errorf("update registration status if: %w", res.Error)
	}
	return res.RowsAffected, nil
}

// GetByIDForUpdate 行锁查询报名记录，供审核/拒绝事务使用。
func (r *RegistrationRepository) GetByIDForUpdate(tx *gorm.DB, id int64) (*model.Registration, error) {
	var reg model.Registration
	if err := tx.Clauses(lockedClause()).First(&reg, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get registration for update: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get registration for update: %w", err)
	}
	return &reg, nil
}

// CountOccupied 统计活动当前占用名额的报名数（待审核 + 已通过 + 已完成）。
// 必须在持有活动行锁的事务内调用，以保证计数与插入的原子性。
func (r *RegistrationRepository) CountOccupied(tx *gorm.DB, activityID int64) (int64, error) {
	var count int64
	if err := tx.Model(&model.Registration{}).
		Where("activity_id = ? AND status IN ?", activityID,
			constants.OccupiedRegistrationStatuses).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count occupied registrations: %w", err)
	}
	return count, nil
}

// GetEarliestWaitlist 行锁取出活动最早提交的候补队伍（按提交时间、id 升序）。
func (r *RegistrationRepository) GetEarliestWaitlist(tx *gorm.DB, activityID int64) (*model.Registration, error) {
	var reg model.Registration
	err := tx.Clauses(lockedClause()).
		Where("activity_id = ? AND status = ?", activityID, constants.RegistrationStatusWaitlist).
		Order("registered_at ASC, id ASC").First(&reg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get earliest waitlist: %w", err)
	}
	return &reg, nil
}

// CountWaitlistAhead 统计指定候补记录前面还有多少候补队伍（按提交时间、id 升序排队）。
func (r *RegistrationRepository) CountWaitlistAhead(tx *gorm.DB, activityID int64, registeredAt time.Time, selfID int64) (int64, error) {
	var count int64
	if err := tx.Model(&model.Registration{}).
		Where("activity_id = ? AND status = ? AND (registered_at < ? OR (registered_at = ? AND id < ?))",
			activityID, constants.RegistrationStatusWaitlist, registeredAt, registeredAt, selfID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count waitlist ahead: %w", err)
	}
	return count, nil
}

// CountWaitlistAheadDB 非事务版本：查询单条候补记录前面还有多少队伍。
func (r *RegistrationRepository) CountWaitlistAheadDB(activityID int64, registeredAt time.Time, selfID int64) (int64, error) {
	return r.CountWaitlistAhead(r.db, activityID, registeredAt, selfID)
}

func (r *RegistrationRepository) Update(tx *gorm.DB, reg *model.Registration) error {
	if err := tx.Save(reg).Error; err != nil {
		return fmt.Errorf("update registration: %w", err)
	}
	return nil
}

func (r *RegistrationRepository) ListByTeam(teamID int64) ([]model.Registration, error) {
	var regs []model.Registration
	if err := r.db.Where("team_id = ?", teamID).Order("id DESC").Find(&regs).Error; err != nil {
		return nil, fmt.Errorf("list registrations by team: %w", err)
	}
	return regs, nil
}

func (r *RegistrationRepository) ListByActivity(activityID int64) ([]model.Registration, error) {
	var regs []model.Registration
	// 按提交顺序展示（待审核/候补排队顺序），排行榜会自行重排。
	if err := r.db.Where("activity_id = ?", activityID).
		Order("registered_at ASC, id ASC").Find(&regs).Error; err != nil {
		return nil, fmt.Errorf("list registrations by activity: %w", err)
	}
	return regs, nil
}

func (r *RegistrationRepository) ListByUserTeams(userID int64) ([]model.Registration, error) {
	var regs []model.Registration
	if err := r.db.
		Joins("JOIN team_members ON team_members.team_id = registrations.team_id").
		Where("team_members.user_id = ?", userID).
		Order("registrations.id DESC").
		Find(&regs).Error; err != nil {
		return nil, fmt.Errorf("list registrations by user teams: %w", err)
	}
	return regs, nil
}

func (r *RegistrationRepository) ListByActivityForUpdate(tx *gorm.DB, activityID int64) ([]model.Registration, error) {
	var regs []model.Registration
	if err := tx.Clauses(lockedClause()).
		Where("activity_id = ?", activityID).Find(&regs).Error; err != nil {
		return nil, fmt.Errorf("list registrations for update: %w", err)
	}
	return regs, nil
}

func (r *RegistrationRepository) StartAll(tx *gorm.DB, activityID int64) error {
	if err := tx.Model(&model.Registration{}).
		Where("activity_id = ? AND status = ?", activityID, "approved").
		Updates(map[string]any{"status": "approved", "start_time": gorm.Expr("NOW()")}).Error; err != nil {
		return fmt.Errorf("start all registrations: %w", err)
	}
	return nil
}

// WithTx 在事务中执行业务函数。
func (r *TeamRepository) WithTx(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// WithTx 在事务中执行业务函数。
func (r *RegistrationRepository) WithTx(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
