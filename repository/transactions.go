package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"vandesar/entity"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db}
}

func (c *TransactionRepository) AddTrans(ctx context.Context, trans entity.TransactionReq) []error {
	var errs []error

	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cartList, _ := json.Marshal(trans.CartList)

		transaction := entity.Transaction{
			UserID:      trans.UserID,
			Debt:        trans.Debt,
			Status:      trans.Status,
			Money:       trans.Money,
			CartList:    cartList,
			TotalPrice:  trans.TotalPrice,
			Notes:       trans.Notes,
			TotalProfit: trans.TotalProfit,
		}

		err := tx.Create(&transaction).Error
		if err != nil {
			return err
		}

		for _, v := range trans.CartList {
			var product entity.Product
			err = tx.Table("products").Where("id = ?", v.ProductID).First(&product).Error
			if err != nil {
				return err
			}

			if product.Stock < v.Quantity {
				errs = append(errs, fmt.Errorf("stock for product id %d not enough", product.ID))
				continue
			}

			err = tx.Table("products").Where("id = ?", v.ProductID).Update("stock", gorm.Expr("stock - ?", v.Quantity)).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (c *TransactionRepository) UpdateTrans(trans entity.TransactionReq, tranId uint) (entity.Transaction, error) {
	var result entity.Transaction

	cartList, _ := json.Marshal(trans.CartList)
	transaction := entity.Transaction{
		Model: gorm.Model{
			ID: tranId,
		},
		UserID:      trans.UserID,
		Debt:        trans.Debt,
		Status:      trans.Status,
		Money:       trans.Money,
		CartList:    cartList,
		TotalPrice:  trans.TotalPrice,
		Notes:       trans.Notes,
		TotalProfit: trans.TotalProfit,
	}

	err := c.db.
		Table("transactions").
		Where("id = ?", tranId).
		Updates(&transaction).
		First(&result).Error
	if err != nil {
		return entity.Transaction{}, err
	}

	return result, nil
}

func (c *TransactionRepository) DeleteTrans(id uint) error {
	err := c.db.Delete(&entity.Transaction{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *TransactionRepository) ReadTransByCashier(userId uint) ([]entity.TransactionReq, error) {
	var transactions []entity.Transaction
	err := c.db.
		Table("transactions").
		Select("*").
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL").
		Scan(&transactions).Error

	if err != nil {
		return nil, err
	}

	resp := make([]entity.TransactionReq, 0, len(transactions))

	for _, v := range transactions {
		var cartList []entity.Prods
		_ = json.Unmarshal(v.CartList, &cartList)

		resp = append(resp, entity.TransactionReq{
			UserID:      v.UserID,
			Debt:        v.Debt,
			Status:      v.Status,
			Money:       v.Money,
			CartList:    cartList,
			TotalPrice:  v.TotalPrice,
			Notes:       v.Notes,
			TotalProfit: v.TotalProfit,
		})
	}

	return resp, nil
}

func (c *TransactionRepository) ReadTransByDateRange(startDate, endDate time.Time, adminId uint) ([]entity.TransactionReq, error) {
	var transactions []entity.Transaction
	err := c.db.Debug().
		Table("transactions").
		Select("transactions.*").
		Where("transactions.created_at BETWEEN ? AND ?", startDate, endDate).
		Joins("JOIN cashiers ON cashiers.id = transactions.user_id").
		Where("cashiers.admin_id = ?", adminId).
		Where("transactions.deleted_at IS NULL").
		Scan(&transactions).Error
	if err != nil {
		return nil, err
	}

	fmt.Println("=======")
	fmt.Println(transactions)
	fmt.Println("=======")

	resp := make([]entity.TransactionReq, 0, len(transactions))

	for _, v := range transactions {
		var cartList []entity.Prods
		_ = json.Unmarshal(v.CartList, &cartList)

		resp = append(resp, entity.TransactionReq{
			UserID:      v.UserID,
			Debt:        v.Debt,
			Status:      v.Status,
			Money:       v.Money,
			CartList:    cartList,
			TotalPrice:  v.TotalPrice,
			Notes:       v.Notes,
			TotalProfit: v.TotalProfit,
		})
	}

	return resp, nil
}

func (c *TransactionRepository) ReadTransByAdmin(adminId uint) ([]entity.TransactionReq, error) {
	var transactions []entity.Transaction
	err := c.db.
		Table("transactions").
		Select("*").
		Joins("JOIN cashiers ON cashiers.id = transactions.user_id").
		Where("cashiers.admin_id = ?", adminId).
		Where("transactions.deleted_at IS NULL").
		Scan(&transactions).Error
	if err != nil {
		return nil, err
	}

	resp := make([]entity.TransactionReq, 0, len(transactions))

	for _, v := range transactions {
		var cartList []entity.Prods
		_ = json.Unmarshal(v.CartList, &cartList)

		resp = append(resp, entity.TransactionReq{
			UserID:      v.UserID,
			Debt:        v.Debt,
			Status:      v.Status,
			Money:       v.Money,
			CartList:    cartList,
			TotalPrice:  v.TotalPrice,
			Notes:       v.Notes,
			TotalProfit: v.TotalProfit,
		})
	}

	return resp, nil
}

func (c *TransactionRepository) ReadTransByAdminDebt(adminId uint) ([]entity.TransactionReq, error) {
	var transactions []entity.Transaction
	debt := "hutang"
	err := c.db.
		Table("transactions").
		Select("*").
		Joins("JOIN cashiers ON cashiers.id = transactions.user_id").
		Where("cashiers.admin_id = ?", adminId).
		Where("transactions.status = ?", debt).
		Where("transactions.deleted_at IS NULL").
		Scan(&transactions).Error
	if err != nil {
		return nil, err
	}

	resp := make([]entity.TransactionReq, 0, len(transactions))

	for _, v := range transactions {
		var cartList []entity.Prods
		_ = json.Unmarshal(v.CartList, &cartList)

		resp = append(resp, entity.TransactionReq{
			UserID:      v.UserID,
			Debt:        v.Debt,
			Status:      v.Status,
			Money:       v.Money,
			CartList:    cartList,
			TotalPrice:  v.TotalPrice,
			Notes:       v.Notes,
			TotalProfit: v.TotalProfit,
		})
	}

	return resp, nil
}

// func (c *TransactionRepository) UpdateTransDebt(status *string, debt *float64, tranId uint) (entity.Transaction, error) {
// 	var result entity.Transaction
// 	statuses := "hutang"
// 	trans := entity.TransactionReq{}

// 	cartList, _ := json.Marshal(trans.CartList)
// 	transaction := entity.Transaction{
// 		Model: gorm.Model{
// 			ID: tranId,
// 		},
// 		UserID:      trans.UserID,
// 		Debt:        trans.Debt,
// 		Status:      trans.Status,
// 		Money:       trans.Money,
// 		CartList:    cartList,
// 		TotalPrice:  trans.TotalPrice,
// 		Notes:       trans.Notes,
// 		TotalProfit: trans.TotalProfit,
// 	}
// 	status = &transaction.Status
// 	debt = &transaction.Debt
// 	//db.Table("users").Where("id = ?", user.ID).Updates(map[string]interface{}{"name": user.Name, "email": user.Email})
// 	err := c.db.
// 		Table("transactions").
// 		Where("id = ?", tranId).
// 		Where("status = ?", statuses).
// 		Updates(map[string]interface{}{"status": status, "debt": debt}).
// 		First(&result).Error
// 	if err != nil {
// 		return entity.Transaction{}, err
// 	}

// 	return result, nil
// }

func (c *TransactionRepository) UpdateTransDebt(trans entity.UpdateTrans, tranId uint) (entity.Transaction, error) {
	var result entity.Transaction

	transaction := entity.Transaction{
		Model: gorm.Model{
			ID: tranId,
		},
		Debt:   trans.Debt,
		Status: trans.Status,
		Money:  trans.Money,
	}

	err := c.db.
		Table("transactions").
		Where("id = ?", tranId).
		Updates(&transaction).
		First(&result).Error
	if err != nil {
		return entity.Transaction{}, err
	}

	return result, nil
}
