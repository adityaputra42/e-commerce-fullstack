package services

import (
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"errors"

	"gorm.io/gorm"
)

type AddressService interface {
	CreateAddress(userId int64, param models.CreateAddress) (*models.AddressResponse, error)
	UpdateAddress(id, userId int64, param models.UpdateAddress) (*models.AddressResponse, error)
	FindAllAddress(param models.AddressListRequest) (*models.AddressListResponse, error)
	FindById(id, userId int64) (*models.AddressResponse, error)
	DeleteAddress(id int64) error
}
type AddressServiceImpl struct {
	addressRepo repository.AddressRepository
	userRepo    repository.UserRepository
}

// CreateAddress implements AddressService.
func (a *AddressServiceImpl) CreateAddress(userId int64, param models.CreateAddress) (*models.AddressResponse, error) {
	var result models.AddressResponse
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		user, err := a.userRepo.FindById(userId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("user not found")
			}
			return err
		}

		address := models.Address{
			UserID:               int64(user.ID),
			RecipientName:        param.RecipientName,
			RecipientPhoneNumber: param.RecipientPhoneNumber,
			FullAddress:          param.FullAddress,
			Village:              param.Village,
			PostalCode:           param.PostalCode,
			District:             param.District,
			City:                 param.City,
			Province:             param.Province,
		}

		newAddress, err := a.addressRepo.Create(address, tx)

		result = *newAddress.ToResponseAddress()

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteAddress implements AddressService.
func (a *AddressServiceImpl) DeleteAddress(id int64) error {
	address, err := a.addressRepo.FindById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("address not found")
		}
		return err
	}

	err = a.addressRepo.Delete(address)

	if err != nil {
		return err
	}
	return nil
}

// FindAllAddress implements AddressService.
func (a *AddressServiceImpl) FindAllAddress(param models.AddressListRequest) (*models.AddressListResponse, error) {
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}
	addresses, err := a.addressRepo.FindAll(param)
	if err != nil {
		return nil, err
	}

	total, err := a.addressRepo.CountByUser(int64(*param.UserId))
	if err != nil {
		return nil, err
	}

	var addressResponses []models.AddressResponse
	for _, addr := range addresses {
		addressResponses = append(addressResponses, *addr.ToResponseAddress())
	}
	totalPages := int((total + int64(param.Limit) - 1) / int64(param.Limit))

	return &models.AddressListResponse{
		Addresses:  addressResponses,
		Page:       param.Page,
		Limit:      param.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (a *AddressServiceImpl) FindById(id, userId int64) (*models.AddressResponse, error) {
	address, err := a.addressRepo.FindById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("address not found")
		}
		return nil, err
	}

	if address.UserID != userId {
		return nil, errors.New("unauthorized: you can only access your own address")
	}

	return address.ToResponseAddress(), nil
}

// UpdateAddress implements AddressService.
func (a *AddressServiceImpl) UpdateAddress(id, userID int64, param models.UpdateAddress) (*models.AddressResponse, error) {
	var result models.AddressResponse

	err := database.DB.Transaction(func(tx *gorm.DB) error {

		address, err := a.addressRepo.FindById(uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("address not found")
			}
			return err
		}

		if address.UserID != userID {
			return errors.New("unauthorized: you can only update your own address")
		}

		paramAddress := models.Address{
			ID:                   address.ID,
			UserID:               address.UserID,
			RecipientName:        param.RecipientName,
			RecipientPhoneNumber: param.RecipientPhoneNumber,
			FullAddress:          param.FullAddress,
			Village:              param.Village,
			PostalCode:           param.PostalCode,
			District:             param.District,
			City:                 param.City,
			Province:             param.Province,
		}
		updatedAddress, err := a.addressRepo.Update(paramAddress, tx)

		if err != nil {
			return err
		}

		result = *updatedAddress.ToResponseAddress()

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func NewAddressService(addressRepo repository.AddressRepository,
	userRepo repository.UserRepository) AddressService {
	return &AddressServiceImpl{}
}
