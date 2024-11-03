# User Registration & Login API (Phase 0)

This project implements basic user registration and login functionalities using Go, Gin, and PostgreSQL as part of Phase 0.

## Features

- **User Registration**: Allows users to create an account by providing their name, email, phone number, and password.
- **User Login**: Enables users to log in using their email and password.

## Tech Stack

- **Backend**: Go, Gin, PostgreSQL
- **Frontend**: Angular

## Setup Instructions

### Prerequisites

- **Go** (>=1.17)
- **PostgreSQL**

### Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/yourusername/Package_Tracking_Backend.git
   cd Package_Tracking_Backend
   
2. **Install Go dependencies**:

   ```bash
   go mod tidy

3. **Database Setup**:

   ```go
   connStr := "user=postgres password=123 dbname=PackageTracking_db sslmode=disable"

4. **Run the application**:

   ```bash
   go run main.go
The backend server will start on http://localhost:8080.

**API Endpoints (Phase 0)**:
POST /register - Registers a new user.
POST /login - Logs in an existing user.
