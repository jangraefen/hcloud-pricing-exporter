//nolint:revive
package e2e_test

import (
	"context"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/jangraefen/hcloud-pricing-exporter/fetcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

const (
	serverTypeName = "cx22"
)

var _ = Describe("For servers", Ordered, Label("servers"), func() {
	sutServer := fetcher.NewServer(&fetcher.PriceProvider{Client: testClient}, "suite")
	sutBackup := fetcher.NewServerBackup(&fetcher.PriceProvider{Client: testClient}, "suite")

	BeforeAll(func(ctx context.Context) {
		location, _, err := testClient.Location.GetByName(ctx, "fsn1")
		Expect(err).NotTo(HaveOccurred())

		serverType, _, err := testClient.ServerType.GetByName(ctx, serverTypeName)
		Expect(err).NotTo(HaveOccurred())

		image, _, err := testClient.Image.GetByNameAndArchitecture(ctx, "ubuntu-24.04", hcloud.ArchitectureX86)
		Expect(err).NotTo(HaveOccurred())

		sshKey, _, err := testClient.SSHKey.Create(ctx, hcloud.SSHKeyCreateOpts{
			Name:      "test-key",
			Labels:    testLabels,
			PublicKey: generatePublicKey(),
		})
		Expect(err).ShouldNot(HaveOccurred())
		DeferCleanup(testClient.SSHKey.Delete, sshKey)

		By("Setting up a server")
		res, _, err := testClient.Server.Create(ctx, hcloud.ServerCreateOpts{
			Name:       "test-server",
			Labels:     testLabels,
			Location:   location,
			ServerType: serverType,
			Image:      image,
			SSHKeys:    []*hcloud.SSHKey{sshKey},
			PublicNet:  &hcloud.ServerCreatePublicNet{EnableIPv4: false, EnableIPv6: true},
		})
		Expect(err).ShouldNot(HaveOccurred())
		DeferCleanup(testClient.Server.DeleteWithResult, context.Background(), res.Server)

		waitUntilActionSucceeds(ctx, res.Action)

		By("Enabling backups for the server")
		Eventually(func() (hcloud.ServerStatus, error) {
			server, _, serverErr := testClient.Server.GetByID(ctx, res.Server.ID)
			if serverErr != nil {
				return hcloud.ServerStatusOff, err
			}

			return server.Status, nil
		}).
			Within(1 * time.Minute).
			ProbeEvery(5 * time.Second).
			Should(Equal(hcloud.ServerStatusRunning))

		action, _, err := testClient.Server.EnableBackup(ctx, res.Server, "")
		Expect(err).ShouldNot(HaveOccurred())

		waitUntilActionSucceeds(ctx, action)
	})

	//nolint:dupl
	When("getting prices", func() {
		It("should fetch them", func() {
			By("Running the price collection")
			Expect(sutServer.Run(testClient)).To(Succeed())
			Expect(sutBackup.Run(testClient)).To(Succeed())
		})

		It("should get prices for correct values", func() {
			By("Checking server prices")
			Expect(testutil.ToFloat64(sutServer.GetHourly().WithLabelValues("test-server", "fsn1", serverTypeName, "e2e_suite_test"))).Should(BeNumerically(">", 0.0))
			Expect(testutil.ToFloat64(sutServer.GetMonthly().WithLabelValues("test-server", "fsn1", serverTypeName, "e2e_suite_test"))).Should(BeNumerically(">", 0.0))

			By("Checking server backup prices")
			Expect(testutil.ToFloat64(sutBackup.GetHourly().WithLabelValues("test-server", "fsn1", serverTypeName, "e2e_suite_test"))).Should(BeNumerically(">", 0.0))
			Expect(testutil.ToFloat64(sutBackup.GetMonthly().WithLabelValues("test-server", "fsn1", serverTypeName, "e2e_suite_test"))).Should(BeNumerically(">", 0.0))
		})

		It("should get zero for incorrect values", func() {
			By("Checking server prices")
			Expect(testutil.ToFloat64(sutServer.GetHourly().WithLabelValues("invalid-name", "fsn1", serverTypeName, "e2e_suite_test"))).Should(BeNumerically("==", 0))
			Expect(testutil.ToFloat64(sutServer.GetHourly().WithLabelValues("test-server", "nbg1", serverTypeName, "e2e_suite_test"))).Should(BeNumerically("==", 0))
			Expect(testutil.ToFloat64(sutServer.GetHourly().WithLabelValues("test-server", "fsn1", "cx21", "e2e_suite_test"))).Should(BeNumerically("==", 0))
			Expect(testutil.ToFloat64(sutServer.GetHourly().WithLabelValues("test-server", "fsn1", serverTypeName, "e3e_suite_test"))).Should(BeNumerically("==", 0))

			By("Checking server backup prices")
			Expect(testutil.ToFloat64(sutBackup.GetHourly().WithLabelValues("invalid-name", "fsn1", serverTypeName, "e2e_suite_test"))).Should(BeNumerically("==", 0))
			Expect(testutil.ToFloat64(sutBackup.GetHourly().WithLabelValues("test-server", "nbg1", serverTypeName, "e2e_suite_test"))).Should(BeNumerically("==", 0))
			Expect(testutil.ToFloat64(sutBackup.GetHourly().WithLabelValues("test-server", "fsn1", "cx21", "e2e_suite_test"))).Should(BeNumerically("==", 0))
			Expect(testutil.ToFloat64(sutBackup.GetHourly().WithLabelValues("test-server", "fsn1", serverTypeName, "e3e_suite_test"))).Should(BeNumerically("==", 0))
		})
	})
})
