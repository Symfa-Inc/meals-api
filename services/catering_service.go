package services

import (
	"github.com/Aiscom-LLC/meals-api/api/types"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// CateringService struct
type CateringService struct{}

// NewCateringService returns pointer to Auth struct
// with all methods
func NewCateringService() *CateringService {
	return &CateringService{}
}

func (cs *CateringService) Get(claims jwt.MapClaims, query *types.PaginationQuery) ([]domain.Catering, int, int, error) {
	cateringUserRepo := repository.NewCateringUserRepo()
	cateringRepo := repository.NewCateringRepo()
	var cateringID string

	id := claims["id"].(string)

	user, _ := cateringUserRepo.GetByKey("id", id)

	if user.CateringID == uuid.Nil {
		cateringID = ""
	} else {
		cateringID = user.CateringID.String()
	}

	caterings, total, err := cateringRepo.Get(cateringID, *query)

	if err != nil {
		return nil, 0, http.StatusBadRequest, err
	}

	if query.Page == 0 {
		query.Page = 1
	}

	return caterings, total, 0, nil
}
