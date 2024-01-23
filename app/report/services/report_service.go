package services

import (
	"gdsc/baro/app/report/models"
	"gdsc/baro/app/report/repositories"
	"gdsc/baro/app/report/types"
	usermodel "gdsc/baro/app/user/models"
	"gdsc/baro/global/utils"
	"log"
	"os"
	"strings"

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
	HandleRequest(url string, user usermodel.User, input types.RequestAnalysis)
	ParseHTML(n *html.Node) string
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

	go service.HandleRequest(u.String(), *user, input)

	return message, nil
}

func (service *ReportService) HandleRequest(url string, user usermodel.User, input types.RequestAnalysis) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	report := models.Report{
		UserID:       user.ID,
		AlertCount:   input.AlertCount,
		AnalysisTime: input.AnalysisTime,
		Type:         input.Type,
	}

	report.Predict = service.ParseHTML(doc)

	_, err = service.ReportRepository.Save(&report)
	if err != nil {
		log.Println(err)
		return
	}
}

func (service *ReportService) ParseHTML(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "p" {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := service.ParseHTML(c)
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

	reports, err := service.ReportRepository.FindByYearAndMonth(user.ID, yearAndMonth)
	if err != nil {
		return nil, err
	}

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
	reports, err := service.ReportRepository.FindAll()
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
