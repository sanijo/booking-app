package dbrepo

import (
	"context"
	"time"

	"github.com/sanijo/rent-app/internal/models"
)


func (m *postgresDbRepo) AllUsers() bool {
    return true
}

// InsertRent inserts a rent into the database after data is obtained from the
// form.
func (m *postgresDbRepo) InsertRent(rent models.Rent) (int, error) {
    // Create a context with a timeout of 3 seconds which will be used to
    // kill the query if it takes too long.
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    var newID int

    query := `insert into rent (first_name, last_name, email, phone, start_date,
            end_date, model_id, created_at, updated_at)
            values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

    err := m.DB.QueryRowContext(
        ctx,
        query,
        rent.FirstName, 
        rent.LastName,
        rent.Email,
        rent.Phone,
        rent.StartDate,
        rent.EndDate,
        rent.ModelID,
        time.Now(),
        time.Now(),
    ).Scan(&newID)
    
    if err != nil {
        return 0, err
    }

    return newID, nil
}

// InsertRentRestriction inserts a rent restriction into the database after data 
// is obtained from the form.
func (m *postgresDbRepo) InsertRentRestriction(rentRestriction models.RentRestriction) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := `insert into rent_restrictions (start_date, end_date, model_id, 
            rent_id, restriction_id, created_at, updated_at)
            values ($1, $2, $3, $4, $5, $6, $7)`

    _, err := m.DB.ExecContext(
        ctx,
        query,
        rentRestriction.StartDate,
        rentRestriction.EndDate,
        rentRestriction.ModelID,
        rentRestriction.RentID,
        rentRestriction.RestrictionID,
        time.Now(),
        time.Now(),
    )
    
    if err != nil {
        return err
    }

    return nil
}

// SearchAvailabilityByDatesByModelID returns true if availability exists for
// modelID, and false if no availability exists.
func (m *postgresDbRepo) SearchAvailabilityByDatesAndModelID(start, end time.Time, modelID int) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := `
        select 
            count(id) 
        from 
            rent_restrictions 
        where 
            model_id = $1 and $2 < end_date and $3 > start_date;`

    var numRows int

    // do query
    queryResult := m.DB.QueryRowContext(ctx, query, modelID, start, end)
    // scan the result into the address of numRows variable
    err := queryResult.Scan(&numRows)
    if err != nil {
        return false, err
    }

    if numRows == 0 {
        return true, nil
    }

    return false, nil
}
    
// SearchAvailabilityForAllModels returns a slice of available models if any,
// for given start and end dates.
func (m *postgresDbRepo) SearchAvailabilityForAllModels(start, end time.Time) ([]models.Model, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    var availableCarModels []models.Model

    query := `
        select 
            m.id, m.model_name
        from 
            models m
        where 
            m.id not in 
            (select 
                rr.model_id 
            from 
                rent_restrictions rr
            where 
                $1 < rr.end_date and $2 > rr.start_date);`

    rows, err := m.DB.QueryContext(ctx, query, start, end)
    if err != nil {
        return availableCarModels, err
    }
    defer rows.Close()

    for rows.Next() {
        var model models.Model
        err = rows.Scan(
            &model.ID,
            &model.ModelName,
        )
        if err != nil {
            return availableCarModels, err
        }

        availableCarModels = append(availableCarModels, model)
    }

    if err = rows.Err(); err != nil {
        return availableCarModels, err
    }

    return availableCarModels, nil
}

// GetModelByID returns a model by id.
func (m *postgresDbRepo) GetModelByID(id int) (models.Model, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    var model models.Model

    query := `
        select 
            id, model_name, created_at, updated_at 
        from 
            models 
        where 
            id = $1`

    row := m.DB.QueryRowContext(ctx, query, id)

    err := row.Scan(
        &model.ID,
        &model.ModelName,
        &model.CreatedAt,
        &model.UpdatedAt,
    )
    if err != nil {
        return model, err
    }

    return model, nil
}
