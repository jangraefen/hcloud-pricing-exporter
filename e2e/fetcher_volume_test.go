package e2e_test

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/jangraefen/hcloud-pricing-exporter/fetcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var _ = Describe("For volumes", Ordered, Label("volumes"), func() {
	sut := fetcher.NewVolume(&fetcher.PriceProvider{Client: testClient}, "suite")

	BeforeAll(func(ctx context.Context) {
		location, _, err := testClient.Location.GetByName(ctx, "fsn1")
		Expect(err).NotTo(HaveOccurred())

		res, _, err := testClient.Volume.Create(ctx, hcloud.VolumeCreateOpts{
			Name:     ("test-volume"),
			Labels:   testLabels,
			Location: location,
			Size:     10,
		})
		Expect(err).ShouldNot(HaveOccurred())
		DeferCleanup(testClient.Volume.Delete, res.Volume)

		waitUntilActionSucceeds(ctx, res.Action)
	})

	//nolint:dupl
	When("getting prices", func() {
		It("should fetch them", func() {
			Expect(sut.Run(testClient)).To(Succeed())
		})

		It("should get prices for correct values", func() {
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-volume", "fsn1", "10", "e2e_suite_test"))).Should(BeNumerically(">", 0.0))
			Expect(testutil.ToFloat64(sut.GetMonthly().WithLabelValues("test-volume", "fsn1", "10", "e2e_suite_test"))).Should(BeNumerically(">", 0.0))
		})

		It("should get zero for incorrect values", func() {
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("invalid-name", "fsn1", "10", "e2e_suite_test"))).Should(BeNumerically("==", 0))
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-volume", "nbg1", "10", "e2e_suite_test"))).Should(BeNumerically("==", 0))
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-volume", "fsn1", "99", "e2e_suite_test"))).Should(BeNumerically("==", 0))
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-volume", "fsn1", "10", "e3e_suite_test"))).Should(BeNumerically("==", 0))
		})
	})
})
