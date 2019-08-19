package scwcp

import (
	"context"
	"errors"
	"fmt"
	"io"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	instance_api "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	cloudprovider "k8s.io/cloud-provider"
)

const (
	ProviderName string = "scaleway"
	metadata_url string = "http://169.254.42.42/conf?format=json"
)

type cloud struct {
	client    *instance_api.MetadataAPI
	instances cloudprovider.Instances
}

func (c *cloud) Instances() (cloudprovider.Instances, bool) {
	return c, true
}

func newScalewayCloud() (*cloud, error) {
	return &cloud{
		client: instance_api.NewMetadataAPI(),
	}, nil
}

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, func(io.Reader) (cloudprovider.Interface, error) {
		return newScalewayCloud()
	})
}

func (c *cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
}

func (c *cloud) Zones() (cloudprovider.Zones, bool) {
	return nil, false
}

func (c *cloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

func (c *cloud) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

func (c *cloud) ProviderName() string {
	return ProviderName
}

func (c *cloud) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nil, nil
}

func (c *cloud) HasClusterID() bool {
	return false
}

func (c *cloud) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	metadata, err := c.client.GetMetadata()
	if err != nil {
		return nil, fmt.Errorf("Could not get metadata because : %v", err)
	}
	public_address := v1.NodeAddress{Type: v1.NodeExternalIP, Address: metadata.PublicIP.Address}
	private_address := v1.NodeAddress{Type: v1.NodeInternalIP, Address: metadata.PrivateIP}

	return []v1.NodeAddress{public_address, private_address}, nil
}

func (c *cloud) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return fmt.Errorf("Not implemented")
}

func (c *cloud) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (c *cloud) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	return "", fmt.Errorf("Not implemented")
}

func (c *cloud) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	metadata, err := c.client.GetMetadata()
	if err != nil {
		return "", fmt.Errorf("Could not get metadata because : %v", err)
	}
	return metadata.ID, nil
}

func (c *cloud) InstanceType(ctx context.Context, nodeName types.NodeName) (string, error) {
	metadata, err := c.client.GetMetadata()
	if err != nil {
		return "", fmt.Errorf("Could not get metadata because : %v", err)
	}
	return metadata.CommercialType, nil
}

func (c *cloud) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

// NodeAddressesByProviderID returns the node addresses of an instances with the specified unique providerID
// This method will not be called from the node that is requesting this ID. i.e. metadata service
// and other local methods cannot be used here
func (c *cloud) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	return []v1.NodeAddress{}, errors.New("unimplemented")
}

// CurrentNodeName returns the name of the node we are currently running on
func (c *cloud) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (c *cloud) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
}
