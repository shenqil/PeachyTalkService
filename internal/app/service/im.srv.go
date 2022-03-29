package service

import (
	"context"
	"ginAdmin/internal/app/model/gormx/repo"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/errors"
	"ginAdmin/pkg/util/hash"

	"github.com/google/wire"
)

// IMSet 注入IM
var IMSet = wire.NewSet(wire.Struct(new(IM), "*"))

// IM 管理
type IM struct {
	UserModel *repo.User
}

// Verify 登陆验证
func (a *IM) Verify(ctx context.Context, userName, password string) error {
	//	检查是否是超级用户
	root := schema.GetRootUser()
	if userName == root.UserName && root.Password == password {
		return nil
	}

	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserName: userName,
		Status:   1,
	})
	if err != nil {
		return err
	} else if len(result.Data) == 0 {
		return errors.ErrInvalidUserName
	}

	item := result.Data[0]
	if item.Password != hash.SHA1String(password) {
		return errors.ErrInvalidPassword
	} else if item.Status != 1 {
		return errors.ErrUserDisable
	}

	return nil
}

// 判断是否为超级用户
func (a *IM) IsSuperUser(userName string) error {
	//	检查是否是超级用户
	root := schema.GetRootUser()
	if userName == root.UserName {
		return nil
	}

	return errors.ErrInvalidUser
}

// alc 校验
func (a *IM) AclVerify(ctx context.Context, username, topic string, access schema.AccessType) error {
	//if b, err := a.CasbinEnforcer.Enforce(username, topic, access); err != nil {
	//	return err
	//} else if !b {
	//	return errors.ErrNoPerm
	//}
	return nil
}
