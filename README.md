# Project Overview

This is a personal project focused on simulates a food delivery for a single restaurant. This app allows customers to browse the menu, place orders, track their orders, and leave reviews. Additionally, employees can manage menu items, handle orders, and apply promotional offers. Admins have access to user management features, including assigning employee roles.

## Features

### User Management

| HTTP Method | URL                     | Description                             | Access        |
| ----------- | ----------------------- | --------------------------------------- | ------------- |
| `POST`      | `/api/v1/auth/register` | Register the first admin, or customer   | No Auth       |
| `POST`      | `/api/v1/auth/login`    | User login and token generation         | No Auth       |
| `POST`      | `/api/v1/auth/logout`   | Logout a user                           | Authenticated |
| `PUT`       | `/api/v1/user/:id/role` | Update a user's role to employee        | Admin Only    |
| `DELETE`    | `/api/v1/user/:id`      | Delete a user                           | Admin Only    |
| `PUT`       | `/api/v1/user`          | Update authenticated user's information | Authenticated |
| `GET`       | `/api/v1/user`          | Get all users                           | Admin Only    |

### Menu Management

| HTTP Method | URL                | Description                  | Access   |
| ----------- | ------------------ | ---------------------------- | -------- |
| `POST`      | `/api/v1/menu`     | Add a new menu item          | Employee |
| `GET`       | `/api/v1/menu`     | Get all menu items           | No Auth  |
| `PUT`       | `/api/v1/menu/:id` | Update an existing menu item | Employee |
| `DELETE`    | `/api/v1/menu/:id` | Delete a menu item           | Employee |

### Balance Management

| HTTP Method | URL               | Description                     | Access   |
| ----------- | ----------------- | ------------------------------- | -------- |
| `POST`      | `/api/v1/balance` | Add balance to customer account | Customer |
| `GET`       | `/api/v1/balance` | Get customer’s current balance  | Customer |

### Promotions Management

| HTTP Method | URL                       | Description                      | Access   |
| ----------- | ------------------------- | -------------------------------- | -------- |
| `POST`      | `/api/v1/promo`           | Add a new promotion              | Employee |
| `GET`       | `/api/v1/promo`           | Get all promotions               | Employee |
| `GET`       | `/api/v1/available-promo` | Get available promo for customer | Customer |
| `DELETE`    | `/api/v1/promo/:id`       | Delete a promotion               | Employee |

### Order Management

| HTTP Method | URL                        | Description                                     | Access   |
| ----------- | -------------------------- | ----------------------------------------------- | -------- |
| `POST`      | `/api/v1/order`            | Place a new order (requires sufficient balance) | Customer |
| `GET`       | `/api/v1/unfinish-order`   | Track order status for specific customer        | Customer |
| `PUT`       | `/api/v1/order-status/:id` | Update order status                             | Employee |
| `GET`       | `/api/v1/order`            | Get all customer's orders                       | Employee |
| `GET`       | `/api/v1/finish-order`     | Get order history for specific customer         | Customer |

### Reviews Management

| HTTP Method | URL                  | Description                       | Access   |
| ----------- | -------------------- | --------------------------------- | -------- |
| `POST`      | `/api/v1/review`     | Add a review                      | Customer |
| `GET`       | `/api/v1/review`     | Get all reviews                   | No Auth  |
| `PUT`       | `/api/v1/review/:id` | Update specific customer's review | Customer |
| `DELETE`    | `/api/v1/review/:id` | Delete specific customer's review | Customer |

- [Installation Guide](docs/installation.md)
- [API Documentation](http://localhost:8080/swagger/index.html)
