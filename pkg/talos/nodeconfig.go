package talos

import (
	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
	taloscfg "github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/generate"
	"github.com/siderolabs/talos/pkg/machinery/config/machine"
)

func GenerateNodeConfigBytes(node *config.Nodes, input *generate.Input) ([]byte, error) {
	cfg, err := generateNodeConfig(node, input)
	if err != nil {
		return nil, err
	}
	return cfg.Bytes()
}

func generateNodeConfig(node *config.Nodes, input *generate.Input) (taloscfg.Provider, error) {
	var c taloscfg.Provider
	var err error

	nodeInput, err := patchNodeInput(node, input)
	if err != nil {
		return nil, err
	}

	switch node.ControlPlane {
	case true:
		c, err = nodeInput.Config(machine.TypeControlPlane)
		if err != nil {
			return nil, err
		}
	case false:
		c, err = nodeInput.Config(machine.TypeWorker)
		if err != nil {
			return nil, err
		}
	}

	// https://github.com/budimanjojo/talhelper/issues/81
	if input.Options.VersionContract.SecretboxEncryptionSupported() && input.Options.SecretsBundle.Secrets.AESCBCEncryptionSecret != "" {
		c.RawV1Alpha1().ClusterConfig.ClusterAESCBCEncryptionSecret = input.Options.SecretsBundle.Secrets.AESCBCEncryptionSecret
	}

	cfg := applyNodeOverride(node, c)

	return *cfg, nil
}

func applyNodeOverride(node *config.Nodes, cfg taloscfg.Provider) *taloscfg.Provider {
	cfg.RawV1Alpha1().MachineConfig.MachineNetwork.NetworkHostname = node.Hostname

	if len(node.Nameservers) != 0 {
		cfg.RawV1Alpha1().MachineConfig.MachineNetwork.NameServers = node.Nameservers
	}

	if node.DisableSearchDomain {
		cfg.RawV1Alpha1().MachineConfig.MachineNetwork.NetworkDisableSearchDomain = &node.DisableSearchDomain
	}

	if len(node.NetworkInterfaces) != 0 {
		cfg.RawV1Alpha1().MachineConfig.MachineNetwork.NetworkInterfaces = node.NetworkInterfaces
	}

	if node.InstallDiskSelector != nil {
		cfg.RawV1Alpha1().MachineConfig.MachineInstall.InstallDiskSelector = node.InstallDiskSelector
	}

	if len(node.KernelModules) != 0 {
		cfg.RawV1Alpha1().MachineConfig.MachineKernel = &v1alpha1.KernelConfig{}
		cfg.RawV1Alpha1().MachineConfig.MachineKernel.KernelModules = node.KernelModules
	}

	if node.NodeLabels != nil {
		cfg.RawV1Alpha1().MachineConfig.MachineNodeLabels = node.NodeLabels
	}

	return &cfg
}

func patchNodeInput(node *config.Nodes, input *generate.Input) (*generate.Input, error) {
	nodeInput := input
	if node.InstallDisk != "" {
		nodeInput.Options.InstallDisk = node.InstallDisk
	}

	if len(node.MachineDisks) > 0 {
		nodeInput.Options.MachineDisks = node.MachineDisks
	}

	return nodeInput, nil
}
