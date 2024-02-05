package repositories_test

import (
	"gdsc/baro/app/report/models"
	"gdsc/baro/app/report/repositories"
	usermodel "gdsc/baro/app/user/models"
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
		UserID:            1,
		AlertCount:        30,
		AnalysisTime:      3600,
		Type:              "Study",
		Predict:           "Good",
		Score:             "90.000",
		NormalRatio:       "90.000",
		NeckAngles:        "angle",
		Distances:         "distance",
		StatusFrequencies: "status",
		CreatedAt:         time.Now(),
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(report.UserID, report.AlertCount, report.AnalysisTime, report.Type, report.Predict, report.Score, report.NormalRatio, report.NeckAngles, report.Distances, report.StatusFrequencies, report.CreatedAt).
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
	assert.Equal(t, report.Type, createdReport.Type)
	assert.Equal(t, report.Predict, createdReport.Predict)
	assert.Equal(t, report.Score, createdReport.Score)
	assert.Equal(t, report.NormalRatio, createdReport.NormalRatio)
	assert.Equal(t, report.NeckAngles, createdReport.NeckAngles)
	assert.Equal(t, report.Distances, createdReport.Distances)
	assert.Equal(t, report.StatusFrequencies, createdReport.StatusFrequencies)
	assert.Equal(t, report.CreatedAt, createdReport.CreatedAt)
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
		UserID:            1,
		AlertCount:        30,
		AnalysisTime:      3600,
		Type:              "Study",
		Predict:           "Good",
		Score:             "90.000",
		NormalRatio:       "90.000",
		NeckAngles:        "angle",
		Distances:         "distance",
		StatusFrequencies: "status",
		CreatedAt:         time.Now(),
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(report.UserID, report.AlertCount, report.AnalysisTime, report.Type, report.Predict, report.Score, report.NormalRatio, report.NeckAngles, report.Distances, report.StatusFrequencies, report.CreatedAt).
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
		UserID:            1,
		AlertCount:        30,
		AnalysisTime:      3600,
		Type:              "Study",
		Predict:           "Good",
		Score:             "90.000",
		NormalRatio:       "90.000",
		NeckAngles:        "angle",
		Distances:         "distance",
		StatusFrequencies: "status",
		CreatedAt:         time.Now(),
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(report.UserID, report.AlertCount, report.AnalysisTime, report.Type, report.Predict, report.Score, report.NormalRatio, report.NeckAngles, report.Distances, report.StatusFrequencies, report.CreatedAt).
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
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "alert_count", "analysis_time", "type", "predict", "score", "normal_ratio", "neck_angles", "distances", "status_frequencies", "created_at"}).
			AddRow(savedReport.ID, savedReport.UserID, savedReport.AlertCount, savedReport.AnalysisTime, savedReport.Type, savedReport.Predict, savedReport.Score, savedReport.NormalRatio, savedReport.NeckAngles, savedReport.Distances, savedReport.StatusFrequencies, savedReport.CreatedAt))

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
	assert.Equal(t, savedReport.Type, reports[0].Type)
	assert.Equal(t, savedReport.Predict, reports[0].Predict)
	assert.Equal(t, savedReport.Score, reports[0].Score)
	assert.Equal(t, savedReport.NormalRatio, reports[0].NormalRatio)
	assert.Equal(t, savedReport.NeckAngles, reports[0].NeckAngles)
	assert.Equal(t, savedReport.Distances, reports[0].Distances)
	assert.Equal(t, savedReport.StatusFrequencies, reports[0].StatusFrequencies)
	assert.Equal(t, savedReport.CreatedAt, reports[0].CreatedAt)
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
		UserID:            1,
		AlertCount:        30,
		AnalysisTime:      3600,
		Type:              "Study",
		Predict:           "Good",
		Score:             "90.000",
		NormalRatio:       "90.000",
		NeckAngles:        "angle",
		Distances:         "distance",
		StatusFrequencies: "status",
		CreatedAt:         time.Now(),
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(report.UserID, report.AlertCount, report.AnalysisTime, report.Type, report.Predict, report.Score, report.NormalRatio, report.NeckAngles, report.Distances, report.StatusFrequencies, report.CreatedAt).
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
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "alert_count", "analysis_time", "type", "predict", "score", "normal_ratio", "neck_angles", "distances", "status_frequencies", "created_at"}).
			AddRow(savedReport.ID, savedReport.UserID, savedReport.AlertCount, savedReport.AnalysisTime, savedReport.Type, savedReport.Predict, savedReport.Score, savedReport.NormalRatio, savedReport.NeckAngles, savedReport.Distances, savedReport.StatusFrequencies, savedReport.CreatedAt))

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
	assert.Equal(t, savedReport.Type, report.Type)
	assert.Equal(t, savedReport.Predict, report.Predict)
	assert.Equal(t, savedReport.Score, report.Score)
	assert.Equal(t, savedReport.NormalRatio, report.NormalRatio)
	assert.Equal(t, savedReport.NeckAngles, report.NeckAngles)
	assert.Equal(t, savedReport.Distances, report.Distances)
	assert.Equal(t, savedReport.StatusFrequencies, report.StatusFrequencies)
	assert.Equal(t, savedReport.CreatedAt, report.CreatedAt)
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
		UserID:            1,
		AlertCount:        30,
		AnalysisTime:      3600,
		Type:              "Study",
		Predict:           "Good",
		Score:             "90.000",
		NormalRatio:       "90.000",
		NeckAngles:        "angle",
		Distances:         "distance",
		StatusFrequencies: "status",
		CreatedAt:         time.Now(),
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(report.UserID, report.AlertCount, report.AnalysisTime, report.Type, report.Predict, report.Score, report.NormalRatio, report.NeckAngles, report.Distances, report.StatusFrequencies, report.CreatedAt).
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
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "alert_count", "analysis_time", "type", "predict", "score", "normal_ratio", "neck_angles", "distances", "status_frequencies", "created_at"}).
			AddRow(savedReport.ID, savedReport.UserID, savedReport.AlertCount, savedReport.AnalysisTime, savedReport.Type, savedReport.Predict, savedReport.Score, savedReport.NormalRatio, savedReport.NeckAngles, savedReport.Distances, savedReport.StatusFrequencies, savedReport.CreatedAt))

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
	assert.Equal(t, savedReport.Type, reports[0].Type)
	assert.Equal(t, savedReport.Predict, reports[0].Predict)
	assert.Equal(t, savedReport.Score, reports[0].Score)
	assert.Equal(t, savedReport.NormalRatio, reports[0].NormalRatio)
	assert.Equal(t, savedReport.NeckAngles, reports[0].NeckAngles)
	assert.Equal(t, savedReport.Distances, reports[0].Distances)
	assert.Equal(t, savedReport.StatusFrequencies, reports[0].StatusFrequencies)
	assert.Equal(t, savedReport.CreatedAt, reports[0].CreatedAt)
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
		UserID:            1,
		AlertCount:        30,
		AnalysisTime:      3600,
		Type:              "Study",
		Predict:           "Good",
		Score:             "90.000",
		NormalRatio:       "90.000",
		NeckAngles:        "angle",
		Distances:         "distance",
		StatusFrequencies: "status",
		CreatedAt:         time.Now(),
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(report.UserID, report.AlertCount, report.AnalysisTime, report.Type, report.Predict, report.Score, report.NormalRatio, report.NeckAngles, report.Distances, report.StatusFrequencies, report.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method under test
	savedReport, err := reportRepository.Save(&report)
	if err != nil {
		t.Fatalf("Error creating report: %v", err)
	}

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT \\* FROM `reports`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "alert_count", "analysis_time", "type", "predict", "score", "normal_ratio", "neck_angles", "distances", "status_frequencies", "created_at"}).
			AddRow(savedReport.ID, savedReport.UserID, savedReport.AlertCount, savedReport.AnalysisTime, savedReport.Type, savedReport.Predict, savedReport.Score, savedReport.NormalRatio, savedReport.NeckAngles, savedReport.Distances, savedReport.StatusFrequencies, savedReport.CreatedAt))

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
	assert.Equal(t, savedReport.Type, reports[0].Type)
	assert.Equal(t, savedReport.Predict, reports[0].Predict)
	assert.Equal(t, savedReport.Score, reports[0].Score)
	assert.Equal(t, savedReport.NormalRatio, reports[0].NormalRatio)
	assert.Equal(t, savedReport.NeckAngles, reports[0].NeckAngles)
	assert.Equal(t, savedReport.Distances, reports[0].Distances)
	assert.Equal(t, savedReport.StatusFrequencies, reports[0].StatusFrequencies)
	assert.Equal(t, savedReport.CreatedAt, reports[0].CreatedAt)
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

