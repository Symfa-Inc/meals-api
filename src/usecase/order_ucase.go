package usecase

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/Aiscom-LLC/meals-api/src/schemes/request"
	"github.com/Aiscom-LLC/meals-api/src/types"
	"github.com/Aiscom-LLC/meals-api/src/utils"

	"github.com/gin-gonic/gin"
)

// Order struct
type Order struct{}

// NewOrder returns pointer to client struct
// with all methods
func NewOrder() *Order {
	return &Order{}
}

var orderRepo = repository.NewOrderRepo()

// Add creates order for client user
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags users orders
// @Param id path string false "User ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Param body body request.OrderRequest false "User order"
// @Success 201 {object} response.UserOrder false "Order for user"
// @Failure 400 {object} types.Error "Error"
// @Router /users/{id}/orders [post]
func (o Order) Add(c *gin.Context) {
	var path types.PathID
	var order request.OrderRequest
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&order, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	for i, dish := range order.Items {
		if dish.Amount == 0 {
			utils.CreateError(http.StatusBadRequest, "can't add dish with 0 amount", c)
			return
		}
		for j := i + 1; j < len(order.Items); j++ {
			if dish.DishID == order.Items[j].DishID {
				utils.CreateError(http.StatusBadRequest, "can't add 2 same dishes, please increment amount field instead", c)
				return
			}
		}
	}

	date, err := time.Parse(time.RFC3339, query.Date)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	difference := date.Sub(time.Now().Truncate(time.Hour * 24)).Hours()

	if difference < 0 {
		utils.CreateError(http.StatusBadRequest, "can't add order to previous date", c)
		return
	}

	userOrder, err := orderRepo.Add(path.ID, date, order)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, userOrder)
}

// CancelOrder changes status of order to canceled
// @Summary Returns error or 204 status code if success
// @Produce json
// @Accept json
// @Tags users orders
// @Param id path string false "User ID"
// @Param orderId path string false "Order ID"
// @Success 204 "Successfully canceled"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /users/{id}/orders/{orderId} [delete]
func (o Order) CancelOrder(c *gin.Context) {
	var path types.PathOrder

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	code, err := orderRepo.CancelOrder(path.ID, path.OrderID)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserOrder returns orders of provided user
// @Summary returns orders of provided user
// @Tags users orders
// @Produce json
// @Param id path string true "User ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {array} response.UserOrder false "Orders for user"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /users/{id}/orders [get]
func (o Order) GetUserOrder(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	_, err := time.Parse(time.RFC3339, query.Date)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	userOrders, code, err := orderRepo.GetUserOrder(path.ID, query.Date)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, userOrders)
}

// GetClientOrders returns orders of provided client
// @Summary returns orders of provided client
// @Tags clients orders
// @Produce json
// @Param id path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {object} response.SummaryOrdersResponse false "Orders for clients"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id}/orders [get]
func (o Order) GetClientOrders(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	_, err := time.Parse(time.RFC3339, query.Date)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	client := types.CompanyTypesEnum.Client
	result, code, err := orderRepo.GetOrders("", path.ID, query.Date, client)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCateringClientOrders returns orders of provided client
// @Summary returns orders of provided client
// @Tags caterings orders
// @Produce json
// @Param id path string true "Catering ID"
// @Param clientId path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {object} response.SummaryOrdersResponse false "Orders for clients"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/clients/{clientId}/orders [get]
func (o Order) GetCateringClientOrders(c *gin.Context) {
	var path types.PathClient
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	_, err := time.Parse(time.RFC3339, query.Date)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	client := types.CompanyTypesEnum.Catering
	result, code, err := orderRepo.GetOrders(path.ID, path.ClientID, query.Date, client)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, result)
}

// ApproveOrders changes status of orders for provided day
// to approved
// @Summary approve user orders
// @Tags clients orders
// @Produce json
// @Param id path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 "Successfully Approved orders"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id}/orders [put]
func (o Order) ApproveOrders(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	if err := orderRepo.ApproveOrders(path.ID, query.Date); err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.Status(http.StatusOK)
}

// GetOrderStatus returns status of order
// @Summary returns status of order
// @Tags clients orders
// @Produce json
// @Param id path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {object} response.OrderStatus "order status"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id}/order-status [get]
func (o Order) GetOrderStatus(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	status := orderRepo.GetOrdersStatus(path.ID, query.Date)

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}

