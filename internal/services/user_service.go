package services

import (
	"errors"

	"starter-kit-restapi-gofiber/internal/dto"
	"starter-kit-restapi-gofiber/internal/models"
	"starter-kit-restapi-gofiber/pkg/utils"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*models.User, error) {
	if s.IsEmailTaken(req.Email) {
		return nil, errors.New("email already taken")
	}

	hashedPwd, _ := utils.HashPassword(req.Password)
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPwd,
		Role:     req.Role,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) QueryUsers(pagination *utils.Pagination, name string, role string) (*utils.Pagination, error) {
	var users []models.User
	
	db := s.DB.Model(&models.User{})
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if role != "" {
		db = db.Where("role = ?", role)
	}

	db.Scopes(utils.Paginate(users, pagination, db)).Find(&users)
	
	pagination.Rows = users
	return pagination, nil
}

func (s *UserService) GetUserById(id string) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUserById(userId string, req *dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	if req.Email != "" && req.Email != user.Email && s.IsEmailTaken(req.Email) {
		return nil, errors.New("email already taken")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
		user.IsEmailVerified = false
	}
	if req.Password != "" {
		hash, _ := utils.HashPassword(req.Password)
		user.Password = hash
	}

	s.DB.Save(&user)
	return user, nil
}

func (s *UserService) DeleteUserById(userId string) error {
	user, err := s.GetUserById(userId)
	if err != nil {
		return err
	}
	return s.DB.Delete(&user).Error
}

func (s *UserService) IsEmailTaken(email string) bool {
	var count int64
	s.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}