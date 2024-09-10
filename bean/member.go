package bean

import "github.com/gogf/gf/v2/os/gtime"

type MerchantMember struct {
	Id          uint64                    `json:"id"         description:"userId"`                // userId
	GmtCreate   *gtime.Time               `json:"gmtCreate"  description:"create time"`           // create time
	GmtModify   *gtime.Time               `json:"gmtModify"  description:"update time"`           // update time
	MerchantId  uint64                    `json:"merchantId" description:"merchant id"`           // merchant id
	IsDeleted   int                       `json:"isDeleted"  description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	Password    string                    `json:"password"   description:"password"`              // password
	UserName    string                    `json:"userName"   description:"user name"`             // user name
	Mobile      string                    `json:"mobile"     description:"mobile"`                // mobile
	Email       string                    `json:"email"      description:"email"`                 // email
	FirstName   string                    `json:"firstName"  description:"first name"`            // first name
	LastName    string                    `json:"lastName"   description:"last name"`             // last name
	CreateTime  int64                     `json:"createTime" description:"create utc time"`       // create utc time
	Role        string                    `json:"role"       description:"role"`                  // role
	Status      int                       `json:"status"     description:"0-Active, 2-Suspend"`   // 0-Active, 2-Suspend
	Permissions []*MerchantRolePermission `json:"permissions"       description:"Permissions"`    // Permissions
	IsOwner     bool                      `json:"isOwner"       description:"IsOwner"`            // role
}

type MerchantRolePermission struct {
	Group       string   `json:"group"           description:"Group"`             // group
	Permissions []string `json:"permissions"           description:"Permissions"` // group
}
