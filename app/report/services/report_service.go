package services

import (
	"gdsc/baro/app/report/models"
	"gdsc/baro/app/report/repositories"
	"gdsc/baro/app/report/types"
	"gdsc/baro/global/utils"
	"log"
	"os"
	"strings"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

type ReportService struct {
	ReportRepository *repositories.ReportRepository
	UserUtil         *utils.UserUtil
}

func NewReportService(reportRepository *repositories.ReportRepository, userUtil *utils.UserUtil) *ReportService {
	return &ReportService{
		ReportRepository: reportRepository,
		UserUtil:         userUtil,
	}
}

func (service *ReportService) Predict(c *gin.Context, input types.RequestPredict) (string, error) {
	user, err := service.UserUtil.FindCurrentUser(c)
	if err != nil {
		return "", err
	}

	REQUEST_URL := os.Getenv("AI_SERVER_API_URL")

	u, err := url.Parse(REQUEST_URL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("video_url", input.VideoURL)
	u.RawQuery = q.Encode()

	message := "Video submitted successfully"

	go func() {
		req, err := http.NewRequest("POST", u.String(), nil)
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
			UserID: user.ID,
		}

		// TODO: refactor, AI Server Change to return JSON
		var f func(*html.Node) string
		f = func(n *html.Node) string {
			if n.Type == html.ElementNode && n.Data == "p" {
				return n.FirstChild.Data
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				result := f(c)
				if strings.TrimSpace(result) != "" {
					return result
				}
			}
			return ""
		}

		report.Predict = f(doc)
		err = service.ReportRepository.Save(&report)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	return message, nil
}

func (service *ReportService) FindReportByCurrentUser(c *gin.Context) ([]types.ResponsePredict, error) {
	user, err := service.UserUtil.FindCurrentUser(c)
	if err != nil {
		return nil, err
	}

	reports, err := service.ReportRepository.FindByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	var responseReports []types.ResponsePredict
	for _, report := range reports {
		responseReport := types.ResponsePredict{
			ID:        report.ID,
			UserID:    report.UserID,
			Predict:   report.Predict,
			CreatedAt: report.CreatedAt,
		}
		responseReports = append(responseReports, responseReport)
	}

	return responseReports, nil
}

func (service *ReportService) FindAll() ([]types.ResponsePredict, error) {
	reports, err := service.ReportRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var responseReports []types.ResponsePredict
	for _, report := range reports {
		responseReport := types.ResponsePredict{
			ID:        report.ID,
			UserID:    report.UserID,
			Predict:   report.Predict,
			CreatedAt: report.CreatedAt,
		}
		responseReports = append(responseReports, responseReport)
	}

	return responseReports, nil
}
