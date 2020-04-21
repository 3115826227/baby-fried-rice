package model

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/jinzhu/gorm"
	"time"
)

func init() {
	Sync(db.DB)
}

func Sync(engine *gorm.DB) {
	err := engine.AutoMigrate(
		new(AdminPermission),
		new(AdminRole),
		new(AdminRolePermissionRelation),
		new(AccountAdminRoleRelation),
		new(AccountRoot),
		new(AccountAdmin),
		new(AccountUser),
		new(AccountClient),
		new(AccountUserDetail),
		new(AccountUserSchoolDetail),
		new(School),
		new(AccountSchoolOrganize),
		new(AccountSchoolStudent),
		new(SchoolDepartment),
		new(SchoolCommunity),
		new(SchoolUserCertification),
		new(ClientSchoolRelation),
		new(Area),
	).Error
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}

type CommonField struct {
	ID        string    `gorm:"column:id;type:char(36);primary_key;not null"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp" json:"-"`
}

type CommonIntField struct {
	ID        int       `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp" json:"-"`
}

type AdminPermission struct {
	ID       int    `gorm:"column:id;primary_key;not null" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Path     string `gorm:"column:path" json:"path"`
	Method   string `gorm:"column:method" json:"method"`
	Types    int    `gorm:"column:types" json:"types"`
	ParentId int    `gorm:"column:parent_id" json:"parent_id"`
}

func (table *AdminPermission) TableName() string {
	return "admin_permission"
}

type AdminRole struct {
	CommonIntField

	Name     string `gorm:"column:name;" json:"name"`
	SchoolId string `gorm:"column:school_id" json:"-"`
}

func (table *AdminRole) TableName() string {
	return "admin_role"
}

type AdminRolePermissionRelation struct {
	CommonIntField

	RoleId       int
	PermissionId int
}

func (table *AdminRolePermissionRelation) TableName() string {
	return "admin_role_permission_relation"
}

type AccountRoot struct {
	CommonField

	LoginName  string `gorm:"column:login_name;type:varchar(255);"`
	Username   string
	Password   string
	EncodeType string
	ReqId      string `gorm:"column:req_id;type:varchar(255);"`
}

type AccountAdminRoleRelation struct {
	CommonIntField

	AdminId string
	RoleId  int
}

type AccountAdmin struct {
	CommonField

	LoginName  string `gorm:"column:login_name;type:varchar(255);"`
	Username   string
	Password   string
	Name       string `gorm:"column:name;not null"`
	Super      bool   `gorm:"column:super"`
	SchoolId   string `gorm:"column:school_id"`
	EncodeType string
	ReqId      string `gorm:"column:req_id;type:varchar(255);"`
}

type AccountClient struct {
	CommonField

	Name   string
	Origin string
	Status int
}

func (client *AccountClient) TableName() string {
	return "account_client"
}

type ClientSchoolRelation struct {
	SchoolId string `gorm:"column:school_id;not null"`
	ClientId string `gorm:"column:client_id;not null"`
}

func (relation *ClientSchoolRelation) TableName() string {
	return "account_rel_client_school"
}

type School struct {
	ID       string `gorm:"column:id" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Province string `gorm:"column:province" json:"province"`
	City     string `gorm:"column:city" json:"city"`
}

func (table *School) TableName() string {
	return "account_school"
}

type SchoolDepartment struct {
	CommonField

	SchoolId string `gorm:"school_id"`
	Name     string `gorm:"name"`
	FullName string `gorm:"full_name"`
	IsLeaf   int    `gorm:"is_leaf"`
	ParentId string `gorm:"parent_id"`
}

func (table *SchoolDepartment) TableName() string {
	return "account_school_department"
}

type SchoolCommunity struct {
	CommonField

	SchoolId      string    `gorm:"school_id"`
	Name          string    `gorm:"name"`
	Origin        string    `gorm:"origin"`
	EstablishTime time.Time `gorm:"column:establish_time;type:timestamp"`
}

func (table *SchoolCommunity) TableName() string {
	return "account_school_community"
}

type SchoolUserCertification struct {
	CommonField

	Identify           string `gorm:"column:identify;type:varchar(255);unique"`
	Name               string `gorm:"column:name;type:varchar(255)"`
	SchoolDepartmentId string `gorm:"column:school_department_id;type:varchar(255)"`
}

func (table *SchoolUserCertification) TableName() string {
	return "account_school_user_certification"
}

type AccountUser struct {
	CommonField

	LoginName  string `gorm:"column:login_name;type:varchar(255);"`
	Password   string `gorm:"column:password;type:varchar(255);"`
	EncodeType string `gorm:"column:encode_type"`
}

type AccountUserDetail struct {
	CommonField

	Username string `gorm:"column:username"`
	SchoolId string `gorm:"column:school_id"`
	Verify   bool   `gorm:"column:verify"`
	Birthday string `gorm:"column:birthday"`
	Gender   bool   `gorm:"column:gender"`
	Age      int    `gorm:"column:age"`
	//HeadImgUrl string `gorm:"column:head_img_url"`
	Phone    string `gorm:"column:phone"`
	Wx       string `gorm:"column:wx"`
	QQ       string `gorm:"column:qq"`
	Addr     string `gorm:"column:addr"`
	Hometown string `gorm:"column:hometown"`
	Ethnic   string `gorm:"column:ethnic"`
}

func (table *AccountUserDetail) TableName() string {
	return "account_user_detail"
}

type AccountUserSchoolDetail struct {
	CommonField

	Name     string `gorm:"column:name"`
	Identify string `gorm:"column:identify"`
	Number   string `gorm:"column:number"`
	OrgId    string `gorm:"column:org_id"`
}

func (table *AccountUserSchoolDetail) TableName() string {
	return "account_user_school_detail"
}

type Area struct {
	Code       string
	Name       string
	ParentCode string
	Level      int
}

func (table *Area) TableName() string {
	return "area"
}

type AccountSchoolOrganize struct {
	CommonField
	Label    string
	ParentId string
	SchoolId string
	Status   bool
	Count    int `gorm:"-"`
}

func (table *AccountSchoolOrganize) TableName() string {
	return "account_school_organize"
}

type AccountSchoolStudent struct {
	CommonField

	Name     string
	Identify string
	Status   bool
	Number   string
	Phone    string
	OrgId    string
}

func (table *AccountSchoolStudent) TableName() string {
	return "account_school_student"
}
