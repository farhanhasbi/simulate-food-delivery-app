package repository

import (
	"database/sql"
	"fmt"
	"food-delivery-apps/config"
	"food-delivery-apps/entity"
	"food-delivery-apps/shared/model"
	"math"
	"time"

	"github.com/lib/pq"
)

type menuRepository struct{
	db *sql.DB
}

type MenuRepository interface{
	AddMenu(payload entity.Menu) (entity.MenuResponse, error)
	GetAllMenu(page, size int, mtype, mname string) ([]entity.MenuResponse, model.Paging, error)
	GetMenubyId(id string) (entity.MenuResponse, error)
	UpdateMenu(payload entity.MenuResponse) (entity.MenuResponse, error)
	DeleteMenu(id string) error
	GetMenubyName(name string) (entity.Menu, error)
}

func (r *menuRepository) AddMenu(payload entity.Menu) (entity.MenuResponse, error){
	// Insert the value for menus.
	err := r.db.QueryRow(config.CreateMenuQuery, payload.Name, payload.Type,
		payload.Desc, payload.UnitType, payload.Price, payload.CreatedBy,
		payload.UpdatedAt).Scan(&payload.Id, &payload.CreatedAt, &payload.CreatedBy)
	
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
				case "23505": // Unique violation
				if pqErr.Constraint == "unique_menu_name" {
					return entity.MenuResponse{}, fmt.Errorf("menu with name %s already exists", payload.Name)
				}
				if pqErr.Constraint == "unique_menu_description" {
					return entity.MenuResponse{}, fmt.Errorf("menu with description '%s' already exists", payload.Desc)
				}
			}
		}
		return entity.MenuResponse{}, fmt.Errorf("failed to create new menu: %v", err.Error())
	}

	// Retrieve the customer's username based on CustomerId.
	var username string
	query := "SELECT username FROM users WHERE id = $1"
	if err := r.db.QueryRow(query, payload.CreatedBy).Scan(&username); err != nil{
		return entity.MenuResponse{}, fmt.Errorf("failed to retrive username: %v", err.Error())
	}	
	payload.CreatedBy = username

	// Format CreatedAt and UpdatedAt for the response in a readable format.
	formattedCreatedAt := payload.CreatedAt.Format("January 02, 2006 03:04 PM")
	formattedUpdatedAt := payload.UpdatedAt.Format("January 02, 2006 03:04 PM")

	// Construct the response object with formatted data.
	response := entity.MenuResponse{
		Id: payload.Id,
		Name: payload.Name,
		Type: payload.Type,
		Desc: payload.Desc,
		UnitType: payload.UnitType,
		Price: payload.Price,
		Rating: payload.Rating,
		CreatedBy: payload.CreatedBy,
		CreatedAt: formattedCreatedAt,
		UpdatedAt: formattedUpdatedAt,
	}

	return response, nil
}

func (r *menuRepository) GetAllMenu(page, size int, mtype, mname string) ([]entity.MenuResponse, model.Paging, error){
	var menus []entity.MenuResponse

	// Calculate the offset for pagination based on the current page and page size.
	offset := (page - 1) *size

	var rows *sql.Rows
	var err error
	
	// Retrieve the Menus with pagination, otherwise include the filter by name and type
	if mtype != "" && mname != ""{
		rows, err = r.db.Query(config.GetAllMenuWithAllFilterQuery, size, offset, mtype, mname)
	} else if mname != ""{
		rows, err = r.db.Query(config.GetAllMenuWithFilterNameQuery, size, offset, mname)
	} else if mtype != ""{
		rows, err = r.db.Query(config.GetAllMenuWithFilterTypeQuery, size, offset, mtype)
	} else {
		rows, err = r.db.Query(config.GetAllMenuQuery, size, offset)
	}

	if err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to retrieve menu: %v", err.Error())
	}
	defer rows.Close()

	// Iterate over the rows from the database, scanning each row into a menu object.
	for rows.Next(){
		var menu entity.MenuResponse
		var createdAt, updateAt time.Time

		// Scan menu data into struct fields, including timestamps for creation and update.
		if err := rows.Scan(&menu.Id, &menu.Name, &menu.Type, &menu.Desc, &menu.UnitType,
			&menu.Price, &menu.Rating, &menu.CreatedBy, &createdAt, &updateAt); err != nil{
				 return nil, model.Paging{}, fmt.Errorf("failed to scan menu: %v", err.Error())
			}

		// Format the timestamps for the response in a readable format.
		menu.CreatedAt = createdAt.Format("January 02, 2006 03:04 PM")
		menu.UpdatedAt = updateAt.Format("January 02, 2006 03:04 PM")

		// Append the menu object to the menus slice.
		menus = append(menus, menu)
	}

	// Count the total number of menus to set up paging information.
	totalRowsMenu := 0
	if err := r.db.QueryRow(config.CountMenuQuery).Scan(&totalRowsMenu); err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to count menu: %v", err.Error())
	}
	
	// Construct the paging object based on the total rows, page, and size.
	paging := model.Paging{
		Page: page,
		RowsPerPage: size,
		TotalRows: totalRowsMenu,
		TotalPages: int(math.Ceil(float64(totalRowsMenu) / float64(size))),
	}

	return menus, paging, nil
}

