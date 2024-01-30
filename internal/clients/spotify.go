/*
Copyright 2021 Upbound Inc.
*/

package clients

import (
	"context"
	"encoding/json"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/upjet/pkg/terraform"

	"github.com/tampakrap/provider-spotify/apis/v1beta1"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal spotify credentials as JSON"

	// provider config variables
	keyAPIKey     = "api_key"
	keyAuthServer = "auth_server"
	keyTokenID    = "token_id"
	keyUsername   = "username"
)

type spotifyConfig struct {
	APIKey     *string `json:"api_key,omitempty"`
	AuthServer *string `json:"auth_server,omitempty"`
	TokenID    *string `json:"token_id,omitempty"`
	Username   *string `json:"username,omitempty"`
}

func terraformProviderConfigurationBuilder(creds spotifyConfig) (terraform.ProviderConfiguration, error) {
	cnf := terraform.ProviderConfiguration{}

	if creds.APIKey != nil {
		cnf[keyAPIKey] = *creds.APIKey
	}

	if creds.AuthServer != nil {
		cnf[keyAuthServer] = *creds.AuthServer
	}

	if creds.TokenID != nil {
		cnf[keyTokenID] = *creds.TokenID
	}

	if creds.Username != nil {
		cnf[keyUsername] = *creds.Username
	}

	return cnf, nil
}

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return ps, errors.New(errNoProviderConfig)
		}

		pc := &v1beta1.ProviderConfig{}
		if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
			return ps, errors.Wrap(err, errGetProviderConfig)
		}

		t := resource.NewProviderConfigUsageTracker(client, &v1beta1.ProviderConfigUsage{})
		if err := t.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}

		data, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, client, pc.Spec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}

		creds := spotifyConfig{}
		if err := json.Unmarshal(data, &creds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		// Set credentials in Terraform provider configuration.
		ps.Configuration, err = terraformProviderConfigurationBuilder(creds)
		if err != nil {
			return ps, errors.Wrap(err, errProviderConfigurationBuilder)
		}

		return ps, nil
	}
}
