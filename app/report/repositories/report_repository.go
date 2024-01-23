package repositories

import (
	"gdsc/baro/app/report/models"

	"gorm.io/gorm"
)

type ReportRepositoryInterface interface {
	Save(report *models.Report) error
	FindByUserID(userID uint) ([]models.Report, error)
	FindById(id uint) (models.Report, error)
	FindByYearAndMonth(userID uint, month string) ([]models.Report, error)
	FindAll() ([]models.Report, error)
}

type ReportRepository struct {
	DB *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{
		DB: db,
	}
}

func (repo *ReportRepository) Save(report *models.Report) error {
	if err := repo.DB.Create(report).Error; err != nil {
		return err
	}
	return nil
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
