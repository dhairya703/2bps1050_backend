Redis and all Are working locally but since my aws was billed previously so instead of s3 i used cloudinary

Installation
Clone the repository:

bash
Copy
Edit
git clone https://github.com/dhairya703/2bps1050_backend.git
cd 2bps1050_backend
Install Go dependencies:

bash
Copy
Edit
go mod tidy
Set up the PostgreSQL database:

Create a PostgreSQL database and configure the credentials in the .env file (see Environment Variables).

Set up environment variables in .env:

Create a .env file in the root directory and add the following variables:

env
Copy
Edit
DB_HOST=localhost
DB_PORT=5432
DB_NAME=file_sharing
DB_USER=your_db_user
DB_PASSWORD=your_db_password
JWT_SECRET=your_jwt_secret
CLOUDINARY_URL=your_cloudinary_url
REDIS_URL=your_redis_url
Environment Variables
DB_HOST: The host of the PostgreSQL database.

DB_PORT: The port where the database is running (default is 5432).

DB_NAME: The name of the PostgreSQL database.

DB_USER: Database username.

DB_PASSWORD: Database password.

JWT_SECRET: Secret key used for JWT signing.

CLOUDINARY_URL: URL for Cloudinary for file storage.

REDIS_URL: URL for Redis instance.

Running the Project
To run the application locally, use the following command:

bash
Copy
Edit
go run main.go
This will start the server on http://localhost:8080.


Public Routes
POST /api/register: Register a new user.

Request body: { "email": "user@example.com", "password": "password123" }

Response: { "message": "User registered successfully" }

![image](https://github.com/user-attachments/assets/35c17a00-22a3-43bd-b672-91ce543c9be2)


POST /api/login: Login a user and return a JWT token.

Request body: { "email": "user@example.com", "password": "password123" }

Response: { "access_token": "your_jwt_token" }

![image](https://github.com/user-attachments/assets/06bfd93e-b3fc-47f1-9125-0fde3aa062dc)
Session Token Storage
![image](https://github.com/user-attachments/assets/d24dd2f3-57d1-47d2-8e85-b18959108bf6)


Now Only Valid Jwt user can access other apis 

Protected Routes (Requires JWT Authentication)
POST /api/upload: Upload a file.

Request body: multipart/form-data

Response: { "message": "File uploaded successfully", "file_id": 123 }
![image](https://github.com/user-attachments/assets/b33ece12-4af7-4318-bf6a-a6db46648fb8)
![image](https://github.com/user-attachments/assets/5032af72-3317-4765-b5cf-8ffde67e3118)

GET /api/profile: View user profile.

Response: { "email": "user@example.com", "user_id": 1 }

GET /api/files: List all files uploaded by the user.

Response: [ { "file_name": "file1.jpg", "size": 123456 }, { "file_name": "file2.pdf", "size": 789101 } ]

GET /api/share/:file_id: Share a file by its ID.

Response: { "share_url": "http://example.com/file/123" }

GET /api/search: Search files by name, size, or date.

Response: [ { "file_name": "file1.jpg", "size": 123456 }, { "file_name": "file2.pdf", "size": 789101 } ]

Testing
To run tests for your project:

Write test cases under the /tests folder.

Run the tests using the following command:

bash
Copy
Edit
go test ./tests -v
