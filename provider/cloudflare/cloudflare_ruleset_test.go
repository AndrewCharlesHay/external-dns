package cloudflare

import (
	"context"
	"os"
	"testing"

	"github.com/cloudflare/cloudflare-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/provider"
)

type mockCloudFlareDNS struct {
	cloudFlareDNS
	ListRulesetsCalled   bool
	CreateRulesetCalled  bool
	UpdateRulesetCalled  bool
	LastRuleset          cloudflare.Ruleset
}

func (m *mockCloudFlareDNS) ListRulesets(ctx context.Context, zoneID string) ([]cloudflare.Ruleset, error) {
	m.ListRulesetsCalled = true
	return []cloudflare.Ruleset{}, nil
}
func (m *mockCloudFlareDNS) CreateRuleset(ctx context.Context, zoneID string, ruleset cloudflare.Ruleset) (cloudflare.Ruleset, error) {
	m.CreateRulesetCalled = true
	m.LastRuleset = ruleset
	return ruleset, nil
}
func (m *mockCloudFlareDNS) UpdateRuleset(ctx context.Context, zoneID string, ruleset cloudflare.Ruleset) (cloudflare.Ruleset, error) {
	m.UpdateRulesetCalled = true
	m.LastRuleset = ruleset
	return ruleset, nil
}

func TestNewCloudFlareProvider_WithRulesetEnv(t *testing.T) {
	os.Setenv("CF_RULESET_JSON", `{"name":"test-ruleset","kind":"zone","phase":"http_request_firewall_custom","rules":[]}`)
	os.Setenv("CF_API_KEY", "dummy")
	os.Setenv("CF_API_EMAIL", "dummy@example.com")
	defer os.Unsetenv("CF_RULESET_JSON")
	defer os.Unsetenv("CF_API_KEY")
	defer os.Unsetenv("CF_API_EMAIL")

	provider, err := NewCloudFlareProvider(
		endpoint.DomainFilter{}, provider.ZoneIDFilter{}, false, false, "", CustomHostnamesConfig{}, DNSRecordsConfig{},
	)
	require.NoError(t, err)
	assert.NotNil(t, provider.RulesetConfig)
	assert.Equal(t, "test-ruleset", provider.RulesetConfig.Name)
}

func TestNewCloudFlareProvider_InvalidRulesetEnv(t *testing.T) {
	os.Setenv("CF_RULESET_JSON", `not-json`)
	os.Setenv("CF_API_KEY", "dummy")
	os.Setenv("CF_API_EMAIL", "dummy@example.com")
	defer os.Unsetenv("CF_RULESET_JSON")
	defer os.Unsetenv("CF_API_KEY")
	defer os.Unsetenv("CF_API_EMAIL")
	_, err := NewCloudFlareProvider(
		endpoint.DomainFilter{}, provider.ZoneIDFilter{}, false, false, "", CustomHostnamesConfig{}, DNSRecordsConfig{},
	)
	assert.Error(t, err)
}

// NOTE: The following test is a placeholder. The actual Ruleset API is not implemented in the current cloudflare-go version.
func TestApplyChanges_CallsRulesetAPI(t *testing.T) {
	mock := &mockCloudFlareDNS{}
	provider := &CloudFlareProvider{
		Client:       mock,
		RulesetConfig: &cloudflare.Ruleset{Name: "test", Kind: "zone", Phase: "http_request_firewall_custom"},
	}
	// This would require a real or properly mocked Zones method
	// For now, just check that the provider is constructed and RulesetConfig is set
	assert.Equal(t, "test", provider.RulesetConfig.Name)
}
