// Code generated by protoc-gen-go. DO NOT EDIT.
// source: entity.proto

package com_salty_protos

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GrpcReq_Language int32

const (
	GrpcReq_CHINESE GrpcReq_Language = 0
	GrpcReq_ENGLISH GrpcReq_Language = 1
)

var GrpcReq_Language_name = map[int32]string{
	0: "CHINESE",
	1: "ENGLISH",
}

var GrpcReq_Language_value = map[string]int32{
	"CHINESE": 0,
	"ENGLISH": 1,
}

func (x GrpcReq_Language) String() string {
	return proto.EnumName(GrpcReq_Language_name, int32(x))
}

func (GrpcReq_Language) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{0, 0}
}

type GrpcReq_OS int32

const (
	GrpcReq_UNKNOWN GrpcReq_OS = 0
	GrpcReq_ANDROID GrpcReq_OS = 1
	GrpcReq_IOS     GrpcReq_OS = 2
	GrpcReq_WINDOWS GrpcReq_OS = 3
)

var GrpcReq_OS_name = map[int32]string{
	0: "UNKNOWN",
	1: "ANDROID",
	2: "IOS",
	3: "WINDOWS",
}

var GrpcReq_OS_value = map[string]int32{
	"UNKNOWN": 0,
	"ANDROID": 1,
	"IOS":     2,
	"WINDOWS": 3,
}

func (x GrpcReq_OS) String() string {
	return proto.EnumName(GrpcReq_OS_name, int32(x))
}

func (GrpcReq_OS) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{0, 1}
}

type UserProfile_Sex int32

const (
	UserProfile_NOT_SET UserProfile_Sex = 0
	UserProfile_MALE    UserProfile_Sex = 1
	UserProfile_FEMALE  UserProfile_Sex = 2
)

var UserProfile_Sex_name = map[int32]string{
	0: "NOT_SET",
	1: "MALE",
	2: "FEMALE",
}

var UserProfile_Sex_value = map[string]int32{
	"NOT_SET": 0,
	"MALE":    1,
	"FEMALE":  2,
}

func (x UserProfile_Sex) String() string {
	return proto.EnumName(UserProfile_Sex_name, int32(x))
}

func (UserProfile_Sex) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{2, 0}
}

type RegisterReq_RegisterType int32

const (
	RegisterReq_TELEPHONE RegisterReq_RegisterType = 0
	RegisterReq_EMAIL     RegisterReq_RegisterType = 1
)

var RegisterReq_RegisterType_name = map[int32]string{
	0: "TELEPHONE",
	1: "EMAIL",
}

var RegisterReq_RegisterType_value = map[string]int32{
	"TELEPHONE": 0,
	"EMAIL":     1,
}

func (x RegisterReq_RegisterType) String() string {
	return proto.EnumName(RegisterReq_RegisterType_name, int32(x))
}

func (RegisterReq_RegisterType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{3, 0}
}

type LoginReq_LoginType int32

const (
	LoginReq_TELEPHONE LoginReq_LoginType = 0
	LoginReq_EMAIL     LoginReq_LoginType = 1
)

var LoginReq_LoginType_name = map[int32]string{
	0: "TELEPHONE",
	1: "EMAIL",
}

var LoginReq_LoginType_value = map[string]int32{
	"TELEPHONE": 0,
	"EMAIL":     1,
}

func (x LoginReq_LoginType) String() string {
	return proto.EnumName(LoginReq_LoginType_name, int32(x))
}

func (LoginReq_LoginType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{5, 0}
}

type SMSReq_CodeType int32

const (
	SMSReq_REGISTER SMSReq_CodeType = 0
	SMSReq_LOGIN    SMSReq_CodeType = 1
)

var SMSReq_CodeType_name = map[int32]string{
	0: "REGISTER",
	1: "LOGIN",
}

var SMSReq_CodeType_value = map[string]int32{
	"REGISTER": 0,
	"LOGIN":    1,
}

func (x SMSReq_CodeType) String() string {
	return proto.EnumName(SMSReq_CodeType_name, int32(x))
}

