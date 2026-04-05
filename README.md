# Finance Management Dashboard Backend

A powerful, modular Go-based backend for financial tracking and management. This application provides robust user management, financial record-keeping, and insightful analytics.

## 🚀 Features

-   **Modular Architecture**: Clean separation of concerns with handlers, views, and database layers.
-   **JWT Authentication**: Secure user sessions using JSON Web Tokens.
-   **Dynamic RBAC**: Role-Based Access Control managed through the database (`Admin`, `Analyst`, `Viewer`).
-   **Financial Records**: Comprehensive CRUD for incomes and expenses with categories and soft-delete capabilities.
-   **Insightful Summary**: Real-time totals, category breakdowns, and monthly trends.
-   **Security**: Rate limiting and secure BCRYPT password hashing.

## 🛠️ Technology Stack

-   **Language**: Go (v1.25.4)
-   **Web Framework**: Gin
-   **Database**: PostgreSQL (using `pgx` driver)
-   **Authentication**: JWT (JSON Web Tokens)
-   **Config**: Environment variables with `.env` support

## 🔧 Getting Started

### Prerequisites

-   **Go**: 1.25.4+ installed.
-   **PostgreSQL**: A running instance with a database for the project.

### Installation

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/radhikabhut/Finance-Management.git
    cd finance-dashboard-backend
    ```

2.  **Environment Setup**:
    Create a `.env` file in the root directory and add the following:
    ```env
    SERVER_PORT=8080
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=your_db_password
    DB_NAME=finance
    JWT_SECRET=your_jwt_secret
    ```

3.  **Database Initialization**:
    Run the SQL script found in `db/schema.sql` to set up the tables and initial permissions.

4.  **Install Dependencies**:
    ```bash
    go mod download
    ```

5.  **Run the Server**:
    ```bash
    go run main.go router.go
    ```

## 🔐 Default Credentials (for testing)

| Role     | Username  | Password     | Permissions                                  |
| :------- | :-------- | :----------- | :------------------------------------------- |
| **Admin** | `admin`   | `admin2204`  | Full access (User & Record Management)        |
| **Analyst**| `analyst` | `analyst2003`| Record Listing & Analytics View              |
| **Viewer** | `viewer`  | `viewer2022` | Record Listing                               |

## 📄 API Documentation

Detailed endpoint documentation can be found in [API_DOCUMENTATION.md](API_DOCUMENTATION.md).

## 🗃️ Database Schema

The core schema includes:
-   `users`: Manages authentication and roles.
-   `financial_records`: Stores individual income and expense entries.
-   `role_permissions`: Defines the dynamic RBAC permissions.

