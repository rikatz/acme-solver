package main

import (
	"os"

	"k8s.io/client-go/rest"

	"github.com/jetstack/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/jetstack/cert-manager/pkg/acme/webhook/cmd"

	"k8s.io/klog/v2"
)

// GroupName defines the webhook groupName
var GroupName = os.Getenv("GROUP_NAME")

func main() {
	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	cmd.RunWebhookServer(GroupName,
		&customDNSProviderSolver{},
	)
}

type customDNSProviderSolver struct {
}

func (c *customDNSProviderSolver) Name() string {
	return "nullsolver"
}

func (c *customDNSProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	klog.Infof("%s of %s got key %s", ch.ResolvedFQDN, ch.DNSName, ch.Key)
	return nil
}

func (c *customDNSProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	return nil
}

func (c *customDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	return nil
}
