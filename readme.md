# Sales Record Application

## How to Run the App

1. **Prerequisites:**
   - Make sure you have Go installed on your machine. You can download it from [golang.org](https://golang.org/dl/).
   - Ensure that you have MongoDB installed and running. You can find installation instructions on the [MongoDB website](https://www.mongodb.com/try/download/community). You can use local instance or create an account with MongoDB Atlas

   - Setup your path to csv and mongodb connection details in .env file

2. **Clone the Repository:**
   ```bash
   git clone https://github.com/RishanKP/sales-record-app.git
   cd sales-record-app
   go mod tidy
   go run main.go

## API DOCUMENTATION

| Method | Endpoint          | Description                                                     |  
|--------|-------------------|-----------------------------------------------------------------|
| POST   | /refresh          | Initiates a background refresh of sales data from a CSV file.   | 
| GET    | /revenue          | Retrieves total revenue for a specified date range.             | 
| GET    | /revenue/product  | Retrieves total revenue by product for a specified date range.  | 
| GET    | /revenue/category | Retrieves total revenue by category for a specified date range. | 

## Sample Requests
    ```bash
    curl -X POST http://localhost:8080/refresh

    curl -X GET "http://localhost:8080/revenue?starDate=2023-01-01&endDate=2023-12-31"




