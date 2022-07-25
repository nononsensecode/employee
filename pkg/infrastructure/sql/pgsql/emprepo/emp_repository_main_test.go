package emprepo_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	pgUsername = "kaushik"
	pgPassword = "password"
	pgDb       = "employee"
)

var (
	hostAndPort string
	pgPort      string
	dbUrl       string
)

func Test_Main(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", pgPassword),
			fmt.Sprintf("POSTGRES_USER=%s", pgUsername),
			fmt.Sprintf("POSTGRES_DB=%s", pgDb),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort = resource.GetHostPort("5432/tcp")
	pgPort = resource.GetPort("5432/tcp")
	dbUrl = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pgUsername, pgPassword, hostAndPort, pgDb)

	resource.Expire(120)
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		pgPool, err := pgxpool.Connect(context.Background(), dbUrl)
		if err != nil {
			return err
		}
		defer func() {
			pgPool.Close()
		}()
		return pgPool.Ping(context.Background())
	}); err != nil {
		log.Fatalf("cannot connect to pool: %s", err)
	}

	exitCode := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("cannot purge the resource: %s", err)
	}

	os.Exit(exitCode)
}
