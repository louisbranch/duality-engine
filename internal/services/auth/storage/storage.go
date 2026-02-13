package storage

import (
	"context"
	"time"

	"github.com/louisbranch/fracturing.space/internal/platform/errors"
	"github.com/louisbranch/fracturing.space/internal/services/auth/user"
)

// ErrNotFound indicates a requested record is missing.
var ErrNotFound = errors.New(errors.CodeNotFound, "record not found")

// UserStore persists auth user records.
type UserStore interface {
	PutUser(ctx context.Context, u user.User) error
	GetUser(ctx context.Context, userID string) (user.User, error)
	ListUsers(ctx context.Context, pageSize int, pageToken string) (UserPage, error)
}

// UserPage describes a page of user records.
type UserPage struct {
	Users         []user.User
	NextPageToken string
}

// AuthStatistics contains aggregate counts across auth data.
type AuthStatistics struct {
	UserCount int64
}

// StatisticsStore provides aggregate auth statistics.
type StatisticsStore interface {
	// GetAuthStatistics returns aggregate counts.
	// When since is nil, counts are for all time.
	GetAuthStatistics(ctx context.Context, since *time.Time) (AuthStatistics, error)
}
