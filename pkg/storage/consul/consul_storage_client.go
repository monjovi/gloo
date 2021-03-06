package consul

import (
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/storage"
	"github.com/solo-io/gloo/pkg/storage/base"
)

//go:generate go run ${GOPATH}/src/github.com/solo-io/gloo/pkg/storage/generate/generate_clients.go -f ${GOPATH}/src/github.com/solo-io/gloo/pkg/storage/consul/client_template.go.tmpl -o ${GOPATH}/src/github.com/solo-io/gloo/pkg/storage/consul/
type Client struct {
	v1 *v1client
}

// TODO: support basic auth and tls
func NewStorage(cfg *api.Config, rootPath string, syncFrequency time.Duration) (storage.Interface, error) {
	cfg.WaitTime = syncFrequency

	// Get a new client
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "creating consul client")
	}

	return &Client{
		v1: &v1client{
			upstreams: &upstreamsClient{
				base: base.NewConsulStorageClient(rootPath+"/upstreams", client),
			},
			virtualServices: &virtualServicesClient{
				base: base.NewConsulStorageClient(rootPath+"/virtualservices", client),
			},
			roles: &rolesClient{
				base: base.NewConsulStorageClient(rootPath+"/roles", client),
			},
			attributes: &attributesClient{
				base: base.NewConsulStorageClient(rootPath+"/attributes", client),
			},
		},
	}, nil
}

func (c *Client) V1() storage.V1 {
	return c.v1
}

type v1client struct {
	upstreams       *upstreamsClient
	virtualServices *virtualServicesClient
	roles           *rolesClient
	attributes      *attributesClient
}

func (c *v1client) Register() error {
	return nil
}

func (c *v1client) Upstreams() storage.Upstreams {
	return c.upstreams
}

func (c *v1client) VirtualServices() storage.VirtualServices {
	return c.virtualServices
}

func (c *v1client) Roles() storage.Roles {
	return c.roles
}

func (c *v1client) Attributes() storage.Attributes {
	return c.attributes
}