func (r *menuRepository) GetMenubyId(id string) (entity.MenuResponse, error){
	var menu entity.MenuResponse

	// Retrieve menu by id
	err := r.db.QueryRow(config.GetMenubyIdQuery, id).Scan(&menu.Id, &menu.Name,
		&menu.Type, &menu.Desc, &menu.UnitType, &menu.Price, &menu.CreatedBy, &menu.CreatedAt, &menu.UpdatedAt)

	// Handle potential errors from the query
	if err != nil{
		// If no rows are found, return a specific "menu not found" error message
		if err == sql.ErrNoRows{
			return entity.MenuResponse{}, fmt.Errorf("menu with id %s is not found: %v", id, err.Error())
		}
		// For other errors, return a general retrieval failure message
		return entity.MenuResponse{}, fmt.Errorf("failed to retrieve menu: %v", err.Error())
	}

	// Parse and format the createdAt for the response in a readable format.
	parsedCreatedAt, err := time.Parse(time.RFC3339, menu.CreatedAt)
	if err != nil {
		return entity.MenuResponse{}, fmt.Errorf("error parsing created_at time: %v", err.Error())
	}
	formattedCreatedAt := parsedCreatedAt.Format("January 02, 2006 03:04 PM")

	// Parse and format the updatedAt for the response in a readable format.
	parseUpdatedAt, err := time.Parse(time.RFC3339, menu.UpdatedAt)
	if err != nil {
		return entity.MenuResponse{}, fmt.Errorf("error parsing updated_at time: %v", err.Error())
	}
	formattedUpdatedAt := parseUpdatedAt.Format("January 02, 2006 03:04 PM")

	// Construct the response object with formatted data.
	response := entity.MenuResponse{
		Id: menu.Id,
		Name: menu.Name,
		Type: menu.Type,
		Desc: menu.Desc,
		UnitType: menu.UnitType,
		Price: menu.Price,
		Rating: menu.Rating,
		CreatedBy: menu.CreatedBy,
		CreatedAt: formattedCreatedAt,
		UpdatedAt: formattedUpdatedAt,
	}

	return response, nil
}

func (r *menuRepository) UpdateMenu(payload entity.MenuResponse) (entity.MenuResponse, error){
	_, err := r.db.Exec(config.UpdateMenuQuery, payload.Id, payload.Name, payload.Type,
		payload.Desc, payload.UnitType, payload.Price, payload.UpdatedAt)
	if err != nil{
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
				case "23505": // Unique violation
				if pqErr.Constraint == "unique_menu_name" {
					return entity.MenuResponse{}, fmt.Errorf("menu with name %s already exists", payload.Name)
				}
				if pqErr.Constraint == "unique_menu_description" {
					return entity.MenuResponse{}, fmt.Errorf("menu with description '%s' already exists", payload.Desc)
				}
			}
		}
		return entity.MenuResponse{}, fmt.Errorf("failed to update menu: %v", err.Error())
	}

	// Retrieve the customer's username based on CustomerId.
	var username string
	query := "SELECT username FROM users WHERE id = $1"
	if err := r.db.QueryRow(query, payload.CreatedBy).Scan(&username); err != nil{
		return entity.MenuResponse{}, fmt.Errorf("failed to retrieve username: %v", err.Error())
	}
	payload.CreatedBy = username

	return payload, nil
}

func (r *menuRepository) DeleteMenu(id string) error{
	_, err := r.db.Exec(config.DeleteMenuQuery, id)
	if err != nil{
		return fmt.Errorf("failed to delete menu: %v", err.Error())
	}

	return nil
}

func (r *menuRepository) GetMenubyName(name string) (entity.Menu, error){
	var menu entity.Menu

	// Retrieve menu by name
	err := r.db.QueryRow(config.GetMenubyNameQuery, name).Scan(&menu.Id, &menu.Name, &menu.Price)

	// Handle potential errors from the query
	if err != nil{
		// If no rows are found, return a specific "menu not found" error message
		if err == sql.ErrNoRows{
			return entity.Menu{}, fmt.Errorf("menu with name %s is not found: %v", name, err.Error())
		}
		// For other errors, return a general retrieval failure message
		return entity.Menu{}, fmt.Errorf("failed to retrieve menu: %v", err.Error())
	}

	return menu, nil
}

func NewMenuRepository(db *sql.DB) MenuRepository{
	return &menuRepository{db: db}
}


