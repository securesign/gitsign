// Copyright 2024 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package sigstoreroot loads the Sigstore trusted root via the sigstore-go TUF client.
package sigstoreroot

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sigstore/cosign/v3/pkg/cosign"
	"github.com/sigstore/sigstore-go/pkg/root"
	"github.com/sigstore/sigstore-go/pkg/tuf"
	sigstoretuf "github.com/sigstore/sigstore/pkg/tuf"
)

// SigningConfig represents the service URLs from signing_config.json.
type SigningConfig struct {
	FulcioURL    string
	RekorURL     string
	OIDCURL      string
	TimestampURL string
}

// CacheExists returns true if the TUF cache exists and has been initialized.
func CacheExists() bool {
	opts := tuf.DefaultOptions()
	if mirror, err := readRemoteHint(opts.CachePath); err == nil && mirror != "" {
		signingConfigPath := filepath.Join(opts.CachePath, tuf.URLToPath(mirror), "targets", "signing_config.v0.2.json")
		if _, err := os.Stat(signingConfigPath); err == nil {
			return true
		}
	}
	return false
}

// GetSigningConfig loads the signing configuration from the TUF cache.
// Returns nil if the cache doesn't exist or signing config is not available.
func GetSigningConfig() (*SigningConfig, error) {
	opts := tuf.DefaultOptions()
	mirror, err := readRemoteHint(opts.CachePath)
	if err != nil {
		return nil, nil
	}
	if mirror == "" {
		mirror = tuf.DefaultMirror
	}

	signingConfigPath := filepath.Join(opts.CachePath, tuf.URLToPath(mirror), "targets", "signing_config.v0.2.json")
	data, err := os.ReadFile(signingConfigPath)
	if err != nil {
		return nil, nil
	}

	var raw signingConfigJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing signing config: %w", err)
	}

	config := &SigningConfig{}
	if url := selectValidURL(raw.CAUrls); url != "" {
		config.FulcioURL = url
	}
	if url := selectValidURL(raw.RekorTlogUrls); url != "" {
		config.RekorURL = url
	}
	if url := selectValidURL(raw.OIDCUrls); url != "" {
		config.OIDCURL = url
	}
	if url := selectValidURL(raw.TSAUrls); url != "" {
		config.TimestampURL = url
	}

	return config, nil
}

type signingConfigJSON struct {
	CAUrls        []serviceURL `json:"caUrls"`
	OIDCUrls      []serviceURL `json:"oidcUrls"`
	RekorTlogUrls []serviceURL `json:"rekorTlogUrls"`
	TSAUrls       []serviceURL `json:"tsaUrls"`
}

type serviceURL struct {
	URL      string   `json:"url"`
	ValidFor validity `json:"validFor"`
}

type validity struct {
	Start string `json:"start"`
	End   string `json:"end,omitempty"`
}

func selectValidURL(urls []serviceURL) string {
	now := time.Now()
	for _, u := range urls {
		start, err := time.Parse(time.RFC3339, u.ValidFor.Start)
		if err != nil {
			continue
		}
		if now.Before(start) {
			continue
		}
		if u.ValidFor.End != "" {
			end, err := time.Parse(time.RFC3339, u.ValidFor.End)
			if err != nil {
				continue
			}
			if now.After(end) {
				continue
			}
		}
		return u.URL
	}
	if len(urls) > 0 {
		return urls[0].URL
	}
	return ""
}

// TUFOptions returns sigstore-go TUF options, reading the mirror URL and
// cached root.json from the local TUF cache when a custom mirror is configured.
func TUFOptions() *tuf.Options {
	opts := tuf.DefaultOptions()
	if mirror, err := readRemoteHint(opts.CachePath); err == nil && mirror != "" {
		opts.RepositoryBaseURL = mirror
	}
	if opts.RepositoryBaseURL != tuf.DefaultMirror {
		cachedRoot := filepath.Join(opts.CachePath, tuf.URLToPath(opts.RepositoryBaseURL), "root.json")
		if rootBytes, err := os.ReadFile(cachedRoot); err == nil {
			opts.Root = rootBytes
		}
	}
	return opts
}

// FetchTrustedRoot loads the Sigstore trusted root from the TUF cache.
func FetchTrustedRoot() (*root.TrustedRoot, error) {
	return root.FetchTrustedRootWithOptions(TUFOptions())
}

// GetCTLogPubs returns CT log public keys from the trusted root.
func GetCTLogPubs(trustedRoot *root.TrustedRoot) (*cosign.TrustedTransparencyLogPubKeys, error) {
	return transparencyLogPubKeys(trustedRoot.CTLogs())
}

// GetRekorPubs returns Rekor transparency log public keys from the trusted root.
func GetRekorPubs(trustedRoot *root.TrustedRoot) (*cosign.TrustedTransparencyLogPubKeys, error) {
	return transparencyLogPubKeys(trustedRoot.RekorLogs())
}

// FulcioCertificates extracts Fulcio root and intermediate certificates from the trusted root.
func FulcioCertificates(trustedRoot *root.TrustedRoot) ([]*x509.Certificate, error) {
	cas := trustedRoot.FulcioCertificateAuthorities()
	if len(cas) == 0 {
		return nil, fmt.Errorf("no Fulcio certificate authorities found in trusted root")
	}

	var certs []*x509.Certificate
	for _, ca := range cas {
		fca, ok := ca.(*root.FulcioCertificateAuthority)
		if !ok {
			continue
		}
		if fca.Root != nil {
			certs = append(certs, fca.Root)
		}
		certs = append(certs, fca.Intermediates...)
	}
	if len(certs) == 0 {
		return nil, fmt.Errorf("no Fulcio certificates found in trusted root")
	}
	return certs, nil
}

// GetTSACertificates extracts TSA (Timestamp Authority) certificates from the trusted root.
func GetTSACertificates(trustedRoot *root.TrustedRoot) ([]*x509.Certificate, error) {
	tsas := trustedRoot.TimestampingAuthorities()
	if len(tsas) == 0 {
		return nil, nil
	}

	var certs []*x509.Certificate
	for _, tsa := range tsas {
		sta, ok := tsa.(*root.SigstoreTimestampingAuthority)
		if !ok {
			continue
		}
		if sta.Root != nil {
			certs = append(certs, sta.Root)
		}
		certs = append(certs, sta.Intermediates...)
		if sta.Leaf != nil {
			certs = append(certs, sta.Leaf)
		}
	}
	return certs, nil
}

func transparencyLogPubKeys(logs map[string]*root.TransparencyLog) (*cosign.TrustedTransparencyLogPubKeys, error) {
	pubKeys := cosign.NewTrustedTransparencyLogPubKeys()

	for logID, log := range logs {
		if log.PublicKey == nil {
			continue
		}

		status := sigstoretuf.Active
		if !log.ValidityPeriodEnd.IsZero() && time.Now().After(log.ValidityPeriodEnd) {
			status = sigstoretuf.Expired
		}

		pubKeys.Keys[logID] = cosign.TransparencyLogPubKey{
			PubKey: log.PublicKey,
			Status: status,
		}
	}

	if len(pubKeys.Keys) == 0 {
		return nil, fmt.Errorf("no transparency log public keys found")
	}
	return &pubKeys, nil
}

func readRemoteHint(cachePath string) (string, error) {
	data, err := os.ReadFile(filepath.Join(cachePath, "remote.json"))
	if err != nil {
		return "", err
	}
	var remote struct {
		Mirror string `json:"mirror"`
	}
	if err := json.Unmarshal(data, &remote); err != nil {
		return "", err
	}
	return remote.Mirror, nil
}
