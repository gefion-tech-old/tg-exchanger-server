package mocks

var LOG_RECORD__ADMIN = map[string]interface{}{
	"username": USER_IN_BOT_REGISTRATION_REQ["username"],
	"info":     "some error text here",
	"service":  300,
	"module":   "HTTP SERVER",
}

var LOG_RECORD = map[string]interface{}{
	"info":    "some error text here",
	"service": 100,
	"module":  "HTTP SERVER",
}
