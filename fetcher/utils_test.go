package fetcher_test

import (
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/gomega"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func hcloudAPITokenFromENV() string {
	if token, ok := os.LookupEnv("HCLOUD_API_TOKEN"); ok {
		return token
	}

	panic(fmt.Errorf("environment variable HCLOUD_API_TOKEN not set, but required"))
}

func waitUntilActionSucceeds(ctx context.Context, actionID int) {
	Eventually(func() (hcloud.ActionStatus, error) {
		action, _, err := testClient.Action.GetByID(ctx, actionID)
		if err != nil {
			return hcloud.ActionStatusError, err
		}

		return action.Status, nil
	}).
		WithOffset(1).
		WithTimeout(1 * time.Minute).
		WithPolling(5 * time.Second).
		Should(Equal(hcloud.ActionStatusSuccess))
}
