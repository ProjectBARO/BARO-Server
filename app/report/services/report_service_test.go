package services_test

import (
	"errors"
	"gdsc/baro/app/report/models"
	"gdsc/baro/app/report/services"
	"gdsc/baro/app/report/types"
	usermodel "gdsc/baro/app/user/models"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/html"
)

type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) Save(report *models.Report) (models.Report, error) {
	args := m.Called(report)
	return *args.Get(0).(*models.Report), args.Error(1)
}

func (m *MockReportRepository) FindByUserID(userID uint) ([]models.Report, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Report), args.Error(1)
}

func (m *MockReportRepository) FindById(id uint) (models.Report, error) {
	args := m.Called(id)
	if tmp := args.Get(0); tmp != nil {
		return tmp.(models.Report), args.Error(1)
	}
	return models.Report{}, args.Error(1)
}

func (m *MockReportRepository) FindByYearAndMonth(userID uint, month string) ([]models.Report, error) {
	args := m.Called(userID, month)
	if tmp := args.Get(0); tmp != nil {
		return tmp.([]models.Report), args.Error(1)
	}
	return []models.Report{}, args.Error(1)
}

func (m *MockReportRepository) FindAll() ([]models.Report, error) {
	args := m.Called()
	return args.Get(0).([]models.Report), args.Error(1)
}

type MockUserUtil struct {
	mock.Mock
}

func (m *MockUserUtil) FindCurrentUser(c *gin.Context) (*usermodel.User, error) {
	args := m.Called(c)
	return args.Get(0).(*usermodel.User), args.Error(1)
}

func TestAnalysis(t *testing.T) {
	// Mock ReportRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	reportService := services.NewReportService(mockReportRepository, mockUserUtil)

	// Set up sample user for the test
	user := usermodel.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
	}

	// Set up expectations for the mock repository and util
	mockUserUtil.On("FindCurrentUser", mock.Anything).Return(&user, nil)

	// Create a test context
	c, _ := gin.CreateTestContext(nil)

	request := types.RequestAnalysis{
		VideoURL:     "test",
		AlertCount:   10,
		AnalysisTime: 1800,
		Type:         "Study",
	}

	// Call the service
	responseMessage, err := reportService.Analysis(c, request)
	assert.NoError(t, err)

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)

	// Check the results
	assert.Equal(t, "Video submitted successfully", responseMessage)
}

func TestAnalysis_NoUser(t *testing.T) {
	// Mock ReportRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	reportService := services.NewReportService(mockReportRepository, mockUserUtil)

	// Set up expectations for the mock repository and util
	mockUserUtil.On("FindCurrentUser", mock.Anything).Return((*usermodel.User)(nil), errors.New("record not found"))

	// Create a test context
	c, _ := gin.CreateTestContext(nil)

	request := types.RequestAnalysis{
		VideoURL:     "test",
		AlertCount:   10,
		AnalysisTime: 1800,
		Type:         "Study",
	}

	// Call the service
	_, err := reportService.Analysis(c, request)
	assert.Error(t, err)

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)
}

func TestParseHTML(t *testing.T) {
	// Create service
	service := services.NewReportService(nil, nil)

	// Set up test data
	testHTML := "<html><body><p>Test</p></body></html>"
	doc, err := html.Parse(strings.NewReader(testHTML))
	assert.Nil(t, err)

	// Execute method under test
	result := service.ParseHTML(doc)

	// Assert result
	assert.Equal(t, "Test", result)
}

func TestParseHTML_NoData(t *testing.T) {
	// Create service
	service := services.NewReportService(nil, nil)

	// Set up test data
	testHTML := "<html><body></body></html>"
	doc, err := html.Parse(strings.NewReader(testHTML))
	assert.Nil(t, err)

	// Execute method under test
	result := service.ParseHTML(doc)

	// Assert result
	assert.Equal(t, "", result)
}

