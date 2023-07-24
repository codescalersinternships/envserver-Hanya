# Envserver
# Overivew
Envserver is a simple web application that displays environment variables on a local server
# How to Install and Run
1. Run this command in your terminal
    ```go
    $ go get -d github.com/codescalersinternships/envserver-Hanya
    ```
2. Import the package's internal functions
    ```go
    import envserver "github.com/codescalersinternships/envserver-Hanya/internal"
    ```

3. Start a new server in your main function and run it
    ```go
	server, err := envserver.NewServer(3000)
    err = server.Run()
    ```
4. Build the package
    ```go
    $ go build -o envserver cmd/main.go
    ```
5. Run your program
    ```go
    $ ./envserver
    ```
# Routes
1. /env : displays all environment variables
2. /env/key : displays the value of key
# Errors
- ErrInvalidPort : raised when a reserved or out-of-range port number is used 
# How to Test
This package includes unit tests for 100% of the internal functions, run them using this command:
```go
$ go test ./...
```
You should see an output similar to this
```
ok      github.com/codescalersinternships/envserver-Hanya/internal 
```
If a test fails the error will provide relevant information