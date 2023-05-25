package main

import (
	"errors"
	"net/http"

	"github.com/WrastAct/lospoloshermanos/internal/data"
	"gorm.io/gorm"
)

func (app *application) createPropertyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	property := data.Properties{
		Name: input.Name,
	}

	result := app.db.Create(&property)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"property": property}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showPropertyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var property data.Properties

	err = app.db.First(&property, id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"property": property}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listPropertiesHandler(w http.ResponseWriter, r *http.Request) {
	var properties []data.Properties

	result := app.db.Find(&properties)

	err := result.Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"properties": properties}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePropertyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var properties data.Properties

	err = app.db.First(&properties, id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name *string `json:"name"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		properties.Name = *input.Name
	}

	err = app.db.Save(&properties).Error
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"property": properties}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePropertyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var property data.Properties

	err = app.db.Delete(&property, id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "order property successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
