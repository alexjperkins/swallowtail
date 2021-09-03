package handler

import (
	"context"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.discord/client"
	"swallowtail/s.discord/domain"
	"swallowtail/s.discord/marshaling"
	discordproto "swallowtail/s.discord/proto"
)

// RemoveUserRole ...
func (s *DiscordService) RemoveUserRole(
	ctx context.Context, in *discordproto.RemoveUserRoleRequest,
) (*discordproto.RemoveUserRoleResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_param.user_id", nil)
	case in.GetRole().RoleId == "":
		return nil, gerrors.BadParam("missing_param.role.role_id", nil)
	}

	errParams := map[string]string{
		"user_id":           in.UserId,
		"actor_id":          in.ActorId,
		"to_remove_role_id": in.Role.RoleId,
	}

	// Confirm we have a valid actor.
	actorValid, err := isValidActor(ctx, in.ActorId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_remove_user_roles", errParams)
	}
	if !actorValid {
		return nil, gerrors.Unimplemented("failed_to_remove_user_role.unauthorized_actor", errParams)
	}

	rsp, err := (&discordproto.ReadUserRolesRequest{
		UserId: in.UserId,
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_remove_user_role", errParams)
	}

	currentRoles := marshaling.RolesProtoToDomain(rsp.Roles)

	// Filter out the role we want to remove from our current roles.
	updatedRoles := []*domain.Role{}
	for _, cr := range currentRoles {
		if cr.ID == in.Role.RoleId {
			continue
		}

		updatedRoles = append(updatedRoles, cr)
	}

	if err := client.SetRoles(ctx, in.UserId, updatedRoles); err != nil {
		return nil, gerrors.Augment(err, "failed_to_remove_user_role", errParams)
	}

	return &discordproto.RemoveUserRoleResponse{}, nil
}