// GetClientOrdersExcel returns excel file of provided client
// @Summary returns excel file of orders of provided client
// @Tags clients orders
// @Produce json
// @Param id path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id}/orders-file [get]
func (o Order) GetClientOrdersExcel(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	_, err := time.Parse(time.RFC3339, query.Date)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	client := types.CompanyTypesEnum.Client
	result, code, err := orderRepo.GetOrders("", path.ID, query.Date, client)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	f := excelize.NewFile()

	dir, _ := os.Getwd()
	style, _ := f.NewStyle(`{"alignment":{"horizontal": "center", "vertical": "center"}}`)
	commentStyle, _ := f.NewStyle(`{"alignment":{"horizontal": "left", "vertical": "top", "wrap_text": true}}`)
	headers := map[string]string{"A1": "Имя", "B1": "Этаж", "C1": "Заказ", "D1": "Комментарий", "E1": "Сумма"}
	for k, v := range headers {
		f.SetCellValue("Sheet1", k, v)
		f.SetCellStyle("Sheet1", k, k, style)
	}

	f.SetColWidth("Sheet1", "A", "A", 25)
	f.SetColWidth("Sheet1", "D", "D", 25)
	f.SetColWidth("Sheet1", "B", "B", 10)
	f.SetColWidth("Sheet1", "C", "C", 25)
	f.SetColWidth("Sheet1", "E", "E", 15)
	f.SetColWidth("Sheet1", "G", "G", 15)
	f.SetColWidth("Sheet1", "H", "H", 25)
	f.SetColWidth("Sheet1", "I", "I", 10)
	f.SetColWidth("Sheet1", "K", "K", 30)
	f.SetCellValue("Sheet1", "K1", "Общая сумма заказов")
	f.SetCellValue("Sheet1", "K2", result.Total)
	f.SetCellStyle("Sheet1", "K1", "K2", style)

	for index, order := range result.UserOrders {
		var startLine int
		for i, ord := range result.UserOrders {
			if i < index {
				startLine += len(ord.Items)
			}
		}
		if index > 0 {
			start := startLine + 2
			for idx, dish := range order.Items {
				f.SetCellValue("Sheet1", "C"+strconv.Itoa(start+idx), dish.Name+" "+strconv.Itoa(dish.Amount))
				f.SetCellStyle("Sheet1", "C"+strconv.Itoa(start+idx), "C"+strconv.Itoa(start+idx), style)
			}
			st := strconv.Itoa(start)
			end := strconv.Itoa(start + len(order.Items) - 1)
			f.SetCellValue("Sheet1", "A"+st, order.Name)
			f.SetCellValue("Sheet1", "B"+st, order.Floor)

			f.SetCellValue("Sheet1", "E"+st, order.Total)
			if order.Comment != "" {
				f.SetCellValue("Sheet1", "D"+st, order.Comment)
			} else {
				f.SetCellValue("Sheet1", "D"+st, "Нет комментария к заказу")
			}
			f.MergeCell("Sheet1", "A"+st, "A"+end)
			f.MergeCell("Sheet1", "B"+st, "B"+end)
			f.MergeCell("Sheet1", "D"+st, "D"+end)
			f.MergeCell("Sheet1", "E"+st, "E"+end)
			f.SetCellStyle("Sheet1", "A"+st, "A"+end, style)
			f.SetCellStyle("Sheet1", "B"+st, "B"+end, style)
			f.SetCellStyle("Sheet1", "D"+st, "D"+end, commentStyle)
			f.SetCellStyle("Sheet1", "E"+st, "E"+end, style)
		} else {
			for idx, dish := range order.Items {
				f.SetCellValue("Sheet1", "C"+strconv.Itoa(2+idx), dish.Name+" "+strconv.Itoa(dish.Amount))
				f.SetCellStyle("Sheet1", "C"+strconv.Itoa(2+idx), "C"+strconv.Itoa(2+idx), style)
			}
			f.SetCellValue("Sheet1", "A2", order.Name)
			f.SetCellValue("Sheet1", "B2", order.Floor)
			f.SetCellValue("Sheet1", "E2", order.Total)
			if order.Comment != "" {
				f.SetCellValue("Sheet1", "D2", order.Comment)
			} else {
				f.SetCellValue("Sheet1", "D2", "Нет комментария к заказу")
			}
			end := strconv.Itoa(1 + len(order.Items))
			f.MergeCell("Sheet1", "A2", "A"+end)
			f.MergeCell("Sheet1", "B2", "B"+end)
			f.MergeCell("Sheet1", "D2", "D"+end)
			f.MergeCell("Sheet1", "E2", "E"+end)
			f.SetCellStyle("Sheet1", "A2", "A"+end, style)
			f.SetCellStyle("Sheet1", "B2", "B"+end, style)
			f.SetCellStyle("Sheet1", "D2", "D"+end, commentStyle)
			f.SetCellStyle("Sheet1", "E2", "E"+end, style)
		}
	}

	for index, order := range result.SummaryOrders {
		var startLine int
		for i, ord := range result.SummaryOrders {
			if i < index {
				startLine += len(ord.Items)
			}
		}
		if index > 0 {
			start := startLine + 1
			for idx, dish := range order.Items {
				f.SetCellValue("Sheet1", "H"+strconv.Itoa(start+idx), dish.Name)
				f.SetCellValue("Sheet1", "I"+strconv.Itoa(start+idx), dish.Amount)
				f.SetCellStyle("Sheet1", "H"+strconv.Itoa(start+idx), "H"+strconv.Itoa(start+idx), style)
				f.SetCellStyle("Sheet1", "I"+strconv.Itoa(start+idx), "I"+strconv.Itoa(start+idx), style)
			}
			st := strconv.Itoa(start)
			end := strconv.Itoa(start + len(order.Items) - 1)
			f.SetCellValue("Sheet1", "G"+st, order.CategorySummaryOrder.Name)
			f.MergeCell("Sheet1", "G"+st, "G"+end)
			f.SetCellStyle("Sheet1", "G"+st, "G"+end, style)
		} else {
			for idx, dish := range order.Items {
				f.SetCellValue("Sheet1", "H"+strconv.Itoa(1+idx), dish.Name)
				f.SetCellValue("Sheet1", "I"+strconv.Itoa(1+idx), dish.Amount)
				f.SetCellStyle("Sheet1", "H"+strconv.Itoa(1+idx), "H"+strconv.Itoa(1+idx), style)
				f.SetCellStyle("Sheet1", "I"+strconv.Itoa(1+idx), "I"+strconv.Itoa(1+idx), style)
			}
			end := strconv.Itoa(len(order.Items))
			f.SetCellValue("Sheet1", "G1", order.CategorySummaryOrder.Name)
			f.MergeCell("Sheet1", "G1", "G"+end)
			f.SetCellStyle("Sheet1", "G1", "G"+end, style)
		}
	}

	pathDir := filepath.Join(dir, "src", "static", "Book1.xlsx")

	if err := f.SaveAs(pathDir); err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.File(pathDir)
}
