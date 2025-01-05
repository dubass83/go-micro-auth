package api

import (
	"os"
	"testing"

	data "github.com/dubass83/go-micro-auth/data/sqlc"
	"github.com/dubass83/go-micro-auth/util"
)

func NewTestServer(t *testing.T, store data.Store) *Server {
	config := util.Config{
		LogService: "localhost:8080",
		Enviroment: "test",
	}
	server := CreateNewServer(config, store)
	// require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
