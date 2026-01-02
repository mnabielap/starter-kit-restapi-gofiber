package services

import (
	"errors"
	"strings"

	"starter-kit-restapi-gofiber/internal/dto"
	"starter-kit-restapi-gofiber/internal/models"
	"starter-kit-restapi-gofiber/pkg/utils"

	"github.com/google/uuid"
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

func (s *UserService) QueryUsers(params dto.UserQueryParams) (*utils.Pagination, error) {
	var users []models.User
	db := s.DB.Model(&models.User{})

	// 1. Filter by Role
	if params.Role != "" {
		db = db.Where("role = ?", params.Role)
	}

	// 2. Search Logic
	if params.Search != "" {
		searchStr := "%" + strings.ToLower(params.Search) + "%"
		
		// Create a sub-query for the grouping to avoid "unsupported type func" error
		// We use s.DB.Session(&gorm.Session{}) to ensure a clean state builder
		subQuery := s.DB.Session(&gorm.Session{NewDB: true})
		var searchGroup *gorm.DB

		switch params.Scope {
		case "name":
			searchGroup = subQuery.Where("LOWER(name) LIKE ?", searchStr)
		case "email":
			searchGroup = subQuery.Where("LOWER(email) LIKE ?", searchStr)
		case "id":
			if _, err := uuid.Parse(params.Search); err == nil {
				searchGroup = subQuery.Where("id = ?", params.Search)
			} else {
				// Impossible condition if ID is invalid
				searchGroup = subQuery.Where("1 = 0")
			}
		case "all":
			fallthrough
		default:
			// Equivalent to: AND (name LIKE ? OR email LIKE ? [OR id = ?])
			searchGroup = subQuery.Where("LOWER(name) LIKE ?", searchStr).
				Or("LOWER(email) LIKE ?", searchStr)
			
			if _, err := uuid.Parse(params.Search); err == nil {
				searchGroup = searchGroup.Or("id = ?", params.Search)
			}
		}

		// Apply the grouped condition
		db = db.Where(searchGroup)
	}

	// 3. Prepare Pagination & Sorting
	pagination := &utils.Pagination{
		Limit: params.Limit,
		Page:  params.Page,
		Sort:  params.SortBy,
	}

	// 4. Execute Query with Pagination Scope
	err := db.Scopes(utils.Paginate(&users, pagination, db)).Find(&users).Error
	if err != nil {
		return nil, err
	}
	
	// 5. Explicitly assign the filled slice to Results
	pagination.Results = users
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

	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}
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