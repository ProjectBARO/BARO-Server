package services

import (
	"encoding/json"
	"fmt"
	"gdsc/baro/app/report/models"
	"gdsc/baro/app/report/repositories"
	"gdsc/baro/app/report/types"
	usermodel "gdsc/baro/app/user/models"
	"gdsc/baro/global/fcm"
	"gdsc/baro/global/utils"
	"io"
	"os"
	"time"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type ReportServiceInterface interface {
	Analysis(c *gin.Context, input types.RequestAnalysis) (string, error)
	FindReportByCurrentUser(c *gin.Context) ([]types.ResponseReport, error)
	FindById(c *gin.Context, id uint) (types.ResponseReport, error)
	FindReportSummaryByMonth(c *gin.Context, yearAndMonth string) ([]types.ResponseReportSummary, error)
	FindAll() ([]types.ResponseReport, error)
	FindRankAtAgeAndGender(c *gin.Context) (types.ResponseRank, error)
}

type ReportService struct {
	ReportRepository repositories.ReportRepositoryInterface
	UserUtil         utils.UserUtilInterface
}

func NewReportService(reportRepository repositories.ReportRepositoryInterface, userUtil utils.UserUtilInterface) *ReportService {
	return &ReportService{
		ReportRepository: reportRepository,
		UserUtil:         userUtil,
	}
}

var REQUEST_URL string
var client *http.Client

func init() {
	REQUEST_URL = os.Getenv("AI_SERVER_API_URL")
	client = &http.Client{}
}

func (service *ReportService) Analysis(c *gin.Context, input types.RequestAnalysis) (string, error) {
	user, err := service.UserUtil.FindCurrentUser(c)
	if err != nil {
		return "Not Found User", err
	}

	u, _ := url.Parse(REQUEST_URL)

	q := u.Query()
	q.Add("video_url", input.VideoURL)
	u.RawQuery = q.Encode()

	message := "Video submitted successfully"

	errCh := make(chan error, 1)
	go func() {
		errCh <- Predict(*service, u.String(), *user, input)
	}()

	defer func() {
		if err := <-errCh; err != nil {
			fmt.Println(err, " [", user.ID, ", ", input.VideoURL, "]")
		}
	}()

	return message, nil
}

func Predict(service ReportService, url string, user usermodel.User, input types.RequestAnalysis) error {
	response, err := HandleRequest(url)
	if err != nil {
		return err
	}

	result, scores, nomalRatio, statusFrequencies, distances, landmarksInfo := ParseAnalysis(&response)
	score := CalculateScores(result, scores)

	report := models.Report{
		UserID:            user.ID,
		AlertCount:        input.AlertCount,
		AnalysisTime:      input.AnalysisTime,
		Type:              input.Type,
		Predict:           fmt.Sprintf("%v", result),
		Score:             score,
		NormalRatio:       nomalRatio,
		StatusFrequencies: statusFrequencies,
		Distances:         distances,
		NeckAngles:        landmarksInfo,
	}

	savedReport, _ := service.ReportRepository.Save(&report)

	title, body, _ := GenerateMessage(savedReport.CreatedAt.String())

	err = fcm.SendPushNotification(user.FcmToken, title, body)
	if err != nil {
		return err
	}

	return nil
}

func HandleRequest(url string) (types.ResponseAnalysis, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.ResponseAnalysis{}, err
	}

	response, err := client.Do(req)
	if err != nil {
		return types.ResponseAnalysis{}, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return types.ResponseAnalysis{}, err
	}

	var data types.ResponseAnalysis
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println("응답 파싱 에러:", err)
		return types.ResponseAnalysis{}, err
	}

	return data, nil
}

func ParseAnalysis(response *types.ResponseAnalysis) ([]int, []float64, string, string, string, string) {
	// 결과 및 신뢰도
	result := response.Result
	scores := response.Scores

	// 정상 비율
	nomalRatio := response.NormalRatio

	// 빈도수
	statusFrequencies := []int{
		response.StatusFrequencies["Fine"],
		response.StatusFrequencies["Danger"],
		response.StatusFrequencies["Serious"],
		response.StatusFrequencies["Very Serious"],
	}

	// 없는 필드는 0으로 초기화
	for i := range statusFrequencies {
		if statusFrequencies[i] == 0 {
			statusFrequencies[i] = 0
		}
	}

	// 길이 및 각도
	var distances []float64
	var angles []float64
	for i := range response.LandmarksInfo {
		// response.LandmarksInfo[i][2] = 길이
		// response.LandmarksInfo[i][3] = 각도
		distances = append(distances, response.LandmarksInfo[i][2].(float64))
		angles = append(angles, response.LandmarksInfo[i][3].(float64))
	}

	return result, scores, fmt.Sprintf("%.3f", nomalRatio), fmt.Sprintf("%v", statusFrequencies), fmt.Sprintf("%.3f", distances), fmt.Sprintf("%.3f", angles)
}

