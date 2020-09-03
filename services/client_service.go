package services

import (
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/schemes/response"
	"github.com/Aiscom-LLC/meals-api/types"
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

func (cl *Client) Get(query types.PaginationQuery, cateringID string, claims jwt.MapClaims) ([]response.Client, int, types.PaginationQuery, int, error) {

	id := claims["id"].(string)

	cateringUser, _ := cateringUserRepo.GetByKey("user_id", id)
	user, err := userRepo.GetByID(id)

	if err != nil {
		return []response.Client{}, 0, query, http.StatusBadRequest, err
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
