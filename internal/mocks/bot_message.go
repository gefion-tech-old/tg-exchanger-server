package mocks

var BOT_MESSAGE_REQ = map[string]interface{}{
	"connector":    "text_connector",
	"message_text": "some message text here",
	"created_by":   MANAGER_IN_ADMIN_REQ["username"],
}
