package user

import (
	"echo-framework/internal/infrastructure/db"
	"echo-framework/internal/model"

	"gorm.io/gorm"
)

func getDB() (*gorm.DB, error) {
	return db.GetDB()
}

// FindAll 查询全部用户
func FindAll() ([]*model.SysUser, error) {
	dbCli, err := getDB()
	if err != nil {
		return nil, err
	}
	var users []*model.SysUser
	if err := dbCli.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// FindByID 根据ID查询用户
func FindByID(id int) (*model.SysUser, error) {
	dbCli, err := getDB()
	if err != nil {
		return nil, err
	}
	var u model.SysUser
	if err := dbCli.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByUserName 根据用户名查询
func FindByUserName(userName string) (*model.SysUser, error) {
	dbCli, err := getDB()
	if err != nil {
		return nil, err
	}
	var u model.SysUser
	if err := dbCli.Where("user_name = ?", userName).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Create 创建用户
func Create(user *model.SysUser) error {
	dbCli, err := getDB()
	if err != nil {
		return err
	}
	return dbCli.Create(user).Error
}

// UpdateByID 根据ID更新用户
func UpdateByID(id int, user *model.SysUser) error {
	dbCli, err := getDB()
	if err != nil {
		return err
	}
	return dbCli.Model(&model.SysUser{}).Where("id = ?", id).Updates(user).Error
}

// DeleteByID 根据ID删除用户
func DeleteByID(id int) error {
	dbCli, err := getDB()
	if err != nil {
		return err
	}
	return dbCli.Delete(&model.SysUser{}, id).Error
}
