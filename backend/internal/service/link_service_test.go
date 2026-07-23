package service

import (
	"errors"
	"testing"

	"go-links/internal/db"
	"go-links/internal/models"
	"go-links/internal/repository"
)

func newTestService(t *testing.T) *LinkService {
	t.Helper()

	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	repo := repository.NewLinkRepository(conn)
	return NewLinkService(repo)
}

func TestValidateShortcut(t *testing.T) {
	tests := []struct {
		name      string
		shortcut  string
		wantValid bool
	}{
		{"valid alphanumeric", "oncall", true},
		{"valid with hyphens", "design-system", true},
		{"valid with numbers", "team42", true},
		{"empty", "", false},
		{"contains space", "on call", false},
		{"contains slash", "on/call", false},
		{"contains underscore", "on_call", false},
		{"contains dot", "on.call", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateShortcut(tt.shortcut)
			if tt.wantValid && err != nil {
				t.Errorf("ValidateShortcut(%q) = %v, want nil", tt.shortcut, err)
			}
			if !tt.wantValid && err == nil {
				t.Errorf("ValidateShortcut(%q) = nil, want error", tt.shortcut)
			}
		})
	}
}

func TestValidateDestination(t *testing.T) {
	tests := []struct {
		name        string
		destination string
		wantValid   bool
	}{
		{"valid https", "https://example.com/docs", true},
		{"valid http", "http://example.com", true},
		{"empty", "", false},
		{"missing scheme", "example.com", false},
		{"unsupported scheme", "ftp://example.com", false},
		{"not a url", "not a url", false},
		{"scheme only", "https://", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDestination(tt.destination)
			if tt.wantValid && err != nil {
				t.Errorf("ValidateDestination(%q) = %v, want nil", tt.destination, err)
			}
			if !tt.wantValid && err == nil {
				t.Errorf("ValidateDestination(%q) = nil, want error", tt.destination)
			}
		})
	}
}

func TestCreateLink_Success(t *testing.T) {
	svc := newTestService(t)

	link, err := svc.CreateLink(models.CreateLinkRequest{
		Shortcut:    "oncall",
		Destination: "https://example.com/oncall",
	})
	if err != nil {
		t.Fatalf("CreateLink() error = %v, want nil", err)
	}
	if link.Shortcut != "oncall" {
		t.Errorf("link.Shortcut = %q, want %q", link.Shortcut, "oncall")
	}
	if link.Destination != "https://example.com/oncall" {
		t.Errorf("link.Destination = %q, want %q", link.Destination, "https://example.com/oncall")
	}
	if link.ClickCount != 0 {
		t.Errorf("link.ClickCount = %d, want 0", link.ClickCount)
	}
	if link.ID == 0 {
		t.Error("link.ID = 0, want a nonzero id")
	}
}

func TestCreateLink_InvalidShortcut(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.CreateLink(models.CreateLinkRequest{
		Shortcut:    "bad shortcut!",
		Destination: "https://example.com",
	})
	if err == nil {
		t.Fatal("CreateLink() error = nil, want validation error")
	}

	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("CreateLink() error = %v, want *ValidationError", err)
	}
	if valErr.Field != "shortcut" {
		t.Errorf("valErr.Field = %q, want %q", valErr.Field, "shortcut")
	}
}

func TestCreateLink_InvalidDestination(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.CreateLink(models.CreateLinkRequest{
		Shortcut:    "docs",
		Destination: "not-a-url",
	})
	if err == nil {
		t.Fatal("CreateLink() error = nil, want validation error")
	}

	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("CreateLink() error = %v, want *ValidationError", err)
	}
	if valErr.Field != "destination" {
		t.Errorf("valErr.Field = %q, want %q", valErr.Field, "destination")
	}
}

func TestCreateLink_DuplicateShortcut(t *testing.T) {
	svc := newTestService(t)

	req := models.CreateLinkRequest{Shortcut: "docs", Destination: "https://example.com/docs"}
	if _, err := svc.CreateLink(req); err != nil {
		t.Fatalf("first CreateLink() error = %v, want nil", err)
	}

	_, err := svc.CreateLink(models.CreateLinkRequest{
		Shortcut:    "docs",
		Destination: "https://example.com/other-docs",
	})
	if err == nil {
		t.Fatal("second CreateLink() error = nil, want duplicate error")
	}
	if !errors.Is(err, ErrShortcut) {
		t.Errorf("CreateLink() error = %v, want ErrDuplicateShortcut", err)
	}
}

func TestListLinks_OrderedByMostRecent(t *testing.T) {
	svc := newTestService(t)

	first, err := svc.CreateLink(models.CreateLinkRequest{Shortcut: "first", Destination: "https://example.com/1"})
	if err != nil {
		t.Fatalf("CreateLink(first) error = %v", err)
	}
	second, err := svc.CreateLink(models.CreateLinkRequest{Shortcut: "second", Destination: "https://example.com/2"})
	if err != nil {
		t.Fatalf("CreateLink(second) error = %v", err)
	}

	links, err := svc.ListLinks()
	if err != nil {
		t.Fatalf("ListLinks() error = %v", err)
	}
	if len(links) != 2 {
		t.Fatalf("ListLinks() returned %d links, want 2", len(links))
	}
	if links[0].ID != second.ID || links[1].ID != first.ID {
		t.Errorf("ListLinks() order = [%d, %d], want most-recent-first [%d, %d]",
			links[0].ID, links[1].ID, second.ID, first.ID)
	}
}

func TestDeleteLink(t *testing.T) {
	svc := newTestService(t)

	link, err := svc.CreateLink(models.CreateLinkRequest{Shortcut: "temp", Destination: "https://example.com"})
	if err != nil {
		t.Fatalf("CreateLink() error = %v", err)
	}

	if err := svc.DeleteLink(link.ID); err != nil {
		t.Fatalf("DeleteLink() error = %v, want nil", err)
	}

	_, err = svc.GetLink(link.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("GetLink() after delete error = %v, want ErrNotFound", err)
	}
}

func TestDeleteLink_NotFound(t *testing.T) {
	svc := newTestService(t)

	err := svc.DeleteLink(9999)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("DeleteLink() error = %v, want ErrNotFound", err)
	}
}

func TestResolveShortcut_IncrementsClickCount(t *testing.T) {
	svc := newTestService(t)

	link, err := svc.CreateLink(models.CreateLinkRequest{Shortcut: "docs", Destination: "https://example.com/docs"})
	if err != nil {
		t.Fatalf("CreateLink() error = %v", err)
	}

	resolved, err := svc.ResolveShortcut("docs")
	if err != nil {
		t.Fatalf("ResolveShortcut() error = %v, want nil", err)
	}
	if resolved.ClickCount != link.ClickCount+1 {
		t.Errorf("resolved.ClickCount = %d, want %d", resolved.ClickCount, link.ClickCount+1)
	}

	resolved, err = svc.ResolveShortcut("docs")
	if err != nil {
		t.Fatalf("second ResolveShortcut() error = %v, want nil", err)
	}
	if resolved.ClickCount != 2 {
		t.Errorf("resolved.ClickCount = %d, want 2", resolved.ClickCount)
	}
}

func TestResolveShortcut_NotFound(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.ResolveShortcut("nope")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("ResolveShortcut() error = %v, want ErrNotFound", err)
	}
}
