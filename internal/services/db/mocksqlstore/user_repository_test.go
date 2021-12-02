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
}
