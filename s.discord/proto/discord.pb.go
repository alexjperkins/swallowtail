// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.17.3
// source: s.discord/proto/discord.proto

package discordproto

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type SendMsgToChannelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId      string `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	SenderId       string `protobuf:"bytes,2,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	Content        string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	IdempotencyKey string `protobuf:"bytes,4,opt,name=idempotency_key,json=idempotencyKey,proto3" json:"idempotency_key,omitempty"`
	Force          bool   `protobuf:"varint,5,opt,name=force,proto3" json:"force,omitempty"`
}

func (x *SendMsgToChannelRequest) Reset() {
	*x = SendMsgToChannelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_discord_proto_discord_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMsgToChannelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMsgToChannelRequest) ProtoMessage() {}

func (x *SendMsgToChannelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_s_discord_proto_discord_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMsgToChannelRequest.ProtoReflect.Descriptor instead.
func (*SendMsgToChannelRequest) Descriptor() ([]byte, []int) {
	return file_s_discord_proto_discord_proto_rawDescGZIP(), []int{0}
}

func (x *SendMsgToChannelRequest) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

func (x *SendMsgToChannelRequest) GetSenderId() string {
	if x != nil {
		return x.SenderId
	}
	return ""
}

func (x *SendMsgToChannelRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *SendMsgToChannelRequest) GetIdempotencyKey() string {
	if x != nil {
		return x.IdempotencyKey
	}
	return ""
}

func (x *SendMsgToChannelRequest) GetForce() bool {
	if x != nil {
		return x.Force
	}
	return false
}

type SendMsgToChannelResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendMsgToChannelResponse) Reset() {
	*x = SendMsgToChannelResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_discord_proto_discord_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMsgToChannelResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMsgToChannelResponse) ProtoMessage() {}

func (x *SendMsgToChannelResponse) ProtoReflect() protoreflect.Message {
	mi := &file_s_discord_proto_discord_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMsgToChannelResponse.ProtoReflect.Descriptor instead.
func (*SendMsgToChannelResponse) Descriptor() ([]byte, []int) {
	return file_s_discord_proto_discord_proto_rawDescGZIP(), []int{1}
}

type SendMsgToPrivateChannelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId         string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	SenderId       string `protobuf:"bytes,2,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	Content        string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	IdempotencyKey string `protobuf:"bytes,4,opt,name=idempotency_key,json=idempotencyKey,proto3" json:"idempotency_key,omitempty"`
	Force          bool   `protobuf:"varint,5,opt,name=force,proto3" json:"force,omitempty"`
}

func (x *SendMsgToPrivateChannelRequest) Reset() {
	*x = SendMsgToPrivateChannelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_discord_proto_discord_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMsgToPrivateChannelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMsgToPrivateChannelRequest) ProtoMessage() {}

func (x *SendMsgToPrivateChannelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_s_discord_proto_discord_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMsgToPrivateChannelRequest.ProtoReflect.Descriptor instead.
func (*SendMsgToPrivateChannelRequest) Descriptor() ([]byte, []int) {
	return file_s_discord_proto_discord_proto_rawDescGZIP(), []int{2}
}

func (x *SendMsgToPrivateChannelRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *SendMsgToPrivateChannelRequest) GetSenderId() string {
	if x != nil {
		return x.SenderId
	}
	return ""
}

func (x *SendMsgToPrivateChannelRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *SendMsgToPrivateChannelRequest) GetIdempotencyKey() string {
	if x != nil {
		return x.IdempotencyKey
	}
	return ""
}

func (x *SendMsgToPrivateChannelRequest) GetForce() bool {
	if x != nil {
		return x.Force
	}
	return false
}

type SendMsgToPrivateChannelResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendMsgToPrivateChannelResponse) Reset() {
	*x = SendMsgToPrivateChannelResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_discord_proto_discord_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMsgToPrivateChannelResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMsgToPrivateChannelResponse) ProtoMessage() {}

func (x *SendMsgToPrivateChannelResponse) ProtoReflect() protoreflect.Message {
	mi := &file_s_discord_proto_discord_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMsgToPrivateChannelResponse.ProtoReflect.Descriptor instead.
func (*SendMsgToPrivateChannelResponse) Descriptor() ([]byte, []int) {
	return file_s_discord_proto_discord_proto_rawDescGZIP(), []int{3}
}

type Role struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoleId   string `protobuf:"bytes,1,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
	RoleName string `protobuf:"bytes,2,opt,name=role_name,json=roleName,proto3" json:"role_name,omitempty"`
}

