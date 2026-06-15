package db

import (
	"echo-framework/internal/model"

	"go.uber.org/zap"
)

type TableModel interface {
	TableName() string
}

var models = []TableModel{
	new(model.SysUser),
}

// SyncDatabase 同步数据库
func SyncDatabase() error {
	db, err := GetDB()
	if err != nil {
		zap.L().Error("failed to get db", zap.Error(err))
		return err
	}

	for _, m := range models {
		err = db.AutoMigrate(m)
		if err != nil {
			zap.L().Error("failed to sync table", zap.Error(err), zap.String("table", m.TableName()))
			return err
		}
		zap.L().Info("table sync successfully", zap.String("table", m.TableName()))
	}
	return nil
}
