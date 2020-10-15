package repository

import (
	"errors"
	"net/http"
	"time"

	"github.com/Aiscom-LLC/meals-api/repository/models"

	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository/enums"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// OrderRepo struct
type OrderRepo struct{}

// NewOrderRepo returns pointer to order repository
// with all methods
func NewOrderRepo() *OrderRepo {
	return &OrderRepo{}
}

// Add adds order for provided user id
func (o OrderRepo) Add(userID string, date time.Time, newOrder models.OrderRequest) (models.UserOrder, error) {
	var orderExist int
	var order domain.Order
	var userOrder domain.UserOrders
	var total float32
	var userOrderResponse models.UserOrder

	config.DB.
		Model(&domain.UserOrders{}).
		Select("user_orders.*").
		Joins("left join orders o on user_orders.order_id = o.id").
		Where("user_orders.user_id = ? AND o.date = ? AND o.status != ?",
			userID, date, enums.OrderStatusTypesEnum.Canceled).
		Count(&orderExist)

	if orderExist != 0 {
		return models.UserOrder{}, errors.New("order for current day already created")
	}

	config.DB.Create(&order)

	for _, dish := range newOrder.Items {
		var price []float32

		orderDish := domain.OrderDishes{
			OrderID: order.ID,
			DishID:  dish.DishID,
			Amount:  dish.Amount,
		}

		if err := config.DB.Create(&orderDish).Error; err != nil {
			return models.UserOrder{}, err
		}

		config.DB.
			Model(&domain.Dish{}).
			Where("id = ?", dish.DishID).
			Pluck("price", &price)

		total += price[0] * float32(dish.Amount)

		order.Total = &total
		order.Date = date
		order.Status = &enums.OrderStatusTypesEnum.Pending
		order.Comment = &newOrder.Comment

		parsedUserID, _ := uuid.FromString(userID)
		userOrder = domain.UserOrders{
			UserID:  parsedUserID,
			OrderID: order.ID,
		}
	}

	config.DB.
		Model(&order).
		Update(&order)

	if err := config.DB.
		Create(&userOrder).
		Error; err != nil {
		return models.UserOrder{}, err
	}

	if err := o.getDishesForOrder(userOrder.OrderID, &userOrderResponse.Items); err != nil {
		return models.UserOrder{}, err
	}

	userOrderResponse.OrderID = userOrder.OrderID
	userOrderResponse.Total = total
	userOrderResponse.Status = *order.Status

	return userOrderResponse, nil
}

// CancelOrder changes status of order to canceled
func (o OrderRepo) CancelOrder(userID, orderID string) (int, error) {
	if err := config.DB.
		Model(&domain.UserOrders{}).
		Where("user_id = ? AND order_id = ?", userID, orderID).
		Find(&domain.UserOrders{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return http.StatusNotFound, err
		}
		return http.StatusBadRequest, err
	}

	if err := config.DB.
		Model(&domain.Order{}).
		Where("id = ? AND status = ?", orderID, enums.OrderStatusTypesEnum.Pending).
		Find(&domain.Order{}).
		Update("status", enums.OrderStatusTypesEnum.Canceled).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return http.StatusNotFound, errors.New("order not found or already canceled/approved")
		}
		return http.StatusBadRequest, err
	}

	return 0, nil
}

// GetUserOrder returns order for provided date for certain user
func (o OrderRepo) GetUserOrder(userID, date string) (models.UserOrder, int, error) {
	var userOrder models.UserOrder

	if err := config.DB.
		Model(&domain.UserOrders{}).
		Select("distinct on (o.id) o.id as order_id, o.total, o.status").
		Joins("left join orders o on user_orders.order_id = o.id").
		Where("user_orders.user_id = ? AND o.date = ? AND o.status != ?", userID, date, enums.OrderStatusTypesEnum.Canceled).
		Scan(&userOrder).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return models.UserOrder{}, http.StatusNotFound, err
		}
		return models.UserOrder{}, http.StatusBadRequest, err
	}

	if err := o.getDishesForOrder(userOrder.OrderID, &userOrder.Items); err != nil {
		return models.UserOrder{}, http.StatusBadRequest, err
	}

	return userOrder, 0, nil
}

