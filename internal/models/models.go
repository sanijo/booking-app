package models

import "time"

// User holds database users data
type User struct {
    ID int
    FirstName string
    LastName string
    Email string
    Password string
    CreatedAt time.Time
    UpdatedAt time.Time
    AccessLevel int
}

// Model holds database models data
type Model struct {
    ID int
    ModelName string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// RestrictionType holds database restriction types data
type RestrictionType struct {
    ID int
    RestrictionName string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Rent holds database rent data
type Rent struct {
    ID int
    FirstName string
    LastName string
    Email string
    Phone string
    StartDate time.Time
    EndDate time.Time
    ModelID int
    CreatedAt time.Time
    UpdatedAt time.Time
    Model Model
}

// RentRestriction holds database rent restrictions data
type RentRestriction struct {
    ID int
    StartDate time.Time
    EndDate time.Time
    ModelID int
    RentID int
    CreatedAt time.Time
    UpdatedAt time.Time
    RestrictionID int
    Model Model
    Rent Rent
    Restriction RestrictionType
}
