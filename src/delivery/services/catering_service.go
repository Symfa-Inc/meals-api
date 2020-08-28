package services

import (
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/Aiscom-LLC/meals-api/src/schemes/response"
	"github.com/Aiscom-LLC/meals-api/src/types"
	"github.com/Aiscom-LLC/meals-api/src/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// Catering struct
type Catering struct{}

// NewCatering returns pointer to catering struct
// with all methods
func NewCatering() *Catering {
	return &Catering{}
}

var cateringRepo = repository.NewCateringRepo()

func (ca Catering) AddService(catering domain.Catering, c *gin.Context) {

	c.JSON(http.StatusCreated, catering)
}

func (ca Catering) DeleteService(path types.PathID, c *gin.Context) {

	c.Status(http.StatusNoContent)
}

func (ca Catering) GetByIdService(r domain.Catering, c *gin.Context) {

	c.JSON(http.StatusOK, r)
}

func (ca Catering) GetService(user domain.CateringUser, query types.PaginationQuery, c *gin.Context) {
	var cateringID string

	if user.CateringID == uuid.Nil {
		cateringID = ""
	} else {
		cateringID = user.CateringID.String()
	}

	result, total, err := cateringRepo.Get(cateringID, query)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}

	c.JSON(http.StatusOK, response.GetCaterings{
		Items: result,
		Page:  query.Page,
		Total: total,
	})
}

func (ca Catering) Update(path types.PathID, c *gin.Context, cm domain.Catering) {

	if code, err := cateringRepo.Update(path.ID, cm); err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}
