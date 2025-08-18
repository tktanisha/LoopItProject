package commands

import (
	"context"
	"fmt"
	"loopit/cli/initializer"
	"loopit/cli/utils"
	"loopit/internal/config"
	"loopit/internal/constants"
	"loopit/internal/models"
	validators "loopit/internal/utils"
	"loopit/pkg/logger"
)

var log = logger.GetLogger()

func loginUtil(email, password string, ctx *context.Context) bool {
	log.Info(fmt.Sprintf("Login CLI initiated for email: %s", email))

	token, user, err := initializer.AuthService.Login(email, password)
	if err != nil {
		log.Warning(fmt.Sprintf("Login CLI failed for email: %s, error: %v", email, err))
		fmt.Println(config.Red+"Login failed:"+config.Reset, err)
		return false
	}

	_, err = validators.ValidateJWT(token)
	if err != nil {
		log.Error(fmt.Sprintf("Token validation failed for email: %s, error: %v", email, err))
		fmt.Println(config.Red+"Invalid token:"+config.Reset, err)
		return false
	}

	userCtx := &models.UserContext{
		ID:   user.ID,
		Name: user.FullName,
		Role: user.Role,
	}
	*ctx = context.WithValue(*ctx, constants.UserCtxKey, userCtx)

	log.Info(fmt.Sprintf("Login CLI successful for email: %s", email))
	fmt.Println(config.Green + "Logged in successfully" + config.Reset)

	return true
}

func AuthLogin(ctx *context.Context) bool {
	log.Info("AuthLogin CLI command started")
	fmt.Println("\nLogin")
	email := utils.InputWithValidation("Enter Email", validators.ValidateEmail)
	password := utils.InputPassword("Enter Password", validators.ValidatePassword)

	return loginUtil(email, password, ctx)
}

func AuthRegister(ctx *context.Context) bool {
	log.Info("AuthRegister CLI command started")
	fmt.Println("\nCreate Account")
	fullname := utils.InputWithValidation("Enter Full Name", validators.ValidateFullName)
	email := utils.InputWithValidation("Enter Email", validators.ValidateEmail)
	password := utils.InputPassword("Enter Password", validators.ValidatePassword)
	phoneNumber := utils.InputWithValidation("Enter Phone Number", validators.ValidatePhoneNumber)
	address := utils.InputWithValidation("Enter Address", validators.ValidateAddress)

	err := initializer.AuthService.Register(&models.User{
		FullName:     fullname,
		Email:        email,
		PasswordHash: password,
		PhoneNumber:  phoneNumber,
		Address:      address,
	})

	if err != nil {
		log.Warning(fmt.Sprintf("Registration CLI failed for email: %s, error: %v", email, err))
		fmt.Println(config.Red+"Registration failed:"+config.Reset, err)
		return false
	}

	log.Info(fmt.Sprintf("Registration CLI successful for email: %s", email))
	return loginUtil(email, password, ctx)
}

func AuthLogout(ctx *context.Context) {
	log.Info("Logout CLI command executed")
	*ctx = nil
	fmt.Println("Logged out successfully.")
}
