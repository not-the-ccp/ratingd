package dunning

import (
	"context"
)

// Store persists dunning data.
type Store interface {
	LoadNotice(ctx context.Context, id int64) (*Notice, error)
	ListNotices(ctx context.Context, limit int, offset int) ([]*Notice, error)
	ListNoticesByStatus(ctx context.Context, status NoticeStatus) ([]*Notice, error)
	InsertNotice(ctx context.Context, notice *Notice) error
	UpdateNotice(ctx context.Context, notice *Notice) error
}
