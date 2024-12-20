package global

import (
	"go.uber.org/zap"
	"sync"
)

// 全局日志实例
var (
	Sugar *zap.SugaredLogger
	mutex sync.RWMutex // 创建一个 RWMutex 实例
)

func GetSugar() *zap.SugaredLogger {
	mutex.RLock()
	defer mutex.RUnlock()

	return Sugar
}

func SetSugar(newSugar *zap.SugaredLogger) {
	mutex.Lock()         // 获取写锁
	defer mutex.Unlock() // 确保释放锁

	Sugar = newSugar
}
