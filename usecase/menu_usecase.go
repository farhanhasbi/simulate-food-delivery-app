package usecase

import (
	"food-delivery-apps/entity"
	"food-delivery-apps/repository"
	"food-delivery-apps/shared/model"
	"time"
)

type menuUseCase struct{
	repo repository.MenuRepository
}

type MenuUseCase interface{
	CreateNewMenu(payload entity.Menu) (entity.MenuResponse, error)
	GetAllMenu(page, size int, mtype, mname string) ([]entity.MenuResponse, model.Paging, error)
	UpdateMenu(payload entity.Menu) (entity.MenuResponse, error)
	DeleteMenu(id string) error
}

func (uc *menuUseCase) CreateNewMenu(payload entity.Menu) (entity.MenuResponse, error){
	// Validate the fields provided in the payload
	if err := payload.Validate(); err != nil{
		return entity.MenuResponse{}, err
	}
	
	payload.UpdatedAt = time.Now()

	return uc.repo.AddMenu(payload)
}

func (uc *menuUseCase) GetAllMenu(page, size int, mtype, mname string) ([]entity.MenuResponse, model.Paging, error){
	return uc.repo.GetAllMenu(page, size, mtype, mname)
}

func (uc *menuUseCase) UpdateMenu(payload entity.Menu) (entity.MenuResponse, error){
	// Retrieve the current menu by id
	menu, err := uc.repo.GetMenubyId(payload.Id)
	if err != nil{
		return entity.MenuResponse{}, err
	}

	// Validate the fields provided in the payload
	if err := payload.ValidateUpdate(); err != nil{
		return entity.MenuResponse{}, err
	}

	// Check if fields are present before updating them
	if payload.Name != ""{
		menu.Name = payload.Name
	}
	if payload.Type != ""{
		menu.Type = payload.Type
	}
	if payload.Desc != ""{
		menu.Desc = payload.Desc
	}
	if payload.Price != 0{
		menu.Price = payload.Price 
	}
	
	menu.UpdatedAt = time.Now().Format("January 02, 2006 03:04 PM")

	return uc.repo.UpdateMenu(menu)
}

func (uc *menuUseCase) DeleteMenu(id string) error{
	// Retrieve the current menu by id
	_, err := uc.repo.GetMenubyId(id)
	if err != nil{
		return err
	}

	return uc.repo.DeleteMenu(id)
}

func NewMenuUseCase(repo repository.MenuRepository) MenuUseCase{
	return &menuUseCase{repo: repo}
}

// parseUpdateTime, err := time.Parse(time.RFC3339, menu.UpdatedAt)
	// parseUpdateTime = time.Now()
	// if err != nil{
	// 	return entity.MenuResponse{}, fmt.Errorf("error parsing time: %v", err.Error())
	// }
	
	// formattedUpdatedAt := parseUpdateTime.Format("January 02, 2006 03:04 PM")