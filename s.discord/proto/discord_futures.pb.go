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
			ctx:  ctx,
			errc: errc,
			closer: func() error {
				if conn == nil {
					return nil
				}
				return conn.Close()
			},
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
			ctx:  ctx,
			errc: errc,
			closer: func() error {
				if conn != nil {
					return conn.Close()
				}
				return nil
			},
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

// --- Read User Roles --- //

type ReadUserRolesFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ReadUserRolesResponse
	ctx     context.Context
}

func (a *ReadUserRolesFuture) Response() (*ReadUserRolesResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "read_user_roles", err)
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

func (r *ReadUserRolesRequest) Send(ctx context.Context) *ReadUserRolesFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ReadUserRolesRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ReadUserRolesFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ReadUserRolesResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-discord:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_discord_connection_failed", nil)
		return &ReadUserRolesFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewDiscordClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ReadUserRoles(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_read_user_roles", nil)
			return
		}
		resultc <- rsp
	}()

	return &ReadUserRolesFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Update User Roles --- //

type UpdateUserRolesFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *UpdateUserRolesResponse
	ctx     context.Context
}

func (a *UpdateUserRolesFuture) Response() (*UpdateUserRolesResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "update_user_roles", err)
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

func (r *UpdateUserRolesRequest) Send(ctx context.Context) *UpdateUserRolesFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *UpdateUserRolesRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *UpdateUserRolesFuture {
	errc := make(chan error, 1)
	resultc := make(chan *UpdateUserRolesResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-discord:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_discord_connection_failed", nil)
		return &UpdateUserRolesFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewDiscordClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.UpdateUserRoles(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_update_user_roles", nil)
			return
		}
		resultc <- rsp
	}()

	return &UpdateUserRolesFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Update User Roles --- //

type RemoveUserRoleFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *RemoveUserRoleResponse
	ctx     context.Context
}

func (a *RemoveUserRoleFuture) Response() (*RemoveUserRoleResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "remove_user_role", err)
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

func (r *RemoveUserRoleRequest) Send(ctx context.Context) *RemoveUserRoleFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *RemoveUserRoleRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *RemoveUserRoleFuture {
	errc := make(chan error, 1)
	resultc := make(chan *RemoveUserRoleResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-discord:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_discord_connection_failed", nil)
		return &RemoveUserRoleFuture{
			ctx:     ctx,
			errc:    errc,
			closer:  conn.Close,
			resultc: resultc,
		}
	}
	c := NewDiscordClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.RemoveUserRole(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_remove_user_role", nil)
			return
		}
		resultc <- rsp
	}()

	return &RemoveUserRoleFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}

// --- Read Message Reactions --- //

type ReadMessageReactionsFuture struct {
	closer  func() error
	errc    chan error
	resultc chan *ReadMessageReactionsResponse
	ctx     context.Context
}

func (a *ReadMessageReactionsFuture) Response() (*ReadMessageReactionsResponse, error) {
	defer func() {
		if err := a.closer(); err != nil {
			slog.Critical(context.Background(), "Failed to close %s grpc connection: %v", "read_message_reactions", err)
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

func (r *ReadMessageReactionsRequest) Send(ctx context.Context) *ReadMessageReactionsFuture {
	return r.SendWithTimeout(ctx, 10*time.Second)
}

func (r *ReadMessageReactionsRequest) SendWithTimeout(ctx context.Context, timeout time.Duration) *ReadMessageReactionsFuture {
	errc := make(chan error, 1)
	resultc := make(chan *ReadMessageReactionsResponse, 1)

	conn, err := grpc.DialContext(ctx, "swallowtail-s-discord:8000", grpc.WithInsecure())
	if err != nil {
		errc <- gerrors.Augment(err, "swallowtail_s_discord_connection_failed", nil)
		return &ReadMessageReactionsFuture{
			ctx:  ctx,
			errc: errc,
			closer: func() error {
				if conn != nil {
					return conn.Close()
				}
				return nil
			},
			resultc: resultc,
		}
	}
	c := NewDiscordClient(conn)

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		rsp, err := c.ReadMessageReactions(ctx, r)
		if err != nil {
			errc <- gerrors.Augment(err, "failed_to_read_message_reactions", nil)
			return
		}
		resultc <- rsp
	}()

	return &ReadMessageReactionsFuture{
		ctx: ctx,
		closer: func() error {
			cancel()
			return conn.Close()
		},
		errc:    errc,
		resultc: resultc,
	}
}
