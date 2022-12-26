package fetcher_test

import (
	"context"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/jangraefen/hcloud-pricing-exporter/fetcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var _ = Describe("For servers", func() {
	sutServer := fetcher.NewServer(&fetcher.PriceProvider{Client: testClient})
	sutBackup := fetcher.NewServer(&fetcher.PriceProvider{Client: testClient})

	BeforeEach(func(ctx context.Context) {
		location, _, err := testClient.Location.GetByName(ctx, "fsn1")
		Expect(err).NotTo(HaveOccurred())

		serverType, _, err := testClient.ServerType.GetByName(ctx, "cx11")
		Expect(err).NotTo(HaveOccurred())

		res, _, err := testClient.Server.Create(ctx, hcloud.ServerCreateOpts{
			Name:       "test-server",
			Labels:     testLabels,
			Location:   location,
			ServerType: serverType,
		})
		Expect(err).ShouldNot(HaveOccurred())
		DeferCleanup(testClient.Server.Delete, res.Server)

		Eventually(func() (hcloud.ServerStatus, error) {
			server, _, serverErr := testClient.Server.GetByID(ctx, res.Server.ID)
			if serverErr != nil {
				return hcloud.ServerStatusOff, err
			}

			return server.Status, nil
		}).
			WithTimeout(1 * time.Minute).
			WithPolling(5 * time.Second).
			Should(Equal(hcloud.ServerStatusRunning))

		action, _, err := testClient.Server.EnableBackup(ctx, res.Server, "22-02")
		waitUntilActionSucceeds(ctx, action.ID)

		Expect(action.Status).Should(Equal(hcloud.ActionStatusSuccess))
	})

	When("getting prices", func() {
		It("should fetch them", func() {
			Expect(sutServer.Run(testClient)).To(Succeed())
		})

		It("should get prices for correct values", func() {
			Eventually(testutil.ToFloat64(sutServer.GetHourly().WithLabelValues("test-server", "fsn1", "cx11"))).Should(BeNumerically(">", 0.0))
			Eventually(testutil.ToFloat64(sutServer.GetMonthly().WithLabelValues("test-server", "fsn1", "cx11"))).Should(BeNumerically(">", 0.0))

			Eventually(testutil.ToFloat64(sutBackup.GetHourly().WithLabelValues("test-server", "fsn1", "cx11"))).Should(BeNumerically(">", 0.0))
			Eventually(testutil.ToFloat64(sutBackup.GetMonthly().WithLabelValues("test-server", "fsn1", "cx11"))).Should(BeNumerically(">", 0.0))
		})

		It("should get zero for incorrect values", func() {
			Eventually(testutil.ToFloat64(sutServer.GetHourly().WithLabelValues("invalid-name", "fsn1", "cx11"))).Should(BeNumerically("==", 0))
			Eventually(testutil.ToFloat64(sutServer.GetHourly().WithLabelValues("test-server", "nbg1", "cx11"))).Should(BeNumerically("==", 0))
			Eventually(testutil.ToFloat64(sutServer.GetHourly().WithLabelValues("test-server", "fsn1", "cx21"))).Should(BeNumerically("==", 0))

			Eventually(testutil.ToFloat64(sutBackup.GetHourly().WithLabelValues("invalid-name", "fsn1", "cx11"))).Should(BeNumerically("==", 0))
			Eventually(testutil.ToFloat64(sutBackup.GetHourly().WithLabelValues("test-server", "nbg1", "cx11"))).Should(BeNumerically("==", 0))
			Eventually(testutil.ToFloat64(sutBackup.GetHourly().WithLabelValues("test-server", "fsn1", "cx21"))).Should(BeNumerically("==", 0))
		})
	})
})
