package service

import (
	"PeachyTalkService/internal/app/model/gormx/repo"
	"PeachyTalkService/internal/app/model/miniox/bucket"
	"PeachyTalkService/internal/app/schema"
	"PeachyTalkService/pkg/errors"
	"PeachyTalkService/pkg/logger"
	"PeachyTalkService/pkg/util/font2Img"
	"PeachyTalkService/pkg/util/hash"
	"PeachyTalkService/pkg/util/uuid"
	"context"
	"github.com/google/wire"
)

// UserSet 注入User
var UserSet = wire.NewSet(wire.Struct(new(User), "*"))

// User 用户管理
type User struct {
	TransModel  *repo.Trans
	UserModel   *repo.User
	AvatarModel *bucket.Avatar
}

// Query 查询数据
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	return a.UserModel.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *User) Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	item, err := a.UserModel.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

// GetByUserName 根据用户名称查询指定数据
func (a *User) GetByUserName(ctx context.Context, userName string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	item, err := a.UserModel.GetByUserName(ctx, userName, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

// Create 创建数据
func (a *User) Create(ctx context.Context, item schema.User) (*schema.IDResult, error) {
	err := a.checkUserName(ctx, item)
	if err != nil {
		return nil, err
	}

	item.Password = hash.SHA1String(item.Password)
	item.ID = uuid.MustString()

	if item.Avatar == "" {
		reader, err := font2Img.GetReader(item.RealName)
		if err == nil {
			info, err := a.AvatarModel.Upload(ctx, item.UserName, reader, -1, "image/png")
			if err == nil {
				item.Avatar = info.Key
			} else {
				logger.WithContext(ctx).Errorf("[user.srv][Create] UploadAvatar err = %s", err.Error())
			}
		} else {
			logger.WithContext(ctx).Errorf("[user.srv][Create] GetReader err = %s", err.Error())
		}
	}

	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		return a.UserModel.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

func (a *User) checkUserName(ctx context.Context, item schema.User) error {
	if item.UserName == schema.GetRootUser().UserName {
		return errors.New400Response("用户名不合法")
	}

	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		UserName:        item.UserName,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("用户名已经存在")
	}
	return nil
}

// Update 更新数据
func (a *User) Update(ctx context.Context, id string, item schema.User) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.UserName != item.UserName {
		err := a.checkUserName(ctx, item)
		if err != nil {
			return err
		}
	}

	if item.Password != "" {
		item.Password = hash.SHA1String(item.Password)
	} else {
		item.Password = oldItem.Password
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		return a.UserModel.Update(ctx, id, item)
	})
	if err != nil {
		return err
	}

	return nil
}

// Delete 删除数据
func (a *User) Delete(ctx context.Context, id string) error {
	oldItem, err := a.UserModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		return a.UserModel.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	return nil
}

// BatchDelete 批量删除
func (a *User) BatchDelete(ctx context.Context, ids []string) error {

	err := a.TransModel.Exec(ctx, func(ctx context.Context) error {
		return a.UserModel.BatchDelete(ctx, ids)
	})
	if err != nil {
		return err
	}

	return nil
}

// UpdateStatus 更新状态
func (a *User) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.UserModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}
	oldItem.Status = status

	err = a.UserModel.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}
	return nil
}