// GetOrders return list of orders for catering or client
func (o OrderRepo) GetOrders(cateringID, clientID, date, companyType string) (models.SummaryOrderResult, int, error) {
	var result models.SummaryOrderResult

	if companyType == enums.CompanyTypesEnum.Client {
		result.Status = o.GetOrdersStatus(clientID, date)

		if err := config.DB.
			Model(&domain.User{}).
			Select("distinct on (c.id) c.name as category_name, c.id as category_id").
			Joins("left join client_users cu on cu.user_id = users.id").
			Joins("left join user_orders uo on uo.user_id = users.id").
			Joins("left join orders o on uo.order_id = o.id").
			Joins("left join order_dishes od on od.order_id = o.id").
			Joins("left join dishes d on od.dish_id = d.id").
			Joins("left join categories c on c.id = d.category_id").
			Where("cu.client_id = ? AND users.company_type = ? AND o.date = ?"+
				" AND o.status != ?", clientID, enums.CompanyTypesEnum.Client, date, enums.OrderStatusTypesEnum.Canceled).
			Scan(&result.SummaryOrders).
			Error; err != nil {
			return models.SummaryOrderResult{}, http.StatusBadRequest, err
		}

		for i := range result.SummaryOrders {
			if err := config.DB.
				Model(&domain.User{}).
				Select("d.name, sum(od.amount) as amount").
				Joins("left join client_users cu on cu.user_id = users.id").
				Joins("left join user_orders uo on uo.user_id = users.id").
				Joins("left join orders o on uo.order_id = o.id").
				Joins("left join order_dishes od on od.order_id = o.id").
				Joins("left join dishes d on od.dish_id = d.id").
				Joins("left join categories c on c.id = d.category_id").
				Where("cu.client_id = ? AND users.company_type = ? AND o.date = ?"+
					" AND c.id = ? and o.status != ?",
					clientID, enums.CompanyTypesEnum.Client, date, result.SummaryOrders[i].ID, enums.OrderStatusTypesEnum.Canceled).
				Group("d.name").
				Scan(&result.SummaryOrders[i].Items).
				Error; err != nil {
				return models.SummaryOrderResult{}, http.StatusBadRequest, err
			}
		}
		if err := config.DB.
			Model(&domain.User{}).
			Select("concat_ws(' ', users.last_name, users.first_name) as full_name,"+
				" cu.floor, users.id, o.total, o.comment").
			Joins("left join client_users cu on cu.user_id = users.id").
			Joins("left join user_orders uo on uo.user_id = users.id").
			Joins("left join orders o on uo.order_id = o.id").
			Where("cu.client_id = ? AND users.company_type = ? AND o.date = ?"+
				" AND o.status != ?", clientID, enums.CompanyTypesEnum.Client, date, enums.OrderStatusTypesEnum.Canceled).
			Scan(&result.UserOrders).
			Error; err != nil {
			return models.SummaryOrderResult{}, http.StatusBadRequest, err
		}

		for i := range result.UserOrders {
			if err := config.DB.
				Model(&domain.User{}).
				Select("d.name, od.amount").
				Joins("left join client_users cu on cu.user_id = users.id").
				Joins("left join user_orders uo on uo.user_id = users.id").
				Joins("left join orders o on uo.order_id = o.id").
				Joins("left join order_dishes od on od.order_id = o.id").
				Joins("left join dishes d on od.dish_id = d.id").
				Where("cu.client_id = ? AND users.company_type = ? AND o.date = ?"+
					" AND uo.user_id = ? AND o.status != ?",
					clientID, enums.CompanyTypesEnum.Client, date, result.UserOrders[i].ID, enums.OrderStatusTypesEnum.Canceled).
				Scan(&result.UserOrders[i].Items).
				Error; err != nil {
				return models.SummaryOrderResult{}, http.StatusBadRequest, err
			}
			result.Total += result.UserOrders[i].Total
		}

		return result, 0, nil
	}

	if cateringExist := config.DB.
		Where("id = ?", cateringID).
		Find(&domain.Catering{}).
		RowsAffected; cateringExist == 0 {
		return models.SummaryOrderResult{}, http.StatusNotFound, errors.New("catering not found")
	}

	if err := config.DB.
		Model(&domain.User{}).
		Select("distinct on (c.id) c.name as category_name, c.id as category_id").
		Joins("left join client_users cu on cu.user_id = users.id").
		Joins("left join user_orders uo on uo.user_id = users.id").
		Joins("left join orders o on uo.order_id = o.id").
		Joins("left join order_dishes od on od.order_id = o.id").
		Joins("left join dishes d on od.dish_id = d.id").
		Joins("left join categories c on c.id = d.category_id").
		Where("cu.client_id = ? AND users.company_type = ? AND o.date = ? AND o.status = ?",
			clientID, enums.CompanyTypesEnum.Client, date, enums.OrderStatusTypesEnum.Approved).
		Scan(&result.SummaryOrders).
		Error; err != nil {
		return models.SummaryOrderResult{}, http.StatusBadRequest, err
	}

	for i := range result.SummaryOrders {
		if err := config.DB.
			Model(&domain.User{}).
			Select("d.name, sum(od.amount) as amount").
			Joins("left join client_users cu on cu.user_id = users.id").
			Joins("left join user_orders uo on uo.user_id = users.id").
			Joins("left join orders o on uo.order_id = o.id").
			Joins("left join order_dishes od on od.order_id = o.id").
			Joins("left join dishes d on od.dish_id = d.id").
			Joins("left join categories c on c.id = d.category_id").
			Where("cu.client_id = ? AND users.company_type = ? AND o.date = ?"+
				" AND c.id = ? AND o.status = ?",
				clientID, enums.CompanyTypesEnum.Client, date, result.SummaryOrders[i].ID, enums.OrderStatusTypesEnum.Approved).
			Group("d.name").
			Scan(&result.SummaryOrders[i].Items).
			Error; err != nil {
			return models.SummaryOrderResult{}, http.StatusBadRequest, err
		}
	}

	if err := config.DB.
		Model(&domain.User{}).
		Select("concat_ws(' ', users.last_name, users.first_name) as full_name,"+
			" cu.floor, users.id, o.total, o.comment").
		Joins("left join client_users cu on cu.user_id = users.id").
		Joins("left join user_orders uo on uo.user_id = users.id").
		Joins("left join orders o on uo.order_id = o.id").
		Where("cu.client_id = ? AND users.company_type = ? AND o.date = ?"+
			" AND o.status = ?", clientID, enums.CompanyTypesEnum.Client, date, enums.OrderStatusTypesEnum.Approved).
		Scan(&result.UserOrders).
		Error; err != nil {
		return models.SummaryOrderResult{}, http.StatusBadRequest, err
	}

	for i := range result.UserOrders {
		if err := config.DB.
			Model(&domain.User{}).
			Select("d.name, od.amount").
			Joins("left join client_users cu on cu.user_id = users.id").
			Joins("left join user_orders uo on uo.user_id = users.id").
			Joins("left join orders o on uo.order_id = o.id").
			Joins("left join order_dishes od on od.order_id = o.id").
			Joins("left join dishes d on od.dish_id = d.id").
			Where("cu.client_id = ? AND users.company_type = ? AND o.date = ?"+
				" AND uo.user_id = ? AND o.status = ?", clientID, enums.CompanyTypesEnum.Client, date, result.UserOrders[i].ID, enums.OrderStatusTypesEnum.Approved).
			Scan(&result.UserOrders[i].Items).
			Error; err != nil {
			return models.SummaryOrderResult{}, http.StatusBadRequest, err
		}
		result.Total += result.UserOrders[i].Total
	}

	if len(result.SummaryOrders) != 0 {
		result.Status = &enums.OrderStatusTypesEnum.Approved
	}

	return result, 0, nil
}

