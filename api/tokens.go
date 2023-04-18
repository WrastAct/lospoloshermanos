package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/WrastAct/maestro/internal/data"
	"github.com/WrastAct/maestro/internal/validator"
	"gorm.io/gorm"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var user data.User

	result := app.db.Where("email = ?", input.Email).Take(&user)
	if result.Error != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, result.Error)
		}
		return
	}

	password := data.Password{Hash: user.Password}

	match, err := password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, tokenPlaintext, err := data.GenerateToken(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	result = app.db.Create(&token)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": tokenPlaintext, "expiry": token.Expiry}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
