package ctypes

import "fmt"

type AppStep string

// Start steps
const (
	StartStepDbConnection    AppStep = "Database connection"
	StartStepNsqConnection   AppStep = "NSQ connection"
	StartStepRedisConnection AppStep = "Redis connection"
	StartStepSqlStoreInit    AppStep = "SQL store init"
	StartStepNsqStoreInit    AppStep = "NSQ store init"
	StartStepPluginsInit     AppStep = "Plugins init"
	StartStepLoggerInit      AppStep = "Logger init"
)

// Listener steps
const (
	ListenerStepGetAllMerachants AppStep = "Get all merchants from database"
	ListenerStepDecodeParams     AppStep = "Decode options params for account"
)

func SprintfStep(format string, args ...interface{}) AppStep {
	str := fmt.Sprintf(format, args...)
	return AppStep(str)
}
