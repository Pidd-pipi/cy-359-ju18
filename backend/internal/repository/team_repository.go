package repository

import (
	"errors"
	"fmt"

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

// GetByIDTx 事务内读取报名记录（不加锁，用于加锁前取得 activity_id）。
func (r *RegistrationRepository) GetByIDTx(tx *gorm.DB, id int64) (*model.Registration, error) {
	var reg model.Registration
	if err := tx.First(&reg, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get registration by id in tx: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get registration by id in tx: %w", err)
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

// GetByTeamAndActivityForUpdate 事务内查询某团队对某活动的报名（行锁），用于并发下的重复报名判定。
func (r *RegistrationRepository) GetByTeamAndActivityForUpdate(tx *gorm.DB, teamID, activityID int64) (*model.Registration, error) {
	var reg model.Registration
	if err := tx.Clauses(lockedClause()).
		Where("team_id = ? AND activity_id = ?", teamID, activityID).First(&reg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get registration by team/activity for update: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("get registration by team/activity for update: %w", err)
	}
	return &reg, nil
}

// GetByIDForUpdate 事务内查询报名记录（行锁）。
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

// FirstWaitlistedForUpdate 锁定并返回活动候补队列中提交最早的一队（按 id 升序，id 由序列生成即提交顺序）。
// 找不到时返回包装后的 ErrNotFound。
func (r *RegistrationRepository) FirstWaitlistedForUpdate(tx *gorm.DB, activityID int64) (*model.Registration, error) {
	var reg model.Registration
	if err := tx.Clauses(lockedClause()).
		Where("activity_id = ? AND status = ?", activityID, constants.RegistrationStatusWaitlisted).
		Order("id ASC").First(&reg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("first waitlisted: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("first waitlisted: %w", err)
	}
	return &reg, nil
}

// CountWaitlistedBefore 统计指定候补报名之前还有多少候补队伍（id 更小即提交更早）。
// tx 为 nil 时使用默认连接（事务内传 tx，列表展示传 nil，与 CheckinRepository 惯例一致）。
func (r *RegistrationRepository) CountWaitlistedBefore(tx *gorm.DB, activityID, regID int64) (int64, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	var count int64
	if err := db.Model(&model.Registration{}).
		Where("activity_id = ? AND status = ? AND id < ?",
			activityID, constants.RegistrationStatusWaitlisted, regID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count waitlisted before: %w", err)
	}
	return count, nil
}

// CountWaitlistedBeforeIDs 批量计算给定报名 id 各自在候补队列中的排位（前面还有几队），
// 排位按全局候补队列计算，调用方即使只拿到过滤后的子集也能得到正确位置。
// 结果扫描进 out（元素须含 `ID int64` 与 `Ahead int64` 字段）。
func (r *RegistrationRepository) CountWaitlistedBeforeIDs(tx *gorm.DB, activityID int64, ids []int64, out any) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	if len(ids) == 0 {
		return nil
	}
	if err := db.Model(&model.Registration{}).
		Select("id, (SELECT COUNT(*) FROM registrations AS w2 "+
			"WHERE w2.activity_id = registrations.activity_id "+
			"AND w2.status = ? AND w2.id < registrations.id) AS ahead",
			constants.RegistrationStatusWaitlisted).
		Where("activity_id = ? AND id IN ?", activityID, ids).
		Find(out).Error; err != nil {
		return fmt.Errorf("count waitlisted before by ids: %w", err)
	}
	return nil
}

// UpdateStatus CAS 更新：仅当当前状态在 expectStatuses 中时才把状态改为 target。
// 返回受影响行数；0 表示状态已被其他事务改变（如两个管理员同时审核同一条报名）。
func (r *RegistrationRepository) UpdateStatusCAS(
	tx *gorm.DB, id int64, target string, expectStatuses ...string,
) (int64, error) {
	res := tx.Model(&model.Registration{}).
		Where("id = ? AND status IN ?", id, expectStatuses).
		Update("status", target)
	if res.Error != nil {
		return 0, fmt.Errorf("cas update registration status: %w", res.Error)
	}
	return res.RowsAffected, nil
}

func (r *RegistrationRepository) UpdateStatus(tx *gorm.DB, id int64, status string) error {
	if err := tx.Model(&model.Registration{}).Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("update registration status: %w", err)
	}
	return nil
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
	// 报名管理列表按提交顺序展示（候补队列的先后与 id 升序一致）。
	if err := r.db.Where("activity_id = ?", activityID).Order("registered_at ASC, id ASC").Find(&regs).Error; err != nil {
		return nil, fmt.Errorf("list registrations by activity: %w", err)
	}
	return regs, nil
}

// NamesByIDs 按 id 批量查询团队名。
func (r *TeamRepository) NamesByIDs(ids []int64) (map[int64]string, error) {
	names := make(map[int64]string, len(ids))
	if len(ids) == 0 {
		return names, nil
	}
	var teams []model.Team
	if err := r.db.Select("id", "name").Where("id IN ?", ids).Find(&teams).Error; err != nil {
		return nil, fmt.Errorf("list team names: %w", err)
	}
	for i := range teams {
		names[teams[i].ID] = teams[i].Name
	}
	return names, nil
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
