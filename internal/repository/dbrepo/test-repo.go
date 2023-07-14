package dbrepo

import (
	"errors"
	"time"

	"github.com/sanijo/rent-app/internal/models"
)


func (m *testDBRepo) AllUsers() bool {
    return true
}

// InsertRent inserts a rent into the database after data is obtained from the
// form.
func (m *testDBRepo) InsertRent(rent models.Rent) (int, error) {
    // If the rent modelID is greater than 3, return an error    
    if rent.ModelID == 3 {
        return 0, errors.New("some error")
    }
    return 1, nil
}

// InsertRentRestriction inserts a rent restriction into the database after data 
// is obtained from the form.
func (m *testDBRepo) InsertRentRestriction(rentRestriction models.RentRestriction) error {
    // If the rent restrictionID is greater than 3, return an error 
    if rentRestriction.ModelID == 4 {
        return errors.New("some error")
    }

    return nil
}

// SearchAvailabilityByDatesByModelID returns true if availability exists for
// modelID, and false if no availability exists.
func (m *testDBRepo) SearchAvailabilityByDatesAndModelID(start, end time.Time, modelID int) (bool, error) {
    // If the start date is equal to 2021-01-01, return false and nil
    layout := "2006-01-02"
    date, _ := time.Parse(layout, "2021-01-01")
    if start == date {
        return false, nil
    }

    return true, nil
}
    
// SearchAvailabilityForAllModels returns a slice of available models if any,
// for given start and end dates.
func (m *testDBRepo) SearchAvailabilityForAllModels(start, end time.Time) ([]models.Model, error) {
    var availableCarModels []models.Model
    // If the start date is equal to 2021-01-01, return an error 
    layout := "2006-01-02"
    date, _ := time.Parse(layout, "2021-01-01")
    if start == date {
        return availableCarModels, errors.New("date equal to 2021-01-01, test error")
    }
    // If the start date is equal to 2022-01-02, return availableCarModels with 
    // modelID 1 and 2
    date, _ = time.Parse(layout, "2022-01-02")
    if start == date {
        availableCarModels = append(availableCarModels, models.Model{
            ID: 1,
        })
        availableCarModels = append(availableCarModels, models.Model{
            ID: 2,
        })
        return availableCarModels, nil
    }

    return availableCarModels, nil
}

// GetModelByID returns a model by id.
func (m *testDBRepo) GetModelByID(id int) (models.Model, error) {
    var model models.Model

    if id > 2 {
        return model, errors.New("some error")
    }

    return model, nil
}
