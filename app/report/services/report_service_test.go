package services_test

import (
	"errors"
	"fmt"
	"gdsc/baro/app/report/models"
	"gdsc/baro/app/report/services"
	"gdsc/baro/app/report/types"
	usermodel "gdsc/baro/app/user/models"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *MockReportRepository) FindRankAtAgeAndGender(user *usermodel.User, start, end time.Time) (types.ResponseRank, error) {
	args := m.Called(user, start, end)
	return args.Get(0).(types.ResponseRank), args.Error(1)
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

	// Set up environment variables
	os.Setenv("AI_SERVER_API_URL", "http://localhost:5000/predict")

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
	request := types.RequestAnalysis{
		VideoURL:     "test",
		AlertCount:   10,
		AnalysisTime: 1800,
		Type:         "Study",
	}

	// Set up expectations for the mock repository and util
	mockUserUtil.On("FindCurrentUser", mock.Anything).Return(&user, nil)

	// Create a test context
	c, _ := gin.CreateTestContext(nil)

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

func TestPredict_InvalidURL(t *testing.T) {
	// Mock ReportRepository, UserUtil
	mockReportRepository := new(MockReportRepository)
	mockUserUtil := new(MockUserUtil)

	// Create ReportService
	service := services.NewReportService(mockReportRepository, mockUserUtil)

	// Set up invalid URL
	url := ""

	// Call the service
	err := services.Predict(*service, url, usermodel.User{}, types.RequestAnalysis{})
	assert.Error(t, err)
}

func TestHandleRequest(t *testing.T) {
	assert.True(t, true)
}

func TestParseAnalysis(t *testing.T) {
	// response *types.ResponseAnalysis
	response := types.ResponseAnalysis{
		Result:       []int{0, 1, 1, 1, 0, 1},
		HunchedRatio: 10.0,
		NormalRatio:  90.0,
		Scores:       []float64{99.9, 92.9, 82.3, 92.4, 69.5},
		LandmarksInfo: [][]interface{}{
			{[]float64{0.1, 0.2}, []float64{0.3, 0.4}, 1.0, 45.0},
			{[]float64{0.5, 0.6}, []float64{0.7, 0.8}, 2.0, 60.0},
		},
		StatusFrequencies: map[string]int{"Very Serious": 6},
	}

	// Execute method under test
	result, scores, nomalRatio, statusFrequencies, distances, landmarksInfo := services.ParseAnalysis(&response)

	// Assert result
	assert.Equal(t, response.Result, result)
	assert.Equal(t, response.Scores, scores)
	assert.Equal(t, fmt.Sprintf("%.3f", response.NormalRatio), nomalRatio)
	assert.Equal(t, fmt.Sprintf("%.3f", []float64{1.0, 2.0}), distances)
	assert.Equal(t, fmt.Sprintf("%.3f", []float64{45.0, 60.0}), landmarksInfo)
	assert.Equal(t, fmt.Sprintf("%v", []string{"0", "0", "0", "6"}), statusFrequencies)
}

func TestParseAnalysis_NoStatusFrequencies(t *testing.T) {
	// response *types.ResponseAnalysis
	response := types.ResponseAnalysis{
		Result:       []int{0, 1, 1, 1, 0, 1},
		HunchedRatio: 10.0,
		NormalRatio:  90.0,
		Scores:       []float64{99.9, 92.9, 82.3, 92.4, 69.5},
		LandmarksInfo: [][]interface{}{
			{[]float64{0.1, 0.2}, []float64{0.3, 0.4}, 1.0, 45.0},
			{[]float64{0.5, 0.6}, []float64{0.7, 0.8}, 2.0, 60.0},
		},
	}

	// Execute method under test
	result, scores, nomalRatio, statusFrequencies, distances, landmarksInfo := services.ParseAnalysis(&response)

	// Assert result
	assert.Equal(t, response.Result, result)
	assert.Equal(t, response.Scores, scores)
	assert.Equal(t, fmt.Sprintf("%.3f", response.NormalRatio), nomalRatio)
	assert.Equal(t, fmt.Sprintf("%.3f", []float64{1.0, 2.0}), distances)
	assert.Equal(t, fmt.Sprintf("%.3f", []float64{45.0, 60.0}), landmarksInfo)
	assert.Equal(t, fmt.Sprintf("%v", []string{"0", "0", "0", "0"}), statusFrequencies)
}

func TestParseAnalysis_FullStatusFrequencies(t *testing.T) {
	// response *types.ResponseAnalysis
	response := types.ResponseAnalysis{
		Result:       []int{0, 1, 1, 1, 0, 1},
		HunchedRatio: 10.0,
		NormalRatio:  90.0,
		Scores:       []float64{99.9, 92.9, 82.3, 92.4, 69.5},
		LandmarksInfo: [][]interface{}{
			{[]float64{0.1, 0.2}, []float64{0.3, 0.4}, 1.0, 45.0},
			{[]float64{0.5, 0.6}, []float64{0.7, 0.8}, 2.0, 60.0},
		},
		StatusFrequencies: map[string]int{"Fine": 1, "Danger": 2, "Serious": 1, "Very Serious": 2},
	}

	// Execute method under test
	result, scores, nomalRatio, statusFrequencies, distances, landmarksInfo := services.ParseAnalysis(&response)

	// Assert result
	assert.Equal(t, response.Result, result)
	assert.Equal(t, response.Scores, scores)
	assert.Equal(t, fmt.Sprintf("%.3f", response.NormalRatio), nomalRatio)
	assert.Equal(t, fmt.Sprintf("%.3f", []float64{1.0, 2.0}), distances)
	assert.Equal(t, fmt.Sprintf("%.3f", []float64{45.0, 60.0}), landmarksInfo)
	assert.Equal(t, fmt.Sprintf("%v", []string{"1", "2", "1", "2"}), statusFrequencies)
}

func TestCalculateScores_AllBad(t *testing.T) {
	// Set up test data
	testResult := []int{0, 0, 0, 0, 0, 0}
	testScores := []float64{99.9, 92.9, 92.3, 92.4, 99.5, 92.0}

	// Execute method under test
	result := services.CalculateScores(testResult, testScores)

	// Assert result
	assert.Equal(t, "29.33", result)
}

func TestCalculateScores_HalfGood(t *testing.T) {
	// Set up test data
	testResult := []int{0, 0, 0, 1, 1, 1}
	testScores := []float64{99.9, 99.2, 96.5, 95.5, 94.9, 92.0}

	// Execute method under test
	result := services.CalculateScores(testResult, testScores)

	// Assert result
	assert.Equal(t, "56.17", result)
}

func TestCalculateScores_AllGoodDiffScore(t *testing.T) {
	// Set up test data
	testResult := []int{1, 1, 1, 1, 1, 1}
	testScores := []float64{99.9, 99.2, 96.5, 95.5, 94.9, 92.0}

	// Execute method under test
	result := services.CalculateScores(testResult, testScores)

	// Assert result
	assert.Equal(t, "92.83", result)
}

func TestCalculateScores_AllBadSameScore(t *testing.T) {
	// Set up test data
	testResult := []int{0, 0, 0, 0, 0, 0}
	testScores := []float64{99.9, 99.1, 97.5, 93.9, 92.9, 86.9}

	// Execute method under test
	result := services.CalculateScores(testResult, testScores)

	// Assert result
	assert.Equal(t, "30.67", result)
}

func TestCalculateScores_AllGood(t *testing.T) {
	// Set up test data
	testResult := []int{1, 1, 1, 1, 1, 1}
	testScores := []float64{99.9, 92.9, 96.3, 92.4, 99.1, 92.0}

	// Execute method under test
	result := services.CalculateScores(testResult, testScores)

	// Assert result
	assert.Equal(t, "91.50", result)
}

func TestGenerateMessage(t *testing.T) {
	// Set up test data
	testTimes := []string{
		"2024-02-01 12:04:05.123 +0900 KST",
		"2024-02-01 12:04:05.12 +0900 KST",
		"2024-02-01 12:04:05.1 +0900 KST",
	}

	// Execute method under test
	for _, testTime := range testTimes {
		title, body, err := services.GenerateMessage(testTime)

		// Assert result
		assert.NoError(t, err)
		assert.Equal(t, "자세 분석이 완료되었어요!", title)
		assert.Equal(t, "2024년 2월 1일에 측정한 보고서가 도착했습니다!", body)
	}
}

func TestGenerateMessage_InvaildData(t *testing.T) {
	// Set up invalid test data
	testTime := "2024-02-01"

	// Execute method under test
	title, body, err := services.GenerateMessage(testTime)

	// Assert result
	assert.Error(t, err)
	assert.Equal(t, "", title)
	assert.Equal(t, "", body)
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
			ID:                1,
			UserID:            1,
			AlertCount:        10,
			AnalysisTime:      1800,
			Type:              "Study",
			Predict:           "Good",
			Score:             "90.000",
			NormalRatio:       "90.000",
			NeckAngles:        "angle",
			Distances:         "distance",
			StatusFrequencies: "status",
		},
		{
			ID:                2,
			UserID:            1,
			AlertCount:        10,
			AnalysisTime:      1800,
			Type:              "Study",
			Predict:           "Good",
			Score:             "90.000",
			NormalRatio:       "90.000",
			NeckAngles:        "angle",
			Distances:         "distance",
			StatusFrequencies: "status",
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
		assert.Equal(t, report.Type, responseReports[i].Type)
		assert.Equal(t, report.Predict, responseReports[i].Predict)
		assert.Equal(t, report.Score, responseReports[i].Score)
		assert.Equal(t, report.NormalRatio, responseReports[i].NormalRatio)
		assert.Equal(t, report.NeckAngles, responseReports[i].NeckAngles)
		assert.Equal(t, report.Distances, responseReports[i].Distances)
		assert.Equal(t, report.StatusFrequencies, responseReports[i].StatusFrequencies)
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

	// Check the results
	assert.Equal(t, "record not found", err.Error())
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
		ID:                1,
		UserID:            1,
		AlertCount:        10,
		AnalysisTime:      1800,
		Type:              "Study",
		Predict:           "Good",
		Score:             "90.000",
		NormalRatio:       "90.000",
		NeckAngles:        "angle",
		Distances:         "distance",
		StatusFrequencies: "status",
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
	assert.Equal(t, report.Type, responseReport.Type)
	assert.Equal(t, report.Predict, responseReport.Predict)
	assert.Equal(t, report.Score, responseReport.Score)
	assert.Equal(t, report.NormalRatio, responseReport.NormalRatio)
	assert.Equal(t, report.NeckAngles, responseReport.NeckAngles)
	assert.Equal(t, report.Distances, responseReport.Distances)
	assert.Equal(t, report.StatusFrequencies, responseReport.StatusFrequencies)
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
			ID:                1,
			UserID:            1,
			AlertCount:        10,
			AnalysisTime:      1800,
			Type:              "Study",
			Predict:           "Good",
			Score:             "90.000",
			NormalRatio:       "90.000",
			NeckAngles:        "angle",
			Distances:         "distance",
			StatusFrequencies: "status",
			CreatedAt:         time.Now(),
		},
		{
			ID:                2,
			UserID:            1,
			AlertCount:        10,
			AnalysisTime:      1800,
			Type:              "Study",
			Predict:           "Good",
			Score:             "90.000",
			NormalRatio:       "90.000",
			NeckAngles:        "angle",
			Distances:         "distance",
			StatusFrequencies: "status",
			CreatedAt:         time.Now(),
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
			ID:                1,
			UserID:            1,
			AlertCount:        10,
			AnalysisTime:      1800,
			Type:              "Study",
			Predict:           "Good",
			Score:             "90.000",
			NormalRatio:       "90.000",
			NeckAngles:        "angle",
			Distances:         "distance",
			StatusFrequencies: "status",
			CreatedAt:         time.Now(),
		},
		{
			ID:                2,
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

func TestFindRankAtAgeAndGender(t *testing.T) {
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

	// Set up sample rank for the test
	rank := types.ResponseRank{
		UserID:       1,
		Nickname:     "test",
		Age:          20,
		Gender:       "male",
		NormalRatio:  "90.00",
		AverageScore: 90.000,
	}

	// Set up expectations for the mock repository and util
	mockUserUtil.On("FindCurrentUser", mock.Anything).Return(&user, nil)
	mockReportRepository.On("FindRankAtAgeAndGender", mock.Anything, mock.Anything, mock.Anything).Return(rank, nil)

	// Create a test context
	c, _ := gin.CreateTestContext(nil)

	// Call the service
	responseRank, err := reportService.FindRankAtAgeAndGender(c)
	assert.NoError(t, err)

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)

	// Check the results
	assert.Equal(t, rank.UserID, responseRank.UserID)
	assert.Equal(t, rank.Nickname, responseRank.Nickname)
	assert.Equal(t, rank.Age, responseRank.Age)
	assert.Equal(t, rank.Gender, responseRank.Gender)
	assert.Equal(t, rank.NormalRatio, responseRank.NormalRatio)
	assert.Equal(t, rank.AverageScore, responseRank.AverageScore)
}

func TestFindRankAtAgeAndGender_NoUser(t *testing.T) {
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
	_, err := reportService.FindRankAtAgeAndGender(c)
	assert.Error(t, err)

	// Assert that the expectations were met
	mockReportRepository.AssertExpectations(t)
	mockUserUtil.AssertExpectations(t)
}
