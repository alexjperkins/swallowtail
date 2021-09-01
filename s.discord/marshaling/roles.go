package marshaling

import (
	"swallowtail/s.discord/domain"
	discordproto "swallowtail/s.discord/proto"
)

// RolesToProto ...
func RolesToProto(roles []*domain.Role) []*discordproto.Role {
	protoRoles := []*discordproto.Role{}
	for _, role := range roles {
		protoRoles = append(protoRoles, &discordproto.Role{
			RoleId:   role.ID,
			RoleName: role.Name,
		})
	}

	return protoRoles
}
