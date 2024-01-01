package repositories

import (
	"gdsc/baro/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (repo *UserRepository) Create(user *models.User) (models.User, error) {
	if err := repo.DB.Create(user).Error; err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func (repo *UserRepository) FindByID(id uint) (models.User, error) {
	var user models.User
	result := repo.DB.Where("id = ?", id).First(&user)
	return user, result.Error
}

func (repo *UserRepository) FindOrCreateByEmail(user *models.User) (*models.User, error) {
	err := repo.DB.Where(models.User{Email: user.Email}).Attrs(user).FirstOrCreate(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := repo.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) Update(user *models.User) (models.User, error) {
	if err := repo.DB.Save(user).Error; err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func (repo *UserRepository) Delete(user *models.User) error {
	return repo.DB.Delete(user).Error
}
