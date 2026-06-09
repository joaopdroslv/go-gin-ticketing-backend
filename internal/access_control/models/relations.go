package acmodels

import "time"

type UserRole struct {
	ID        int64
	UserID    int64
	RoleID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RolePermission struct {
	ID           int64
	RoleID       int64
	PermissionID int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