func (SMSReq_CodeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{7, 0}
}

// Basic Request
type GrpcReq struct {
	DeviceID             string           `protobuf:"bytes,1,opt,name=deviceID,proto3" json:"deviceID,omitempty"`
	Version              string           `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Language             GrpcReq_Language `protobuf:"varint,3,opt,name=language,proto3,enum=com.salty.protos.GrpcReq_Language" json:"language,omitempty"`
	Os                   GrpcReq_OS       `protobuf:"varint,4,opt,name=os,proto3,enum=com.salty.protos.GrpcReq_OS" json:"os,omitempty"`
	Token                string           `protobuf:"bytes,5,opt,name=token,proto3" json:"token,omitempty"`
	Data                 *any.Any         `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GrpcReq) Reset()         { *m = GrpcReq{} }
func (m *GrpcReq) String() string { return proto.CompactTextString(m) }
func (*GrpcReq) ProtoMessage()    {}
func (*GrpcReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{0}
}

func (m *GrpcReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GrpcReq.Unmarshal(m, b)
}
func (m *GrpcReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GrpcReq.Marshal(b, m, deterministic)
}
func (m *GrpcReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GrpcReq.Merge(m, src)
}
func (m *GrpcReq) XXX_Size() int {
	return xxx_messageInfo_GrpcReq.Size(m)
}
func (m *GrpcReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GrpcReq.DiscardUnknown(m)
}

var xxx_messageInfo_GrpcReq proto.InternalMessageInfo

func (m *GrpcReq) GetDeviceID() string {
	if m != nil {
		return m.DeviceID
	}
	return ""
}

func (m *GrpcReq) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *GrpcReq) GetLanguage() GrpcReq_Language {
	if m != nil {
		return m.Language
	}
	return GrpcReq_CHINESE
}

func (m *GrpcReq) GetOs() GrpcReq_OS {
	if m != nil {
		return m.Os
	}
	return GrpcReq_UNKNOWN
}

func (m *GrpcReq) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *GrpcReq) GetData() *any.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

// Basic Response
type GrpcResp struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data                 *any.Any `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GrpcResp) Reset()         { *m = GrpcResp{} }
func (m *GrpcResp) String() string { return proto.CompactTextString(m) }
func (*GrpcResp) ProtoMessage()    {}
func (*GrpcResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{1}
}

func (m *GrpcResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GrpcResp.Unmarshal(m, b)
}
func (m *GrpcResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GrpcResp.Marshal(b, m, deterministic)
}
func (m *GrpcResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GrpcResp.Merge(m, src)
}
func (m *GrpcResp) XXX_Size() int {
	return xxx_messageInfo_GrpcResp.Size(m)
}
func (m *GrpcResp) XXX_DiscardUnknown() {
	xxx_messageInfo_GrpcResp.DiscardUnknown(m)
}

var xxx_messageInfo_GrpcResp proto.InternalMessageInfo

func (m *GrpcResp) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *GrpcResp) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *GrpcResp) GetData() *any.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

type UserProfile struct {
	UserID               string          `protobuf:"bytes,1,opt,name=userID,proto3" json:"userID,omitempty"`
	Telephone            string          `protobuf:"bytes,2,opt,name=telephone,proto3" json:"telephone,omitempty"`
	Email                string          `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Username             string          `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
	Nickname             string          `protobuf:"bytes,5,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Avatar               string          `protobuf:"bytes,6,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Description          string          `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty"`
	Sex                  UserProfile_Sex `protobuf:"varint,8,opt,name=sex,proto3,enum=com.salty.protos.UserProfile_Sex" json:"sex,omitempty"`
	Birthday             int64           `protobuf:"varint,9,opt,name=birthday,proto3" json:"birthday,omitempty"`
	Location             string          `protobuf:"bytes,10,opt,name=location,proto3" json:"location,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *UserProfile) Reset()         { *m = UserProfile{} }
func (m *UserProfile) String() string { return proto.CompactTextString(m) }
func (*UserProfile) ProtoMessage()    {}
func (*UserProfile) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{2}
}

func (m *UserProfile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserProfile.Unmarshal(m, b)
}
func (m *UserProfile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserProfile.Marshal(b, m, deterministic)
}
func (m *UserProfile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserProfile.Merge(m, src)
}
func (m *UserProfile) XXX_Size() int {
	return xxx_messageInfo_UserProfile.Size(m)
}
func (m *UserProfile) XXX_DiscardUnknown() {
	xxx_messageInfo_UserProfile.DiscardUnknown(m)
}

var xxx_messageInfo_UserProfile proto.InternalMessageInfo

func (m *UserProfile) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *UserProfile) GetTelephone() string {
	if m != nil {
		return m.Telephone
	}
	return ""
}

func (m *UserProfile) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UserProfile) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *UserProfile) GetNickname() string {
	if m != nil {
		return m.Nickname
	}
	return ""
}

func (m *UserProfile) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

func (m *UserProfile) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *UserProfile) GetSex() UserProfile_Sex {
	if m != nil {
		return m.Sex
	}
	return UserProfile_NOT_SET
}

func (m *UserProfile) GetBirthday() int64 {
	if m != nil {
		return m.Birthday
	}
	return 0
}

func (m *UserProfile) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

type RegisterReq struct {
	RegisterType         RegisterReq_RegisterType `protobuf:"varint,1,opt,name=registerType,proto3,enum=com.salty.protos.RegisterReq_RegisterType" json:"registerType,omitempty"`
	Password             string                   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	VerificationCode     string                   `protobuf:"bytes,3,opt,name=verificationCode,proto3" json:"verificationCode,omitempty"`
	UserProfile          *UserProfile             `protobuf:"bytes,4,opt,name=userProfile,proto3" json:"userProfile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *RegisterReq) Reset()         { *m = RegisterReq{} }
