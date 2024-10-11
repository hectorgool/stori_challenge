
# Stori Challenge

## Description

This API is built using Docker and utilizes a MySQL database. To function correctly, your computer must have the following ports available:

- **Port 3306** for the MySQL database
- **Port 8081** for the API application

The API offers two endpoints:

1. **Test Endpoint** (`/`): This endpoint is used for testing purposes and simply returns the message: *"Hello, Stori."*
   
2. **Email Summary Endpoint** (`/csv`): This endpoint receives the email address of the recipient who will receive the summary information and the `.csv` file containing the records to be processed. By default, the system is capable of processing files smaller than 1 MB.

The summary email contains the following information:

1. Total balance: 39.74
2. Number of transactions in July: 2
3. Number of transactions in August: 2
4. Average debit amount: -15.38
5. Average credit amount: 35.25

### Requirements

- You need to have Docker Engine installed to use Docker Compose.

### Installation Steps

1. Clone the repository:

   ```sh
   git clone https://github.com/hectorgool/stori_challenge.git
   ```

2. Navigate to the project folder:

   ```sh
   cd stori_challenge
   ```

3. Build and run the containers:

   ```sh
   docker-compose up --build
   ```

   Please wait approximately 30 seconds for the following response indicating that the build is successful:

   ```sh
   docker_function  | running...
   docker_function  | [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.
   docker_function  | 
   docker_function  | [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
   docker_function  |  - using env:	export GIN_MODE=release
   docker_function  |  - using code:	gin.SetMode(gin.ReleaseMode)
   docker_function  | 
   docker_function  | [GIN-debug] GET    /                         --> main.main.func1 (3 handlers)
   docker_function  | [GIN-debug] POST   /csv                      --> stori_challenge/internal/handlers.HandleCSVUpload (3 handlers)
   docker_function  | [GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
   docker_function  | Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
   docker_function  | [GIN-debug] Listening and serving HTTP on :8081
   ```

### Available Endpoints

1. **Test Endpoint**

   Send a GET request to:

   ```sh
   GET http://localhost:8080
   ```

2. **Email Summary**

   This endpoint sends an email to the provided address (e.g., hectorgool@gmail.com) with the following information:

   - Total balance: 39.74
   - Average debit amount: -15.38
   - Number of transactions in July: 2
   - Average credit amount: 35.25
   - Number of transactions in August: 2

   The transaction file (`txns.csv`) is stored in a MySQL database in the `sql_document` table. Be sure to check the provided email address for the report.

   ```sh
   curl -X POST http://localhost:8081/csv \
   -F "email=coolorvibes@gmail.com" \
   -F "file=@path/to/file/txns.csv"
   ```

### Running Tests with `test.sh`

You can use the `test.sh` script to run tests on the API. This script contains a `curl` command that sends an email and a `.csv` file to the `/sendmail` endpoint. To run the script, execute:

```sh
./test.sh
```

This will automatically trigger a request to send the summary email and process the attached `.csv` file into the database.

### Stopping the Containers

In a separate terminal, you can stop the containers with:

```sh
docker-compose down
```

### Accessing the Database

To open the MySQL database, run:

```sh
docker exec -it docker_db mysql -u storiuser -pasdf -h localhost storidb
```

### Email Configuration

To send emails, configure the Gmail password by editing the `.env` file:

```sh
SMTP_PASSWD=""
```

### Support

For any issues or inquiries, please contact Héctor at [hectorgool@gmail.com](mailto:hectorgool@gmail.com).
