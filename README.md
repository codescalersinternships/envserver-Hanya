# Envserver
# Overivew
Envserver is a simple web application that displays environment variables on a local server
# How to Install and Run
1. Download the latest release for the server [from here](https://github.com/codescalersinternships/envserver-Hanya/releases/latest)
2. Run the binary file
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