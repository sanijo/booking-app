package dbrepo

import (
	"time"

	"github.com/sanijo/rent-app/internal/models"
)


func (m *testDBRepo) AllUsers() bool {
    return true
}

// InsertRent inserts a rent into the database after data is obtained from the
// form.
func (m *testDBRepo) InsertRent(rent models.Rent) (int, error) {
    return 1, nil
}

// InsertRentRestriction inserts a rent restriction into the database after data 
// is obtained from the form.
func (m *testDBRepo) InsertRentRestriction(rentRestriction models.RentRestriction) error {
    return nil
}

// SearchAvailabilityByDatesByModelID returns true if availability exists for
// modelID, and false if no availability exists.
func (m *testDBRepo) SearchAvailabilityByDatesAndModelID(start, end time.Time, modelID int) (bool, error) {
    return false, nil
}
    
// SearchAvailabilityForAllModels returns a slice of available models if any,
// for given start and end dates.
func (m *testDBRepo) SearchAvailabilityForAllModels(start, end time.Time) ([]models.Model, error) {
    var availableCarModels []models.Model
    return availableCarModels, nil
}

// GetModelByID returns a model by id.
func (m *testDBRepo) GetModelByID(id int) (models.Model, error) {
    var model models.Model
    return model, nil
}
