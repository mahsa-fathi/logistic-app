package configs

import "time"

type intKey int
type AuthStatusVal int

// context keys
const (
	UserIDKey intKey = iota
	AuthStatusKey
)

// values for AuthStatusKey in context
const (
	AuthStatusValUnauthorized AuthStatusVal = iota
	AuthStatusValAuthorized
)

var TokenExpiration = time.Duration(intEnv("TOKEN_EXPIRATION", 24)) * time.Hour
var SecretKey = stringEnv("SECRET_KEY", "random_secret_key")
var ServerURL = stringEnv("SERVER_URL", "localhost:8080")
var LogError = boolEnv("LOG_ERROR", true)

var JWTDefaults = map[string]any{
	"AUTH_HEADER_TYPES": []string{"Bearer"},
	"AUTH_HEADER_NAME":  "Authorization",
	"USER_ID_FIELD":     "id",
	"USER_ID_CLAIM":     "user_id",
}

var DBAddress = stringEnv("DB_ADDRESS", "localhost")
var DBPort = stringEnv("DB_PORT", "5432")
var DBUser = stringEnv("DB_USER", "user")
var DBPassword = stringEnv("DB_PASSWORD", "password")
var DBName = stringEnv("DB_NAME", "db")

var DBTestAddress = stringEnv("DB_TEST_ADDRESS", "localhost")
var DBTestPort = stringEnv("DB_TEST_PORT", "5431")
var DBTestUser = stringEnv("DB_TEST_USER", "postgres")
var DBTestPassword = stringEnv("DB_TEST_PASSWORD", "postgres")
var DBTestName = stringEnv("DB_TEST_NAME", "postgres")

var OrderUpdatePeriod = time.Duration(intEnv("ORDER_UPDATE_PERIOD", 24*60*60)) * time.Second
var PeriodicTaskMaxConcurrency = intEnv("PERIODIC_TASK_MAX_CONCURRENCY", 10)
