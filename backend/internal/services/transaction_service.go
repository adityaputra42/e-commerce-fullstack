package services

import (
	"context"
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/utils"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, param models.CreateTransaction) (*models.TransactionResponse, error)
	UpdateTransaction(param models.UpdateTransaction) (*models.TransactionResponse, error)
	FindAllTransaction(param models.TransactionListRequest) (*[]models.TransactionResponse, error)
	FindTransactionById(txid string) (*models.TransactionResponse, error)
	deleteTransaction(txid string) error
}
type TransactionServiceImpl struct {
	transactionRepo repository.TransactionRepository
	shippingRepo    repository.ShippingRepository
	addressRepo     repository.AddressRepository
	orderRepo       repository.OrderRepository
	paymentRepo     repository.PaymentMethodRepository
	productRepo     repository.ProductRepository
}

func (t *TransactionServiceImpl) CreateTransaction(ctx context.Context, param models.CreateTransaction) (*models.TransactionResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var (
		address       *models.Address
		shipping      *models.Shipping
		paymentMethod *models.PaymentMethod
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		addr, err := t.addressRepo.FindById(uint(param.AddressID))
		if err != nil {
			return fmt.Errorf("address not found: %w", err)
		}
		address = addr
		return nil
	})

	g.Go(func() error {
		ship, err := t.shippingRepo.FindById(uint(param.ShippingID))
		if err != nil {
			return fmt.Errorf("shipping not found: %w", err)
		}
		shipping = ship
		return nil
	})

	g.Go(func() error {
		pm, err := t.paymentRepo.FindById(uint(param.PaymentMethodID))
		if err != nil {
			return fmt.Errorf("payment method not found: %w", err)
		}
		paymentMethod = pm
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	tx := database.DB.WithContext(gctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	transactionParam := models.Transaction{
		TxID:            utils.Generate("TRX-"),
		AddressID:       param.AddressID,
		ShippingID:      param.ShippingID,
		PaymentMethodID: param.PaymentMethodID,
		ShippingPrice:   param.ShippingPrice,
		Status:          utils.WaitingPayment,
	}

	transactionResult, err := t.transactionRepo.Create(transactionParam, tx)
	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	var orders []models.OrderResponse
	var total float64

	for idx, po := range param.ProductOrders {
		if err := gctx.Err(); err != nil {
			return nil, fmt.Errorf("context canceled: %w", err)
		}

		product, err := t.productRepo.FindProductById(po.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found at index %d: %w", idx, err)
		}

		colorVariant, err := t.productRepo.FindColorVarianById(po.ColorVarianID)
		if err != nil {
			return nil, fmt.Errorf("color variant not found at index %d: %w", idx, err)
		}

		sizeVariant, err := t.productRepo.FindSizeVarianLocked(tx, uint(po.SizeVarianID))
		if err != nil {
			return nil, fmt.Errorf("size variant not found at index %d: %w", idx, err)
		}

		if sizeVariant.Stock < po.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %d at index %d", po.ProductID, idx)
		}

		unitPrice := product.Price
		subtotal := float64(po.Quantity) * unitPrice

		order := &models.Order{
			ID:            utils.Generate("TXO"),
			TransactionID: transactionResult.TxID,
			ProductID:     product.ID,
			ColorVarianID: colorVariant.ID,
			SizeVarianID:  sizeVariant.ID,
			UnitPrice:     unitPrice,
			Quantity:      po.Quantity,
			Subtotal:      subtotal,
			Status:        utils.Pending,
		}

		orderResult, err := t.orderRepo.Create(*order, tx)
		if err != nil {
			return nil, fmt.Errorf("create order index %d: %w", idx, err)
		}

		sizeVariant.Stock -= po.Quantity
		if _, err := t.productRepo.UpdateSizeVarian(*sizeVariant, tx); err != nil {
			return nil, fmt.Errorf("update stock index %d: %w", idx, err)
		}

		total += subtotal
		orders = append(orders, orderResult.ToOrderResponse())
	}

	transactionResult.TotalPrice = total
	txResult, err := t.transactionRepo.Update(transactionResult, tx)
	if err != nil {
		return nil, fmt.Errorf("update transaction total: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	committed = true

	response := &models.TransactionResponse{
		TxID:          txResult.TxID,
		Address:       *address.ToResponseAddress(),
		Shipping:      *shipping.ToResponseShipping(),
		PaymentMethod: *paymentMethod.ToResponsePaymentMethod(),
		ShippingPrice: txResult.ShippingPrice,
		TotalPrice:    txResult.TotalPrice,
		Status:        txResult.Status,
		Orders:        orders,
		CreatedAt:     txResult.CreatedAt,
		UpdatedAt:     txResult.UpdatedAt,
	}

	return response, nil
}

// FindAllTransaction implements [TransactionService].
func (t *TransactionServiceImpl) FindAllTransaction(param models.TransactionListRequest) (*[]models.TransactionResponse, error) {

	transactions, err := t.transactionRepo.FindAll(param)
	if err != nil {
		return nil, err
	}

	var responses []models.TransactionResponse
	for _, transaction := range transactions {
	
		response := models.TransactionResponse{
			TxID:          transaction.TxID,
			ShippingPrice: transaction.ShippingPrice,
			TotalPrice:    transaction.TotalPrice,
			Status:        transaction.Status,
			CreatedAt:     transaction.CreatedAt,
			UpdatedAt:     transaction.UpdatedAt,
		}


		responses = append(responses, response)
	}

	return &responses, nil
}

// FindTransactionById implements [TransactionService].
func (t *TransactionServiceImpl) FindTransactionById(txid string) (*models.TransactionResponse, error) {
	transaction, err := t.transactionRepo.FindById(txid)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	var orderResponse []models.OrderResponse
	for _, v := range transaction.Orders {
		orderResponse = append(orderResponse, v.ToOrderResponse())
	}
	response := &models.TransactionResponse{
		TxID:          transaction.TxID,
		Address:       *transaction.Address.ToResponseAddress(),
		Shipping:      *transaction.Shipping.ToResponseShipping(),
		PaymentMethod: *transaction.PaymentMethod.ToResponsePaymentMethod(),
		ShippingPrice: transaction.ShippingPrice,
		TotalPrice:    transaction.TotalPrice,
		Status:        transaction.Status,
		Orders:        orderResponse,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}

	return response, nil
}

// UpdateTransaction implements [TransactionService].
func (t *TransactionServiceImpl) UpdateTransaction(param models.UpdateTransaction) (*models.TransactionResponse, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transaction, err := t.transactionRepo.FindByIdLocking(tx, param.TxID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("transaction not found")
	}

	if !utils.IsValidStatusTransition(transaction.Status, param.Status) {
		tx.Rollback()
		return nil, errors.New("invalid status transition")
	}

	transaction.Status = param.Status

	updatedTransaction, err := t.transactionRepo.Update(*transaction, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var orderResponse []models.OrderResponse
	for _, v := range updatedTransaction.Orders {
		orderResponse = append(orderResponse, v.ToOrderResponse())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	response := &models.TransactionResponse{
		TxID:          transaction.TxID,
		Address:       *transaction.Address.ToResponseAddress(),
		Shipping:      *transaction.Shipping.ToResponseShipping(),
		PaymentMethod: *transaction.PaymentMethod.ToResponsePaymentMethod(),
		ShippingPrice: transaction.ShippingPrice,
		TotalPrice:    transaction.TotalPrice,
		Status:        transaction.Status,
		Orders:        orderResponse,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}

	return response, nil
}

// deleteTransaction implements [TransactionService].
func (t *TransactionServiceImpl) deleteTransaction(txid string) error {
	panic("unimplemented")
}

func NewTransactionService(TransactionRepo repository.TransactionRepository,
	ShippingRepo repository.ShippingRepository,
	AddressRepo repository.AddressRepository,
	PaymentRepo repository.PaymentMethodRepository,
	OrderRepo repository.OrderRepository,
	ProductRepo repository.ProductRepository) TransactionService {
	return &TransactionServiceImpl{
		transactionRepo: TransactionRepo,
		shippingRepo:    ShippingRepo,
		addressRepo:     AddressRepo,
		orderRepo:       OrderRepo,
		paymentRepo:     PaymentRepo,
		productRepo:     ProductRepo,
	}
}
