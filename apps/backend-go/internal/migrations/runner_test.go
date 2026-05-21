package migrations

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReturnsSQLMigrationsSorted(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"20250102000000_second.sql": "SELECT 2;",
		"README.md":                 "ignore",
		"20250101000000_first.sql":  "SELECT 1;",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatalf("write fixture %s: %v", name, err)
		}
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("load migrations: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 migrations, got %d", len(loaded))
	}
	if loaded[0].Version != "20250101000000_first" || loaded[1].Version != "20250102000000_second" {
		t.Fatalf("unexpected order: %#v", loaded)
	}
}

func TestLoadMissingDirectory(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing"))
	if err != ErrNoMigrationDir {
		t.Fatalf("expected ErrNoMigrationDir, got %v", err)
	}
}
