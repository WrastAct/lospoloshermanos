package main

import (
	"errors"
	"net/http"

	"github.com/WrastAct/maestro/internal/data"
	"gorm.io/gorm"
)

func (app *application) createPackageStatsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Mean              float64 `json:"mean"`
		Variance          float64 `json:"variance"`
		StandardDeviation float64 `json:"standard_deviation"`
		Prediction        float64 `json:"prediction"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	packageStats := data.PackageStats{
		Mean:              input.Mean,
		Variance:          input.Variance,
		StandardDeviation: input.StandardDeviation,
		Prediction:        input.Prediction,
	}

	result := app.db.Create(&packageStats)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"package": packageStats}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showPackageStatsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var packageStats data.PackageStats

	err = app.db.First(&packageStats, id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"package": packageStats}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listPackageStatsHandler(w http.ResponseWriter, r *http.Request) {
	var packageStats []data.PackageStats

	result := app.db.Find(&packageStats)

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

	err = app.writeJSON(w, http.StatusOK, envelope{"packages": packageStats}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePackageStatsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var packageStats data.PackageStats

	err = app.db.First(&packageStats, id).Error
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
		Mean              *float64 `json:"mean"`
		Variance          *float64 `json:"variance"`
		StandardDeviation *float64 `json:"standard_deviation"`
		Prediction        *float64 `json:"prediction"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Mean != nil {
		packageStats.Mean = *input.Mean
	}

	if input.Variance != nil {
		packageStats.Variance = *input.Variance
	}

	if input.StandardDeviation != nil {
		packageStats.StandardDeviation = *input.StandardDeviation
	}

	if input.Prediction != nil {
		packageStats.Prediction = *input.Prediction
	}

	err = app.db.Save(&packageStats).Error
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"package": packageStats}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePackageStatsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var packageStats data.PackageStats

	err = app.db.Delete(&packageStats, id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "package successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
