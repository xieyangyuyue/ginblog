package model

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"time"
)

// GormLogger 自定义GORM日志器（关联请求上下文）
type GormLogger struct {
	LogLevel logger.LogLevel
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &GormLogger{LogLevel: level}
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logrus.WithContext(ctx).Infof(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logrus.WithContext(ctx).Warnf(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logrus.WithContext(ctx).Errorf(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	sql, rows := fc()
	fields := logrus.Fields{
		"RequestID": getRequestID(ctx),
		"SQL":       sql,
		"Rows":      rows,
		"Duration":  time.Since(begin).String(),
	}
	if err != nil {
		logrus.WithContext(ctx).WithFields(fields).Error("SQL Error")
	} else {
		logrus.WithContext(ctx).WithFields(fields).Debug("SQL Executed")
	}
}

// 从上下文获取请求ID
func getRequestID(ctx context.Context) string {
	if id, ok := ctx.Value("RequestID").(string); ok {
		return id
	}
	return ""
}
