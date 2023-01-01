module github.com/jakecoffman/crud

go 1.16

// It's worth noting here that only the router that YOU depend on will
// end up in the binary. So if you use gorilla/mux since it has no
// additional dependencies, you are NOT also including gin and echo
// since they are listed here. Go will not include them in your binary.

require (
	github.com/gin-gonic/gin v1.8.1
	github.com/gorilla/mux v1.8.0
	github.com/labstack/echo/v4 v4.10.0
)
