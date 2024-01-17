package e2e_test

import (
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	testClient = hcloud.NewClient(hcloud.WithToken(hcloudAPITokenFromENV()))
	testLabels = map[string]string{
		"test":  "github.com_jangraefen_hcloud-pricing-exporter",
		"suite": "e2e_suite_test",
	}
)

func TestFetcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}
