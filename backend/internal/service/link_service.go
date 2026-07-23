package service

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"go-links/internal/models"
	"go-links/internal/repository"
)

var shortcutPattern = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

// ValidationError represents a rejected input with a field-specific message.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

var ErrShortcut = repository.ErrShortcut
var ErrNotFound = repository.NoError

type LinkService struct {
	repo *repository.LinkRepository
}

func NewLinkService(repo *repository.LinkRepository) *LinkService {
	return &LinkService{repo: repo}
}

// ValidateShortcut checks that a shortcut is non-empty and contains only
// alphanumeric characters and hyphens.
func ValidateShortcut(shortcut string) error {
	if shortcut == "" {
		return &ValidationError{Field: "shortcut", Message: "shortcut is required"}
	}
	if !shortcutPattern.MatchString(shortcut) {
		return &ValidationError{Field: "shortcut", Message: "shortcut must contain only letters, numbers, and hyphens"}
	}
	return nil
}

// ValidateDestination checks that a destination is a well-formed http(s) URL.
func ValidateDestination(destination string) error {
	if destination == "" {
		return &ValidationError{Field: "destination", Message: "destination is required"}
	}
	u, err := url.ParseRequestURI(destination)
	if err != nil {
		return &ValidationError{Field: "destination", Message: "destination must be a valid URL"}
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return &ValidationError{Field: "destination", Message: "destination must use http or https"}
	}
	if u.Host == "" {
		return &ValidationError{Field: "destination", Message: "destination must be a valid URL"}
	}
	return nil
}

func (s *LinkService) CreateLink(req models.CreateLinkRequest) (*models.Link, error) {
	if err := ValidateShortcut(req.Shortcut); err != nil {
		return nil, err
	}
	if err := ValidateDestination(req.Destination); err != nil {
		return nil, err
	}

	link, err := s.repo.Create(req.Shortcut, req.Destination)
	if err != nil {
		if errors.Is(err, repository.ErrShortcut) {
			return nil, fmt.Errorf("shortcut %q is already in use: %w", req.Shortcut, ErrShortcut)
		}
		return nil, err
	}
	return link, nil
}

func (s *LinkService) ListLinks() ([]*models.Link, error) {
	return s.repo.List()
}

func (s *LinkService) GetLink(id int64) (*models.Link, error) {
	return s.repo.GetByID(id)
}

func (s *LinkService) DeleteLink(id int64) error {
	return s.repo.Delete(id)
}

// ResolveShortcut looks up a shortcut's destination and increments its
// click count, for use by the redirect handler.
func (s *LinkService) ResolveShortcut(shortcut string) (*models.Link, error) {
	return s.repo.IncrementClickCount(shortcut)
}
