package repositories

import (
	"fmt"
	"gdsc/baro/app/report/models"
	"gdsc/baro/app/report/types"
	users "gdsc/baro/app/user/models"
	"time"

	"gorm.io/gorm"
)

type ReportRepositoryInterface interface {
	Save(report *models.Report) (models.Report, error)
	FindByUserID(userID uint) ([]models.Report, error)
	FindById(id uint) (models.Report, error)
	FindByYearAndMonth(userID uint, month string) ([]models.Report, error)
	FindAll() ([]models.Report, error)
	FindRankAtAgeAndGender(user *users.User, start, end time.Time) (types.ResponseRank, error)
}

type ReportRepository struct {
	DB *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{
		DB: db,
	}
}

func (repo *ReportRepository) Save(report *models.Report) (models.Report, error) {
	if err := repo.DB.Create(report).Error; err != nil {
		return models.Report{}, err
	}
	return *report, nil
}

func (repo *ReportRepository) FindByUserID(userID uint) ([]models.Report, error) {
	var reports []models.Report
	result := repo.DB.Where("user_id = ?", userID).Find(&reports)
	return reports, result.Error
}

func (repo *ReportRepository) FindById(id uint) (models.Report, error) {
	var report models.Report
	result := repo.DB.Where("id = ?", id).First(&report)
	return report, result.Error
}

func (repo *ReportRepository) FindByYearAndMonth(userID uint, yearAndMonth string) ([]models.Report, error) {
	var reports []models.Report
	result := repo.DB.Where("user_id = ? AND DATE_FORMAT(created_at, '%Y%m') = ?", userID, yearAndMonth).Find(&reports)
	return reports, result.Error
}

func (repo *ReportRepository) FindAll() ([]models.Report, error) {
	var reports []models.Report
	result := repo.DB.Find(&reports)
	return reports, result.Error
}

func (repo *ReportRepository) FindRankAtAgeAndGender(user *users.User, start, end time.Time) (types.ResponseRank, error) {
	var userAvgScore, allAvgScore float64
	var totalUsers, rank int64
	ageGroup := user.Age / 10 * 10

	// calculate average score for the user
	err := repo.DB.Table("reports").
		Select("avg(score)").
		Joins("inner join users on users.id = reports.user_id").
		Where("reports.user_id = ?", user.ID).
		Where("reports.created_at BETWEEN ? AND ?", start, end).
		Scan(&userAvgScore).Error

	if err != nil {
		return types.ResponseRank{}, err
	}

	// calculate total users in the same age group and gender
	err = repo.DB.Table("users").
		Where("age >= ? AND age < ?", ageGroup, ageGroup+10).
		Where("gender = ?", user.Gender).
		Count(&totalUsers).Error

	if err != nil {
		return types.ResponseRank{}, err
	}

	// calculate average score for all users in the same age group and gender
	err = repo.DB.Table("reports").
		Select("avg(score)").
		Joins("inner join users on users.id = reports.user_id").
		Where("users.age >= ? AND users.age < ?", ageGroup, ageGroup+10).
		Where("users.gender = ?", user.Gender).
		Where("reports.created_at BETWEEN ? AND ?", start, end).
		Scan(&allAvgScore).Error

	if err != nil {
		return types.ResponseRank{}, err
	}

	// calculate rank
	sql := `
	SELECT COUNT(*) as rank_count
	FROM (
		SELECT reports.user_id, AVG(score) as average_score
		FROM reports
		INNER JOIN users on users.id = reports.user_id
		WHERE users.age >= ? AND users.age < ?
		AND users.gender = ?
		AND reports.created_at BETWEEN ? AND ?
		GROUP BY reports.user_id
	) as subquery
	WHERE average_score > ?
	`

	err = repo.DB.Raw(sql, ageGroup, ageGroup+10, user.Gender, start, end, userAvgScore).Scan(&rank).Error
	if err != nil {
		return types.ResponseRank{}, err
	}

	normalRatio := fmt.Sprintf("%.2f", (float64(rank+1)/float64(totalUsers))*100)
	averageScore := fmt.Sprintf("%.2f", userAvgScore)
	allAvgScoreStr := fmt.Sprintf("%.2f", allAvgScore)

	return types.ResponseRank{
		UserID:          user.ID,
		Nickname:        user.Nickname,
		Age:             user.Age,
		Gender:          user.Gender,
		NormalRatio:     normalRatio,
		AverageScore:    averageScore,
		AllAverageScore: allAvgScoreStr,
	}, nil
}
