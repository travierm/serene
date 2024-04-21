package main

import (
	"github.com/labstack/echo/v4"
)

// CustomContext is a custom context type that embeds context.Context
type SereneContext struct {
	echo.Context
	Databases []Database
}

// WithUserID returns a new CustomContext with the given user ID
// func WithUserID(ctx context.Context, userID string) *CustomContext {
// 	return &CustomContext{
// 		Context: ctx,
// 		UserID:  userID,
// 	}
// }
