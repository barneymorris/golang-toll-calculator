package client

import (
	"context"

	"github.com/betelgeusexru/golang-toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}