func (m *RegisterReq) String() string { return proto.CompactTextString(m) }
func (*RegisterReq) ProtoMessage()    {}
func (*RegisterReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{3}
}

func (m *RegisterReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterReq.Unmarshal(m, b)
}
func (m *RegisterReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterReq.Marshal(b, m, deterministic)
}
func (m *RegisterReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterReq.Merge(m, src)
}
func (m *RegisterReq) XXX_Size() int {
	return xxx_messageInfo_RegisterReq.Size(m)
}
func (m *RegisterReq) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterReq.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterReq proto.InternalMessageInfo

func (m *RegisterReq) GetRegisterType() RegisterReq_RegisterType {
	if m != nil {
		return m.RegisterType
	}
	return RegisterReq_TELEPHONE
}

func (m *RegisterReq) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *RegisterReq) GetVerificationCode() string {
	if m != nil {
		return m.VerificationCode
	}
	return ""
}

func (m *RegisterReq) GetUserProfile() *UserProfile {
	if m != nil {
		return m.UserProfile
	}
	return nil
}

type RegisterResp struct {
	Profile              *UserProfile `protobuf:"bytes,1,opt,name=profile,proto3" json:"profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *RegisterResp) Reset()         { *m = RegisterResp{} }
func (m *RegisterResp) String() string { return proto.CompactTextString(m) }
func (*RegisterResp) ProtoMessage()    {}
func (*RegisterResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{4}
}

func (m *RegisterResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterResp.Unmarshal(m, b)
}
func (m *RegisterResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterResp.Marshal(b, m, deterministic)
}
func (m *RegisterResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterResp.Merge(m, src)
}
func (m *RegisterResp) XXX_Size() int {
	return xxx_messageInfo_RegisterResp.Size(m)
}
func (m *RegisterResp) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterResp.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterResp proto.InternalMessageInfo

func (m *RegisterResp) GetProfile() *UserProfile {
	if m != nil {
		return m.Profile
	}
	return nil
}

type LoginReq struct {
	LoginType LoginReq_LoginType `protobuf:"varint,1,opt,name=loginType,proto3,enum=com.salty.protos.LoginReq_LoginType" json:"loginType,omitempty"`
	// Types that are valid to be assigned to TargetAccount:
	//	*LoginReq_Telephone
	//	*LoginReq_Email
	TargetAccount        isLoginReq_TargetAccount `protobuf_oneof:"targetAccount"`
	Password             string                   `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *LoginReq) Reset()         { *m = LoginReq{} }