// ApproveOrders changes status of orders to approved
func (o OrderRepo) ApproveOrders(clientID, date string) error {
	var orderIDs []struct {
		ID string
	}

	if areOrdersExist := config.DB.
		Table("orders as o").
		Select("o.id").
		Joins("left join user_orders uo on uo.order_id = o.id").
		Joins("left join users u on uo.user_id = u.id").
		Joins("left join client_users cu on uo.user_id = cu.user_id").
		Where("cu.client_id = ? AND u.company_type = ? AND o.date = ?"+
			" AND o.status != ?", clientID, enums.CompanyTypesEnum.Client, date, enums.OrderStatusTypesEnum.Canceled).
		Scan(&orderIDs).
		RowsAffected; areOrdersExist == 0 {
		return errors.New("client id is not found or no orders to approve for provided day")
	}

	for _, order := range orderIDs {
		config.DB.
			Model(&domain.Order{}).
			Where("id = ? and status != ?", order.ID, enums.OrderStatusTypesEnum.Canceled).
			Update("status", enums.OrderStatusTypesEnum.Approved)
	}

	return nil
}

func (o OrderRepo) getDishesForOrder(orderID uuid.UUID, dishes *[]models.OrderItem) error {
	if err := config.DB.
		Model(&domain.OrderDishes{}).
		Select("distinct on (d.id) d.name, d.price, d.id as dish_id, i.path as path, order_dishes.amount").
		Joins("left join dishes d on order_dishes.dish_id = d.id").
		Joins("left join image_dishes id on d.id = id.dish_id").
		Joins("left join images i on id.image_id = i.id").
		Where("order_dishes.order_id = ?", orderID).
		Scan(dishes).
		Error; err != nil {
		return err
	}
	return nil
}

// GetOrdersStatus return order status for provided client
func (o OrderRepo) GetOrdersStatus(clientID, date string) *string {
	var ordersStatus []string

	config.DB.
		Model(&domain.User{}).
		Select("o.status").
		Joins("left join client_users cu on cu.user_id = users.id").
		Joins("left join user_orders uo on uo.user_id = users.id").
		Joins("left join orders o on uo.order_id = o.id").
		Where("cu.client_id = ? AND users.company_type = ? AND o.date = ?"+
			" AND o.status != ?", clientID, enums.CompanyTypesEnum.Client, date, enums.OrderStatusTypesEnum.Canceled).
		Pluck("o.status", &ordersStatus)

	if len(ordersStatus) == 0 {
		return nil
	}

	return &ordersStatus[0]
}
