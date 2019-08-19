package scwcp

import (
	"fmt"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
)

var cloudProvider *cloud

func registerResponders() {

	httpmock.RegisterResponder("GET", "http://169.254.42.42/conf?format=json",
		httpmock.NewStringResponder(200, `{"tags": ["lolnet", "head", "k3s"], "state_detail": "booted", "public_ip": {"dynamic": false, "id": "b4d3e520-e1c0-47bd-a6f3-40b639c63947", "address": "51.158.183.184"}, "ssh_public_keys": [{"description": null, "modification_date": "2018-10-09T16:45:03.649769+00:00", "ip": null, "creation_date": "2018-10-09T16:45:03.649769+00:00", "port": 0, "email": null, "user_agent": {}, "key": "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA02BfnOCerDukb+tFTePbh7efhmylYaqwaZsWtUIm+dUDrJDk1KikZx3gw70I0GuW/j7+onnAQsG09TMVQIQzRpQeyjp7YIak6u5S0OnmdhDxEc4M6vNhA7La2iRlCmkSNDrbOnWjsmS5rG6C5BRjHzqHIufeRbprt9MYw287rIIFuP1j67Suv7OOLC5tqoPfmQ9BCth72thdiBuOkEgL9Szf4RzYPn8uozQLVmxbLflBdzdbkKD7tO+St8rgUNvBgzKEkL26h6dIoUJU4mOAQmgbMVKncrBAeuF0OBHMADcjr45GHZyvW9/SguLhmywjLz6hrKnFmxZN0JKZXoXz6Q== guillaume@localhost.localdomain", "fingerprint": "2048 80:e6:34:2f:4b:d1:64:59:d8:e2:29:4b:c7:04:3d:49  guillaume@localhost.localdomain (RSA)", "id": "d20810eb-fc00-4ebe-9222-80ac2ff701da"}], "private_ip": "10.6.44.45", "timezone": null, "id": "3b800d09-b747-4c66-8fac-a1b916cb12f8", "extra_networks": [], "name": "lolnet-k3s-head", "hostname": "lolnet-k3s-head", "bootscript": {"kernel": "http://169.254.42.24/kernel/x86_64-mainline-lts-4.9-4.9.93-rev1/vmlinuz-4.9.93", "title": "x86_64 mainline 4.9.93 rev1", "default": false, "dtb": "", "public": false, "initrd": "http://169.254.42.24/initrd/initrd-Linux-x86_64-v3.14.6.gz", "bootcmdargs": "LINUX_COMMON scaleway boot=local nbd.max_part=16", "architecture": "x86_64", "organization": "11111111-1111-4111-8111-111111111111", "id": "54ee2857-8ffb-4323-abba-964f55fea4a2"}, "location": {"platform_id": "21", "node_id": "13", "cluster_id": "29", "zone_id": "ams1", "chassis_id": "1"}, "volumes": {"1": {"name": "lolnet-k3s-head-50", "modification_date": "2019-08-11T18:11:11.174723+00:00", "export_uri": "nbd://10.6.134.129:4481", "volume_type": "l_ssd", "creation_date": "2019-08-11T18:11:09.502241+00:00", "state": "available", "organization": "5cc6d17b-43bb-4bab-85c4-cf5ee19c06c2", "server": {"id": "3b800d09-b747-4c66-8fac-a1b916cb12f8", "name": "lolnet-k3s-head"}, "id": "f81b218b-31df-4bc8-b8dd-bce0531585f5", "size": 50000000000}, "0": {"name": "snapshot-de728daa-0bf6-4c64-abf5-a9477e791c83-2019-03-05_10:13", "modification_date": "2019-08-11T18:11:11.174691+00:00", "export_uri": "nbd://10.6.134.129:4480", "volume_type": "l_ssd", "creation_date": "2019-08-11T18:11:10.689914+00:00", "state": "available", "organization": "5cc6d17b-43bb-4bab-85c4-cf5ee19c06c2", "server": {"id": "3b800d09-b747-4c66-8fac-a1b916cb12f8", "name": "lolnet-k3s-head"}, "id": "530a3249-4188-4e76-8517-eecd71bc4ae8", "size": 50000000000}}, "ipv6": {"netmask": "127", "gateway": "2001:bc8:4700:2100::302c", "address": "2001:bc8:4700:2100::302d"}, "organization": "5cc6d17b-43bb-4bab-85c4-cf5ee19c06c2", "commercial_type": "C2S"}`))
}

func TestMain(m *testing.M) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerResponders()

	cp, err := newScalewayCloud()
	if err != nil {
		fmt.Printf("Could not initialize cloud provider : %v", err)
		os.Exit(1)
	}
	cloudProvider = cp

	os.Exit(m.Run())
}

func TestMetadataNodeAddresses(t *testing.T) {
	metadata, err := cloudProvider.client.GetMetadata()
	if err != nil {
		t.Fatalf("Could not get instance metadata because : %v", err)
	}
	if metadata.PublicIP.Address != "51.158.183.184" {
		t.Errorf("Public IP Address is different from what was expected (got %v)", metadata.PublicIP.Address)
	}
	if metadata.PrivateIP != "10.6.44.45" {
		t.Errorf("Got the wrong private IP Address")
	}
}

func TestNodeAddresses(t *testing.T) {
	instances, implemented := cloudProvider.Instances()
	if implemented != true {
		t.Fatalf("Instance is not implemented")
	}
	addresses, err := instances.NodeAddresses(nil, "A node")
	if err != nil {
		t.Errorf("Coudl not get addresses: %v", err)
	}
	fmt.Printf("Got %v", addresses)
}
