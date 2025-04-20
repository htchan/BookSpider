package repo

import (
	"fmt"
	"log"

	"github.com/htchan/BookSpider/internal/config/v2"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func CreatePsqlContainer(
	name string, conf config.DatabaseConfig, initDB func() error,
) (close func() error, err error) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not construct pool: %w", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not connect to Docker: %w", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       name,
		Repository: "postgres",
		Tag:        "14.0-alpine",
		Env: []string{
			"POSTGRES_USER=" + conf.User,
			"POSTGRES_PASSWORD=" + conf.Password,
			"POSTGRES_DB=" + conf.Name,
		},
		ExposedPorts: []string{"5432/tcp"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostIP: "localhost", HostPort: conf.Port + "/tcp"}},
		},
	}, func(conf *docker.HostConfig) {
		conf.AutoRemove = true
	})
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %w", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if retryErr := pool.Retry(initDB); retryErr != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return func() error { return CloseResource(pool, resource) }, nil
}

func CloseResource(pool *dockertest.Pool, resource *dockertest.Resource) error {
	fmt.Println("close resource")
	if err := pool.Purge(resource); err != nil {
		log.Printf("could not purge resource: %s", err)

		return err
	}

	return nil
}
