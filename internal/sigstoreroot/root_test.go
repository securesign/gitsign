// Copyright 2026 The Sigstore Authors.
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

package sigstoreroot

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sigstore/sigstore-go/pkg/tuf"
)

func TestTUFOptions(t *testing.T) {
	opts := TUFOptions()
	if opts == nil {
		t.Fatal("TUFOptions() returned nil")
	}
	if opts.CachePath == "" {
		t.Fatal("TUFOptions() CachePath is empty")
	}
}

func TestReadRemoteHint(t *testing.T) {
	tmpDir := t.TempDir()

	remoteHint := struct {
		Mirror string `json:"mirror"`
	}{
		Mirror: "https://custom.mirror.example.com",
	}
	data, err := json.Marshal(remoteHint)
	if err != nil {
		t.Fatalf("failed to marshal remote hint: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "remote.json"), data, 0644); err != nil {
		t.Fatalf("failed to write remote.json: %v", err)
	}

	mirror, err := readRemoteHint(tmpDir)
	if err != nil {
		t.Fatalf("readRemoteHint() error = %v", err)
	}
	if mirror != "https://custom.mirror.example.com" {
		t.Errorf("readRemoteHint() = %q, want %q", mirror, "https://custom.mirror.example.com")
	}
}

func TestTUFOptionsCustomMirrorCachedRoot(t *testing.T) {
	tmpDir := t.TempDir()

	customMirror := "https://tuf.custom.example.com"
	fakeRoot := []byte(`{"signed":{"_type":"root","version":1},"signatures":[]}`)

	remoteHint := struct {
		Mirror string `json:"mirror"`
	}{Mirror: customMirror}
	data, err := json.Marshal(remoteHint)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "remote.json"), data, 0644); err != nil {
		t.Fatal(err)
	}

	mirrorDir := filepath.Join(tmpDir, tuf.URLToPath(customMirror))
	if err := os.MkdirAll(mirrorDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(mirrorDir, "root.json"), fakeRoot, 0644); err != nil {
		t.Fatal(err)
	}

	origHome := os.Getenv("HOME")
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)
	defer func() {
		if origHome != "" {
			os.Setenv("HOME", origHome)
		}
	}()

	homeCache := filepath.Join(fakeHome, ".sigstore", "root")
	if err := os.MkdirAll(homeCache, 0755); err != nil {
		t.Fatal(err)
	}
	remoteData, _ := json.Marshal(remoteHint)
	os.WriteFile(filepath.Join(homeCache, "remote.json"), remoteData, 0644)
	homeMirrorDir := filepath.Join(homeCache, tuf.URLToPath(customMirror))
	os.MkdirAll(homeMirrorDir, 0755)
	os.WriteFile(filepath.Join(homeMirrorDir, "root.json"), fakeRoot, 0644)

	opts := TUFOptions()
	if opts.RepositoryBaseURL != customMirror {
		t.Errorf("RepositoryBaseURL = %q, want %q", opts.RepositoryBaseURL, customMirror)
	}
	if !bytes.Equal(opts.Root, fakeRoot) {
		t.Errorf("Root was not set to cached custom mirror root.json; got default embedded root instead")
	}
}

func TestTUFOptionsDefaultMirrorKeepsEmbeddedRoot(t *testing.T) {
	origHome := os.Getenv("HOME")
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)
	defer func() {
		if origHome != "" {
			os.Setenv("HOME", origHome)
		}
	}()

	opts := TUFOptions()
	defaultRoot := tuf.DefaultRoot()
	if !bytes.Equal(opts.Root, defaultRoot) {
		t.Error("With default mirror, Root should be the embedded production root")
	}
	if opts.RepositoryBaseURL != tuf.DefaultMirror {
		t.Errorf("RepositoryBaseURL = %q, want default mirror %q", opts.RepositoryBaseURL, tuf.DefaultMirror)
	}
}

func TestReadRemoteHintMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := readRemoteHint(tmpDir)
	if err == nil {
		t.Fatal("readRemoteHint() expected error for missing file, got nil")
	}
}

func TestReadRemoteHintInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "remote.json"), []byte("not json"), 0644); err != nil {
		t.Fatalf("failed to write remote.json: %v", err)
	}

	_, err := readRemoteHint(tmpDir)
	if err == nil {
		t.Fatal("readRemoteHint() expected error for invalid JSON, got nil")
	}
}
