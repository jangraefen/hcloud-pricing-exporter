package fetcher_test

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/jangraefen/hcloud-pricing-exporter/fetcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var _ = Describe("For floating IPs", func() {
	sut := fetcher.NewFloatingIP(&fetcher.PriceProvider{Client: testClient})

	BeforeEach(func(ctx context.Context) {
		location, _, err := testClient.Location.GetByName(ctx, "fsn1")
		Expect(err).NotTo(HaveOccurred())

		res, _, err := testClient.FloatingIP.Create(ctx, hcloud.FloatingIPCreateOpts{
			Name:         hcloud.String("test-floatingip"),
			Labels:       testLabels,
			HomeLocation: location,
			Type:         hcloud.FloatingIPTypeIPv6,
		})
		Expect(err).ShouldNot(HaveOccurred())
		DeferCleanup(testClient.FloatingIP.Delete, res.FloatingIP)
	})

	When("getting prices", func() {
		It("should fetch them", func() {
			Expect(sut.Run(testClient)).To(Succeed())
		})

		It("should get prices for correct values", func() {
			Eventually(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-floatingip", "fsn1"))).Should(BeNumerically(">", 0.0))
			Eventually(testutil.ToFloat64(sut.GetMonthly().WithLabelValues("test-floatingip", "fsn1"))).Should(BeNumerically(">", 0.0))
		})

		It("should get zero for incorrect values", func() {
			Eventually(testutil.ToFloat64(sut.GetHourly().WithLabelValues("invalid-name", "fsn1"))).Should(BeNumerically("==", 0))
			Eventually(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-floatingip", "nbg1"))).Should(BeNumerically("==", 0))
		})
	})
})
