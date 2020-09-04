package services

import (
	"github.com/Aiscom-LLC/meals-api/api/swagger"
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// Client struct
type Client struct{}

// NewClient return pointer to client struct
// with all methods
func NewClient() *Client {
	return &Client{}
}

var cateringUserRepo = repository.NewCateringUserRepo()
var clientRepo = repository.NewClientRepo()

func (cl *Client) Get(query url.PaginationQuery, claims jwt.MapClaims) ([]swagger.Client, int, url.PaginationQuery, int, error) {
	var cateringID string

	id := claims["id"].(string)

	cateringUser, _ := cateringUserRepo.GetByKey("user_id", id)
	user, err := userRepo.GetByID(id)

	if err != nil {
		return []swagger.Client{}, 0, query, http.StatusBadRequest, err
	}

	if cateringUser.CateringID == uuid.Nil {
		cateringID = ""
	} else {
		cateringID = cateringUser.CateringID.String()
	}

	result, total, err := clientRepo.Get(query, cateringID, user.Role)

	if query.Page == 0 {
		query.Page = 1
	}

	return result, total, query, 0, err
}
