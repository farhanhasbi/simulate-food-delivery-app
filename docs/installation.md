# Installation Guide

Follow these steps to set up and run the project on your local environment.

# Prerequisites

Ensure you have the following installed on your machine:

**Go:** Version 1.22.3 or later
**PostgreSQL:** Version 16.3 or later
**Insomnia or Postman:** for API testing

# Steps

## 1. Clone the Repository

```bash
git clone https://github.com/yourusername/yourproject.git
cd yourproject
```

## 2. Set Up Environment Variables

1. Rename .env.example to .env in the root directory.
2. Add the following variables (modify the values as needed):

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_database_user
DB_PASSWORD=your_secure_password
DB_NAME=your_database_name
DB_DRIVER=postgres
API_PORT=8080
TOKEN_EXPIRE=60
TOKEN_ISSUER=your_app_name
TOKEN_SECRET=your_super_secret_key
```

## 3. Install Dependencies

Run the following command to download the necessary dependencies:

```bash
go mod tidy
```

## 4. Set Up the Database

Create a PostgreSQL database using the name specified in the .env file.

## 5. Run the Application

Start the application by running:

```bash
go run main.go
```

## 6. Test the API

[API Documentation](http://localhost:8080/swagger/index.html).
