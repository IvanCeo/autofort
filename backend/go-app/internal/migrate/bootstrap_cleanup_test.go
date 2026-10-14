package migrate

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestBootstrapCleanupAfterCreationFailures(t *testing.T) {
	const name, token = "autofort-bootstrap-injected", "owned-nonce"
	id := strings.Repeat("a", 64)
	for _, test := range []struct {
		name   string
		stdout string
		err    error
	}{
		{"committed creation with client EOF", "", io.EOF},
		{"startup failed after ID returned", id + "\n", errors.New("startup failed")},
		{"missing ID", "", nil},
		{"invalid ID", "not-an-id", nil},
		{"successful creation followed by validation failure", id + "\n", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			var cleanup func() error
			created, removed := false, false
			run := func(ctx context.Context, args ...string) (string, error) {
				switch args[0] {
				case "run":
					if cleanup == nil {
						t.Fatal("cleanup was not registered before creation")
					}
					created = true // The server committed a container before failing.
					return test.stdout, test.err
				case "inspect":
					if !created || args[len(args)-1] != name || ctx.Err() != nil {
						t.Fatalf("cleanup must inspect its own name using a live context: %v", args)
					}
					return cleanupInspection(t, name, token, id), nil
				case "rm":
					if want := []string{"rm", "--force", "--volumes", id}; !reflect.DeepEqual(args, want) {
						t.Fatalf("unexpected removal: %v", args)
					}
					removed = true
					return id, nil
				default:
					t.Fatalf("unexpected Docker call: %v", args)
					return "", nil
				}
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			returnedID, err := createBootstrapContainer(ctx, run, name, token, []string{"run"}, func(fn func() error) {
				if cleanup != nil {
					t.Fatal("cleanup registered twice")
				}
				cleanup = fn
			})
			if test.err != nil || !isDockerContainerID(strings.TrimSpace(test.stdout)) {
				if err == nil || returnedID != "" {
					t.Fatalf("expected creation/ID failure: id=%q err=%v", returnedID, err)
				}
			} else if err != nil || returnedID != id {
				t.Fatalf("expected valid captured ID: id=%q err=%v", returnedID, err)
			}
			cancel() // Cleanup must work even after the creation context expires.
			if err := cleanup(); err != nil {
				t.Fatal(err)
			}
			if !removed {
				t.Fatal("committed disposable container was leaked")
			}
		})
	}
}

func TestBootstrapCleanupRefusesWrongOwnership(t *testing.T) {
	const name, token = "autofort-bootstrap-injected", "owned-nonce"
	id := strings.Repeat("a", 64)
	for _, test := range []struct {
		name       string
		actualName string
		label      string
		actualID   string
		expectedID string
	}{
		{"wrong name", "different-container", token, id, ""},
		{"wrong owner", name, "other-nonce", id, ""},
		{"missing owner", name, "", id, ""},
		{"returned ID mismatch", name, token, strings.Repeat("b", 64), id},
		{"invalid inspected ID", name, token, "invalid", ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			run := func(_ context.Context, args ...string) (string, error) {
				if !reflect.DeepEqual(args, []string{"inspect", "--type", "container", name}) {
					t.Fatalf("cleanup must refuse removal without matching ownership: %v", args)
				}
				return cleanupInspection(t, test.actualName, test.label, test.actualID), nil
			}
			if err := cleanupBootstrapContainer(context.Background(), run, name, token, test.expectedID); err == nil {
				t.Fatal("expected ownership refusal")
			}
		})
	}
}

func TestBootstrapCleanupDistinguishesAbsenceFromInspectionFailure(t *testing.T) {
	const name = "autofort-bootstrap-injected"
	for _, test := range []struct {
		name       string
		listing    string
		listingErr error
		wantErr    bool
	}{
		{"creation never committed", "", nil, false},
		{"container exists but cannot be inspected", strings.Repeat("a", 64), nil, true},
		{"daemon unavailable", "", io.EOF, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			run := func(_ context.Context, args ...string) (string, error) {
				calls++
				if calls == 1 && reflect.DeepEqual(args, []string{"inspect", "--type", "container", name}) {
					return "", errors.New("inspect failed")
				}
				want := []string{"container", "ls", "--all", "--filter", "name=^/" + name + "$", "--no-trunc", "--format", "{{.ID}}"}
				if calls != 2 || !reflect.DeepEqual(args, want) {
					t.Fatalf("unexpected lookup or removal: %v", args)
				}
				return test.listing, test.listingErr
			}
			err := cleanupBootstrapContainer(context.Background(), run, name, "owned-nonce", "")
			if (err != nil) != test.wantErr || calls != 2 {
				t.Fatalf("absence check: err=%v calls=%d", err, calls)
			}
		})
	}
}

func TestBootstrapCleanupPreservesExplicitVolumeMounts(t *testing.T) {
	const name, token = "autofort-bootstrap-injected", "owned-nonce"
	id := strings.Repeat("a", 64)
	for _, field := range []string{"Binds", "Mounts"} {
		t.Run(field, func(t *testing.T) {
			var inspected []disposableContainer
			if err := json.Unmarshal([]byte(cleanupInspection(t, name, token, id)), &inspected); err != nil {
				t.Fatal(err)
			}
			if field == "Binds" {
				inspected[0].HostConfig.Binds = []string{"existing-volume:/data"}
			} else {
				inspected[0].HostConfig.Mounts = []json.RawMessage{json.RawMessage(`{"Type":"volume","Source":"existing-volume","Target":"/data"}`)}
			}
			payload, err := json.Marshal(inspected)
			if err != nil {
				t.Fatal(err)
			}
			removed := false
			run := func(_ context.Context, args ...string) (string, error) {
				if args[0] == "inspect" {
					return string(payload), nil
				}
				if !reflect.DeepEqual(args, []string{"rm", "--force", id}) {
					t.Fatalf("existing mounted volumes must never be deleted: %v", args)
				}
				removed = true
				return id, nil
			}
			if err := cleanupBootstrapContainer(context.Background(), run, name, token, id); err != nil || !removed {
				t.Fatalf("cleanup: removed=%v err=%v", removed, err)
			}
		})
	}
}

func TestBootstrapCleanupReportsRemovalFailure(t *testing.T) {
	const name, token = "autofort-bootstrap-injected", "owned-nonce"
	id := strings.Repeat("a", 64)
	run := func(_ context.Context, args ...string) (string, error) {
		if args[0] == "inspect" {
			return cleanupInspection(t, name, token, id), nil
		}
		return "", io.EOF
	}
	if err := cleanupBootstrapContainer(context.Background(), run, name, token, id); !errors.Is(err, io.EOF) {
		t.Fatalf("cleanup error was lost: %v", err)
	}
}

func cleanupInspection(t *testing.T, name, token, id string) string {
	t.Helper()
	container := disposableContainer{ID: id, Name: "/" + name}
	container.Config.Labels = map[string]string{"autofort.disposable-bootstrap": token}
	output, err := json.Marshal([]disposableContainer{container})
	if err != nil {
		t.Fatal(err)
	}
	return string(output)
}
