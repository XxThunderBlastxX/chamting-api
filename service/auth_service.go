package service

import (
	"github.com/XxThunderBlastxX/chamting-api/models"
	"github.com/XxThunderBlastxX/chamting-api/repository"
)

//AuthService is an interface from which our api module can access repository
type AuthService interface {
	AddUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type authService struct {
	AuthRepo repository.AuthRepo
}

//NewAuthService is used to create a single instance of the service
func NewAuthService(a repository.AuthRepo) AuthService {
	return &authService{
		AuthRepo: a,
	}
}

//AddUser is a service layer that helps to add new user
func (a *authService) AddUser(user *models.User) (*models.User, error) {
	return a.AuthRepo.AddUser(user)
}

//GetUserByEmail is a service layer that helps to check user by their email
func (a *authService) GetUserByEmail(email string) (*models.User, error) {
	return a.AuthRepo.GetUserByEmail(email)
}
