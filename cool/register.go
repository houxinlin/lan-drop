package cool

import "cool-transmission/common"

var taskMap map[string]common.TaskInfo

func init() {
	taskMap = make(map[string]common.TaskInfo)
}
func RegisterTask(taskId string, info common.TaskInfo) {
	taskMap[taskId] = info
}

func GetTaskInfo(taskId string) common.TaskInfo {
	return taskMap[taskId]
}
