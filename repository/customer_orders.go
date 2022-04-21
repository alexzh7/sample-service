package repository

import "github.com/alexzh7/sample-service/models"

// GetOrders gets orders for provided customer id
// TODO: pass models.Customer to fill orders? o_O
func (p *pgRepo) GetOrders(customerId int) []*models.Order {
	return nil
}

// GetOrder gets order by order id
// TODO: API provides only CUSTOMERS orders, not all
// Возможно вынести методы работы с базой в закрытые методы и
// как-то дать доступ к ним через общий обьект..?
func (p *pgRepo) GetOrder(orderId int) *models.Order {
	return nil
}

// AddOrder returning id, sum

// DeleteOrder
