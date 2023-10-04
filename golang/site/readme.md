## About

This code sets up a basic web server that serves these HTML templates and handles login form submissions. 
When you run the Go program, you can access the webpages at:
http://localhost:8080/ and http://localhost:8080/login.

Please note that this is for demonstration purposes. 
In using this, you should implement proper user authentication and security measures to protect user data and credentials.

## How To Run

1. Install Go (if not already installed):
If you haven't already installed Go, you can download and install it from the official website: https://golang.org/dl/

2. Open a terminal/command prompt, navigate to your project directory, and run the following command to start the Go application:
go run main.go

3. Once the application is running, open a web browser and go to the following URLs:

Home Page: http://localhost:8080/
Login Page: http://localhost:8080/login
You should see the home page with a "Login" link. Clicking on the "Login" link will take you to the login page with a username and password form.

4. Enter a username and password (these are not validated in the provided example) and click the "Login" button. The application will print the entered credentials to the console and redirect you to the home page.

## Testing
```
$ make test
# or for html profiling
$ make html
```