func CalculateScores(result []int, scores []float64) string {
	normalCases := make([]float64, 0, len(result))
	abnormalCases := make([]float64, 0, len(result))

	for i, r := range result {
		if r == 1 {
			normalCases = append(normalCases, scores[i])
		} else {
			abnormalCases = append(abnormalCases, scores[i])
		}
	}

	totalCases := len(normalCases) + len(abnormalCases)
	caseScore := 100.0 / float64(totalCases)

	var totalScore float64

	for _, score := range normalCases {
		switch {
		case score >= 99.8:
			totalScore += 1.0 * caseScore
		case score >= 99:
			totalScore += 0.97 * caseScore
		case score >= 96:
			totalScore += 0.94 * caseScore
		case score >= 93:
			totalScore += 0.9 * caseScore
		default:
			totalScore += 0.86 * caseScore
		}
	}

	for _, score := range abnormalCases {
		switch {
		case score >= 99.5:
			totalScore += 0.18 * caseScore
		case score >= 98:
			totalScore += 0.23 * caseScore
		case score >= 95:
			totalScore += 0.3 * caseScore
		case score >= 90:
			totalScore += 0.35 * caseScore
		default:
			totalScore += 0.43 * caseScore
		}

	}

	return fmt.Sprintf("%.2f", totalScore)
}

func GenerateMessage(date string) (string, string, error) {
	timeFormats := []string{
		"2006-01-02 15:04:05.000 -0700 MST",
		"2006-01-02 15:04:05.00 -0700 MST",
		"2006-01-02 15:04:05.0 -0700 MST",
	}

	var t time.Time
	var err error

	for _, format := range timeFormats {
		t, err = time.Parse(format, date)
		if err == nil {
			break
		}
	}

	if err != nil {
		return "", "", err
	}

	title := "자세 분석이 완료되었어요!"
	body := fmt.Sprintf("%d년 %d월 %d일에 측정한 보고서가 도착했습니다!", t.Year(), t.Month(), t.Day())

	return title, body, nil
}

func (service *ReportService) FindReportByCurrentUser(c *gin.Context) ([]types.ResponseReport, error) {
	user, err := service.UserUtil.FindCurrentUser(c)
	if err != nil {
		return nil, err
	}

	reports, err := service.ReportRepository.FindByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	var responseReports []types.ResponseReport
	for _, report := range reports {
		responseReport := types.ResponseReport{
			ID:                report.ID,
			UserID:            report.UserID,
			AlertCount:        report.AlertCount,
			AnalysisTime:      report.AnalysisTime,
			Type:              report.Type,
			Predict:           report.Predict,
			Score:             report.Score,
			NormalRatio:       report.NormalRatio,
			NeckAngles:        report.NeckAngles,
			Distances:         report.Distances,
			StatusFrequencies: report.StatusFrequencies,
			CreatedAt:         report.CreatedAt,
		}
		responseReports = append(responseReports, responseReport)
	}

	return responseReports, nil
}

func (service *ReportService) FindById(c *gin.Context, id uint) (types.ResponseReport, error) {
	report, err := service.ReportRepository.FindById(id)
	if err != nil {
		return types.ResponseReport{}, err
	}

	responseReport := types.ResponseReport{
		ID:                report.ID,
		UserID:            report.UserID,
		AlertCount:        report.AlertCount,
		AnalysisTime:      report.AnalysisTime,
		Type:              report.Type,
		Predict:           report.Predict,
		Score:             report.Score,
		NormalRatio:       report.NormalRatio,
		NeckAngles:        report.NeckAngles,
		Distances:         report.Distances,
		StatusFrequencies: report.StatusFrequencies,
		CreatedAt:         report.CreatedAt,
	}

	return responseReport, nil
}

func (service *ReportService) FindReportSummaryByMonth(c *gin.Context, yearAndMonth string) ([]types.ResponseReportSummary, error) {
	user, err := service.UserUtil.FindCurrentUser(c)
	if err != nil {
		return nil, err
	}

	reports, _ := service.ReportRepository.FindByYearAndMonth(user.ID, yearAndMonth)

	var responseReports []types.ResponseReportSummary
	for _, report := range reports {
		responseReport := types.ResponseReportSummary{
			ID:        report.ID,
			CreatedAt: report.CreatedAt,
		}
		responseReports = append(responseReports, responseReport)
	}

	return responseReports, nil
}

func (service *ReportService) FindAll() ([]types.ResponseReport, error) {
	reports, _ := service.ReportRepository.FindAll()

	var responseReports []types.ResponseReport
	for _, report := range reports {
		responseReport := types.ResponseReport{
			ID:                report.ID,
			UserID:            report.UserID,
			AlertCount:        report.AlertCount,
			AnalysisTime:      report.AnalysisTime,
			Type:              report.Type,
			Predict:           report.Predict,
			Score:             report.Score,
			NormalRatio:       report.NormalRatio,
			NeckAngles:        report.NeckAngles,
			Distances:         report.Distances,
			StatusFrequencies: report.StatusFrequencies,
			CreatedAt:         report.CreatedAt,
		}
		responseReports = append(responseReports, responseReport)
	}

	return responseReports, nil
}

func (service *ReportService) FindRankAtAgeAndGender(c *gin.Context) (types.ResponseRank, error) {
	user, err := service.UserUtil.FindCurrentUser(c)
	if err != nil {
		return types.ResponseRank{}, err
	}

	rank, _ := service.ReportRepository.FindRankAtAgeAndGender(user, time.Now().AddDate(0, 0, -30), time.Now())

	responseRank := types.ResponseRank{
		UserID:       rank.UserID,
		Nickname:     user.Nickname,
		Age:          rank.Age,
		Gender:       rank.Gender,
		NormalRatio:  rank.NormalRatio,
		AverageScore: rank.AverageScore,
	}

	return responseRank, nil
}
