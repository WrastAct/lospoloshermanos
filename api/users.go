package main

import (
	"crypto/sha256"
	"errors"
	"net/http"
	"time"

	"github.com/WrastAct/maestro/internal/data"
	"github.com/WrastAct/maestro/internal/validator"
	"gorm.io/gorm"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Birthday string `json:"birthday"`
		Address  string `json:"address"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Language string `json:"lang"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var password data.Password
	err = password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := data.User{
		Name:      input.Name,
		Birthday:  input.Birthday,
		Address:   input.Address,
		Email:     input.Email,
		Password:  password.Hash,
		Activated: false,
	}

	v := validator.New()

	if data.ValidateUser(v, &user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	result := app.db.Create(&user)
	if result.Error != nil {
		switch {
		case result.Error.Error() == `duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`:
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, result.Error)
		}
		return
	}

	token, tokenPlaintext, err := data.GenerateToken(user.ID, 5*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	result = app.db.Create(&token)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	app.background(func() {
		data := map[string]interface{}{
			"activationToken": tokenPlaintext,
			"userID":          user.ID,
		}

		switch input.Language {
		case "en":
			err = app.mailer.Send(user.Email, "user_welcome_en.tmpl", data)
		case "ua":
			err = app.mailer.Send(user.Email, "user_welcome_ua.tmpl", data)
		default:
			err = app.mailer.Send(user.Email, "user_welcome_en.tmpl", data)
		}
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var user data.User

	tokenHash := sha256.Sum256([]byte(input.TokenPlaintext))
	result := app.db.Table("users").Joins("left join tokens ON tokens.user_id = users.id").
		Where("tokens.hash = ?", tokenHash[:]).
		Where("tokens.scope = ?", data.ScopeActivation).
		Where("tokens.expiry > ?", time.Now()).
		Take(&user)

	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, result.Error)
		}
		return
	}

	user.Activated = true
	user.Permissions = []data.Permissions{{ID: 2, Code: "user"}}

	result = app.db.Save(&user)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	result = app.db.Where("scope = ?", data.ScopeActivation).
		Where("user_id = ?", user.ID).
		Delete(&data.Token{})

	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
