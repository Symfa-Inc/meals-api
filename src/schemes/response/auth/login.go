package auth

import (
	uuid "github.com/satori/go.uuid"
	"go_api/src/types"
	"time"
)

// IsAuthenticated model for /is-authenticated response route
type IsAuthenticated struct {
	ID        uuid.UUID      `json:"id,omitempty"`
	FirstName string         `json:"firstName,omitempty"`
	LastName  string         `json:"lastName,omitempty"`
	Email     string         `json:"email,omitempty"`
	Role      types.UserRole `json:"role,omitempty"`
}//@name IsAuthenticatedResponse

//Response model for /login response route
type LoginResponse struct {
	Code   int       `json:"code"`
	Expire time.Time `json:"expire"`
	Token  string    `json:"token"`
}//@name LoginResponse
