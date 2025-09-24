module restaurant-service

go 1.24

require (
	github.com/gin-gonic/gin v1.10.1
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/redis/go-redis/v9 v9.12.1
	github.com/segmentio/kafka-go v0.4.48
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.30.1
)

require github.com/rogpeppe/go-internal v1.14.1 // indirect
