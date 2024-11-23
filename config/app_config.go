package config

const ApiGroup = "/api/v1"

// User Route
const (
	Register   = "/auth/register"
	Login      = "/auth/login"
	Logout     = "/auth/logout"
	Role       = "/user/:id/role"
	DeleteUser = "/user/:id"
	UpdateUser = "/user"
	GetAllUser = "/user"
)

// Menu Route
const (
	AddMenu    = "/menu"
	GetMenu    = "/menu"
	UpdateMenu = "/menu/:id"
	DeleteMenu = "/menu/:id"
)

// balance Route
const (
	CreateBalance = "/balance"
	GetBalance    = "/balance"
)

// Order Route
const (
	AddOrder          = "/order"
	GetUnfinishOrder  = "/unfinish-order"
	UpdateOrderStatus = "/order-status/:id"
	GetAllOrder       = "/order"
	GetOrderHistory   = "/finish-order"
)

// Promo Route
const (
	AddPromo     = "/promo"
	GetPromo     = "/promo"
	GetPromoCust = "/available-promo"
	DeletePromo  = "/promo/:id"
)

// Review Route
const (
	AddReview    = "/review"
	GetReview    = "/review"
	UpdateReview = "/review/:id"
	DeleteReview = "/review/:id"
)