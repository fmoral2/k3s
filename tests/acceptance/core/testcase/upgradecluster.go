package testcase

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/k3s-io/k3s/tests/acceptance/core/service/customflag"
	"github.com/k3s-io/k3s/tests/acceptance/shared/util"
	. "github.com/onsi/ginkgo/v2"
)

// TestUpgradeClusterManually upgrades the cluster "manually"
func TestUpgradeClusterManually(version string) error {
	if version == "" {
		return fmt.Errorf("please provide a non-empty k3s version or commit to upgrade to")
	}

	serverIPs := strings.Split(util.ServerIPs, ",")
	agentIPs := strings.Split(util.AgentIPs, ",")

	if util.NumServers == 0 && util.NumAgents == 0 {
		return fmt.Errorf("no nodes found to upgrade")
	}

	if util.NumServers > 0 {
		if err := upgradeServer(version, serverIPs); err != nil {
			return err
		}
	}

	if util.NumAgents > 0 {
		if err := upgradeAgent(version, agentIPs); err != nil {
			return err
		}
	}

	return nil
}

// upgradeServer upgrades servers in the cluster.
func upgradeServer(installType string, serverIPs []string) error {
	var wg sync.WaitGroup
	errCh := make(chan error)

	for _, ip := range serverIPs {
		switch {
		case customflag.InstallType.Version != "":
			installType = fmt.Sprintf("INSTALL_K3S_VERSION=%s", customflag.InstallType.Version)
		case customflag.InstallType.Commit != "":
			installType = fmt.Sprintf("INSTALL_K3S_COMMIT=%s", customflag.InstallType.Commit)
		}
		upgradeCommand := fmt.Sprintf(util.InstallK3sServer, installType)
		wg.Add(1)
		go func(ip, installFlagServer string) {
			defer wg.Done()
			defer GinkgoRecover()

			fmt.Printf("\nUpgrading server to:  " + upgradeCommand)
			if _, err := util.RunCmdOnNode(upgradeCommand, ip); err != nil {
				fmt.Printf("Error upgrading server %s: %v\n\n", ip, err)
				errCh <- err
				close(errCh)
				return
			}
			if _, err := util.RestartCluster(ip); err != nil {
				fmt.Printf("Error restarting server %s: %v\n\n", ip, err)
			}
			time.Sleep(30 * time.Second)
		}(ip, installType)
	}
	wg.Wait()

	return nil
}

// upgradeAgent upgrades agents in the cluster.
func upgradeAgent(installType string, agentIPs []string) error {
	var wg sync.WaitGroup
	errCh := make(chan error)

	for _, ip := range agentIPs {
		switch {
		case customflag.InstallType.Version != "":
			installType = fmt.Sprintf("INSTALL_K3S_VERSION=%s", customflag.InstallType.Version)
		case customflag.InstallType.Commit != "":
			installType = fmt.Sprintf("INSTALL_K3S_COMMIT=%s", customflag.InstallType.Commit)
		}
		upgradeCommand := fmt.Sprintf(util.InstallK3sAgent, installType)
		fmt.Println("\nUpgrading agent to: " + upgradeCommand)
		wg.Add(1)
		go func(ip, installFlagAgent string) {
			defer wg.Done()
			defer GinkgoRecover()

			if _, err := util.RunCmdOnNode(upgradeCommand, ip); err != nil {
				fmt.Printf("Error upgrading agent %s: %v\n\n", ip, err)
				errCh <- err
				close(errCh)
				return
			}
		}(ip, installType)
	}
	wg.Wait()

	return nil
}