func (m *LoginReq) String() string { return proto.CompactTextString(m) }
func (*LoginReq) ProtoMessage()    {}
func (*LoginReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{5}
}

func (m *LoginReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginReq.Unmarshal(m, b)
}
func (m *LoginReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginReq.Marshal(b, m, deterministic)
}
func (m *LoginReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginReq.Merge(m, src)
}
func (m *LoginReq) XXX_Size() int {
	return xxx_messageInfo_LoginReq.Size(m)
}
func (m *LoginReq) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginReq.DiscardUnknown(m)
}

var xxx_messageInfo_LoginReq proto.InternalMessageInfo

func (m *LoginReq) GetLoginType() LoginReq_LoginType {
	if m != nil {
		return m.LoginType
	}
	return LoginReq_TELEPHONE
}

type isLoginReq_TargetAccount interface {
	isLoginReq_TargetAccount()
}

type LoginReq_Telephone struct {
	Telephone string `protobuf:"bytes,2,opt,name=telephone,proto3,oneof"`
}

type LoginReq_Email struct {
	Email string `protobuf:"bytes,3,opt,name=email,proto3,oneof"`
}

func (*LoginReq_Telephone) isLoginReq_TargetAccount() {}

func (*LoginReq_Email) isLoginReq_TargetAccount() {}

func (m *LoginReq) GetTargetAccount() isLoginReq_TargetAccount {
	if m != nil {
		return m.TargetAccount
	}
	return nil
}

func (m *LoginReq) GetTelephone() string {
	if x, ok := m.GetTargetAccount().(*LoginReq_Telephone); ok {
		return x.Telephone
	}
	return ""
}

func (m *LoginReq) GetEmail() string {
	if x, ok := m.GetTargetAccount().(*LoginReq_Email); ok {
		return x.Email
	}
	return ""
}

func (m *LoginReq) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*LoginReq) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*LoginReq_Telephone)(nil),
		(*LoginReq_Email)(nil),
	}
}

type LoginResp struct {
	Profile              *UserProfile `protobuf:"bytes,1,opt,name=profile,proto3" json:"profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *LoginResp) Reset()         { *m = LoginResp{} }
func (m *LoginResp) String() string { return proto.CompactTextString(m) }
func (*LoginResp) ProtoMessage()    {}
func (*LoginResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{6}
}

func (m *LoginResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginResp.Unmarshal(m, b)
}
func (m *LoginResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginResp.Marshal(b, m, deterministic)
}
func (m *LoginResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginResp.Merge(m, src)
}
func (m *LoginResp) XXX_Size() int {
	return xxx_messageInfo_LoginResp.Size(m)
}
func (m *LoginResp) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginResp.DiscardUnknown(m)
}

var xxx_messageInfo_LoginResp proto.InternalMessageInfo

func (m *LoginResp) GetProfile() *UserProfile {
	if m != nil {
		return m.Profile
	}
	return nil
}

type SMSReq struct {
	CodeType SMSReq_CodeType `protobuf:"varint,1,opt,name=codeType,proto3,enum=com.salty.protos.SMSReq_CodeType" json:"codeType,omitempty"`
	// Types that are valid to be assigned to TargetAccount:
	//	*SMSReq_Telephone
	//	*SMSReq_Email
	TargetAccount        isSMSReq_TargetAccount `protobuf_oneof:"targetAccount"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *SMSReq) Reset()         { *m = SMSReq{} }
func (m *SMSReq) String() string { return proto.CompactTextString(m) }
func (*SMSReq) ProtoMessage()    {}
func (*SMSReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{7}
}

