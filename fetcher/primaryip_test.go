package fetcher_test

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/jangraefen/hcloud-pricing-exporter/fetcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var _ = Describe("For primary IPs", func() {
	sut := fetcher.NewPrimaryIP(&fetcher.PriceProvider{Client: testClient})

	BeforeEach(func(ctx context.Context) {
		resv4, _, err := testClient.PrimaryIP.Create(ctx, hcloud.PrimaryIPCreateOpts{
			Name:       "test-primaryipv4",
			Labels:     testLabels,
			Datacenter: "fsn1-dc8",
			Type:       hcloud.PrimaryIPTypeIPv4,
		})
		Expect(err).ShouldNot(HaveOccurred())
		DeferCleanup(testClient.PrimaryIP.Delete, resv4.PrimaryIP)

		resv6, _, err := testClient.PrimaryIP.Create(ctx, hcloud.PrimaryIPCreateOpts{
			Name:       "test-primaryipv6",
			Labels:     testLabels,
			Datacenter: "fsn1-dc8",
			Type:       hcloud.PrimaryIPTypeIPv6,
		})
		Expect(err).ShouldNot(HaveOccurred())
		DeferCleanup(testClient.PrimaryIP.Delete, resv6.PrimaryIP)
	})

	When("getting prices", func() {
		It("should fetch them", func() {
			Expect(sut.Run(testClient)).To(Succeed())
		})

		It("should get prices for correct values for v4", func() {
			Eventually(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-primaryipv4", "fsn1"))).Should(BeNumerically(">", 0.0))
			Eventually(testutil.ToFloat64(sut.GetMonthly().WithLabelValues("test-primaryipv4", "fsn1"))).Should(BeNumerically(">", 0.0))
		})

		It("should get prices for correct values for v6", func() {
			Eventually(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-primaryipv6", "fsn1"))).Should(BeNumerically("==", 0.0))
			Eventually(testutil.ToFloat64(sut.GetMonthly().WithLabelValues("test-primaryipv6", "fsn1"))).Should(BeNumerically("==", 0.0))
		})

		It("should get zero for incorrect values", func() {
			Eventually(testutil.ToFloat64(sut.GetHourly().WithLabelValues("invalid-name", "fsn1"))).Should(BeNumerically("==", 0))
			Eventually(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-primaryip", "nbg1"))).Should(BeNumerically("==", 0))
		})
	})
})
