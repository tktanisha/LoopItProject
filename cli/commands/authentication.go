package commands

import (
	"context"
	"fmt"
	"loopit/cli/utils"
	"loopit/internal/config"
	"loopit/internal/constants"
	"loopit/internal/models"
	validators "loopit/internal/utils"
)

// var (
// 	User_file_repo = user_repo.NewUserFileRepo("data/users.json", "data/sessions.json", "data/lenders.json")
// 	authService    = auth_service.NewAuthService(User_file_repo)
// )

func loginUtil(email, password string, ctx *context.Context) bool {
	token, user, err := AuthService.Login(email, password)
	if err != nil {
		fmt.Println(config.Red+"Login failed:"+config.Reset, err)
		return false
	}

	_, err = validators.ValidateJWT(token)
	if err != nil {
		fmt.Println(config.Red+"Invalid token:"+config.Reset, err)
		return false
	}

	userCtx := &models.UserContext{
		ID:   user.ID,
		Name: user.FullName,
		Role: user.Role,
	}
	*ctx = context.WithValue(*ctx, constants.UserCtxKey, userCtx)

	fmt.Println(config.Green + "Logged in successfully" + config.Reset)

	return true
}

func AuthLogin(ctx *context.Context) bool {
	fmt.Println("\nLogin")
	email := utils.InputWithValidation("Enter Email", validators.ValidateEmail)
	password := utils.InputPassword("Enter Password", validators.ValidatePassword)

	return loginUtil(email, password, ctx)
}

func AuthRegister(ctx *context.Context) bool {
	fmt.Println("\nCreate Account")
	fullname := utils.InputWithValidation("Enter Full Name", validators.ValidateFullName)
	email := utils.InputWithValidation("Enter Email", validators.ValidateEmail)
	password := utils.InputPassword("Enter Password", validators.ValidatePassword)
	phoneNumber := utils.InputWithValidation("Enter Phone Number", validators.ValidatePhoneNumber)
	address := utils.InputWithValidation("Enter Address", validators.ValidateAddress)

	err := AuthService.Register(&models.User{
		FullName:     fullname,
		Email:        email,
		PasswordHash: password,
		PhoneNumber:  phoneNumber,
		Address:      address,
	})

	if err != nil {
		fmt.Println(config.Red+"Registration failed:"+config.Reset, err)
		return false
	}

	return loginUtil(email, password, ctx)
}

func AuthLogout(ctx *context.Context) {
	*ctx = nil // Clear context to remove user session
	fmt.Println("Logged out successfully.")
}
