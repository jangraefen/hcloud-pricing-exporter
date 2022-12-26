package fetcher_test

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/jangraefen/hcloud-pricing-exporter/fetcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var _ = Describe("For volumes", func() {
	sut := fetcher.NewVolume(&fetcher.PriceProvider{Client: testClient})

	BeforeEach(func(ctx context.Context) {
		location, _, err := testClient.Location.GetByName(ctx, "fsn1")
		Expect(err).NotTo(HaveOccurred())

		res, _, err := testClient.Volume.Create(ctx, hcloud.VolumeCreateOpts{
			Name:     ("test-volume"),
			Labels:   testLabels,
			Location: location,
			Size:     10737418240,
		})
		Expect(err).ShouldNot(HaveOccurred())
		DeferCleanup(testClient.Volume.Delete, res.Volume)
	})

	When("getting prices", func() {
		It("should fetch them", func() {
			Expect(sut.Run(testClient)).To(Succeed())
		})

		It("should get prices for correct values", func() {
			Eventually(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-volume", "fsn1", "10737418240"))).Should(BeNumerically(">", 0.0))
			Eventually(testutil.ToFloat64(sut.GetMonthly().WithLabelValues("test-volume", "fsn1", "10737418240"))).Should(BeNumerically(">", 0.0))
		})

		It("should get zero for incorrect values", func() {
			Eventually(testutil.ToFloat64(sut.GetHourly().WithLabelValues("invalid-name", "fsn1", "10737418240"))).Should(BeNumerically("==", 0))
			Eventually(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-volume", "nbg1", "10737418240"))).Should(BeNumerically("==", 0))
			Eventually(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-volume", "fsn1", "11111111111"))).Should(BeNumerically("==", 0))
		})
	})
})
