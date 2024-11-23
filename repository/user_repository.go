package repository

import (
	"database/sql"
	"fmt"
	"food-delivery-apps/config"
	"food-delivery-apps/entity"
	"food-delivery-apps/shared/model"
	"math"
	"time"
)

type userRepository struct{
	db *sql.DB
}

type UserRepository interface{
	CreateNewUser(payload entity.User) (entity.UserResponse, error)
	GetAllUser(page, size int, role string) ([]entity.GetUserResponse, model.Paging, error)
	GetUserbyId(id string) (entity.UserResponse, error)
	UpdateUser(payload entity.UserResponse) (entity.UserResponse, error)
	UpdateRole(payload entity.UserResponse) (entity.UserResponse, error)
	DeleteUser(id string) error
	GetUserbyEmail(email string) (entity.User, error)
	GetUserbyUsername(username string) (entity.User, error)
	CountUser(count *int) error
	BlackListToken(token string) error
	IsTokenBlacklisted(token string) bool
	CleanUpExpiredTokens() (int64, error)
}

func (r *userRepository) CreateNewUser(payload entity.User) (entity.UserResponse, error){
	// Insert the value for user.
	if err := r.db.QueryRow(config.RegisterQuery, payload.Username, payload.Email,
		payload.Password, payload.Role, payload.Gender, payload.UpdatedAt).Scan(&payload.Id, &payload.CreatedAt); err != nil{
			return entity.UserResponse{}, fmt.Errorf("failed to create user: %v", err.Error())
		}

	// Format the timestamps for the response in a readable format.
	formattedCreatedAt := payload.CreatedAt.Format("January 02, 2006 03:04 PM")
	formattedUpdatedAt := payload.UpdatedAt.Format("January 02, 2006 03:04 PM")

	// Construct the response object with formatted data.
	response := entity.UserResponse{
		Id: payload.Id,
		Username: payload.Username,
		Email: payload.Email,
		Password: payload.Password,
		Role: payload.Role,
		Gender: payload.Gender,
		CreatedAt: formattedCreatedAt,
		UpdatedAt: formattedUpdatedAt,
	}

	return response, nil
}

func (r *userRepository) GetAllUser(page, size int, role string) ([]entity.GetUserResponse, model.Paging, error){
	var users []entity.GetUserResponse

	// Calculate the offset for pagination based on the current page and page size.
	offset := (page - 1) * size

	var rows *sql.Rows
	var err error

	// Retrieve users with pagination, otherwise include filter by role
	if role != ""{
		rows, err = r.db.Query(config.GetUserFilterQuery, size, offset, role)
	} else {
		rows, err = r.db.Query(config.GetAllUserQuery, size, offset)
	}
	if err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to retrieve user: %v", err.Error())
	}
	defer rows.Close()

	// Iterate over the rows from the database, scanning each row into a user object.
	for rows.Next(){
		var user entity.GetUserResponse

		var createdAt, updatedAt time.Time

		// Scan user data into struct fields, including timestamps for creation and update.
		if err := rows.Scan(&user.Id, &user.Username, &user.Role,
			&user.Gender, &createdAt, &updatedAt); err != nil{
				return nil, model.Paging{}, fmt.Errorf("failed to scan user: %v", err.Error())
			}
		
		// Format the timestamps for the response in a readable format.
		user.CreatedAt = createdAt.Format("January 02, 2006 03:04 PM")
		user.UpdatedAt = updatedAt.Format("January 02, 2006 03:04 PM")

		// Append the user object to the users slice.
		users = append(users, user)
	}

	// Count the total number of users to set up paging information.
	totalRows := 0
	if err := r.CountUser(&totalRows); err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to count users: %v", err.Error())
	}

	// Construct the paging object based on the total rows, page, and size.
	paging := model.Paging{
		Page: page,
		RowsPerPage: size,
		TotalRows: totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}
	
	return users, paging, nil
}

