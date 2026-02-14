// Package user provides auth user management.
package user

import (
	"fmt"
	"strings"
	"time"

	commonv1 "github.com/louisbranch/fracturing.space/api/gen/go/common/v1"
	apperrors "github.com/louisbranch/fracturing.space/internal/platform/errors"
	platformi18n "github.com/louisbranch/fracturing.space/internal/platform/i18n"
	"github.com/louisbranch/fracturing.space/internal/platform/id"
)

var (
	// ErrEmptyDisplayName indicates a missing user display name.
	ErrEmptyDisplayName = apperrors.New(apperrors.CodeUserEmptyDisplayName, "display name is required")
)

// User represents an authenticated identity record.
type User struct {
	ID          string
	DisplayName string
	Locale      commonv1.Locale
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateUserInput describes the metadata needed to create a user.
type CreateUserInput struct {
	DisplayName string
	Locale      commonv1.Locale
}

// CreateUser creates a new user with a generated ID and timestamps.
func CreateUser(input CreateUserInput, now func() time.Time, idGenerator func() (string, error)) (User, error) {
	if now == nil {
		now = time.Now
	}
	if idGenerator == nil {
		idGenerator = id.NewID
	}

	normalized, err := NormalizeCreateUserInput(input)
	if err != nil {
		return User{}, err
	}

	userID, err := idGenerator()
	if err != nil {
		return User{}, fmt.Errorf("generate user id: %w", err)
	}

	createdAt := now().UTC()
	return User{
		ID:          userID,
		DisplayName: normalized.DisplayName,
		Locale:      normalized.Locale,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}, nil
}

// NormalizeCreateUserInput trims and validates user input metadata.
func NormalizeCreateUserInput(input CreateUserInput) (CreateUserInput, error) {
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	if input.DisplayName == "" {
		return CreateUserInput{}, ErrEmptyDisplayName
	}
	input.Locale = platformi18n.NormalizeLocale(input.Locale)
	return input, nil
}
