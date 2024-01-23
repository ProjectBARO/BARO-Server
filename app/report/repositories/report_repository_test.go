package repositories_test

import (
	"gdsc/baro/app/report/models"
	"gdsc/baro/app/report/repositories"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestReportrRepository_Save(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).
			AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create ReportRepository
	reportRepository := repositories.NewReportRepository(gormDB)

	// Create sample report for the test
	report := models.Report{
		UserID:       1,
		AlertCount:   30,
		AnalysisTime: 3600,
		Predict:      "Good",
		Type:         "Study",
		CreatedAt:    time.Now(),
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(report.UserID, report.AlertCount, report.AnalysisTime, report.Type, report.Predict, report.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method under test
	createdReport, err := reportRepository.Save(&report)
	if err != nil {
		t.Fatalf("Error creating report: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check the result
	assert.Equal(t, report.UserID, createdReport.UserID)
	assert.Equal(t, report.AlertCount, createdReport.AlertCount)
	assert.Equal(t, report.AnalysisTime, createdReport.AnalysisTime)
	assert.Equal(t, report.Predict, createdReport.Predict)
	assert.Equal(t, report.Type, createdReport.Type)
}

func TestReportRepository_Save_Error(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).
			AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create ReportRepository
	reportRepository := repositories.NewReportRepository(gormDB)

	// Create sample report for the test
	report := models.Report{
		UserID:       1,
		AlertCount:   30,
		AnalysisTime: 3600,
		Predict:      "Good",
		Type:         "Study",
		CreatedAt:    time.Now(),
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(report.UserID, report.AlertCount, report.AnalysisTime, report.Type, report.Predict, report.CreatedAt).
		WillReturnError(err)
	mock.ExpectRollback()

	// Call the method under test
	_, err = reportRepository.Save(&report)
	if err == nil {
		t.Fatalf("Error creating report: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_FindByUserID(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).
			AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create ReportRepository
	reportRepository := repositories.NewReportRepository(gormDB)

	// Create sample report for the test
	report := models.Report{
		UserID:       1,
		AlertCount:   30,
		AnalysisTime: 3600,
		Predict:      "Good",
		Type:         "Study",
		CreatedAt:    time.Now(),
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(report.UserID, report.AlertCount, report.AnalysisTime, report.Type, report.Predict, report.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method under test
	savedReport, err := reportRepository.Save(&report)
	if err != nil {
		t.Fatalf("Error creating report: %v", err)
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT \\* FROM `reports` WHERE user_id = \\?").
		WithArgs(report.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "alert_count", "analysis_time", "type", "predict", "created_at"}).
			AddRow(savedReport.ID, savedReport.UserID, savedReport.AlertCount, savedReport.AnalysisTime, savedReport.Type, savedReport.Predict, savedReport.CreatedAt))

	// Call the method under test
	reports, err := reportRepository.FindByUserID(report.UserID)
	if err != nil {
		t.Fatalf("Error finding reports: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check the result
	assert.Equal(t, savedReport.UserID, reports[0].UserID)
	assert.Equal(t, savedReport.AlertCount, reports[0].AlertCount)
	assert.Equal(t, savedReport.AnalysisTime, reports[0].AnalysisTime)
	assert.Equal(t, savedReport.Predict, reports[0].Predict)
	assert.Equal(t, savedReport.Type, reports[0].Type)
}

func TestReportRepository_FindByUserID_Error(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).
			AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create ReportRepository
	reportRepository := repositories.NewReportRepository(gormDB)

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT \\* FROM `reports` WHERE user_id = \\?").
		WithArgs(1).
		WillReturnError(err)

	// Call the method under test
	_, err = reportRepository.FindByUserID(1)
	if err == nil {
		t.Fatalf("Error finding reports: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_FindById(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).
			AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create ReportRepository
	reportRepository := repositories.NewReportRepository(gormDB)

	// Create sample report for the test
	report := models.Report{
		UserID:       1,
		AlertCount:   30,
		AnalysisTime: 3600,
		Predict:      "Good",
		Type:         "Study",
		CreatedAt:    time.Now(),
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(report.UserID, report.AlertCount, report.AnalysisTime, report.Type, report.Predict, report.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method under test
	savedReport, err := reportRepository.Save(&report)
	if err != nil {
		t.Fatalf("Error creating report: %v", err)
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT \\* FROM `reports` WHERE id = \\?").
		WithArgs(savedReport.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "alert_count", "analysis_time", "type", "predict", "created_at"}).
			AddRow(savedReport.ID, savedReport.UserID, savedReport.AlertCount, savedReport.AnalysisTime, savedReport.Type, savedReport.Predict, savedReport.CreatedAt))

	// Call the method under test
	report, err = reportRepository.FindById(savedReport.ID)
	if err != nil {
		t.Fatalf("Error finding report: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check the result
	assert.Equal(t, savedReport.UserID, report.UserID)
	assert.Equal(t, savedReport.AlertCount, report.AlertCount)
	assert.Equal(t, savedReport.AnalysisTime, report.AnalysisTime)
	assert.Equal(t, savedReport.Predict, report.Predict)
	assert.Equal(t, savedReport.Type, report.Type)
}

func TestReportRepository_FindById_Error(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).
			AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create ReportRepository
	reportRepository := repositories.NewReportRepository(gormDB)

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT \\* FROM `reports` WHERE id = \\?").
		WithArgs(1).
		WillReturnError(err)

	// Call the method under test
	_, err = reportRepository.FindById(1)
	if err == nil {
		t.Fatalf("Error finding report: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_FindByYearAndMonth(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).
			AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create ReportRepository
	reportRepository := repositories.NewReportRepository(gormDB)

	// Create sample report for the test
	report := models.Report{
		UserID:       1,
		AlertCount:   30,
		AnalysisTime: 3600,
		Predict:      "Good",
		Type:         "Study",
		CreatedAt:    time.Now(),
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(report.UserID, report.AlertCount, report.AnalysisTime, report.Type, report.Predict, report.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method under test
	savedReport, err := reportRepository.Save(&report)
	if err != nil {
		t.Fatalf("Error creating report: %v", err)
	}

	// YearAndMonth
	yearAndMonth := time.Now().Format("200601")

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT \\* FROM `reports` WHERE user_id = \\? AND DATE_FORMAT\\(created_at, '%Y%m'\\) = \\?").
		WithArgs(savedReport.UserID, yearAndMonth).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "alert_count", "analysis_time", "type", "predict", "created_at"}).
			AddRow(savedReport.ID, savedReport.UserID, savedReport.AlertCount, savedReport.AnalysisTime, savedReport.Type, savedReport.Predict, savedReport.CreatedAt))

	// Call the method under test
	reports, err := reportRepository.FindByYearAndMonth(savedReport.UserID, yearAndMonth)
	if err != nil {
		t.Fatalf("Error finding report: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check the result
	assert.Equal(t, savedReport.UserID, reports[0].UserID)
	assert.Equal(t, savedReport.AlertCount, reports[0].AlertCount)
	assert.Equal(t, savedReport.AnalysisTime, reports[0].AnalysisTime)
	assert.Equal(t, savedReport.Predict, reports[0].Predict)
	assert.Equal(t, savedReport.Type, reports[0].Type)
}

func TestReportRepository_FindByYearAndMonth_Error(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).
			AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create ReportRepository
	reportRepository := repositories.NewReportRepository(gormDB)

	// YearAndMonth
	yearAndMonth := time.Now().Format("200601")

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT \\* FROM `reports` WHERE user_id = \\? AND DATE_FORMAT\\(created_at, '%Y%m'\\) = \\?").
		WithArgs(1, yearAndMonth).
		WillReturnError(err)

	// Call the method under test
	_, err = reportRepository.FindByYearAndMonth(1, yearAndMonth)
	if err == nil {
		t.Fatalf("Error finding report: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_FindAll(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).
			AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create ReportRepository
	reportRepository := repositories.NewReportRepository(gormDB)

	// Create sample report for the test
	report := models.Report{
		UserID:       1,
		AlertCount:   30,
		AnalysisTime: 3600,
		Predict:      "Good",
		Type:         "Study",
		CreatedAt:    time.Now(),
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(report.UserID, report.AlertCount, report.AnalysisTime, report.Type, report.Predict, report.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method under test
	savedReport, err := reportRepository.Save(&report)
	if err != nil {
		t.Fatalf("Error creating report: %v", err)
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT \\* FROM `reports`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "alert_count", "analysis_time", "type", "predict", "created_at"}).
			AddRow(savedReport.ID, savedReport.UserID, savedReport.AlertCount, savedReport.AnalysisTime, savedReport.Type, savedReport.Predict, savedReport.CreatedAt))

	// Call the method under test
	reports, err := reportRepository.FindAll()
	if err != nil {
		t.Fatalf("Error finding report: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check the result
	assert.Equal(t, savedReport.UserID, reports[0].UserID)
	assert.Equal(t, savedReport.AlertCount, reports[0].AlertCount)
	assert.Equal(t, savedReport.AnalysisTime, reports[0].AnalysisTime)
	assert.Equal(t, savedReport.Predict, reports[0].Predict)
	assert.Equal(t, savedReport.Type, reports[0].Type)
}

func TestReportRepository_FindAll_Error(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).
			AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create ReportRepository
	reportRepository := repositories.NewReportRepository(gormDB)

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT \\* FROM `reports`").
		WillReturnError(err)

	// Call the method under test
	_, err = reportRepository.FindAll()
	if err == nil {
		t.Fatalf("Error finding report: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
