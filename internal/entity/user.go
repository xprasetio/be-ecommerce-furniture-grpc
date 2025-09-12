package entity

import "time"


const (
	UserRoleAdmin = "admin"
	UserRoleCustomer = "customer"
)

type UserRole struct {
	Id        int
	Name      string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
	DeletedAt time.Time
	DeletedBy string
	IsDeleted bool
}

type User struct {
	Id        string
	FullName  string
	Email     string
	Password  string
	RoleCode  string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
	DeletedAt time.Time
	DeletedBy string
	IsDeleted bool
}
