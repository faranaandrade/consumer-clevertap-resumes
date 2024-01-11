package config

import (
	"fmt"
	"strings"

	"github.com/occmundial/consumer-clevertap-applies/pkg/envconfig"
	"github.com/occmundial/go-common/logger"

	"context"

	goenvconfig "github.com/sethvargo/go-envconfig"
)

// https://medium.com/@felipedutratine/manage-config-in-golang-to-get-variables-from-file-and-env-variables-33d876887152

func NewConfiguration(log *logger.Log) *Configuration {
	var configuration = Configuration{}
	if err := goenvconfig.Process(context.Background(), &configuration); err != nil {
		log.Fatal("config", "go-envconfig.Process", err)
	}
	if missingEnvVars := envconfig.GetMissingEnvVars(&configuration); len(missingEnvVars) > 0 {
		log.Fatal(
			"config",
			"envconfig.GetMissingEnvVars",
			fmt.Errorf("missing required environment variables: %s", strings.Join(missingEnvVars, ", ")))
	}
	return &configuration
}