func (x *Role) Reset() {
	*x = Role{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_discord_proto_discord_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Role) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Role) ProtoMessage() {}

func (x *Role) ProtoReflect() protoreflect.Message {
	mi := &file_s_discord_proto_discord_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Role.ProtoReflect.Descriptor instead.
func (*Role) Descriptor() ([]byte, []int) {
	return file_s_discord_proto_discord_proto_rawDescGZIP(), []int{4}
}

func (x *Role) GetRoleId() string {
	if x != nil {
		return x.RoleId
	}
	return ""
}

func (x *Role) GetRoleName() string {
	if x != nil {
		return x.RoleName
	}
	return ""
}

type ReadUserRolesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *ReadUserRolesRequest) Reset() {
	*x = ReadUserRolesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_discord_proto_discord_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadUserRolesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadUserRolesRequest) ProtoMessage() {}

func (x *ReadUserRolesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_s_discord_proto_discord_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadUserRolesRequest.ProtoReflect.Descriptor instead.
func (*ReadUserRolesRequest) Descriptor() ([]byte, []int) {
	return file_s_discord_proto_discord_proto_rawDescGZIP(), []int{5}
}

func (x *ReadUserRolesRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type ReadUserRolesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Roles []*Role `protobuf:"bytes,1,rep,name=roles,proto3" json:"roles,omitempty"`
}

func (x *ReadUserRolesResponse) Reset() {
	*x = ReadUserRolesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_discord_proto_discord_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadUserRolesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadUserRolesResponse) ProtoMessage() {}

func (x *ReadUserRolesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_s_discord_proto_discord_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadUserRolesResponse.ProtoReflect.Descriptor instead.
func (*ReadUserRolesResponse) Descriptor() ([]byte, []int) {
	return file_s_discord_proto_discord_proto_rawDescGZIP(), []int{6}
}

func (x *ReadUserRolesResponse) GetRoles() []*Role {
	if x != nil {
		return x.Roles
	}
	return nil
}

type UpdateUserRolesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId            string  `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Roles             []*Role `protobuf:"bytes,2,rep,name=roles,proto3" json:"roles,omitempty"`
	MergeWithExisting bool    `protobuf:"varint,3,opt,name=merge_with_existing,json=mergeWithExisting,proto3" json:"merge_with_existing,omitempty"`
	ActorId           string  `protobuf:"bytes,4,opt,name=actor_id,json=actorId,proto3" json:"actor_id,omitempty"`
}

func (x *UpdateUserRolesRequest) Reset() {
	*x = UpdateUserRolesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_discord_proto_discord_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUserRolesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserRolesRequest) ProtoMessage() {}

func (x *UpdateUserRolesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_s_discord_proto_discord_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserRolesRequest.ProtoReflect.Descriptor instead.
func (*UpdateUserRolesRequest) Descriptor() ([]byte, []int) {
	return file_s_discord_proto_discord_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateUserRolesRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UpdateUserRolesRequest) GetRoles() []*Role {
	if x != nil {
		return x.Roles
	}
	return nil
}

func (x *UpdateUserRolesRequest) GetMergeWithExisting() bool {
	if x != nil {
		return x.MergeWithExisting
	}
	return false
}

func (x *UpdateUserRolesRequest) GetActorId() string {
	if x != nil {
		return x.ActorId
	}
	return ""
}

type UpdateUserRolesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Roles []*Role `protobuf:"bytes,1,rep,name=roles,proto3" json:"roles,omitempty"`
}

func (x *UpdateUserRolesResponse) Reset() {
	*x = UpdateUserRolesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_discord_proto_discord_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUserRolesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserRolesResponse) ProtoMessage() {}

func (x *UpdateUserRolesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_s_discord_proto_discord_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserRolesResponse.ProtoReflect.Descriptor instead.
func (*UpdateUserRolesResponse) Descriptor() ([]byte, []int) {
	return file_s_discord_proto_discord_proto_rawDescGZIP(), []int{8}
}

func (x *UpdateUserRolesResponse) GetRoles() []*Role {
	if x != nil {
		return x.Roles
	}
	return nil
}

type RemoveUserRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ActorId string `protobuf:"bytes,2,opt,name=actor_id,json=actorId,proto3" json:"actor_id,omitempty"`
	Role    *Role  `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *RemoveUserRoleRequest) Reset() {
	*x = RemoveUserRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_discord_proto_discord_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveUserRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveUserRoleRequest) ProtoMessage() {}

func (x *RemoveUserRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_s_discord_proto_discord_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveUserRoleRequest.ProtoReflect.Descriptor instead.
func (*RemoveUserRoleRequest) Descriptor() ([]byte, []int) {
	return file_s_discord_proto_discord_proto_rawDescGZIP(), []int{9}
}

func (x *RemoveUserRoleRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *RemoveUserRoleRequest) GetActorId() string {
	if x != nil {
		return x.ActorId
	}
	return ""
}

func (x *RemoveUserRoleRequest) GetRole() *Role {
	if x != nil {
		return x.Role
	}
	return nil
}

type RemoveUserRoleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RemoveUserRoleResponse) Reset() {
	*x = RemoveUserRoleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_discord_proto_discord_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveUserRoleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveUserRoleResponse) ProtoMessage() {}

func (x *RemoveUserRoleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_s_discord_proto_discord_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveUserRoleResponse.ProtoReflect.Descriptor instead.
func (*RemoveUserRoleResponse) Descriptor() ([]byte, []int) {
	return file_s_discord_proto_discord_proto_rawDescGZIP(), []int{10}
}

var File_s_discord_proto_discord_proto protoreflect.FileDescriptor

var file_s_discord_proto_discord_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x73, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xae, 0x01, 0x0a, 0x17, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x54, 0x6f, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x12, 0x27, 0x0a, 0x0f, 0x69, 0x64, 0x65, 0x6d, 0x70, 0x6f, 0x74, 0x65, 0x6e, 0x63, 0x79,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x69, 0x64, 0x65, 0x6d,
	0x70, 0x6f, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6f,
	0x72, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x66, 0x6f, 0x72, 0x63, 0x65,
	0x22, 0x1a, 0x0a, 0x18, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x54, 0x6f, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xaf, 0x01, 0x0a,
	0x1e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x54, 0x6f, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74,
	0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12,
	0x27, 0x0a, 0x0f, 0x69, 0x64, 0x65, 0x6d, 0x70, 0x6f, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x69, 0x64, 0x65, 0x6d, 0x70, 0x6f,
	0x74, 0x65, 0x6e, 0x63, 0x79, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6f, 0x72, 0x63,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x22, 0x21,
	0x0a, 0x1f, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x54, 0x6f, 0x50, 0x72, 0x69, 0x76, 0x61,
	0x74, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x3c, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6c,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x6c, 0x65,
	0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x6f, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x2f, 0x0a, 0x14, 0x52, 0x65, 0x61, 0x64, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x22, 0x34, 0x0a, 0x15, 0x52, 0x65, 0x61, 0x64, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x05, 0x72, 0x6f, 0x6c,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52,
	0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x22, 0x99, 0x01, 0x0a, 0x16, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x05, 0x72, 0x6f,
	0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x52, 0x6f, 0x6c, 0x65,
	0x52, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x2e, 0x0a, 0x13, 0x6d, 0x65, 0x72, 0x67, 0x65,
	0x5f, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x65, 0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x11, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x57, 0x69, 0x74, 0x68, 0x45,
	0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x63, 0x74, 0x6f, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x63, 0x74, 0x6f, 0x72,
	0x49, 0x64, 0x22, 0x36, 0x0a, 0x17, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a,
	0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x52,
	0x6f, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x22, 0x66, 0x0a, 0x15, 0x52, 0x65,
	0x6d, 0x6f, 0x76, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08,
	0x61, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x61, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f,
	0x6c, 0x65, 0x22, 0x18, 0x0a, 0x16, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xf9, 0x02, 0x0a,
	0x07, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x47, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64,
	0x4d, 0x73, 0x67, 0x54, 0x6f, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x18, 0x2e, 0x53,
	0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x54, 0x6f, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67,
	0x54, 0x6f, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x5c, 0x0a, 0x17, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x54, 0x6f, 0x50, 0x72,
	0x69, 0x76, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x1f, 0x2e, 0x53,
	0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x54, 0x6f, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x43,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e,
	0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x54, 0x6f, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65,
	0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3e, 0x0a, 0x0d, 0x52, 0x65, 0x61, 0x64, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x73,
	0x12, 0x15, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x44, 0x0a, 0x0f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c,
	0x65, 0x73, 0x12, 0x17, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x16, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x17, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x11, 0x5a, 0x0f, 0x2e, 0x2f, 0x3b, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_s_discord_proto_discord_proto_rawDescOnce sync.Once
	file_s_discord_proto_discord_proto_rawDescData = file_s_discord_proto_discord_proto_rawDesc
)

func file_s_discord_proto_discord_proto_rawDescGZIP() []byte {
	file_s_discord_proto_discord_proto_rawDescOnce.Do(func() {
		file_s_discord_proto_discord_proto_rawDescData = protoimpl.X.CompressGZIP(file_s_discord_proto_discord_proto_rawDescData)
	})
	return file_s_discord_proto_discord_proto_rawDescData
}

var file_s_discord_proto_discord_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_s_discord_proto_discord_proto_goTypes = []interface{}{
	(*SendMsgToChannelRequest)(nil),         // 0: SendMsgToChannelRequest
	(*SendMsgToChannelResponse)(nil),        // 1: SendMsgToChannelResponse
	(*SendMsgToPrivateChannelRequest)(nil),  // 2: SendMsgToPrivateChannelRequest
	(*SendMsgToPrivateChannelResponse)(nil), // 3: SendMsgToPrivateChannelResponse
	(*Role)(nil),                            // 4: Role
	(*ReadUserRolesRequest)(nil),            // 5: ReadUserRolesRequest
	(*ReadUserRolesResponse)(nil),           // 6: ReadUserRolesResponse
	(*UpdateUserRolesRequest)(nil),          // 7: UpdateUserRolesRequest
	(*UpdateUserRolesResponse)(nil),         // 8: UpdateUserRolesResponse
	(*RemoveUserRoleRequest)(nil),           // 9: RemoveUserRoleRequest
	(*RemoveUserRoleResponse)(nil),          // 10: RemoveUserRoleResponse
}
var file_s_discord_proto_discord_proto_depIdxs = []int32{
	4,  // 0: ReadUserRolesResponse.roles:type_name -> Role
	4,  // 1: UpdateUserRolesRequest.roles:type_name -> Role
	4,  // 2: UpdateUserRolesResponse.roles:type_name -> Role
	4,  // 3: RemoveUserRoleRequest.role:type_name -> Role
	0,  // 4: discord.SendMsgToChannel:input_type -> SendMsgToChannelRequest
	2,  // 5: discord.SendMsgToPrivateChannel:input_type -> SendMsgToPrivateChannelRequest
	5,  // 6: discord.ReadUserRoles:input_type -> ReadUserRolesRequest
	7,  // 7: discord.UpdateUserRoles:input_type -> UpdateUserRolesRequest
	9,  // 8: discord.RemoveUserRole:input_type -> RemoveUserRoleRequest
	1,  // 9: discord.SendMsgToChannel:output_type -> SendMsgToChannelResponse
	3,  // 10: discord.SendMsgToPrivateChannel:output_type -> SendMsgToPrivateChannelResponse
	6,  // 11: discord.ReadUserRoles:output_type -> ReadUserRolesResponse
	8,  // 12: discord.UpdateUserRoles:output_type -> UpdateUserRolesResponse
	10, // 13: discord.RemoveUserRole:output_type -> RemoveUserRoleResponse
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_s_discord_proto_discord_proto_init() }
func file_s_discord_proto_discord_proto_init() {
	if File_s_discord_proto_discord_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_s_discord_proto_discord_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMsgToChannelRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_s_discord_proto_discord_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMsgToChannelResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_s_discord_proto_discord_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMsgToPrivateChannelRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_s_discord_proto_discord_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMsgToPrivateChannelResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_s_discord_proto_discord_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Role); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_s_discord_proto_discord_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadUserRolesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_s_discord_proto_discord_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadUserRolesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_s_discord_proto_discord_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateUserRolesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_s_discord_proto_discord_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateUserRolesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_s_discord_proto_discord_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveUserRoleRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_s_discord_proto_discord_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveUserRoleResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_s_discord_proto_discord_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_s_discord_proto_discord_proto_goTypes,
		DependencyIndexes: file_s_discord_proto_discord_proto_depIdxs,
		MessageInfos:      file_s_discord_proto_discord_proto_msgTypes,
	}.Build()
	File_s_discord_proto_discord_proto = out.File
	file_s_discord_proto_discord_proto_rawDesc = nil
	file_s_discord_proto_discord_proto_goTypes = nil
	file_s_discord_proto_discord_proto_depIdxs = nil
}