func (r *userRepository) GetUserbyId(id string) (entity.UserResponse, error){
	var user entity.UserResponse

	// Query the database to retrieve user details by ID
	err := r.db.QueryRow(config.GetUserbyIdQuery, id).Scan(&user.Id, &user.Username,
		&user.Email, &user.Password, &user.Role, &user.Gender, &user.CreatedAt, &user.UpdatedAt)

	// Handle potential errors from the query
	if err != nil{
		if err == sql.ErrNoRows{
			// If no rows are found, return a specific "user not found" error message
			return entity.UserResponse{}, fmt.Errorf("user with id %s is not found: %v", id, err.Error())
		}
		// For other errors, return a general retrieval failure message
		return entity.UserResponse{}, fmt.Errorf("failed to retrieve user: %v", err.Error())
	}

	// Parse and format the createdAt for the response in a readable format.
	parsedCreatedAt, err := time.Parse(time.RFC3339, user.CreatedAt)
	if err != nil {
		return entity.UserResponse{}, fmt.Errorf("error parsing created_at time: %v", err.Error())
	}
	formattedCreatedAt := parsedCreatedAt.Format("January 02, 2006 03:04 PM")

	// Parse and format the UpdatedAt for the response in a readable format.
	parsedUpdatedAt, err := time.Parse(time.RFC3339, user.UpdatedAt)
	if err != nil {
		return entity.UserResponse{}, fmt.Errorf("error parsing updated_at time: %v", err.Error())
	}	
	formattedUpdatedAt := parsedUpdatedAt.Format("January 02, 2006 03:04 PM")
	
	// Construct the response object with formatted data.
	response := entity.UserResponse{
		Id: user.Id,
		Username: user.Username,
		Email: user.Email,
		Password: user.Password,
		Role: user.Role,
		Gender: user.Gender,
		CreatedAt: formattedCreatedAt,
		UpdatedAt: formattedUpdatedAt,
	}
	
	return response, nil
}

func (r *userRepository) UpdateUser(payload entity.UserResponse) (entity.UserResponse, error){
	_, err := r.db.Exec(config.UpdateUserQuery, payload.Id, payload.Username,
		payload.Email, payload.Gender, payload.UpdatedAt)
	if err != nil {
		return entity.UserResponse{}, fmt.Errorf("failed to update user: %v", err.Error())
	}

	return payload, nil
}

func (r *userRepository) UpdateRole(payload entity.UserResponse) (entity.UserResponse, error){
	_, err := r.db.Exec(config.AssignToEmployeeQuery, payload.Id, payload.Role, payload.UpdatedAt)
	if err != nil{
		return entity.UserResponse{}, fmt.Errorf("failed to update role: %v", err.Error())
	}

	return payload, nil
}

func (r *userRepository) DeleteUser(id string) error{
	_, err := r.db.Exec(config.DeleteUserQuery, id)
	if err != nil{
		return fmt.Errorf("failed to delete user: %v", err.Error())
	}

	return nil
}

func (r *userRepository) GetUserbyEmail(email string) (entity.User, error){
	var user entity.User

	// Retrieve user by email
	err := r.db.QueryRow(config.LoginQuery, email).Scan(&user.Id, &user.Email, &user.Password, &user.Username, &user.Role)

	// Handle potential errors from the query
	if err != nil{
		// If no rows are found, return a specific "menu not found" error message
		if err == sql.ErrNoRows{
			return entity.User{}, fmt.Errorf("user with email %s is not found: %v", email, err.Error())
		}
		// For other errors, return a general retrieval failure message
		return entity.User{}, fmt.Errorf("failed to retrieve user: %v", err.Error())
	}
	
	return user, nil
}

func (r *userRepository) GetUserbyUsername(username string) (entity.User, error){
	var user entity.User

	// retrieve user by username
	err := r.db.QueryRow(config.GetUserbyUsernameQuery, username).Scan(&user.Username)

	// Handle potential errors from the query
	if err != nil{
		// If no rows are found, return a specific "menu not found" error message
		if err == sql.ErrNoRows{
			return entity.User{}, fmt.Errorf("user with username %s is not found: %v", username, err.Error())
		}
		// For other errors, return a general retrieval failure message
		return entity.User{}, fmt.Errorf("failed to retrieve user: %v", err.Error())
	}
		
	return user, nil
}

func (r *userRepository) CountUser(count *int) error{
	if err := r.db.QueryRow(config.CountUserQuery).Scan(count); err != nil{
		return fmt.Errorf("failed to count user")
	}

	return nil
} 

func (r *userRepository) BlackListToken(token string) error{
	expiration := time.Now().Add(time.Hour)
	_, err := r.db.Exec(config.InsertTokenQuery, token, expiration)
	if err != nil {
			return fmt.Errorf("failed insert token to blacklist: %v", err)
	}
	return nil
}

func (r *userRepository) IsTokenBlacklisted(token string) bool {
	var expiresAt time.Time
	err := r.db.QueryRow(config.GetTokenQuery, token).Scan(&expiresAt)
	if err != nil {
			return err != sql.ErrNoRows
	}
	return true
}

func (r *userRepository) CleanUpExpiredTokens() (int64, error) {
	// Get the current time
	currentTime := time.Now().Add(7 * time.Hour).UTC() // Add 7 hours to get UTC+7
	
	// delete token with the current time to remove expired tokens
	rows, err := r.db.Exec(config.DeleteTokenQuery, currentTime)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired tokens: %v", err.Error())
	}

	// Get the number of rows affected by the delete operation
	affectedRows, err := rows.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %v", err.Error())
	}

	return affectedRows, nil
}


func NewUserRepository(db *sql.DB) UserRepository{
	return &userRepository{db: db} 
}