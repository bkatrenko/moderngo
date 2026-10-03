# moderngo

Bootstrap application for a masterclass "Go 1.26 in 2026".
This is simple and imperfect server app made with Go 1.27, which gives us an ability to try recently released changes in Golang.

# Structure: 
- main.go - application entry point, will use is to setup tracing and other project-wide features
- /controller - handler layer 
- /repository - repository layer

# Build & Run: 
- `go build` from the root of the application
- for a simplicity, application default port is hardcoded as `:8080`

# APIs:
GET http://localhost:8080/weather
GET http://localhost:8080/weather/batch

Please feel free to use any REST client or curl to try it out.

# Work hard, have fun!
