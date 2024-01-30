package repositories_test

import (
	"gdsc/baro/app/user/models"
	"gdsc/baro/app/user/repositories"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUserRepository_Create(t *testing.T) {
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

	// Create UserRepository
	userRepository := repositories.NewUserRepository(gormDB)

	// Create sample user for the test
	user := &models.User{
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
		FcmToken: "test_token",
		Deleted:  gorm.DeletedAt{},
	}

	// Set up expectations for the mock DB to return the sample user
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs(user.Name, user.Nickname, user.Email, user.Age, user.Gender, user.FcmToken, user.Deleted).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method under test
	createdUser, err := userRepository.Create(user)
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check the result
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Age, createdUser.Age)
	assert.Equal(t, user.Gender, createdUser.Gender)
	assert.Equal(t, user.FcmToken, createdUser.FcmToken)
}

func TestUserRepository_FindByID(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create UserRepository
	userRepository := repositories.NewUserRepository(gormDB)

	// Create sample user for the test
	user := &models.User{
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
		FcmToken: "test_token",
		Deleted:  gorm.DeletedAt{},
	}

	// Set up expectations for the mock DB to return the sample user
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs(user.Name, user.Nickname, user.Email, user.Age, user.Gender, user.FcmToken, user.Deleted).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method under test
	createdUser, err := userRepository.Create(user)
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	// Set up expectations for the mock DB to return the user by ID
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE id = \\? AND `users`.`deleted` IS NULL ORDER BY `users`.`id` LIMIT 1").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"ID", "Name", "Nickname", "Email", "Age", "Gender", "FcmToken", "Deleted"}).
			AddRow(createdUser.ID, createdUser.Name, createdUser.Nickname, createdUser.Email, createdUser.Age, createdUser.Gender, createdUser.FcmToken, createdUser.Deleted.Time))

	// Find the user by ID
	foundUser, err := userRepository.FindByID(strconv.Itoa(int(createdUser.ID)))
	if err != nil {
		t.Fatalf("Error finding user by ID: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check the result
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, createdUser.ID, foundUser.ID)
	assert.Equal(t, user.Name, foundUser.Name)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.Age, foundUser.Age)
	assert.Equal(t, user.Gender, foundUser.Gender)
	assert.Equal(t, user.FcmToken, foundUser.FcmToken)
}

func TestUserRepository_FindOrCreateByEmail_First_Login(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create UserRepository
	userRepository := repositories.NewUserRepository(gormDB)

	// Create sample user for the test
	user := &models.User{
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
		FcmToken: "test_token",
		Deleted:  gorm.DeletedAt{},
	}

	// Set up expectations for the mock DB to create the user
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`email` = \\? AND `users`.`deleted` IS NULL ORDER BY `users`.`id` LIMIT 1").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "nickname", "email", "age", "gender", "fcm_token", "deleted"}))

	// Set up expectations for the mock DB to create the user
	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO `users`").
		WithArgs(user.Name, user.Nickname, user.Email, user.Age, user.Gender, user.FcmToken, user.Deleted).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// Call the method under test
	createdUser, err := userRepository.FindOrCreateByEmail(user)
	if err != nil {
		t.Fatalf("Error finding or creating user: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check the result
	assert.NotNil(t, createdUser)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Age, createdUser.Age)
	assert.Equal(t, user.Gender, createdUser.Gender)
	assert.Equal(t, user.FcmToken, createdUser.FcmToken)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create UserRepository
	userRepository := repositories.NewUserRepository(gormDB)

	// Create sample user for the test
	email := "test@gmail.com"
	expectedUser := &models.User{
		Name:     "test",
		Nickname: "test",
		Email:    email,
		Age:      20,
		Gender:   "male",
		FcmToken: "test_token",
		Deleted:  gorm.DeletedAt{},
	}

	// Set up expectations for the mock DB to find the user by email
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\?").WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "nickname", "email", "age", "gender", "fcm_token", "deleted"}).
			AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Nickname, expectedUser.Email, expectedUser.Age, expectedUser.Gender, expectedUser.FcmToken, expectedUser.Deleted))

	// Call the method under test
	resultUser, err := userRepository.FindByEmail(email)
	if err != nil {
		t.Fatalf("Error finding user by email: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check the result
	assert.NotNil(t, resultUser)
	assert.Equal(t, expectedUser.ID, resultUser.ID)
	assert.Equal(t, expectedUser.Name, resultUser.Name)
	assert.Equal(t, expectedUser.Email, resultUser.Email)
	assert.Equal(t, expectedUser.Age, resultUser.Age)
	assert.Equal(t, expectedUser.Gender, resultUser.Gender)
	assert.Equal(t, expectedUser.FcmToken, resultUser.FcmToken)
}

func TestUserRepository_Update(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create UserRepository
	userRepository := repositories.NewUserRepository(gormDB)

	// Create sample user for the test
	user := &models.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
		FcmToken: "test_token",
		Deleted:  gorm.DeletedAt{},
	}

	// Set up expectations for the mock DB to update the user within a transaction
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users`").
		WithArgs(user.Name, user.Nickname, user.Email, user.Age, user.Gender, user.FcmToken, user.Deleted, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method under test
	updatedUser, err := userRepository.Update(user)
	if err != nil {
		t.Fatalf("Error updating user: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check the result
	assert.Equal(t, user.ID, updatedUser.ID)
	assert.Equal(t, user.Name, updatedUser.Name)
	assert.Equal(t, user.Email, updatedUser.Email)
	assert.Equal(t, user.Age, updatedUser.Age)
	assert.Equal(t, user.Gender, updatedUser.Gender)
}

func TestUserRepository_Delete(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Set up expectations for the mock DB (ex: SELECT VERSION())
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION"}).AddRow("8.0.0"))

	// Create gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating gorm.DB: %v", err)
	}

	// Create UserRepository
	userRepository := repositories.NewUserRepository(gormDB)

	// Create sample user for the test
	user := &models.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
		FcmToken: "test_token",
		Deleted:  gorm.DeletedAt{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	// Set up expectations for the mock DB to update the user's Deleted field
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users`").
		WithArgs(sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method under test
	err = userRepository.Delete(user)
	if err != nil {
		t.Fatalf("Error deleting user: %v", err)
	}

	// Check that the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
