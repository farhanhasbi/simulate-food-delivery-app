package config

// User Query
const (
	RegisterQuery = `INSERT INTO users(username, email, password, role, gender, updated_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	LoginQuery             = `SELECT id, email, password, username, role FROM users WHERE email = $1`
	GetUserbyUsernameQuery = `SELECT username FROM users WHERE username = $1`
	GetUserbyIdQuery       = `SELECT id, username, email, password, role, gender, created_at, updated_at FROM users WHERE id = $1`
	CountUserQuery         = `SELECT COUNT(*) FROM users`
	AssignToEmployeeQuery  = `UPDATE users SET role = $2, updated_at = $3 WHERE id = $1`
	DeleteUserQuery        = `DELETE FROM users WHERE id = $1`
	GetAllUserQuery     = `SELECT id, username, role, gender, created_at, updated_at FROM users ORDER BY created_at ASC limit $1 OFFSET $2`
	GetUserFilterQuery  = `SELECT id, username, role, gender, created_at, updated_at FROM users WHERE ROLE = $3 ORDER BY created_at ASC limit $1 OFFSET $2`
	UpdateUserQuery        = `UPDATE users SET username = $2, email = $3, gender = $4, updated_at = $5 WHERE id = $1`
)

// Token Query
const (
	InsertTokenQuery = `INSERT INTO token_blacklists(token, expires_at) VALUES ($1 , $2)`
	GetTokenQuery = `SELECT expires_at FROM token_blacklists WHERE token = $1`
	DeleteTokenQuery = `DELETE FROM token_blacklists WHERE expires_at < $1`
)

// Menu Query
const (
	CreateMenuQuery = `INSERT INTO menus(name, type, description, unit_type, price, created_by, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, created_by`
	GetMenubyNameQuery = "SELECT id, name, price FROM menus WHERE name = $1"
	GetMenubyDescQuery = "SELECT description FROM menus WHERE description = $1"
	GetAllMenuQuery = `
	SELECT m.id, m.name, m.type, m.description, m.unit_type, m.price,
	COALESCE(AVG(r.rating), 0) AS rating, u.username AS created_by,
	m.created_at, m.updated_at FROM menus m
	JOIN users u on m.created_by = u.id
	LEFT JOIN reviews r on m.id = r.menu_id
	GROUP BY m.id, u.username
	ORDER BY rating DESC, created_at ASC
	LIMIT $1 OFFSET $2`
	GetAllMenuWithAllFilterQuery = `SELECT m.id, m.name, m.type, m.description, m.unit_type, m.price,
	COALESCE(AVG(r.rating), 0) AS rating, u.username AS created_by,
	m.created_at, m.updated_at FROM menus m
	JOIN users u on m.created_by = u.id
	LEFT JOIN reviews r on m.id = r.menu_id
	WHERE m.type = $3 AND m.name LIKE '%' || $4 || '%'
	GROUP BY m.id, u.username
	ORDER BY rating DESC, created_at ASC
	LIMIT $1 OFFSET $2`
	GetAllMenuWithFilterNameQuery = `SELECT m.id, m.name, m.type, m.description, m.unit_type, m.price,
	COALESCE(AVG(r.rating), 0) AS rating, u.username AS created_by,
	m.created_at, m.updated_at FROM menus m
	JOIN users u ON m.created_by = u.id
	LEFT JOIN reviews r ON m.id = r.menu_id
	WHERE m.name LIKE '%' || $3 || '%'
	GROUP BY m.id, u.username
	ORDER BY rating DESC, created_at ASC
	LIMIT $1 OFFSET $2`
	GetAllMenuWithFilterTypeQuery = `SELECT m.id, m.name, m.type, m.description, m.unit_type, m.price,
	COALESCE(AVG(r.rating), 0) AS rating, u.username AS created_by,
	m.created_at, m.updated_at FROM menus m
	JOIN users u ON m.created_by = u.id
	LEFT JOIN reviews r ON m.id = r.menu_id
	WHERE m.type = $3
	GROUP BY m.id, u.username
	ORDER BY rating DESC, created_at ASC
	LIMIT $1 OFFSET $2`
	GetMenubyIdQuery = `SELECT id, name, type, description, unit_type, price, created_by, created_at, updated_at FROM menus WHERE id = $1`
	UpdateMenuQuery = `UPDATE menus SET name = $2, type = $3, description = $4, unit_type = $5, price = $6, updated_at = $7 WHERE id = $1`
	DeleteMenuQuery = "DELETE FROM menus WHERE id = $1"
	CountMenuQuery = `SELECT COUNT(*) FROM menus`
)

// Balance Query
const (
	CreateBalanceQuery = `INSERT INTO balances(customer_id, transaction_type, amount, description, balance) VALUES($1, $2, $3, $4, $5) RETURNING id, balance, created_at`
	GetUserBalanceQuery = `SELECT balance FROM balances WHERE customer_id = $1 ORDER BY created_at DESC LIMIT 1`
	GetAllUserBalanceQuery = `SELECT b.id, u.username AS customer_name, b.transaction_type, b.amount, b.description, b.balance, b.created_at FROM balances b JOIN users u ON b.customer_id = u.id WHERE b.customer_id = $3 ORDER BY b.created_at ASC LIMIT $1 OFFSET $2`
	CountUserBalanceQuery = "SELECT COUNT (*) FROM balances WHERE customer_id = $1"
)

// Promo Query
const (
	CreatePromoQuery = `INSERT INTO promos(employee_id, promo_code, discount, is_percentage, start_date, end_date, description, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`
	GetPromobyPromoCodeQuery = `SELECT id, promo_code, discount, is_percentage, start_date, end_date, description FROM promos WHERE promo_code = $1`
	GetAllPromoQuery = `SELECT id, promo_code, discount, is_percentage, start_date, end_date, description, created_at, updated_at FROM promos ORDER BY created_at ASC LIMIT $1 OFFSET $2`
	GetPromoForCustomerQuery = `SELECT id, promo_code, discount, is_percentage, start_date, end_date, description,
	created_at, updated_at FROM promos p
	WHERE NOT EXISTS (SELECT 1 FROM orders o WHERE o.promo_code = p.promo_code AND o.customer_id = $3)
	ORDER BY created_at ASC LIMIT $1 OFFSET $2`
	CountPromoQuery = `SELECT COUNT(*) FROM promos`
	CountPromoForCustomerQuery = `SELECT COUNT(*) FROM promos p 
	WHERE NOT EXISTS (SELECT 1 FROM orders o WHERE o.promo_code = p.promo_code AND o.customer_id = $1)`
	GetPromoByIdQuery = `SELECT id FROM promos where id = $1`
	DeletePromoQuery = `DELETE FROM promos WHERE id = $1`
)

// Order Query
const (
	CreateOrderQuery = `INSERT INTO orders(customer_id, address, promo_code, order_status, note, date, total_price) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`
	CreateOrderItemQuery = `INSERT INTO order_items(order_id, menu_id, quantity) VALUES($1, $2, $3) RETURNING id`
	CountunfinishCustomerOrderQuery = `SELECT COUNT(*) FROM orders WHERE customer_id = $1 AND order_status != 'delivered'`
	GetUnfinishOrderByCustomerIdQuery = `SELECT o.id, o.address, o.promo_code, o.order_status, o.note, o.total_price, o.created_at FROM orders o JOIN users u ON o.customer_id = u.id WHERE o.customer_id = $1 AND order_status != 'delivered' LIMIT 1`
	GetOrderByIdQuery = `SELECT o.id, u.username AS customer_name, o.address, o.promo_code, o.order_status, o.note, o.total_price, o.created_at FROM orders o JOIN users u ON o.customer_id = u.id WHERE o.id = $1`
	GetAllOrderQuery = `SELECT o.id, u.username AS customer_name, o.address, o.promo_code,
	o.order_status, o.note, o.total_price, o.created_at
	FROM orders o JOIN users u ON o.customer_id = u.id
	WHERE o.order_status = ANY($3)
	ORDER BY o.created_at ASC LIMIT $1 OFFSET $2`
	GetOrderItemsByOrderIdQuery = `SELECT oi.id, oi.order_id, m.name AS menu_name, oi.quantity FROM order_items oi JOIN menus m ON oi.menu_id = m.id WHERE oi.order_id = $1`
	UpdateOrderStatusQuery = `UPDATE orders SET order_status = $2 WHERE id = $1`
	CountfinishCustomerOrderQuery = `SELECT COUNT(*) FROM order_items oi
	JOIN orders o ON oi.order_id = o.id
	JOIN menus m ON oi.menu_id = m.id
	WHERE m.name = $1
	AND o.customer_id = $2
	AND o.id = $3
	AND o.order_status = 'delivered'`
	GetCustomerOrderHistoryQuery = `SELECT o.id, o.address, o.promo_code, o.order_status, o.note, o.date, o.total_price,
	o.created_at FROM orders o
	JOIN users u ON o.customer_id = u.id
	WHERE o.customer_id = $3 AND o.order_status = 'delivered' ORDER BY o.created_at ASC LIMIT $1 OFFSET $2`
	GetFilterDateCustomerOrderHistoryQuery = `SELECT o.id, o.address, o.promo_code, o.order_status, o.note, o.date, o.total_price,
	o.created_at FROM orders o
	JOIN users u ON o.customer_id = u.id
	WHERE o.date BETWEEN $3 AND $4 AND o.customer_id = $5 AND o.order_status = 'delivered' ORDER BY o.created_at ASC LIMIT $1 OFFSET $2`
	GetCustomerIdWithFinishOrderQuery = `SELECT customer_id, created_at FROM orders WHERE id = $1 AND order_status = 'delivered'`
	CountAllOrderQuery = `SELECT COUNT(*) FROM orders`
	CountFinishCustomerOrderQuery = `SELECT COUNT(*) FROM orders WHERE customer_id = $1 AND order_status = 'delivered'`
	CountUsagePromoQuery = `SELECT COUNT(*) FROM orders WHERE customer_id = $1 AND promo_code = $2 AND promo_used = 'TRUE'`
	UpdatePromoUsedStatusQuery = `UPDATE orders SET promo_used = TRUE WHERE id = $1`
)

// Review Query
const (
	CreateReviewQuery = `INSERT INTO reviews(customer_id, menu_id, order_id, rating, comment, buy_date, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`
	GetAllReviewQuery = `SELECT r.id, u.username AS customer_name, m.name AS menu_name, r.order_id,
	r.rating, r.comment, r.buy_date, r.created_at, r.updated_at FROM reviews r
	JOIN users u ON r.customer_id = u.id
	JOIN menus m ON r.menu_id = m.id
	ORDER BY r.created_at ASC
	LIMIT $1 OFFSET $2`
	GetReviewByIdQuery = `SELECT r.id, r.customer_id, u.username AS customer_name, m.name AS menu_name, r.order_id,
	r.rating, r.comment, r.buy_date, r.created_at, r.updated_at FROM reviews r
	JOIN users u ON r.customer_id = u.id
	JOIN menus m ON r.menu_id = m.id
	WHERE r.id = $1`
	UpdateReviewQuery = `UPDATE reviews SET rating = $3, comment = $4, updated_at = $5 WHERE id = $1 AND customer_id = $2`
	DeleteReviewQuery = `DELETE FROM reviews WHERE id = $1 AND customer_id = $2`
	CountReviewQuery = `SELECT COUNT(*) FROM reviews`
)
