# Finance Dashboard API Documentation

This document describes the API endpoints for the Finance Dashboard Backend.

## General Information
- **Base URL**: `http://localhost:8080/v1`
- **Method**: All interactive endpoints use **POST**.
- **Content-Type**: `application/json`

---

## Authentication
Most endpoints require a **JSON Web Token (JWT)**.
1.  Obtain a token via the `POST /login` endpoint.
2.  Include the token in the `Authorization` header as a Bearer Token:
    `Authorization: Bearer <your_token>`

---

## Default Credentials
For testing purposes, you can use the following default credentials:
  Admin : 
- **Username**: `admin`
- **Password**: `admin2204`

  Analyst : 
- **Username**: `analyst`
- **Password**: `analyst2003`

  Viewer : 
- **Username**: `viewer`
- **Password**: `viewer2022`

---

## Rate Limiting
- **Global Limit**: 5 requests per second.
- **Burst Limit**: 10 requests.
- **Retry**: On limit exceeded, the server returns `429 Too Many Requests`.

---

## Endpoints

### 1. Authentication
#### `POST /login`
- **Description**: Authenticate a user and receive a JWT token.
- **Request**:
  ```json
  { "username": "admin", "password": "password123" }
  ```
- **Response**:
  ```json
  { "success": true, "token": "...", "role": "Admin" }
  ```

---

### 2. User Management (Admin Only)
#### `POST /users/create`
- **Description**: Create a new user.
- **Permissions**: `user::create`
- **Request**:
  ```json
  { "username": "newuser", "password": "securepassword", "role": "Viewer", "status": "Active" }
  ```

#### `POST /users/list`
- **Description**: List all users (paginated).
- **Permissions**: `user::list`
- **Request**:
  ```json
  { "page": 1, "limit": 10 }
  ```

#### `POST /users/update`
- **Description**: Update user details.
- **Permissions**: `user::update`
- **Request**:
  ```json
  { "id": "uuid", "username": "updated_name", "status": "Inactive" }
  ```

---

### 3. Financial Records
#### `POST /records/create`
- **Description**: Create a new income or expense record.
- **Permissions**: `record::create`
- **Request**:
  ```json
  { "amount": 100.50, "type": "Income", "category": "Salary", "notes": "Monthly salary" }
  ```

#### `POST /records/list`
- **Description**: List your records with optional filters.
- **Permissions**: `record::list`
- **Request Parameters**: `type`, `category`, `startDate`, `endDate`, `page`, `limit`.

#### `POST /records/delete`
- **Description**: Soft delete a record.
- **Permissions**: `record::delete`
- **Request**:
  ```json
  { "id": "uuid" }
  ```

---

### 4. Analytics
#### `POST /summary`
- **Description**: Get your financial summary (totals, trends, and recent activity).
- **Permissions**: `summary::view`
- **Response**:
  ```json
  {
    "success": true,
    "totalIncome": 5000,
    "totalExpense": 2000,
    "netBalance": 3000,
    "categoryTotals": [...],
    "recentActivity": [...],
    "monthlyTrends": [...]
  }
  ```