func TestFindReportByCurrentUser(t *testing.T) {
	// Mock UserRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	reportService := services.NewReportService(mockReportRepository, mockUserUtil)

	// Set up sample user for the test
	user := usermodel.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
	}

	// Set up sample reports for the test
	reports := []models.Report{
		{
			ID:           1,
			UserID:       1,
			AlertCount:   10,
			AnalysisTime: 1800,
			Predict:      "Good",
			Type:         "Study",
		},
		{
			ID:           2,
			UserID:       1,
			AlertCount:   30,
			AnalysisTime: 3600,
			Predict:      "Good",
			Type:         "Study",
		},
	}

	// Set up expectations for the mock repository and util
	mockUserUtil.On("FindCurrentUser", mock.Anything).Return(&user, nil)
	mockReportRepository.On("FindByUserID", mock.Anything).Return(reports, nil)

	// Create a test context
	c, _ := gin.CreateTestContext(nil)

	// Call the service
	responseReports, err := reportService.FindReportByCurrentUser(c)
	assert.NoError(t, err)
	assert.Equal(t, len(reports), len(responseReports))

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)

	// Check the results
	for i, report := range reports {
		assert.Equal(t, report.ID, responseReports[i].ID)
		assert.Equal(t, report.UserID, responseReports[i].UserID)
		assert.Equal(t, report.AlertCount, responseReports[i].AlertCount)
		assert.Equal(t, report.AnalysisTime, responseReports[i].AnalysisTime)
		assert.Equal(t, report.Predict, responseReports[i].Predict)
		assert.Equal(t, report.Type, responseReports[i].Type)
	}
}

func TestFindReportByCurrentUser_NoUser(t *testing.T) {
	// Mock UserRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	reportService := services.NewReportService(mockReportRepository, mockUserUtil)

	// Set up expectations for the mock repository and util
	mockUserUtil.On("FindCurrentUser", mock.Anything).Return((*usermodel.User)(nil), errors.New("record not found"))

	// Create a test context
	c, _ := gin.CreateTestContext(nil)

	// Call the service
	_, err := reportService.FindReportByCurrentUser(c)
	assert.Error(t, err)

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)
}

func TestFindReportByCurrentUser_NoReport(t *testing.T) {
	// Mock UserRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	reportService := services.NewReportService(mockReportRepository, mockUserUtil)

	// Set up sample user for the test
	user := usermodel.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
	}

	// Set up expectations for the mock repository and util
	mockUserUtil.On("FindCurrentUser", mock.Anything).Return(&user, nil)
	mockReportRepository.On("FindByUserID", mock.Anything).Return([]models.Report{}, nil)

	// Create a test context
	c, _ := gin.CreateTestContext(nil)

	// Call the service
	responseReports, err := reportService.FindReportByCurrentUser(c)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(responseReports))

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)
}

func TestFindById(t *testing.T) {
	// Mock UserRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	reportService := services.NewReportService(mockReportRepository, mockUserUtil)

	// Set up sample report for the test
	report := models.Report{
		ID:           1,
		UserID:       1,
		AlertCount:   10,
		AnalysisTime: 1800,
		Predict:      "Good",
		Type:         "Study",
	}

	// Set up expectations for the mock repository
	mockReportRepository.On("FindById", mock.Anything).Return(report, nil)

	// Create a test context
	c, _ := gin.CreateTestContext(nil)

	// Call the service
	responseReport, err := reportService.FindById(c, 1)

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)

	// Check the results
	assert.NoError(t, err)
	assert.Equal(t, report.ID, responseReport.ID)
	assert.Equal(t, report.UserID, responseReport.UserID)
	assert.Equal(t, report.AlertCount, responseReport.AlertCount)
	assert.Equal(t, report.AnalysisTime, responseReport.AnalysisTime)
	assert.Equal(t, report.Predict, responseReport.Predict)
	assert.Equal(t, report.Type, responseReport.Type)
}

func TestFindById_NoReport(t *testing.T) {
	// Mock UserRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	reportService := services.NewReportService(mockReportRepository, mockUserUtil)
	assert.NotNil(t, reportService)

	// Set up expectations for the mock repository
	mockReportRepository.On("FindById", mock.Anything).Return(models.Report{}, errors.New("record not found"))

	// Create a test context
	c, _ := gin.CreateTestContext(nil)

	// Call the service
	_, err := reportService.FindById(c, 1)

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)

	// Check the results
	assert.Error(t, err)
}

