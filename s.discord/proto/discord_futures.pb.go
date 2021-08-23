package discordproto

import (
	context "context"
	"swallowtail/libraries/gerrors"
	"time"

	"github.com/monzo/slog"
	grpc "google.golang.org/grpc"
)

// --- SendMsgToChannel --- //
type SendMsgToChannelFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *SendMsgToChannelResponse
	ctx     context.Context
}

func (a *SendMsgToChannelFuture) Response() (*SendMsgToChannelResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "send_msg_to_channel", err)
		}
	}()

	select {
	case r := <-a.resultc:
		return r, nil
	case <-a.ctx.Done():
		return nil, a.ctx.Err()
	case err := <-a.errc:
		return nil, err
	}
}

func (r *SendMsgToChannelRequest) Send(ctx context.Context) *SendMsgToChannelFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *SendMsgToChannelRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *SendMsgToChannelFuture {
	errc := make(chan error, 1)
	resultc := make(chan *SendMsgToChannelResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-discord:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_discord_connection_failed", nil)
		return &SendMsgToChannelFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewDiscordClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.SendMsgToChannel(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_send_msg_to_channel", nil)
			return
		}
		resultc <- rsp
	}()

	return &SendMsgToChannelFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- SendMsgToPrivateChannel --- //

type SendMsgToPrivateChannelFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *SendMsgToPrivateChannelResponse
	ctx     context.Context
}

func (a *SendMsgToPrivateChannelFuture) Response() (*SendMsgToPrivateChannelResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "send_msg_to_private_channel", err)
		}
	}()

	select {
	case r := <-a.resultc:
		return r, nil
	case <-a.ctx.Done():
		return nil, a.ctx.Err()
	case err := <-a.errc:
		return nil, err
	}
}

func (r *SendMsgToPrivateChannelRequest) Send(ctx context.Context) *SendMsgToPrivateChannelFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *SendMsgToPrivateChannelRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *SendMsgToPrivateChannelFuture {
	errc := make(chan error, 1)
	resultc := make(chan *SendMsgToPrivateChannelResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-discord:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_discord_connection_failed", nil)
		return &SendMsgToPrivateChannelFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewDiscordClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.SendMsgToPrivateChannel(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_send_msg_to_private_channel", nil)
			return
		}
		resultc <- rsp
	}()

	return &SendMsgToPrivateChannelFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
