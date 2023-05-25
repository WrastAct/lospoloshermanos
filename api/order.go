package main

import (
	"errors"
	"net/http"

	"github.com/WrastAct/lospoloshermanos/internal/data"
	"gorm.io/gorm"
)

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SenderID        uint64   `json:"sender_id"`
		SenderAddress   string   `json:"sender_address"`
		ReceiverID      uint64   `json:"receiver_id"`
		ReceiverAddress string   `json:"receiver_address"`
		Mass            float32  `json:"mass"`
		InsuranceID     uint64   `json:"insurance_id"`
		Value           float32  `json:"value"`
		Coverage        float32  `json:"coverage"`
		Properties      []string `json:"properties"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{ID: input.SenderID}

	err = app.db.First(&user).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.senderNotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user = &data.User{ID: input.ReceiverID}

	err = app.db.First(&user).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.receiverNotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	insurance := &data.Insurance{ID: input.InsuranceID}
	err = app.db.First(&insurance).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.insuranceNotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	order := data.Order{
		SenderID:        input.SenderID,
		SenderAddress:   input.SenderAddress,
		ReceiverID:      input.ReceiverID,
		ReceiverAddress: input.ReceiverAddress,
		Mass:            input.Mass,
		InsuranceID:     input.InsuranceID,
		Value:           input.Value,
		Coverage:        input.Coverage,
	}

	for _, v := range input.Properties {
		order.Properties = append(order.Properties, data.Properties{Name: v})
	}

	result := app.db.Create(&order)
	if result.Error != nil {
		app.serverErrorResponse(w, r, result.Error)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"order": order}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var order data.Order

	err = app.db.Model(&data.Order{}).Preload("Properties").First(&order, id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"order": order}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var orders []data.Order

	result := app.db.Find(&orders)

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

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var order data.Order

	err = app.db.First(&order, id).Error
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
		SenderID        *uint64  `json:"sender_id"`
		SenderAddress   *string  `json:"sender_address"`
		ReceiverID      *uint64  `json:"receiver_id"`
		ReceiverAddress *string  `json:"receiver_address"`
		Mass            *float32 `json:"mass"`
		InsuranceID     *uint64  `json:"insurance_id"`
		Value           *float32 `json:"value"`
		Coverage        *float32 `json:"coverage"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.SenderID != nil {
		order.SenderID = *input.SenderID
	}

	if input.SenderAddress != nil {
		order.SenderAddress = *input.SenderAddress
	}

	if input.ReceiverID != nil {
		order.ReceiverID = *input.ReceiverID
	}

	if input.ReceiverAddress != nil {
		order.ReceiverAddress = *input.ReceiverAddress
	}

	if input.Mass != nil {
		order.Mass = *input.Mass
	}

	if input.InsuranceID != nil {
		order.InsuranceID = *input.InsuranceID
	}

	if input.Value != nil {
		order.Value = *input.Value
	}

	if input.Coverage != nil {
		order.Coverage = *input.Coverage
	}

	err = app.db.Save(&order).Error
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"order": order}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var order data.Order

	err = app.db.Delete(&order, id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "order successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