func TestFindReportSummaryByMonth(t *testing.T) {
	// Mock UserRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	reportService := services.NewReportService(mockReportRepository, mockUserUtil)

	// Set up sample user for the test
	user := usermodel.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
	}

	// Set up sample reports for the test
	reports := []models.Report{
		{
			ID:           1,
			UserID:       1,
			AlertCount:   10,
			AnalysisTime: 1800,
			Predict:      "Good",
			Type:         "Study",
			CreatedAt:    time.Now(),
		},
		{
			ID:           2,
			UserID:       1,
			AlertCount:   30,
			AnalysisTime: 3600,
			Predict:      "Good",
			Type:         "Study",
			CreatedAt:    time.Now(),
		},
	}

	// Set up expectations for the mock repository and util
	mockUserUtil.On("FindCurrentUser", mock.Anything).Return(&user, nil)
	mockReportRepository.On("FindByYearAndMonth", mock.Anything, mock.Anything).Return(reports, nil)

	// Create a test context
	c, _ := gin.CreateTestContext(nil)

	// YearAndMonth
	yearAndMonth := time.Now().Format("200601")

	// Call the service
	responseReports, err := reportService.FindReportSummaryByMonth(c, yearAndMonth)
	assert.NoError(t, err)
	assert.Equal(t, len(reports), len(responseReports))

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)

	// Check the results
	for i, report := range reports {
		assert.Equal(t, report.ID, responseReports[i].ID)
		assert.Equal(t, report.CreatedAt, responseReports[i].CreatedAt)
	}
}

func TestFindReportSummaryByMonth_NoUser(t *testing.T) {
	// Mock UserRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	reportService := services.NewReportService(mockReportRepository, mockUserUtil)

	// Set up expectations for the mock repository and util
	mockUserUtil.On("FindCurrentUser", mock.Anything).Return((*usermodel.User)(nil), errors.New("record not found"))

	// Create a test context
	c, _ := gin.CreateTestContext(nil)

	// YearAndMonth
	yearAndMonth := time.Now().Format("200601")

	// Call the service
	_, err := reportService.FindReportSummaryByMonth(c, yearAndMonth)
	assert.Error(t, err)

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)
}

func TestFindReportSummaryByMonth_NoReport(t *testing.T) {
	// Mock UserRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	reportService := services.NewReportService(mockReportRepository, mockUserUtil)

	// Set up sample user for the test
	user := usermodel.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
	}

	// Set up expectations for the mock repository and util
	mockUserUtil.On("FindCurrentUser", mock.Anything).Return(&user, nil)
	mockReportRepository.On("FindByYearAndMonth", mock.Anything, mock.Anything).Return([]models.Report{}, nil)

	// Create a test context
	c, _ := gin.CreateTestContext(nil)

	// YearAndMonth
	yearAndMonth := time.Now().Format("200601")

	// Call the service
	responseReports, err := reportService.FindReportSummaryByMonth(c, yearAndMonth)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(responseReports))

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)
}

func TestFindAll(t *testing.T) {
	// Mock UserRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	reportService := services.NewReportService(mockReportRepository, mockUserUtil)

	// Set up sample reports for the test
	reports := []models.Report{
		{
			ID:           1,
			UserID:       1,
			AlertCount:   10,
			AnalysisTime: 1800,
			Predict:      "Good",
			Type:         "Study",
			CreatedAt:    time.Now(),
		},
		{
			ID:           2,
			UserID:       1,
			AlertCount:   30,
			AnalysisTime: 3600,
			Predict:      "Good",
			Type:         "Study",
			CreatedAt:    time.Now(),
		},
	}

	// Set up expectations for the mock repository and util
	mockReportRepository.On("FindAll").Return(reports, nil)

	// Call the service
	responseReports, err := reportService.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, len(reports), len(responseReports))

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)

	// Check the results
	for i, report := range reports {
		assert.Equal(t, report.ID, responseReports[i].ID)
		assert.Equal(t, report.CreatedAt, responseReports[i].CreatedAt)
	}
}

func TestFindAll_NoReport(t *testing.T) {
	// Mock UserRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	reportService := services.NewReportService(mockReportRepository, mockUserUtil)
	assert.NotNil(t, reportService)

	// Set up expectations for the mock repository
	mockReportRepository.On("FindAll").Return([]models.Report{}, nil)

	// Call the service
	responseReports, err := reportService.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(responseReports))

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)
}
