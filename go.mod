module letter-square-api

go 1.22.2

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/stretchr/testify v1.9.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace filippo.io/edwards25519 => github.com/FiloSottile/edwards25519 v1.1.0

replace gopkg.in/yaml.v3 => github.com/go-yaml/yaml v0.0.0-20220527083530-f6f7691b1fde

replace gopkg.in/check.v1 => github.com/go-check/check v0.0.0-20161208181325-20d25e280405
