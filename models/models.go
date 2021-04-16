package models

import (
	"database/sql"
	"hash"
	"time"
)

type BaseModel struct {
	Id        uint16
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

type UserType string

const (
	Admin    UserType = "admin"
	Customer UserType = "customer"
)

type User struct {
	BaseModel
	UserType UserType
}

type UserProfile struct {
	BaseModel
	UserId      uint16
	FirstName   string
	LastName    string
	Email       string
	Phone       sql.NullString
	DateOfBirth sql.NullTime
	City        string
	Country     string
	Password    hash.Hash
}