func (m *SMSReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SMSReq.Unmarshal(m, b)
}
func (m *SMSReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SMSReq.Marshal(b, m, deterministic)
}
func (m *SMSReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SMSReq.Merge(m, src)
}
func (m *SMSReq) XXX_Size() int {
	return xxx_messageInfo_SMSReq.Size(m)
}
func (m *SMSReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SMSReq.DiscardUnknown(m)
}

var xxx_messageInfo_SMSReq proto.InternalMessageInfo

func (m *SMSReq) GetCodeType() SMSReq_CodeType {
	if m != nil {
		return m.CodeType
	}
	return SMSReq_REGISTER
}

type isSMSReq_TargetAccount interface {
	isSMSReq_TargetAccount()
}

type SMSReq_Telephone struct {
	Telephone string `protobuf:"bytes,2,opt,name=telephone,proto3,oneof"`
}

type SMSReq_Email struct {
	Email string `protobuf:"bytes,3,opt,name=email,proto3,oneof"`
}

func (*SMSReq_Telephone) isSMSReq_TargetAccount() {}

func (*SMSReq_Email) isSMSReq_TargetAccount() {}

func (m *SMSReq) GetTargetAccount() isSMSReq_TargetAccount {
	if m != nil {
		return m.TargetAccount
	}
	return nil
}

func (m *SMSReq) GetTelephone() string {
	if x, ok := m.GetTargetAccount().(*SMSReq_Telephone); ok {
		return x.Telephone
	}
	return ""
}

func (m *SMSReq) GetEmail() string {
	if x, ok := m.GetTargetAccount().(*SMSReq_Email); ok {
		return x.Email
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*SMSReq) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*SMSReq_Telephone)(nil),
		(*SMSReq_Email)(nil),
	}
}

type SMSResp struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SMSResp) Reset()         { *m = SMSResp{} }
func (m *SMSResp) String() string { return proto.CompactTextString(m) }
func (*SMSResp) ProtoMessage()    {}
func (*SMSResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{8}
}

func (m *SMSResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SMSResp.Unmarshal(m, b)
}
func (m *SMSResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SMSResp.Marshal(b, m, deterministic)
}
func (m *SMSResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SMSResp.Merge(m, src)
}
func (m *SMSResp) XXX_Size() int {
	return xxx_messageInfo_SMSResp.Size(m)
}
func (m *SMSResp) XXX_DiscardUnknown() {
	xxx_messageInfo_SMSResp.DiscardUnknown(m)
}

var xxx_messageInfo_SMSResp proto.InternalMessageInfo

type EmptyResp struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyResp) Reset()         { *m = EmptyResp{} }
func (m *EmptyResp) String() string { return proto.CompactTextString(m) }
func (*EmptyResp) ProtoMessage()    {}
func (*EmptyResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{9}
}

func (m *EmptyResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyResp.Unmarshal(m, b)
}
func (m *EmptyResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyResp.Marshal(b, m, deterministic)
}
func (m *EmptyResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyResp.Merge(m, src)
}
func (m *EmptyResp) XXX_Size() int {
	return xxx_messageInfo_EmptyResp.Size(m)
}
func (m *EmptyResp) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyResp.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyResp proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("com.salty.protos.GrpcReq_Language", GrpcReq_Language_name, GrpcReq_Language_value)
	proto.RegisterEnum("com.salty.protos.GrpcReq_OS", GrpcReq_OS_name, GrpcReq_OS_value)
	proto.RegisterEnum("com.salty.protos.UserProfile_Sex", UserProfile_Sex_name, UserProfile_Sex_value)
	proto.RegisterEnum("com.salty.protos.RegisterReq_RegisterType", RegisterReq_RegisterType_name, RegisterReq_RegisterType_value)
	proto.RegisterEnum("com.salty.protos.LoginReq_LoginType", LoginReq_LoginType_name, LoginReq_LoginType_value)
	proto.RegisterEnum("com.salty.protos.SMSReq_CodeType", SMSReq_CodeType_name, SMSReq_CodeType_value)
	proto.RegisterType((*GrpcReq)(nil), "com.salty.protos.GrpcReq")
	proto.RegisterType((*GrpcResp)(nil), "com.salty.protos.GrpcResp")
	proto.RegisterType((*UserProfile)(nil), "com.salty.protos.UserProfile")
	proto.RegisterType((*RegisterReq)(nil), "com.salty.protos.RegisterReq")
	proto.RegisterType((*RegisterResp)(nil), "com.salty.protos.RegisterResp")
	proto.RegisterType((*LoginReq)(nil), "com.salty.protos.LoginReq")
	proto.RegisterType((*LoginResp)(nil), "com.salty.protos.LoginResp")
	proto.RegisterType((*SMSReq)(nil), "com.salty.protos.SMSReq")
	proto.RegisterType((*SMSResp)(nil), "com.salty.protos.SMSResp")
	proto.RegisterType((*EmptyResp)(nil), "com.salty.protos.EmptyResp")
}

