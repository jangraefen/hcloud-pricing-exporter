package e2e_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/gomega"
	"golang.org/x/crypto/ssh"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func hcloudAPITokenFromENV() string {
	if token, ok := os.LookupEnv("HCLOUD_API_TOKEN"); ok {
		return token
	}

	panic(fmt.Errorf("environment variable HCLOUD_API_TOKEN not set, but required"))
}

func waitUntilActionSucceeds(ctx context.Context, actionToTrack *hcloud.Action) {
	if actionToTrack != nil {
		Eventually(func() (hcloud.ActionStatus, error) {
			action, _, err := testClient.Action.GetByID(ctx, actionToTrack.ID)
			if err != nil {
				return hcloud.ActionStatusError, err
			}

			return action.Status, nil
		}).
			WithOffset(1).
			Within(1 * time.Minute).
			ProbeEvery(5 * time.Second).
			Should(Equal(hcloud.ActionStatusSuccess))
	}
}

func generatePublicKey() string {
	public, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	sshKey, err := ssh.NewPublicKey(public)
	if err != nil {
		panic(err)
	}

	return string(ssh.MarshalAuthorizedKey(sshKey))
}
