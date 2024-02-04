package playlist

import "github.com/crossplane/upjet/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("spotify_playlist", func(r *config.Resource) {
		r.ShortGroup = "playlist"
	})
}