func init() { proto.RegisterFile("entity.proto", fileDescriptor_cf50d946d740d100) }

var fileDescriptor_cf50d946d740d100 = []byte{
	// 775 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xdd, 0x6e, 0xe3, 0x44,
	0x14, 0x8e, 0x9d, 0x3f, 0xfb, 0x38, 0xbb, 0x58, 0xa3, 0xd5, 0xca, 0xac, 0x16, 0x14, 0xcc, 0x22,
	0x45, 0x2b, 0xe4, 0x95, 0x52, 0x24, 0xae, 0x00, 0xa5, 0x8d, 0x49, 0x2c, 0x52, 0xbb, 0x1a, 0xa7,
	0xea, 0x25, 0x72, 0x9c, 0x69, 0x6a, 0xd5, 0xf1, 0x18, 0x8f, 0x13, 0x9a, 0xb7, 0xe0, 0x65, 0x78,
	0x08, 0x5e, 0x81, 0x4b, 0x9e, 0x04, 0xcd, 0xf8, 0x27, 0x2e, 0x69, 0x01, 0xc1, 0xdd, 0x7c, 0xe7,
	0xc7, 0xe7, 0x9b, 0xef, 0x3b, 0x63, 0x18, 0x90, 0x24, 0x8f, 0xf2, 0x83, 0x95, 0x66, 0x34, 0xa7,
	0x48, 0x0f, 0xe9, 0xd6, 0x62, 0x41, 0x5c, 0x05, 0xd8, 0x9b, 0x8f, 0x37, 0x94, 0x6e, 0x62, 0xf2,
	0x41, 0xc0, 0xd5, 0xee, 0xf6, 0x43, 0x90, 0x94, 0x39, 0xf3, 0x37, 0x19, 0xfa, 0xb3, 0x2c, 0x0d,
	0x31, 0xf9, 0x09, 0xbd, 0x01, 0x65, 0x4d, 0xf6, 0x51, 0x48, 0x9c, 0xa9, 0x21, 0x0d, 0xa5, 0x91,
	0x8a, 0x6b, 0x8c, 0x0c, 0xe8, 0xef, 0x49, 0xc6, 0x22, 0x9a, 0x18, 0xb2, 0x48, 0x55, 0x10, 0x7d,
	0x0b, 0x4a, 0x1c, 0x24, 0x9b, 0x5d, 0xb0, 0x21, 0x46, 0x7b, 0x28, 0x8d, 0x5e, 0x8e, 0x4d, 0xeb,
	0xaf, 0x0c, 0xac, 0x72, 0x84, 0xb5, 0x28, 0x2b, 0x71, 0xdd, 0x83, 0xbe, 0x04, 0x99, 0x32, 0xa3,
	0x23, 0x3a, 0xdf, 0x3e, 0xdf, 0xe9, 0xf9, 0x58, 0xa6, 0x0c, 0xbd, 0x82, 0x6e, 0x4e, 0xef, 0x49,
	0x62, 0x74, 0x05, 0x8b, 0x02, 0xa0, 0x11, 0x74, 0xd6, 0x41, 0x1e, 0x18, 0xbd, 0xa1, 0x34, 0xd2,
	0xc6, 0xaf, 0xac, 0xe2, 0xbe, 0x56, 0x75, 0x5f, 0x6b, 0x92, 0x1c, 0xb0, 0xa8, 0x30, 0xdf, 0x81,
	0x52, 0x71, 0x40, 0x1a, 0xf4, 0x2f, 0xe6, 0x8e, 0x6b, 0xfb, 0xb6, 0xde, 0xe2, 0xc0, 0x76, 0x67,
	0x0b, 0xc7, 0x9f, 0xeb, 0x92, 0xf9, 0x15, 0xc8, 0x9e, 0xcf, 0x43, 0xd7, 0xee, 0x0f, 0xae, 0x77,
	0xe3, 0x16, 0xf9, 0x89, 0x3b, 0xc5, 0x9e, 0x33, 0xd5, 0x25, 0xd4, 0x87, 0xb6, 0xe3, 0xf9, 0xba,
	0xcc, 0xa3, 0x37, 0x8e, 0x3b, 0xf5, 0x6e, 0x7c, 0xbd, 0x6d, 0xae, 0x40, 0x29, 0xd8, 0xb2, 0x14,
	0x21, 0xe8, 0x84, 0x74, 0x4d, 0x84, 0x8e, 0x5d, 0x2c, 0xce, 0x5c, 0xc3, 0x2d, 0x61, 0x8c, 0x0b,
	0x55, 0x6a, 0x58, 0xc2, 0x9a, 0x7f, 0xfb, 0x1f, 0xf9, 0xff, 0x21, 0x83, 0x76, 0xcd, 0x48, 0x76,
	0x95, 0xd1, 0xdb, 0x28, 0x26, 0xe8, 0x35, 0xf4, 0x76, 0x8c, 0x64, 0xb5, 0x63, 0x25, 0x42, 0x6f,
	0x41, 0xcd, 0x49, 0x4c, 0xd2, 0x3b, 0x9a, 0x54, 0xd3, 0x8e, 0x01, 0xae, 0x22, 0xd9, 0x06, 0x51,
	0x2c, 0x06, 0xaa, 0xb8, 0x00, 0xdc, 0x7f, 0xde, 0x9d, 0x04, 0x5b, 0x22, 0xfc, 0x50, 0x71, 0x8d,
	0x79, 0x2e, 0x89, 0xc2, 0x7b, 0x91, 0x2b, 0xa4, 0xaf, 0x31, 0xe7, 0x10, 0xec, 0x83, 0x3c, 0xc8,
	0x84, 0xfe, 0x2a, 0x2e, 0x11, 0x1a, 0x82, 0xb6, 0x26, 0x2c, 0xcc, 0xa2, 0x34, 0xe7, 0x7b, 0xd3,
	0x17, 0xc9, 0x66, 0x08, 0x9d, 0x41, 0x9b, 0x91, 0x07, 0x43, 0x11, 0xe6, 0x7f, 0x76, 0x6a, 0x7e,
	0xe3, 0xa6, 0x96, 0x4f, 0x1e, 0x30, 0xaf, 0xe6, 0x54, 0x56, 0x51, 0x96, 0xdf, 0xad, 0x83, 0x83,
	0xa1, 0x0e, 0xa5, 0x51, 0x1b, 0xd7, 0x98, 0xe7, 0x62, 0x1a, 0x06, 0x62, 0x1e, 0x14, 0x34, 0x2b,
	0x6c, 0x8e, 0xa0, 0xed, 0x93, 0x07, 0x6e, 0x99, 0xeb, 0x2d, 0x7f, 0xf4, 0xed, 0xa5, 0xde, 0x42,
	0x0a, 0x74, 0x2e, 0x27, 0x0b, 0x5b, 0x97, 0x10, 0x40, 0xef, 0x7b, 0x5b, 0x9c, 0x65, 0xf3, 0x17,
	0x19, 0x34, 0x4c, 0x36, 0x11, 0xcb, 0x49, 0xc6, 0x1f, 0x86, 0x0b, 0x83, 0xac, 0x84, 0xcb, 0x43,
	0x5a, 0x98, 0xfa, 0x72, 0xfc, 0xfe, 0x94, 0x6f, 0xa3, 0xa9, 0x3e, 0xf3, 0x0e, 0xfc, 0xa8, 0x9f,
	0xb3, 0x4c, 0x03, 0xc6, 0x7e, 0xa6, 0xd9, 0xba, 0xf4, 0xa6, 0xc6, 0xe8, 0x3d, 0xe8, 0x7b, 0x92,
	0x45, 0xb7, 0x51, 0xc1, 0xfa, 0x82, 0x2f, 0x51, 0xe1, 0xd2, 0x49, 0x1c, 0x7d, 0x07, 0xda, 0xee,
	0xa8, 0x90, 0xf0, 0x4c, 0x1b, 0x7f, 0xf2, 0xb7, 0x32, 0xe2, 0x66, 0x87, 0x39, 0x82, 0x41, 0x93,
	0x26, 0x7a, 0x01, 0xea, 0xd2, 0x5e, 0xd8, 0x57, 0x73, 0xcf, 0xe5, 0x6f, 0x42, 0x85, 0xae, 0x7d,
	0x39, 0x71, 0x16, 0xba, 0x64, 0xce, 0x8e, 0x95, 0x62, 0xbf, 0xbf, 0x86, 0x7e, 0x5a, 0x8e, 0x95,
	0xfe, 0xcd, 0xd8, 0xaa, 0xda, 0xfc, 0x5d, 0x02, 0x65, 0x41, 0x37, 0x51, 0xc2, 0x85, 0x3d, 0x07,
	0x35, 0xe6, 0xe7, 0x86, 0xaa, 0xef, 0x4e, 0xbf, 0x53, 0x95, 0x17, 0x07, 0xa1, 0xe7, 0xb1, 0x0d,
	0x7d, 0x7a, 0xb2, 0xe9, 0xf3, 0x56, 0x73, 0xd7, 0x5f, 0x3f, 0xda, 0xf5, 0x79, 0xab, 0xb1, 0xed,
	0xb5, 0x09, 0x9d, 0xc7, 0x26, 0x98, 0x5f, 0x80, 0x5a, 0xcf, 0x7a, 0x5e, 0x94, 0xf3, 0x8f, 0xe0,
	0x45, 0x1e, 0x64, 0x1b, 0x92, 0x4f, 0xc2, 0x90, 0xee, 0x92, 0xdc, 0x9c, 0x96, 0x7d, 0xff, 0x4f,
	0xa2, 0x5f, 0x25, 0xe8, 0xf9, 0x97, 0x3e, 0x17, 0xe8, 0x1b, 0x50, 0xf8, 0xaf, 0xa3, 0xa1, 0xcf,
	0x13, 0xaf, 0xa4, 0xa8, 0xb5, 0x2e, 0xca, 0x42, 0x5c, 0xb7, 0xfc, 0x57, 0x6d, 0xcc, 0xcf, 0x41,
	0xa9, 0xbe, 0x86, 0x06, 0xa0, 0x60, 0x7b, 0xe6, 0xf8, 0x4b, 0x1b, 0x17, 0xb7, 0x5f, 0x78, 0x33,
	0xc7, 0x7d, 0xea, 0xf6, 0x2a, 0xf4, 0x05, 0x15, 0x96, 0x9a, 0x1a, 0xa8, 0xf6, 0x36, 0xcd, 0x0f,
	0x1c, 0x9c, 0xcb, 0x57, 0xd2, 0xaa, 0x27, 0xb8, 0x9e, 0xfd, 0x19, 0x00, 0x00, 0xff, 0xff, 0x10,
	0xcd, 0x2e, 0x05, 0xab, 0x06, 0x00, 0x00,
}
