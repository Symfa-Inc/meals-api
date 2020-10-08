package services

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Aiscom-LLC/meals-api/api/url"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/enums"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// OrderService struct
type OrderService struct{}

// NewOrderService return pointer to order struct
// with all methods
func NewOrderService() *OrderService {
	return &OrderService{}
}

var orderRepo = repository.NewOrderRepo()

func (o *OrderService) Add(query string, order models.OrderRequest, claims jwt.MapClaims) (models.UserOrder, int, error) {
	userRepo := repository.NewUserRepo()
	var userID string

	id := claims["id"].(string)

	user, _ := userRepo.GetByKey("id", id)

	if user.ID == uuid.Nil {
		userID = ""
	} else {
		userID = user.ID.String()
	}
	for i, dish := range order.Items {
		if dish.Amount == 0 {
			return models.UserOrder{}, http.StatusBadRequest, errors.New("can't add dish with 0 amount")
		}
		for j := i + 1; j < len(order.Items); j++ {
			if dish.DishID == order.Items[j].DishID {
				return models.UserOrder{}, http.StatusBadRequest, errors.New("can't add 2 same dishes, please increment amount field instead")
			}
		}
	}

	date, err := time.Parse(time.RFC3339, query)

	if err != nil {
		return models.UserOrder{}, http.StatusBadRequest, err
	}

	difference := date.Sub(time.Now().Truncate(time.Hour * 24)).Hours()

	if difference < 0 {
		return models.UserOrder{}, http.StatusBadRequest, errors.New("can't add order to previous date")
	}

	userOrder, err := orderRepo.Add(userID, date, order)

	if err != nil {
		return models.UserOrder{}, http.StatusBadRequest, err
	}

	return userOrder, 0, nil
}

func (o *OrderService) GetClientOrdersExcel(path url.PathID, query url.DateQuery) (string, int, error) {
	client := enums.CompanyTypesEnum.Client
	result, code, err := orderRepo.GetOrders("", path.ID, query.Date, client)

	if err != nil {
		return "", code, err
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

	// nolint:errcheck
	f.SaveAs(pathDir)

	return pathDir, code, err
}
