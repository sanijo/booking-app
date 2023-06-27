package repository

import (
	"time"

	"github.com/sanijo/rent-app/internal/models"
)


type DatabaseRepo interface {
    AllUsers() bool
    InsertRent(rent models.Rent) (int, error)
    InsertRentRestriction(rentRestriction models.RentRestriction) error
    SearchAvailabilityByDatesAndModelID(start, end time.Time, modelID int) (bool, error)
    SearchAvailabilityForAllModels(start, end time.Time) ([]models.Model, error)
    GetModelByID(id int) (models.Model, error) 
}
