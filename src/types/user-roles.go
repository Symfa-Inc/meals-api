package types

// UserRole type
type UserRole string

type enum struct {
	SuperAdmin   UserRole
	CompanyAdmin UserRole
	ClientAdmin  UserRole
	User         UserRole
}

// UserRoleEnum enum
var UserRoleEnum enum = enum{
	SuperAdmin:   "Super administrator",
	CompanyAdmin: "Delivery administrator",
	ClientAdmin:  "Company administrator",
	User:         "User",
}
