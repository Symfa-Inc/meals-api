package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/schemes/response"
	"go_api/src/types"
	"net/http"
	"time"
)

// ClientRepo struct
type ClientRepo struct{}

// NewClientRepo returns pointer to client repository
// with all methods
func NewClientRepo() *ClientRepo {
	return &ClientRepo{}
}

// Add adds client in DB
// returns error if that client name already exists
func (c ClientRepo) Add(cateringID string, client *domain.Client) error {
	if err := config.DB.
		Where("id = ?", cateringID).
		Find(&domain.Catering{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return err
		}

		return err
	}

	if exist := config.DB.
		Where("name = ?", client.Name).
		Find(client).RowsAffected; exist != 0 {
		return errors.New("client with that name already exist")
	}

	err := config.DB.Create(client).Error

	var cateringSchedules []domain.CateringSchedule

	config.DB.
		Where("catering_id = ?", client.CateringID).
		Find(&cateringSchedules)

	for _, schedule := range cateringSchedules {
		clientSchedule := domain.ClientSchedule{
			Day:       schedule.Day,
			Start:     schedule.Start,
			End:       schedule.End,
			IsWorking: schedule.IsWorking,
			ClientID:  client.ID,
		}
		config.DB.Create(&clientSchedule)
	}

	return err
}

// GetCateringClients returns list of catering Clients
func (c ClientRepo) GetCateringClients(cateringID string, query types.PaginationWithDateQuery) ([]response.CateringClient, int, error) {
	var cateringClients []response.CateringClient
	var total int

	page := query.Page
	limit := query.Limit

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	err := config.DB.
		Limit(limit).
		Offset((page-1)*limit).
		Model(&domain.Client{}).
		Select("clients.name, clients.id, concat_ws('/', count(distinct od.order_id), sum(od.amount)) as orders_dishes").
		Joins("left join users u on u.client_id = clients.id").
		Joins("left join user_orders uo on u.id = uo.user_id").
		Joins("left join orders o on uo.order_id = o.id").
		Joins("left join order_dishes od on od.order_id = o.id").
		Where("clients.catering_id = ? AND o.status = ? AND o.date = ?", cateringID, types.OrderStatusTypesEnum.Approved, query.Date).
		Group("clients.name, clients.id").
		Scan(&cateringClients).
		Count(&total).
		Error

	for i := range cateringClients {
		var total []int
		config.DB.
			Limit(limit).
			Offset((page-1)*limit).
			Model(&domain.Client{}).
			Select("o.total as total").
			Joins("left join users u on u.client_id = clients.id").
			Joins("left join user_orders uo on u.id = uo.user_id").
			Joins("left join orders o on uo.order_id = o.id").
			Where("clients.catering_id = ? AND o.status = ? AND o.date = ?", cateringID, types.OrderStatusTypesEnum.Approved, query.Date).
			Group("clients.name, clients.id").
			Pluck("sum(o.total)", &total)

		cateringClients[i].Total = total[0]
	}

	return cateringClients, total, err
}

// Get returns list of clients
func (c ClientRepo) Get(query types.PaginationQuery, cateringID, role string) ([]response.Client, int, error) {
	var clients []response.Client
	var total int
	var err error

	page := query.Page
	limit := query.Limit

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	if role == types.UserRoleEnum.SuperAdmin {
		config.DB.
			Model(&domain.Client{}).
			Count(&total)

		err = config.DB.
			Limit(limit).
			Offset((page - 1) * limit).
			Model(&domain.Client{}).
			Select("clients.*, c.name as catering_name, c.id as catering_id").
			Joins("left join caterings c on c.id = clients.catering_id").
			Order("clients.created_at DESC, clients.name ASC").
			Scan(&clients).
			Error

		return clients, total, err
	}

	config.DB.
		Model(&domain.Client{}).
		Where("catering_id = ?", cateringID).
		Count(&total)

	err = config.DB.
		Limit(limit).
		Offset((page-1)*limit).
		Model(&domain.Client{}).
		Select("clients.*, c.name as catering_name, c.id as catering_id").
		Joins("left join caterings c on c.id = clients.catering_id").
		Where("catering_id = ?", cateringID).
		Scan(&clients).
		Error

	return clients, total, err
}

// Delete soft delete of client
func (c ClientRepo) Delete(id string) error {
	if clientExist := config.DB.
		Where("id = ?", id).
		Delete(&domain.Client{}).
		RowsAffected; clientExist == 0 {
		return errors.New("client not found")
	}

	if userExist := config.DB.
		Model(&domain.User{}).
		Where("client_id = ? AND company_type = ?", id, types.CompanyTypesEnum.Client).
		Update(map[string]interface{}{
			"status":     types.StatusTypesEnum.Deleted,
			"deleted_at": time.Now(),
		}).
		RowsAffected; userExist == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Update updates client with passed args
// returns error and status code
func (c ClientRepo) Update(id string, client domain.Client) (int, error) {
	if clientExist := config.DB.
		Where("name = ? AND id = ?", client.Name, id).
		Find(&client).
		RowsAffected; clientExist == 0 {
		if nameExist := config.DB.
			Where("name = ?", client.Name).
			Find(&client).
			RowsAffected; nameExist != 0 {
			return http.StatusBadRequest, errors.New("client with that name already exist")
		}
	}

	if clientExist := config.DB.
		Model(&client).
		Where("id = ?", id).
		Update(&client).
		RowsAffected; clientExist == 0 {
		return http.StatusNotFound, errors.New("client not found")
	}

	return 0, nil
}

// GetByKey client by provided key value arguments
// Returns client, error
func (c ClientRepo) GetByKey(key, value string) (domain.Client, error) {
	var client domain.Client
	err := config.DB.Where(key+" = ?", value).First(&client).Error
	return client, err
}