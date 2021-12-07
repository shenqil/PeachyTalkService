package service

import (
	"context"
	"ginAdmin/internal/app/model/gormx/repo"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/errors"
	"ginAdmin/pkg/util/uuid"
	"github.com/google/wire"
)

// GroupSet 注入Group
var GroupSet = wire.NewSet(wire.Struct(new(Group)), "*")

// Group 群组管理
type Group struct {
	TransModel       *repo.Trans
	GroupModel       *repo.Group
	GroupMemberModel *repo.GroupMember
}

// 判断是否是群组拥有者
func (a *Group) checkGroupOwner(ctx context.Context, userId, groupID string) error {
	result, err := a.GroupModel.Get(ctx, groupID)
	if err != nil {
		return err
	} else if result == nil {
		return errors.ErrNotFound
	} else if result.Owner != userId {
		return errors.ErrNoPerm
	}

	return nil
}

// Query 查询数据
func (a *Group) Query(ctx context.Context, params schema.GroupQueryParam, opts ...schema.GroupQueryOptions) (*schema.GroupQueryResult, error) {
	// 查询所有群组
	result, err := a.GroupModel.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	} else if result == nil {
		return nil, nil
	}

	// 查询群组中相关的所有成员
	groupMembers, err := a.GroupMemberModel.Query(ctx, schema.GroupMemberQueryParam{
		GroupIDs: result.Data.ToGroupIDs(),
	})

	//	合并成员数据
	m := groupMembers.Data.ToGroupIDMap()
	for _, datum := range result.Data {
		datum.MemberIDs = m[datum.ID].ToMemberIDs()
	}

	return result, nil
}

// Get 查询指定数据
func (a *Group) Get(ctx context.Context, id string, opts ...schema.GroupQueryOptions) (*schema.Group, error) {
	item, err := a.GroupModel.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	groupMembers, err := a.GroupMemberModel.Query(ctx, schema.GroupMemberQueryParam{
		GroupIDs: []string{id},
	})
	if err != nil {
		return nil, err
	}

	item.MemberIDs = groupMembers.Data.ToMemberIDs()

	return item, nil
}

