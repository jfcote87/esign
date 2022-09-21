package ratelimit_test

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/ratelimit"
	"github.com/jfcote87/esign/v2.1/envelopes"
)

func ExampleCredential() {
	ctx := context.TODO()
	apiUserID := "78e5a047-f767-41f8-8dbd-10e3eed65c55"
	cred, err := getCredential(ctx, apiUserID)
	if err != nil {
		log.Fatalf("credential error: %v", err)
	}

	rateLimitCred := &ratelimit.Credential{
		Credential:    cred,
		ReportHandler: globalRateLimitHandler,
	}
	sv := envelopes.New(rateLimitCred)
	envInfo, err := sv.ListStatusChanges().FolderIds("F1", "F2").Do(ctx)
	if err != nil {
		log.Fatalf("listing error %v", err)
	}
	for _, info := range envInfo.Envelopes {
		log.Printf("%s %s", info.Recipients.Signers[0].Name, info.Status)
	}
	if globalRateLimitHandler.MustThrottle() {
		log.Printf("currently throttling esign calls")
	}
	if rlr := globalRateLimitHandler.Report(); rlr != nil {
		log.Printf("Rate Limit: %d", rlr.RateLimit)
		log.Printf("Rate Limit Remaing: %d", rlr.RateRemaining)
		log.Printf("Reset Time: %s", rlr.ResetAt().Format("15:04:05"))
		log.Printf("Rate Limit: %d", rlr.BurstLimit)
		log.Printf("Rate Limit: %d", rlr.BurstRemaining)
	}
}

var globalRateLimitHandler = &RLHandler{}

type RLHandler struct {
	m             sync.Mutex
	lastestReport *ratelimit.Report
	mustThrottle  bool
}

func (rlc *RLHandler) Report() *ratelimit.Report {
	rlc.m.Lock()
	defer rlc.m.Unlock()
	return rlc.lastestReport
}

func (rlc *RLHandler) MustThrottle() bool {
	rlc.m.Lock()
	defer rlc.m.Unlock()
	return rlc.mustThrottle
}

func (rlc *RLHandler) Handle(ctx context.Context, res *http.Response) error {

	rpt := ratelimit.New(res.Header)
	mustThrottle := rpt.RateRemaining < 1000 && rpt.ResetAt().Add(-10*time.Minute).After(time.Now())

	rlc.m.Lock()
	defer rlc.m.Unlock()
	rlc.lastestReport = rpt
	if mustThrottle != rlc.mustThrottle {
		rlc.mustThrottle = mustThrottle

		switch mustThrottle {
		case true:
			log.Printf("send messages to throttle until reset time")
		case false:
			log.Printf("send messages to end throttling")
		}
	}
	return nil
}

func getCredential(ctx context.Context, apiUser string) (*esign.OAuth2Credential, error) {
	cfg := &esign.JWTConfig{
		IntegratorKey: "51d1a791-489c-4622-b743-19e0bd6f359e",
		KeyPairID:     "1f57a66f-cc71-45cd-895e-9aaf39d5ccc4",
		PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
		MIICWwIBAAKBgGeDMVfH1+RBI/JMPfrIJDmBWxJ/wGQRjPUm6Cu2aHXMqtOjK3+u
		.........
		ZS1NWRHL8r7hdJL8lQYZPfNqyYcW7C1RW3vWbCRGMA==
		-----END RSA PRIVATE KEY-----`,
		AccountID: "c23357a7-4f00-47f5-8802-94d2b1fb9a29",
		IsDemo:    true,
	}

	return cfg.Credential(apiUser, nil, nil)

}
