package api

import (
	"os"
	"testing"
	"time"

	db "github.com/DarrelASandbox/go-simple-bank/db/sqlc"
	"github.com/DarrelASandbox/go-simple-bank/db/util"
	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