func TestReportRepository_FindRankAtAgeAndGender(t *testing.T) {
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

	// Create sample user for the test
	user := usermodel.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
	}

	start := time.Now().AddDate(0, 0, -30)
	end := time.Now()

	ageGroup := 20
	userAvgScore := 80.0
	allAvgScore := 90.0

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT avg\\(score\\) FROM `reports` inner join users on users.id = reports.user_id WHERE reports.user_id = \\? AND \\(reports.created_at BETWEEN \\? AND \\?\\)").
		WithArgs(user.ID, start, end).
		WillReturnRows(sqlmock.NewRows([]string{"avg"}).
			AddRow("80.0"))

	// Set up expectations for the mock DB to return total users
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `users` WHERE \\(age >= \\? AND age < \\?\\) AND gender = \\?").
		WithArgs(ageGroup, ageGroup+10, user.Gender).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow("100"))

	// Set up expectations for the mock DB to return average score for all users in the same age group and gender
	mock.ExpectQuery("SELECT avg\\(score\\) FROM `reports` inner join users on users.id = reports.user_id WHERE \\(users.age >= \\? AND users.age < \\?\\) AND users.gender = \\? AND \\(reports.created_at BETWEEN \\? AND \\?\\)").
		WithArgs(ageGroup, ageGroup+10, user.Gender, start, end).
		WillReturnRows(sqlmock.NewRows([]string{"avg"}).
			AddRow(allAvgScore))

	// Set up expectations for the mock DB to return rank
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) as rank_count FROM \\(\\s*SELECT reports.user_id, AVG\\(score\\) as average_score FROM reports INNER JOIN users\\s+on\\s+users.id = reports.user_id WHERE users.age >= \\? AND users.age < \\? AND users.gender = \\? AND reports.created_at BETWEEN \\? AND \\? GROUP BY reports.user_id\\s*\\) as subquery WHERE average_score > \\?").
		WithArgs(ageGroup, ageGroup+10, user.Gender, start, end, userAvgScore).
		WillReturnRows(sqlmock.NewRows([]string{"rank_count"}).
			AddRow("20"))

	// Call the method under test
	responseRank, err := reportRepository.FindRankAtAgeAndGender(&user, start, end)
	if err != nil {
		t.Fatalf("Error finding rank: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check the result
	assert.Equal(t, user.ID, responseRank.UserID)
	assert.Equal(t, user.Nickname, responseRank.Nickname)
	assert.Equal(t, user.Age, responseRank.Age)
	assert.Equal(t, user.Gender, responseRank.Gender)
	assert.Equal(t, "21.00", responseRank.NormalRatio)
	assert.Equal(t, "80.00", responseRank.AverageScore)
	assert.Equal(t, "90.00", responseRank.AllAverageScore)
}

func TestReportRepository_FindRankAtAgeAndGender_Error(t *testing.T) {
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

	// Create sample user for the test
	user := usermodel.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
	}

	start := time.Now().AddDate(0, 0, -30)
	end := time.Now()

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT avg\\(score\\) FROM `reports` inner join users on users.id = reports.user_id WHERE reports.user_id = \\? AND \\(reports.created_at BETWEEN \\? AND \\?\\)").
		WithArgs(user.ID, start, end).
		WillReturnError(err)

	// Call the method under test
	_, err = reportRepository.FindRankAtAgeAndGender(&user, start, end)
	if err == nil {
		t.Fatalf("Error finding rank: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_FindRankAtAgeAndGender_Error2(t *testing.T) {
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

	// Create sample user for the test
	user := usermodel.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
	}

	start := time.Now().AddDate(0, 0, -30)
	end := time.Now()

	ageGroup := 20

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT avg\\(score\\) FROM `reports` inner join users on users.id = reports.user_id WHERE reports.user_id = \\? AND \\(reports.created_at BETWEEN \\? AND \\?\\)").
		WithArgs(user.ID, start, end).
		WillReturnRows(sqlmock.NewRows([]string{"avg"}).
			AddRow("80.0"))

	// Set up expectations for the mock DB to return total users
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `users` WHERE \\(age >= \\? AND age < \\?\\) AND gender = \\?").
		WithArgs(ageGroup, ageGroup+10, user.Gender).
		WillReturnError(err)

	// Call the method under test
	_, err = reportRepository.FindRankAtAgeAndGender(&user, start, end)
	if err == nil {
		t.Fatalf("Error finding rank: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_FindRankAtAgeAndGender_Error3(t *testing.T) {
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

	// Create sample user for the test
	user := usermodel.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
	}

	start := time.Now().AddDate(0, 0, -30)
	end := time.Now()

	ageGroup := 20

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT avg\\(score\\) FROM `reports` inner join users on users.id = reports.user_id WHERE reports.user_id = \\? AND \\(reports.created_at BETWEEN \\? AND \\?\\)").
		WithArgs(user.ID, start, end).
		WillReturnRows(sqlmock.NewRows([]string{"avg"}).
			AddRow("80.0"))

	// Set up expectations for the mock DB to return total users
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `users` WHERE \\(age >= \\? AND age < \\?\\) AND gender = \\?").
		WithArgs(ageGroup, ageGroup+10, user.Gender).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow("100"))

	// Set up expectations for the mock DB to return average score for all users in the same age group and gender
	mock.ExpectQuery("SELECT avg\\(score\\) FROM `reports` inner join users on users.id = reports.user_id WHERE \\(users.age >= \\? AND users.age < \\?\\) AND users.gender = \\? AND \\(reports.created_at BETWEEN \\? AND \\?\\)").
		WithArgs(ageGroup, ageGroup+10, user.Gender, start, end).
		WillReturnError(err)

	// Call the method under test
	_, err = reportRepository.FindRankAtAgeAndGender(&user, start, end)
	if err == nil {
		t.Fatalf("Error finding rank: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReportRepository_FindRankAtAgeAndGender_Error4(t *testing.T) {
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

	// Create sample user for the test
	user := usermodel.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
	}

	start := time.Now().AddDate(0, 0, -30)
	end := time.Now()

	ageGroup := 20
	userAvgScore := 80.0
	allAvgScore := 90.0

	// Set up expectations for the mock DB to return the sample report
	mock.ExpectQuery("SELECT avg\\(score\\) FROM `reports` inner join users on users.id = reports.user_id WHERE reports.user_id = \\? AND \\(reports.created_at BETWEEN \\? AND \\?\\)").
		WithArgs(user.ID, start, end).
		WillReturnRows(sqlmock.NewRows([]string{"avg"}).
			AddRow("80.0"))

	// Set up expectations for the mock DB to return total users
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `users` WHERE \\(age >= \\? AND age < \\?\\) AND gender = \\?").
		WithArgs(ageGroup, ageGroup+10, user.Gender).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow("100"))

	// Set up expectations for the mock DB to return average score for all users in the same age group and gender
	mock.ExpectQuery("SELECT avg\\(score\\) FROM `reports` inner join users on users.id = reports.user_id WHERE \\(users.age >= \\? AND users.age < \\?\\) AND users.gender = \\? AND \\(reports.created_at BETWEEN \\? AND \\?\\)").
		WithArgs(ageGroup, ageGroup+10, user.Gender, start, end).
		WillReturnRows(sqlmock.NewRows([]string{"avg"}).
			AddRow(allAvgScore))

	// Set up expectations for the mock DB to return rank
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) as rank_count FROM \\(\\s*SELECT reports.user_id, AVG\\(score\\) as average_score FROM reports INNER JOIN users\\s+on\\s+users.id = reports.user_id WHERE users.age >= \\? AND users.age < \\? AND users.gender = \\? AND reports.created_at BETWEEN \\? AND \\? GROUP BY reports.user_id\\s*\\) as subquery WHERE average_score > \\?").
		WithArgs(ageGroup, ageGroup+10, user.Gender, start, end, userAvgScore).
		WillReturnError(err)

	// Call the method under test
	_, err = reportRepository.FindRankAtAgeAndGender(&user, start, end)
	if err == nil {
		t.Fatalf("Error finding rank: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
