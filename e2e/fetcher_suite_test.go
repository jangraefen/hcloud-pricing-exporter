package fetcher_test

import (
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	testClient = hcloud.NewClient(hcloud.WithToken(hcloudAPITokenFromENV()))
	testLabels = map[string]string{
		"test":  "github.com_jangraefen_hcloud-pricing-exporter",
		"suite": "fetcher_suite_test",
	}
)

func TestFetcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fetcher Suite")
}
