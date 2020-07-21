package model

import (
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
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
		new(AccountRootLoginLog),
		new(AccountAdmin),
		new(AccountAdminLoginLog),
		new(AccountUser),
		new(AccountUserLoginLog),
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
		new(Ip),
	).Error
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}

type CommonField struct {
	ID        string    `gorm:"column:id;type:char(36);primary_key;not null" json:"id"`
	CreatedAt time.Time `gorm:"column:create_time" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"updated_at"`
}

type CommonIntField struct {
	ID        int       `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `gorm:"column:create_time" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"updated_at"`
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
	SchoolId string `gorm:"column:school_id" json:"school_id"`
	Describe string `gorm:"column:describe" json:"describe"`
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

	LoginName  string `gorm:"column:login_name;type:varchar(255);" json:"login_name"`
	Username   string `json:"username"`
	Password   string `json:"-"`
	EncodeType string `json:"-"`
	Phone      string `json:"phone"`
	ReqId      string `gorm:"column:req_id;type:varchar(255);" json:"req_id"`
}

func (table *AccountRoot) TableName() string {
	return "account_root"
}

type AccountRootLoginLog struct {
	ID         int       `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	RootID     string    `json:"root_id"`
	IP         string    `json:"ip"`
	LoginCount int       `json:"login_count"`
	LoginTime  time.Time `gorm:"column:login_time;type:timestamp" json:"login_time"`
}

func (table *AccountRootLoginLog) TableName() string {
	return "account_root_login_log"
}

type AccountAdminRoleRelation struct {
	CommonIntField

	AdminId string
	RoleId  int
}

type AccountAdmin struct {
	CommonField

	LoginName  string `gorm:"column:login_name;type:varchar(255);" json:"login_name"`
	Username   string `gorm:"column:username" json:"username"`
	Password   string `gorm:"column:password" json:"password"`
	Name       string `gorm:"column:name;not null" json:"name"`
	Super      bool   `gorm:"column:super" json:"super"`
	Phone      string `gorm:"column:phone" json:"phone"`
	SchoolId   string `gorm:"column:school_id" json:"school_id"`
	EncodeType string `gorm:"column:encode_type" json:"encode_type"`
	ReqId      string `gorm:"column:req_id;type:varchar(255);" json:"req_id"`
}

type AccountAdminLoginLog struct {
	ID         int       `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	AdminID    string    `json:"admin_id"`
	IP         string    `json:"ip"`
	LoginCount int       `json:"login_count"`
	LoginTime  time.Time `gorm:"column:login_time;type:timestamp" json:"login_time"`
}

func (table *AccountAdminLoginLog) TableName() string {
	return "account_admin_login_log"
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

	LoginName  string `gorm:"column:login_name;type:varchar(255);" json:"login_name"`
	Password   string `gorm:"column:password;type:varchar(255);" json:"password"`
	EncodeType string `gorm:"column:encode_type" json:"encode_type"`
}

type AccountUserLoginLog struct {
	ID         int       `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	UserID     string    `json:"user_id"`
	LoginCount int       `json:"login_count"`
	IP         string    `json:"ip"`
	LoginTime  time.Time `gorm:"column:login_time;type:timestamp" json:"login_time"`
}

func (table *AccountUserLoginLog) TableName() string {
	return "account_user_login_log"
}

type AccountUserDetail struct {
	CommonField

	AccountID string `gorm:"column:account_id" json:"account_id"`
	Username  string `gorm:"column:username" json:"username"`
	SchoolId  string `gorm:"column:school_id" json:"school_id"`
	Verify    bool   `gorm:"column:verify" json:"verify"`
	Birthday  string `gorm:"column:birthday" json:"birthday"`
	Gender    bool   `gorm:"column:gender" json:"gender"`
	Age       int    `gorm:"column:age" json:"age"`
	//HeadImgUrl string `gorm:"column:head_img_url"`
	Phone    string `gorm:"column:phone" json:"phone"`
	Wx       string `gorm:"column:wx" json:"wx"`
	QQ       string `gorm:"column:qq" json:"qq"`
	Addr     string `gorm:"column:addr" json:"addr"`
	Hometown string `gorm:"column:hometown" json:"hometown"`
	Ethnic   string `gorm:"column:ethnic" json:"ethnic"`
}

func (table *AccountUserDetail) TableName() string {
	return "account_user_detail"
}

type AccountUserSchoolDetail struct {
	CommonField

	Name     string `gorm:"column:name" json:"name"`
	Identify string `gorm:"column:identify" json:"identify"`
	Number   string `gorm:"column:number" json:"number"`
	OrgId    string `gorm:"column:org_id" json:"org_id"`
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

	Name     string `json:"name"`
	Identify string `json:"identify"`
	Status   bool   `json:"status"`
	Number   string `json:"number"`
	Phone    string `json:"phone"`
	OrgId    string `json:"org_id"`
}

func (table *AccountSchoolStudent) TableName() string {
	return "account_school_student"
}

type Ip struct {
	Ip         string    `gorm:"column:ip;primary_key" json:"ip"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	Describe   string    `json:"describe"`
	UpdateTime time.Time `json:"update_time"`
}

func (table *Ip) TableName() string {
	return "account_ip"
}
