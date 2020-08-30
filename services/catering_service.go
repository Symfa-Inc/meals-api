package services

import (
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// AuthService struct
type CateringService struct{}

// NewAuthService returns pointer to Auth struct
// with all methods
func NewCateringService() *CateringService {
	return &CateringService{}
}

func (cs *CateringService) Get(c *gin.Context, query *types.PaginationQuery) ([]domain.Catering, int, int, error) {
	cateringUserRepo := repository.NewCateringUserRepo()
	cateringRepo := repository.NewCateringRepo()
	var cateringID string

	claims, err := middleware.Passport().GetClaimsFromJWT(c)

	if err != nil {
		return nil, 0, http.StatusUnauthorized, err
	}

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