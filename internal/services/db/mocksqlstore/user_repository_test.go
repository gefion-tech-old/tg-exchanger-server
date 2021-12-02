package mocksqlstore_test

import (
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/mocksqlstore"
	"github.com/stretchr/testify/assert"
)

func Test_UserRepository(t *testing.T) {
	s := mocksqlstore.Init()

	u, err := s.User().Create(&mocks.USER_IN_BOT_REGISTRATION_REQUEST)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, mocks.USER_IN_BOT_REGISTRATION_REQUEST.Username, u.Username)

	uf, err := s.User().FindByUsername(mocks.USER_IN_BOT_REGISTRATION_REQUEST.Username)
	assert.NoError(t, err)
	assert.NotNil(t, uf)
	assert.Equal(t, mocks.USER_IN_BOT_REGISTRATION_REQUEST.Username, uf.Username)
}
