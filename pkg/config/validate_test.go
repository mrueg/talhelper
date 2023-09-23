package config

import (
	"testing"
)

func TestValidateFromByte(t *testing.T) {
	data := []byte(`clusterName: "test-cluster"
talosVersion: vv1.0.6
kubernetesVersion: 1.23.6.5
allowSchedulingOnMasters: true
cniConfig:
  name: customo
clusterPodNets:
  - 10.244.0.0/16
clusterSvcNets:
  - 10.244.0.0/16
  - 10.10.0.0/166
nodes:
  - hostname: master1
    ipAddress: 1.2.3.4.5
    installDisk: /dev/sda
    disableSearchDomain: true
    nameservers:
      - 8.8.8.8
    networkInterfaces:
      - addresses:
          - 1.2.3.4
        mtu: 1500
        routes:
          - network: 0.0.0.0/0
            gateway: 1.2.3.4.5.6
    configPatches:
      - op: del
        path: /cluster
  - nodeLabels:
      ra*ck: rack1a
      z***: hahaha
    machineFiles:
      - op: ccreate
`)

	found, err := ValidateFromByte(data)
	if err != nil {
		t.Fatal(err)
	}

	expectedErrors := map[string]bool{
		"clusterName": false,
		"talosVersion": true,
		"kubernetesVersion": true,
		"endpoint": true,
		"cniConfig": true,
		"clusterPodNets": false,
		"clusterSvcNets": true,
		"nodes[0].hostname": false,
		"nodes[0].ipAddress": false,
		"nodes[0].controlPlane": false,
		"nodes[0].installDisk": false,
		"nodes[0].nameservers": false,
		"nodes[0].networkInterfaces": true,
		"nodes[0].configPatches": true,
		"nodes[1].hostname": true,
		"nodes[1].ipAddress": true,
		"nodes[1].installDisk": true,
		"nodes[1].nodeLabels": true,
		"nodes[1].machineFiles": true,
	}

	for k, v := range expectedErrors {
		if found.HasField(k) != v {
			t.Errorf("%s: got %t, want %t", k, found.HasField(k), v)
		}
	}
}