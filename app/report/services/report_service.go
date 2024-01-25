package services

import (
	"fmt"
	"gdsc/baro/app/report/models"
	"gdsc/baro/app/report/repositories"
	"gdsc/baro/app/report/types"
	usermodel "gdsc/baro/app/user/models"
	"gdsc/baro/global/fcm"
	"gdsc/baro/global/utils"
	"os"
	"strings"
	"time"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

type ReportServiceInterface interface {
	Analysis(c *gin.Context, input types.RequestAnalysis) (string, error)
	FindReportByCurrentUser(c *gin.Context) ([]types.ResponseReport, error)
	FindById(c *gin.Context, id uint) (types.ResponseReport, error)
	FindReportSummaryByMonth(c *gin.Context, yearAndMonth string) ([]types.ResponseReportSummary, error)
	FindAll() ([]types.ResponseReport, error)
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

func (service *ReportService) Analysis(c *gin.Context, input types.RequestAnalysis) (string, error) {
	user, err := service.UserUtil.FindCurrentUser(c)
	if err != nil {
		return "Not Found User", err
	}

	REQUEST_URL := os.Getenv("AI_SERVER_API_URL")

	u, _ := url.Parse(REQUEST_URL)

	q := u.Query()
	q.Add("video_url", input.VideoURL)
	u.RawQuery = q.Encode()

	message := "Video submitted successfully"

	go Predict(*service, u.String(), *user, input)

	return message, nil
}

func Predict(service ReportService, url string, user usermodel.User, input types.RequestAnalysis) error {
	response, err := HandleRequest(url)
	if err != nil {
		return err
	}

	report := models.Report{
		UserID:       user.ID,
		AlertCount:   input.AlertCount,
		AnalysisTime: input.AnalysisTime,
		Type:         input.Type,
		Predict:      response,
	}

	savedReport, _ := service.ReportRepository.Save(&report)

	title, body, _ := GenerateMessage(savedReport.CreatedAt.String())

	err = fcm.SendPushNotification(user.FcmToken, title, body)
	if err != nil {
		return err
	}

	return nil
}

func HandleRequest(url string) (string, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		return "", err
	}

	return ParseHTML(doc), nil
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

func ParseHTML(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "p" {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := ParseHTML(c)
		if strings.TrimSpace(result) != "" {
			return result
		}
	}
	return ""
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
			ID:           report.ID,
			UserID:       report.UserID,
			AlertCount:   report.AlertCount,
			AnalysisTime: report.AnalysisTime,
			Predict:      report.Predict,
			Type:         report.Type,
			CreatedAt:    report.CreatedAt,
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
		ID:           report.ID,
		UserID:       report.UserID,
		AlertCount:   report.AlertCount,
		AnalysisTime: report.AnalysisTime,
		Predict:      report.Predict,
		Type:         report.Type,
		CreatedAt:    report.CreatedAt,
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
			ID:           report.ID,
			UserID:       report.UserID,
			AlertCount:   report.AlertCount,
			AnalysisTime: report.AnalysisTime,
			Predict:      report.Predict,
			Type:         report.Type,
			CreatedAt:    report.CreatedAt,
		}
		responseReports = append(responseReports, responseReport)
	}

	return responseReports, nil
}
