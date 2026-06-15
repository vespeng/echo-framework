package user

import (
	"echo-framework/internal/model"
	"errors"
)

// Service 用户业务接口
type Service interface {
	GetAllUsers() ([]*model.SysUser, error)
	GetUserByID(id int) (*model.SysUser, error)
	GetUserByUserName(userName string) (*model.SysUser, error)
	AddUser(user *model.SysUser) error
	ModifyUser(id int, user *model.SysUser) error
	RemoveUser(id int) error
}

// service 用户业务实现
type service struct{}

// NewService 创建用户业务实例
func NewService() Service {
	return &service{}
}

// GetAllUsers 获取全部用户
func (s *service) GetAllUsers() ([]*model.SysUser, error) {
	return FindAll()
}

// GetUserByID 根据ID获取用户
func (s *service) GetUserByID(id int) (*model.SysUser, error) {
	if id <= 0 {
		return nil, errors.New("无效的用户ID")
	}
	u, err := FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return u, nil
}

// GetUserByUserName 根据用户名获取用户
func (s *service) GetUserByUserName(userName string) (*model.SysUser, error) {
	if userName == "" {
		return nil, errors.New("用户名不能为空")
	}
	return FindByUserName(userName)
}

// AddUser 创建用户
func (s *service) AddUser(user *model.SysUser) error {
	if user.UserName == "" {
		return errors.New("用户名不能为空")
	}
	// 检查用户名是否已存在
	_, err := FindByUserName(user.UserName)
	if err == nil {
		return errors.New("用户名已存在")
	}
	return Create(user)
}

// ModifyUser 更新用户
func (s *service) ModifyUser(id int, user *model.SysUser) error {
	if id <= 0 {
		return errors.New("无效的用户ID")
	}
	// 检查用户是否存在
	if _, err := FindByID(id); err != nil {
		return errors.New("用户不存在")
	}
	return UpdateByID(id, user)
}

// RemoveUser 删除用户
func (s *service) RemoveUser(id int) error {
	if id <= 0 {
		return errors.New("无效的用户ID")
	}
	return DeleteByID(id)
}
