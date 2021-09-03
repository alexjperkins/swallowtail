package handler

import (
	"context"
	"fmt"
	"strings"
	"swallowtail/libraries/gerrors"
	accountproto "swallowtail/s.account/proto"
	"swallowtail/s.discord/client"
	"swallowtail/s.discord/domain"
	"swallowtail/s.discord/marshaling"
	discordproto "swallowtail/s.discord/proto"
)

// UpdateUserRoles ...
func (s *DiscordService) UpdateUserRoles(
	ctx context.Context, in *discordproto.UpdateUserRolesRequest,
) (*discordproto.UpdateUserRolesResponse, error) {
	switch {
	case in.UserId == "":
		return nil, gerrors.BadParam("missing_params.user_id", nil)
	case len(in.Roles) == 0 && in.MergeWithExisting == true:
		return &discordproto.UpdateUserRolesResponse{}, nil
	}

	errParams := map[string]string{
		"user_id":  in.UserId,
		"actor_id": in.ActorId,
		"roles":    stringRoles(in.Roles),
	}

	// Confirm we have a valid actor.
	actorValid, err := isValidActor(ctx, in.ActorId)
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_update_user_roles", errParams)
	}
	if !actorValid {
		return nil, gerrors.Unimplemented("failed_to_update_user_roles.unauthorized_actor", errParams)
	}

	newRoles := marshaling.RolesProtoToDomain(in.Roles)

	// If we don't want to merge with existing roles we can just set them and return early.
	if !in.MergeWithExisting {
		if err := client.SetRoles(ctx, in.UserId, newRoles); err != nil {
			return nil, gerrors.Augment(err, "failed_to_update_user_roles", errParams)
		}

		return &discordproto.UpdateUserRolesResponse{}, nil
	}

	// Merge roles
	rsp, err := (&discordproto.ReadUserRolesRequest{
		UserId: in.UserId,
	}).Send(ctx).Response()
	if err != nil {
		return nil, gerrors.Augment(err, "failed_to_update_user_roles", errParams)
	}

	existingRoles := marshaling.RolesProtoToDomain(rsp.Roles)
	mergedRoles := mergeRoles(existingRoles, newRoles)
	mergedRolesProto := marshaling.RolesToProto(mergedRoles)

	if err := client.SetRoles(ctx, in.UserId, newRoles); err != nil {
		errParams["roles"] = stringRoles(mergedRolesProto)
		errParams["old_roles"] = stringRoles(rsp.Roles)
		return nil, gerrors.Augment(err, "failed_to_update_user_roles", errParams)
	}

	return &discordproto.UpdateUserRolesResponse{
		Roles: mergedRolesProto,
	}, nil

}

func mergeRoles(existingRoles, newRoles []*domain.Role) []*domain.Role {
	if len(newRoles) == 0 {
		return existingRoles
	}

	head, tail := newRoles[0], newRoles[1:]
	for _, er := range existingRoles {
		if er.ID == head.ID {
			return mergeRoles(existingRoles, tail)
		}
	}

	return mergeRoles(append(existingRoles, head), tail)
}

func stringRoles(roles []*discordproto.Role) string {
	srs := []string{}
	for _, r := range roles {
		srs = append(srs, fmt.Sprintf("%s:%s", r.RoleName, r.RoleId))
	}

	return strings.Join(srs, ", ")
}

func isValidActor(ctx context.Context, actorID string) (bool, error) {
	switch actorID {
	case "":
		return false, nil
	case discordproto.DiscordRolesUpdateActorPaymentsSystem:
		return true, nil
	}

	account, err := (&accountproto.ReadAccountRequest{
		UserId: actorID,
	}).Send(ctx).Response()
	if err != nil {
		return false, gerrors.Augment(err, "failed_to_list_account_deposits.failed_to_read_account_of_actor", nil)
	}

	return account.GetAccount().IsAdmin, nil
}
