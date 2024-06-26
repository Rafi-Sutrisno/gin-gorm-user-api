package service

import (
	"context"
	"errors"
	"mods/dto"
	"mods/entity"
	"mods/repository"
	"mods/utils"
	"strings"

	"github.com/google/uuid"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	// functional
	CreateUser(ctx context.Context, userDTO dto.CreateUserDTO) (entity.User, error)
	GetAllUser(ctx context.Context) ([]entity.User, error)
	GetAllDoctor(ctx context.Context) ([]entity.User, error)
	DeleteUser(ctx context.Context, id string) error
	IsDuplicateEmail(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	VerifyCredential(ctx context.Context, email string, pass string) (bool, error)
	UpdateUserName(ctx context.Context, updateDTO dto.UpdateNameUserDTO, userID string) (entity.User, error)
	UpdateUserNotelp(ctx context.Context, updateDTO dto.UpdateNotelpUserDTO, userID string) (entity.User, error)
	GetMe(ctx context.Context, id string) (entity.User, error)
	AddDoctor(ctx context.Context, userDTO dto.CreateUserDTO) (entity.User, error)
	ProfilePicture(ctx context.Context, imageDTO dto.PredictImageDTO, uid string) (string, error)
	DefaultCheck(ctx context.Context, uid string) (error)
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepository: ur,
	}
}

func (us *userService) CreateUser(ctx context.Context, userDTO dto.CreateUserDTO) (entity.User, error) {

	id := uuid.NewString()

	newUser := entity.User{
		U_Id:   id,
		Name:   userDTO.Name,
		Email:  userDTO.Email,
		Notelp: userDTO.Notelp,
		Pass:   userDTO.Pass,
		Picture: "https://storage.googleapis.com/example-bucket-test-cc-trw/default.png",
		Role:   "User",
	}

	return us.userRepository.AddUser(ctx, newUser)
}

func (us *userService) IsDuplicateEmail(ctx context.Context, email string) (bool, error) {
	checkUser, err := us.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if checkUser.Email == "" {
		return false, nil
	}

	return true, nil
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	return us.userRepository.GetUserByEmail(ctx, email)
}

func (us *userService) GetAllUser(ctx context.Context) ([]entity.User, error) {
	return us.userRepository.GetAllUser(ctx)
}

func (us *userService) GetAllDoctor(ctx context.Context) ([]entity.User, error) {
	return us.userRepository.GetAllDoctor(ctx)
}

func (us *userService) DeleteUser(ctx context.Context, id string) error {
	return us.userRepository.DeleteUser(ctx, id)
}

func (us *userService) UpdateUserName(ctx context.Context, updateDTO dto.UpdateNameUserDTO, userID string) (entity.User, error) {
	return us.userRepository.UpdateUserName(ctx, updateDTO, userID)
}

func (us *userService) UpdateUserNotelp(ctx context.Context, updateDTO dto.UpdateNotelpUserDTO, userID string) (entity.User, error) {
	return us.userRepository.UpdateUserNotelp(ctx, updateDTO, userID)
}

func (us *userService) VerifyCredential(ctx context.Context, email string, pass string) (bool, error) {
	checkUser, err := us.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	checkPassword, err := utils.PasswordCompare(checkUser.Pass, []byte(pass))
	if err != nil {
		return false, err
	}

	if checkUser.Email == email && checkPassword {
		return true, nil
	}
	return false, nil
}

func (us *userService) GetMe(ctx context.Context, id string) (entity.User, error) {
	return us.userRepository.Me(ctx, id)
}

func (us *userService) AddDoctor(ctx context.Context, userDTO dto.CreateUserDTO) (entity.User, error){
	id := uuid.NewString()

	newUser := entity.User{
		U_Id:   id,
		Name:   userDTO.Name,
		Email:  userDTO.Email,
		Notelp: userDTO.Notelp,
		Pass:   userDTO.Pass,
		Picture: "https://storage.googleapis.com/example-bucket-test-cc-trw/default.png",
		Role:   "Doctor",
	}

	return us.userRepository.AddUser(ctx, newUser)
}

func (us *userService) ProfilePicture(ctx context.Context, imageDTO dto.PredictImageDTO, uid string) (string, error) {
	imageFile := imageDTO.File
	
	img_uuid, err := utils.UploadToBucket(imageFile, "profile_picture")
	if err != nil {
		return "failed to upload to bucket", err
	}

	link := "https://storage.googleapis.com/example-bucket-test-cc-trw/" + img_uuid

	return us.userRepository.ProfilePicture(ctx, link, uid)

}

func (us *userService) DefaultCheck(ctx context.Context, uid string) (error) {
	user, err := us.userRepository.Me(ctx, uid)
	if err != nil {
		return errors.New("failed to get user")
	}

	currPicture := user.Picture

	if (currPicture == "https://storage.googleapis.com/example-bucket-test-cc-trw/default.png"){
		return nil
	}

	parts := strings.Split(currPicture, "/")
	fileName := parts[len(parts)-1]

	delete := utils.DeleteFromBucket("profile_picture", fileName)
	if delete != nil {
		return delete
	}

	return nil
}