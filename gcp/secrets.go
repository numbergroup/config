package gcp

import (
	"context"
	"encoding/json"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/cockroachdb/errors"
	"github.com/numbergroup/cleanenv"
)

func LoadJSONSecretsIntoEnvThenUpdateConfig(ctx context.Context, secretClient *secretmanager.Client, secrets []string, confPtr any) error {
	if len(secrets) == 0 || secretClient == nil {
		return nil
	}
	for _, secretName := range secrets {

		// Create the request to access the secret.
		accessSecretReq := &secretmanagerpb.AccessSecretVersionRequest{
			Name: secretName,
		}
		secret, err := secretClient.AccessSecretVersion(ctx, accessSecretReq)
		if err != nil {
			return err
		}
		// Load into cleanenv
		var values map[string]string
		err = json.Unmarshal(secret.Payload.Data, &values)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal secret %s", secretName)
		}
		for key, val := range values {
			err = os.Setenv(key, val)
			if err != nil {
				return errors.Wrapf(err, "failed to set environment variable %s from secret %s", key, secretName)
			}
		}

	}
	return cleanenv.UpdateEnv(confPtr)
}
