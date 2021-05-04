module github.com/dusansimic/yaas

go 1.15

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/squirrel v1.5.0
	github.com/danielkov/gin-helmet v0.0.0-20171108135313-1387e224435e
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/influxdata/influxdb-client-go/v2 v2.2.2
	github.com/jinzhu/now v1.1.2
	github.com/jkomyno/nanoid v0.0.0-20210124212126-241306812283
	github.com/joho/godotenv v1.3.0
	github.com/komkom/toml v0.0.0-20210317065440-24f427ca88cc
	github.com/lib/pq v1.10.0
	github.com/mssola/user_agent v0.5.2
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
)

replace github.com/gin-contrib/sessions v0.0.3 => github.com/dusansimic/sessions v0.0.4