// Create 创建数据
func (a *Group) Create(ctx context.Context, item schema.Group) (*schema.IDResult, error) {
	item.ID = uuid.MustString()
	err := a.TransModel.Exec(ctx, func(ctx context.Context) error {
		// 创建关联关系
		for _, memberID := range item.MemberIDs {
			err := a.GroupMemberModel.Create(ctx, schema.GroupMember{
				ID:      uuid.MustString(),
				GroupID: item.ID,
				UserID:  memberID,
			})
			if err != nil {
				return err
			}
		}

		// 创建群组
		return a.GroupModel.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

// Update 更新数据
func (a *Group) Update(ctx context.Context, id, userId string, item schema.Group) (*schema.Group, error) {
	// 群主可以更新群信息
	err := a.checkGroupOwner(ctx, userId, id)
	if err != nil {
		return nil, err
	}

	// 获取旧数据
	oldItem, err := a.GroupModel.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	}

	OldGroupMembers, err := a.GroupMemberModel.Query(ctx, schema.GroupMemberQueryParam{
		GroupIDs: []string{id},
	})
	if err != nil {
		return nil, err
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		addIDsList, delList := a.compareGroupMembers(ctx, OldGroupMembers.Data, item.MemberIDs)
		for _, memberId := range addIDsList {
			err := a.GroupMemberModel.Create(ctx, schema.GroupMember{
				ID:      uuid.MustString(),
				GroupID: item.ID,
				UserID:  memberId,
			})
			if err != nil {
				return err
			}
		}

		for _, member := range delList {
			err := a.GroupMemberModel.Delete(ctx, member.ID)
			if err != nil {
				return err
			}
		}

		return a.GroupModel.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (a *Group) compareGroupMembers(ctx context.Context, oldGroupMembers schema.GroupMembers, newGroupMemberIds []string) (addIDsList []string, delList schema.GroupMembers) {

	for _, id := range newGroupMemberIds {
		flag := false
		for _, member := range oldGroupMembers {
			if member.UserID == id {
				flag = true
				continue
			}
		}

		if !flag {
			addIDsList = append(addIDsList, id)
		}
	}

	for _, member := range oldGroupMembers {
		flag := false
		for _, id := range newGroupMemberIds {
			if member.ID == id {
				flag = true
				continue
			}
		}
		if !flag {
			delList = append(delList, member)
		}
	}

	return
}

// Delete 删除数据
func (a *Group) Delete(ctx context.Context, id, userId string) (*schema.Group, error) {
	// 群主可以更新群信息
	err := a.checkGroupOwner(ctx, userId, id)
	if err != nil {
		return nil, err
	}

	// 获取旧数据
	oldItem, err := a.GroupModel.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	}

	OldGroupMembers, err := a.GroupMemberModel.Query(ctx, schema.GroupMemberQueryParam{
		GroupIDs: []string{id},
	})
	if err != nil {
		return nil, err
	}
	oldItem.MemberIDs = OldGroupMembers.Data.ToMemberIDs()

	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		for _, id := range oldItem.MemberIDs {
			err := a.GroupMemberModel.Delete(ctx, id)
			if err != nil {
				return err
			}
		}

		return a.GroupModel.Delete(ctx, id)
	})

	return oldItem, err
}

// AddMembers 添加成员
func (a *Group) AddMembers(ctx context.Context, params schema.GroupMemberChangesParam) ([]*schema.GroupMemberChangesInfo, []string, error) {
	// 拿到旧数据
	oldGroupMembers, err := a.GroupMemberModel.Query(ctx, schema.GroupMemberQueryParam{
		GroupIDs: []string{params.GroupID},
	})
	if err != nil {
		return nil, nil, err
	}
	oldMemberIDs := oldGroupMembers.Data.ToMemberIDs()
	allMemberIDs := oldGroupMembers.Data.ToMemberIDs()

	// 得到需要添加的成员
	addList := make([]*schema.GroupMemberChangesInfo, 0)
	for _, info := range params.List {
		flag := false
		for _, id := range oldMemberIDs {
			if id == info.ID {
				flag = true
				break
			}
		}

		if !flag {
			// 不存在则添加成员
			addList = append(addList, info)
			allMemberIDs = append(allMemberIDs, info.ID)
		}
	}

	// 入库
	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		for _, item := range addList {
			err := a.GroupMemberModel.Create(ctx, schema.GroupMember{
				ID:      uuid.MustString(),
				GroupID: params.GroupID,
				UserID:  item.ID,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return addList, allMemberIDs, nil
}

// DelMembers 删除成员
func (a *Group) DelMembers(ctx context.Context, params schema.GroupMemberChangesParam) ([]*schema.GroupMemberChangesInfo, []string, error) {
	// 拿到旧数据
	oldGroupMembers, err := a.GroupMemberModel.Query(ctx, schema.GroupMemberQueryParam{
		GroupIDs: []string{params.GroupID},
	})
	if err != nil {
		return nil, nil, err
	}

	// 得到需要删除的成员
	delList := make(schema.GroupMembers, 0)
	memberList := make([]*schema.GroupMemberChangesInfo, 0)
	for _, info := range params.List {
		for _, item := range oldGroupMembers.Data {
			if item.UserID == info.ID {
				// 存在则删除成员
				delList = append(delList, item)
				memberList = append(memberList, info)
				break
			}
		}
	}

	//	删除
	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		for _, item := range delList {
			err := a.GroupMemberModel.Delete(ctx, item.ID)
			if err != nil {
				return err
			}
		}
		if len(delList) == len(oldGroupMembers.Data) {
			return a.GroupModel.Delete(ctx, params.GroupID)
		}
		return nil
	})

	return memberList, oldGroupMembers.Data.ToMemberIDs(), nil
}

// ExitGroup 退出群聊
func (a *Group) ExitGroup(ctx context.Context, params schema.GroupMemberChangesParam) ([]string, error) {
	if len(params.List) != 1 {
		return nil, errors.New400Response("参数错误")
	}

	memberItem := params.List[0]

	// 拿到旧数据
	oldGroupMembers, err := a.GroupMemberModel.Query(ctx, schema.GroupMemberQueryParam{
		GroupIDs: []string{params.GroupID},
	})
	if err != nil {
		return nil, err
	}

	//	判断数据是否存在
	var delItem *schema.GroupMember
	for _, item := range oldGroupMembers.Data {
		if item.UserID == memberItem.ID {
			delItem = item
		}
	}
	if delItem == nil {
		return nil, errors.New400Response("不存在该成员")
	}

	err = a.GroupMemberModel.Delete(ctx, delItem.ID)
	if err != nil {
		return nil, err
	}

	return oldGroupMembers.Data.ToMemberIDs(), nil
}
