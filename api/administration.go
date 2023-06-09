package main

import (
	"net/http"

	"github.com/WrastAct/lospoloshermanos/internal/data"
)

func (app *application) customQueryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Query string `json:"query"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	query := input.Query
	var queryResult []map[string]interface{}

	var message envelope
	status := http.StatusAccepted

	result := app.db.Raw(query).Scan(&queryResult)
	if result.Error != nil {
		message = envelope{"error": result.Error.Error()}
		status = http.StatusBadRequest
	} else {
		message = envelope{"result": queryResult}
	}

	err = app.writeJSON(w, status, message, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listTablesHandler(w http.ResponseWriter, r *http.Request) {
	var tables []string

	err := app.db.Table("information_schema.tables").
		Where("table_schema = ?", "public").
		Pluck("table_name", &tables).Error

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tables": tables}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) describeTableHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Table string `json:"table"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var tableDescription []data.ColumnDescription

	rows, err := app.db.Table("information_schema.columns").
		Select("column_name, data_type, CAST (is_nullable AS bool)").
		Where("table_name = ?", input.Table).Rows()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var column data.ColumnDescription

		rows.Scan(&column.Name, &column.Type, &column.IsNullable)
		tableDescription = append(tableDescription, column)
	}

	if tableDescription == nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{input.Table: tableDescription}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
