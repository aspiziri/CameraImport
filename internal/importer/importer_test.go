package importer

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// writeSourceFile creates a file with the given content and a fixed mtime so
// that every file resolves to the same capture date (and thus the same
// destination filename).
func writeSourceFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	mtime := time.Date(2026, 1, 2, 15, 4, 5, 0, time.Local)
	if err := os.Chtimes(path, mtime, mtime); err != nil {
		t.Fatal(err)
	}
	return path
}

func importOne(t *testing.T, s *ImportService, path, destBase string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	file := FileInfo{
		Name: filepath.Base(path),
		Path: path,
		Size: info.Size(),
		Date: "2026-01-02",
		Type: "image",
	}
	if err := s.importFile(file, destBase, "shoot"); err != nil {
		t.Fatalf("importFile(%s): %v", path, err)
	}
}

func readDestDir(t *testing.T, destBase string) map[string]string {
	t.Helper()
	destDir := filepath.Join(destBase, "shoot")
	entries, err := os.ReadDir(destDir)
	if err != nil {
		t.Fatal(err)
	}
	contents := make(map[string]string)
	for _, e := range entries {
		data, err := os.ReadFile(filepath.Join(destDir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		contents[e.Name()] = string(data)
	}
	return contents
}

// Burst shots share a capture timestamp; each must land in its own file with
// no overwrites.
func TestImportFileCollisionDoesNotOverwrite(t *testing.T) {
	s := NewImportService()
	src := t.TempDir()
	dest := t.TempDir()

	paths := []string{
		writeSourceFile(t, src, "a.jpg", "first"),
		writeSourceFile(t, src, "b.jpg", "second, different content"),
		writeSourceFile(t, src, "c.jpg", "third, also different content here"),
	}
	for _, p := range paths {
		importOne(t, s, p, dest)
	}

	contents := readDestDir(t, dest)
	if len(contents) != 3 {
		t.Fatalf("expected 3 files at destination, got %d: %v", len(contents), contents)
	}
	seen := make(map[string]bool)
	for name, body := range contents {
		if seen[body] {
			t.Errorf("duplicate content at destination, a file was overwritten: %v", contents)
		}
		seen[body] = true
		_ = name
	}
}

// Re-importing the same file must be skipped, not duplicated.
func TestImportFileSkipsIdenticalDuplicate(t *testing.T) {
	s := NewImportService()
	src := t.TempDir()
	dest := t.TempDir()

	path := writeSourceFile(t, src, "a.jpg", "same content")
	importOne(t, s, path, dest)
	importOne(t, s, path, dest)

	contents := readDestDir(t, dest)
	if len(contents) != 1 {
		t.Fatalf("expected identical re-import to be skipped, got %d files: %v", len(contents), contents)
	}
}

// Same size but different bytes must NOT be treated as identical.
func TestFilesAreIdenticalComparesContent(t *testing.T) {
	dir := t.TempDir()
	p1 := writeSourceFile(t, dir, "a.jpg", "AAAA")
	p2 := writeSourceFile(t, dir, "b.jpg", "BBBB")

	identical, err := filesAreIdentical(p1, p2)
	if err != nil {
		t.Fatal(err)
	}
	if identical {
		t.Error("same-size files with different content reported as identical")
	}
}

func TestCopyFileRefusesToOverwrite(t *testing.T) {
	dir := t.TempDir()
	src := writeSourceFile(t, dir, "src.jpg", "new content")
	dst := writeSourceFile(t, dir, "dst.jpg", "existing content")

	if err := copyFile(src, dst); err == nil {
		t.Fatal("copyFile overwrote an existing file")
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "existing content" {
		t.Errorf("existing file was modified: %q", data)
	}
}
