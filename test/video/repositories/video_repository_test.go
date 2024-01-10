package repositories

import (
	"gdsc/baro/app/video/repositories"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestVideoRepository_FindAll(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create VideoRepository
	videoRepository := repositories.NewVideoRepository(gormDB)

	// Create sample videos for the test
	rows := sqlmock.NewRows([]string{"VideoID", "Title", "ThumbnailUrl"}).
		AddRow("1", "Video 1", "thumbnail1.jpg").
		AddRow("2", "Video 2", "thumbnail2.jpg")

	// Set up expectations for the SELECT query
	mock.ExpectQuery("SELECT \\* FROM `videos`").WillReturnRows(rows)

	// Call the method under test
	videos, err := videoRepository.FindAll()

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check the result
	assert.NoError(t, err)
	assert.NotNil(t, videos)
	assert.Equal(t, 2, len(videos))
	assert.Equal(t, "Video 1", videos[0].Title)
	assert.Equal(t, "thumbnail1.jpg", videos[0].ThumbnailUrl)
	assert.Equal(t, "Video 2", videos[1].Title)
	assert.Equal(t, "thumbnail2.jpg", videos[1].ThumbnailUrl)
}
