package pb

import (
	"context"
	"fmt"
	"gdsc/baro/app/user/models"
	"gdsc/baro/app/user/repositories"
	"gdsc/baro/global/auth"
	"gdsc/baro/global/utils"

	userpb "gdsc/baro/protos/user"
)

type UserPbApp struct {
	UserRepository repositories.UserRepositoryInterface
	UserUtil       utils.UserUtilInterface
	userpb.UnimplementedUserServiceServer
}

func NewUserPbApp(userRepository repositories.UserRepositoryInterface, userUtil utils.UserUtilInterface) *UserPbApp {
	return &UserPbApp{
		UserRepository: userRepository,
		UserUtil:       userUtil,
	}
}

func (app *UserPbApp) Login(c context.Context, req *userpb.RequestCreateUser) (*userpb.ResponseToken, error) {
	user := models.User{
		Name:     req.Name,
		Nickname: req.Nickname,
		Email:    req.Email,
		Age:      int(req.Age),
		Gender:   req.Gender,
	}

	foundUser, err := app.UserRepository.FindOrCreateByEmail(&user)
	if err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(auth.NewClaim(fmt.Sprint(foundUser.ID)))
	if err != nil {
		return nil, err
	}

	return &userpb.ResponseToken{Token: token}, nil
}

func (app *UserPbApp) GetUserInfo(c context.Context, req *userpb.Empty) (*userpb.ResponseUser, error) {
	userID := c.Value(auth.UserIDKey).(string)

	user, err := app.UserRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &userpb.ResponseUser{
		Id:       uint64(user.ID),
		Name:     user.Name,
		Nickname: user.Nickname,
		Email:    user.Email,
		Age:      int32(user.Age),
		Gender:   user.Gender,
	}, nil
}

func (app *UserPbApp) UpdateUserInfo(c context.Context, req *userpb.RequestUpdateUser) (*userpb.ResponseUser, error) {
	userID := c.Value(auth.UserIDKey).(string)

	user, err := app.UserRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	app.updateUser(&user, req)

	user, err = app.UserRepository.Update(&user)
	if err != nil {
		return nil, err
	}

	return &userpb.ResponseUser{
		Id:       uint64(user.ID),
		Name:     user.Name,
		Nickname: user.Nickname,
		Email:    user.Email,
		Age:      int32(user.Age),
		Gender:   user.Gender,
	}, nil
}

func (app *UserPbApp) updateUser(user *models.User, input *userpb.RequestUpdateUser) {
	if input.Nickname != "" {
		user.Nickname = input.Nickname
	}

	if input.Age != 0 {
		user.Age = int(input.Age)
	}

	if input.Gender != "" {
		user.Gender = input.Gender
	}
}

func (app *UserPbApp) DeleteUser(c context.Context, req *userpb.Empty) (*userpb.Empty, error) {
	userID := c.Value(auth.UserIDKey).(string)

	user, err := app.UserRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	err = app.UserRepository.Delete(&user)
	if err != nil {
		return nil, err
	}

	return &userpb.Empty{}, nil
}

func (s *UserPbApp) mustEmbedUnimplementedVideoServiceServer() {}
