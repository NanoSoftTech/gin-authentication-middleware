# gin-authentication-middleware

Gin authentication middleware sample, read the blogpost at for an explanation [dandua.com/blog/gin-authentication-middleware/](https://dandua.com/blog/gin-authentication-middleware/)

<center><img alt="gopher" src="https://dandua.com/img/blog/blog_1/0.jpg" width="100%"></center>

# Testing the Application
Clone the application from [here](https://github.com/dandua98/gin-authentication-middleware) to your `$GOPATH/src` and then run `dep ensure` and `go run main.go`
```
cd $GOPATH/src

go get github.com/dandua98/gin-authentication-middleware

cd github.com/dandua98/gin-authentication-middleware

go run main.go
```

Now try the endpoints using Postman or cURL with your session generated with the required `authType` from the login endpoint.
