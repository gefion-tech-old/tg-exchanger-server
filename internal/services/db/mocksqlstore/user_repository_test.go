package mocksqlstore_test

import (
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/mocksqlstore"
	"github.com/stretchr/testify/assert"
)

func Test_UserRepository(t *testing.T) {
	s := mocksqlstore.Init()

	u, err := s.User().Create(&models.User{
		ChatID:   int64(mocks.USER_IN_BOT_REGISTRATION_REQ["chat_id"].(int)),
		Username: mocks.USER_IN_BOT_REGISTRATION_REQ["username"].(string),
	})
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, mocks.USER_IN_BOT_REGISTRATION_REQ["username"], u.Username)

	uf, err := s.User().FindByUsername(mocks.USER_IN_BOT_REGISTRATION_REQ["username"].(string))
	assert.NoError(t, err)
	assert.NotNil(t, uf)
	assert.Equal(t, mocks.USER_IN_BOT_REGISTRATION_REQ["username"], uf.Username)
}
