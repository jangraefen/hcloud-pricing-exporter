package e2e_test

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/jangraefen/hcloud-pricing-exporter/fetcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var _ = Describe("For primary IPs", Ordered, Label("primaryips"), func() {
	sut := fetcher.NewPrimaryIP(&fetcher.PriceProvider{Client: testClient})

	BeforeAll(func(ctx context.Context) {
		By("Creating a IPv4 address")
		resv4, _, err := testClient.PrimaryIP.Create(ctx, hcloud.PrimaryIPCreateOpts{
			Name:         "test-primaryipv4",
			Labels:       testLabels,
			Datacenter:   "fsn1-dc14",
			Type:         hcloud.PrimaryIPTypeIPv4,
			AssigneeType: "server",
		})
		Expect(err).ShouldNot(HaveOccurred())
		DeferCleanup(testClient.PrimaryIP.Delete, resv4.PrimaryIP)

		waitUntilActionSucceeds(ctx, resv4.Action)

		By("Creating a IPv6 address")
		resv6, _, err := testClient.PrimaryIP.Create(ctx, hcloud.PrimaryIPCreateOpts{
			Name:         "test-primaryipv6",
			Labels:       testLabels,
			Datacenter:   "fsn1-dc14",
			Type:         hcloud.PrimaryIPTypeIPv6,
			AssigneeType: "server",
		})
		Expect(err).ShouldNot(HaveOccurred())
		DeferCleanup(testClient.PrimaryIP.Delete, resv6.PrimaryIP)

		waitUntilActionSucceeds(ctx, resv6.Action)
	})

	When("getting prices", func() {
		It("should fetch them", func() {
			By("Running the price collection")
			Expect(sut.Run(testClient)).To(Succeed())
		})

		It("should get prices for correct values for v4", func() {
			By("Checking IPv4 prices")
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-primaryipv4", "fsn1-dc14"))).Should(BeNumerically(">", 0.0))
			Expect(testutil.ToFloat64(sut.GetMonthly().WithLabelValues("test-primaryipv4", "fsn1-dc14"))).Should(BeNumerically(">", 0.0))
		})

		It("should get prices for correct values for v6", func() {
			By("Checking IPv6 prices")
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-primaryipv6", "fsn1-dc14"))).Should(BeNumerically("==", 0.0))
			Expect(testutil.ToFloat64(sut.GetMonthly().WithLabelValues("test-primaryipv6", "fsn1-dc14"))).Should(BeNumerically("==", 0.0))
		})

		It("should get zero for incorrect values", func() {
			By("Checking IPv4 prices")
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("invalid-name", "fsn1-dc14"))).Should(BeNumerically("==", 0))

			By("Checking IPv6 prices")
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-primaryip", "nbg1-dc14"))).Should(BeNumerically("==", 0))
		})
	})
})
