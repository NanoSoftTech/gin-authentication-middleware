# gin-authentication-middleware

Gin authentication middleware sample, original blogpost at [dandua.com](https://dandua.com)

<center><img alt="gopher" src="https://dandua.com/img/blog/blog_1/0.jpg" width="100%"></center>

If you have been using Golang to build microservices/web apps for a while, you
have probably heard of Gin. It's a high performance web framework with a pretty powerful and easy to use API.
At heart of Gin's request handling pipeline is support for middlewares at all pre or post route
route level handlers which makes extending the base API fairly easy and can potentially reduce the
development time for your application due to shared resources. This tutorial would explain using
that functionality for access control limiting and route level authentication.

*Note:* The source code for this tutorial can be found at this [repo](https://github.com/dandua98/gin-authentication-middleware).

# The application
This example application has three types of users, `basic`, `subscriber` and `admin`. While it doesn't have any
special functionality, your actual business logic could have multiple user type definitions and various restrictions on endpoints based on what
kind of user is hitting the endpoint. For example, a music application (like Spotify) could restrict download endpoints to
paid subscriber groups only while playback could be available to everyone.

My simple application doesn't have music playback or anything fancy though. It just echoes a message sent through one of the endpoints back to the user with `admin` and `subscriber` users getting their
username and authentication type back in the message. Note that an
admin can't hit subscriber endpoints and vice-versa. Any other user can't hit any of those endpoints. There's a *POST* endpoint (which doesn't do anything different) but can be hit by both `admin` and `subscriber` users. The users specify their authentication type during login.

**Note** that this isn't supposed to represent your business logic and the user authentication type should always be generated on
the backend itself. You can store it in an encrypted session on client side for sure but it's 2018 and storing claims in [JSON web tokens](https://jwt.io/) is a thing and is the recommended way to do it.

## Router groups
Assuming you have already initialized your gin router, you can define the groups as follow:
<script src="https://gist.github.com/dandua98/8ac0d0276a9799b2b466aaa570b2e02a.js"></script>
Note how the groups themselves can have sub-groups with the *relative link* (to the parent group) and that's what we use to apply
middleware to the group specifically (the `AuthenticationRequired` function, more on that later). You can have sub-sub-groups too but
it might mean a middleware encountered earlier aborting the request there on a user being unauthenticated etc. so it's not the best
practice. Also, using `routerGroup.METHOD().Use(middleware)` to define middleware on a route led to the middleware being applied on the whole group so not sure if that's the best practice.

## Authentication Required

<script src="https://gist.github.com/dandua98/7ec4f692eb56b1bfd88d0ad2bcc1e3bd.js"></script>

This is where you actually check the user's claim from the session/JWT and decide if you will let the user access the endpoint
or return a `StatusUnauthorized` (if user is not logged in) or `StatusForbidden` (invalid claim). The function itself returns a `gin.HandlerFunc` which is basically the middleware which would be called wherever you used it (`router.Use(AuthenticationRequired())`) before the request handler. The parameters define
what the valid claims are and after extracting the user's claim, you have to verify if the user has one of the valid claims.

Note that we call `c.Abort()` if the user is unauthenticated/unauthorized. This is because gin calls the next function  in the chain even after you write the header (`c.JSON()`) using `c.Next()`. The next handler is likely your route handler which might
try to rewrite the header which would result in an error. That also implies the last call to `c.Next()` is unnecessary ;)

You can also add another check to ensure the user is in your database (if somehow the secret key for the particular JWT/session) is leaked/found and the user is trying a malicious dictionary attack
to find other valid user IDs/emails signed with that key (a highly unlikely made-up scenario). Use `redis` or something similar to cache
the user check calls to the database in that case since this middleware might be called frequently.

# Testing the Application
Clone the application from [here](https://github.com/dandua98/gin-authentication-middleware) to your `$GOPATH/src` and then run `dep ensure` and `go run main.go`
```
cd $GOPATH/src

git clone https://github.com/dandua98/gin-authentication-middleware.git

// or

go get github.com/dandua98/gin-authentication-middleware

cd github.com/dandua98/gin-authentication-middleware

go run main.go
```

Now try the endpoints using Postman or cURL with your session generated with the required `authType` from the login endpoint.

