package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dbo-test/internal/model"
	"github.com/google/uuid"
)

func (oc *orderController) GetOrderDetail(ctx context.Context, orderID, customerID int) (*model.Order, error) {
	order, err := oc.repo.GetOrderByID(ctx, orderID, customerID)
	if err != nil {
		return nil, err
	}

	orderDetailList, err := oc.repo.GetOrderDetailList(ctx, orderID, customerID)
	if err != nil {
		return nil, err
	}
	var totalAmount float64
	for _, v := range orderDetailList {
		totalAmount += v.TotalPrice
	}
	order.ProductList = orderDetailList
	order.TotalAmount = totalAmount

	return order, nil
}

type OrderDetailRequest struct {
	ProductID int
	Quantity  int
}

func (oc *orderController) CreateOrder(ctx context.Context, customerID int, req []OrderDetailRequest) error {
	customer, err := oc.repo.GetCustomerByID(ctx, customerID)
	if err != nil {
		return err
	}

	// create transaction to keep consistency
	trx, err := oc.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			trx.Rollback()
		}
	}()

	generateUUID := uuid.New()
	currentTime := time.Now()
	invoice := fmt.Sprintf("INV/%v/%v/%v/%s", currentTime.Day(), int(currentTime.Month()), currentTime.Year(), generateUUID.String())
	orderID, err := oc.repo.CreateOrder(ctx, trx, customer.ID, invoice)
	if err != nil {
		return err
	}

	for _, order := range req {
		product, err := oc.repo.GetProductDetailByID(ctx, order.ProductID)
		if err != nil {
			return err
		}

		if order.Quantity > product.Stock {
			return fmt.Errorf("product_id %d out of stock, remaining : %d", product.ID, product.Stock)
		}

		orderProduct := model.OrderProduct{
			ProductID:  product.ID,
			Quantity:   order.Quantity,
			TotalPrice: product.Price * float64(order.Quantity),
		}

		err = oc.repo.CreateOrderDetail(ctx, trx, orderID, &orderProduct)
		if err != nil {
			return err
		}

		err = oc.repo.UpdateProductStock(ctx, trx, product.ID, product.Stock-order.Quantity)
		if err != nil {
			return err
		}
	}

	errCommit := trx.Commit()
	if errCommit != nil {
		return errCommit
	}

	return nil
}

func (oc *orderController) DeleteOrder(ctx context.Context, orderID, customerID int) error {
	order, err := oc.repo.GetOrderByID(ctx, orderID, customerID)
	if err != nil {
		return err
	}

	trx, err := oc.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			trx.Rollback()
		}
	}()

	if order.Status != model.OrderStatusFinished {
		orderList, err := oc.repo.GetOrderDetailList(ctx, orderID, customerID)
		if err != nil {
			return err
		}

		for _, orderDetail := range orderList {
			product, err := oc.repo.GetProductDetailByID(ctx, orderDetail.ProductID)
			if err != nil {
				return err
			}

			if err := oc.repo.UpdateProductStock(ctx, trx, product.ID, product.Stock+orderDetail.Quantity); err != nil {
				return err
			}
		}
	}

	if err := oc.repo.DeleteOrderDetail(ctx, trx, orderID); err != nil {
		return err
	}

	if err := oc.repo.DeleteOrder(ctx, trx, orderID); err != nil {
		return err
	}

	errCommit := trx.Commit()
	if errCommit != nil {
		return errCommit
	}

	return nil
}

func (oc *orderController) UpdateOrder(ctx context.Context, orderID, customerID int, req []OrderDetailRequest) error {
	orderDetailList, err := oc.repo.GetOrderDetailList(ctx, orderID, customerID)
	if err != nil {
		return err
	}

	requestMapToProductID := make(map[int]int, len(req))
	for _, v := range req {
		requestMapToProductID[v.ProductID] = v.Quantity
	}

	trx, err := oc.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			trx.Rollback()
		}
	}()

	for _, order := range orderDetailList {
		newQty, exist := requestMapToProductID[order.ProductID]
		if !exist {
			continue
		}

		product, err := oc.repo.GetProductDetailByID(ctx, order.ProductID)
		if err != nil {
			return err
		}

		orderProduct := model.OrderProduct{
			ProductID:  product.ID,
			Quantity:   newQty,
			TotalPrice: product.Price * float64(newQty),
		}

		err = oc.repo.UpdateOrderDetail(ctx, trx, orderID, &orderProduct)
		if err != nil {
			return err
		}

		err = oc.repo.UpdateProductStock(ctx, trx, product.ID, product.Stock-(newQty-order.Quantity))
		if err != nil {
			return err
		}
	}

	errCommit := trx.Commit()
	if errCommit != nil {
		return errCommit
	}

	return nil
}

func (oc *orderController) SearchOrders(ctx context.Context, customerID int, queryParameters map[string]interface{}) ([]model.Order, bool, error) {
	// adjust pagination
	var (
		page    int
		perPage int
	)
	pageRaw := queryParameters["page"]
	page, valid := pageRaw.(int)
	if !valid {
		return nil, false, errors.New("invalid per_page data type")
	}

	perPageRaw := queryParameters["per_page"]
	perPage, valid = perPageRaw.(int)
	if !valid {
		return nil, false, errors.New("invalid per_page data type")
	}

	page = (page - 1) * perPage
	perPage++
	queryParameters["per_page"] = perPage
	queryParameters["page"] = page

	orders, err := oc.repo.SearchOrder(ctx, customerID, queryParameters)
	if err != nil {
		return nil, false, err
	}

	// check for next page pagination
	hasNext := false
	if len(orders) >= perPage {
		hasNext = true
		orders = orders[:len(orders)-1]
	}

	for idx, order := range orders {
		orderDetailList, err := oc.repo.GetOrderDetailList(ctx, order.ID, order.CustomerID)
		if err != nil {
			return nil, false, err
		}
		var totalAmount float64
		orders[idx].ProductList = orderDetailList
		for _, orderDetail := range orderDetailList {
			totalAmount += orderDetail.TotalPrice
		}
		orders[idx].TotalAmount = totalAmount
	}

	return orders, hasNext, nil
}
