package e2e_test

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/jangraefen/hcloud-pricing-exporter/fetcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var _ = Describe("For loadbalancers", Ordered, Label("loadbalancers"), func() {
	sut := fetcher.NewLoadbalancer(&fetcher.PriceProvider{Client: testClient}, "suite")

	BeforeAll(func(ctx context.Context) {
		location, _, err := testClient.Location.GetByName(ctx, "fsn1")
		Expect(err).NotTo(HaveOccurred())

		lbType, _, err := testClient.LoadBalancerType.GetByName(ctx, "lb11")
		Expect(err).NotTo(HaveOccurred())

		res, _, err := testClient.LoadBalancer.Create(ctx, hcloud.LoadBalancerCreateOpts{
			Name:             "test-loadbalancer",
			Labels:           testLabels,
			Location:         location,
			LoadBalancerType: lbType,
		})
		Expect(err).ShouldNot(HaveOccurred())
		DeferCleanup(testClient.LoadBalancer.Delete, res.LoadBalancer)

		waitUntilActionSucceeds(ctx, res.Action)
	})

	//nolint:dupl
	When("getting prices", func() {
		It("should fetch them", func() {
			Expect(sut.Run(testClient)).To(Succeed())
		})

		It("should get prices for correct values", func() {
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-loadbalancer", "fsn1", "lb11", "e2e_suite_test"))).Should(BeNumerically(">", 0.0))
			Expect(testutil.ToFloat64(sut.GetMonthly().WithLabelValues("test-loadbalancer", "fsn1", "lb11", "e2e_suite_test"))).Should(BeNumerically(">", 0.0))
		})

		It("should get zero for incorrect values", func() {
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("invalid-name", "fsn1", "lb11", "e2e_suite_test"))).Should(BeNumerically("==", 0))
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-loadbalancer", "nbg1", "lb11", "e2e_suite_test"))).Should(BeNumerically("==", 0))
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-loadbalancer", "fsn1", "lb21", "e2e_suite_test"))).Should(BeNumerically("==", 0))
			Expect(testutil.ToFloat64(sut.GetHourly().WithLabelValues("test-loadbalancer", "fsn1", "lb11", "e3e_suite_test"))).Should(BeNumerically("==", 0))
		})
	})
})
