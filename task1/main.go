package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Task 定义了一个任务
type Task struct {
	Time     time.Time
	Message  string
	Callback func()
}

// TaskQueue 管理任务的map
var TaskQueue = make(map[string]*Task)
var mu sync.Mutex // 用于保护TaskQueue的互斥锁

// addTask 添加到任务队列并启动定时器（如果时间已到达则立即执行）
func addTask(task *Task) {
	mu.Lock()
	defer mu.Unlock()

	key := fmt.Sprintf("%d", task.Time.Unix())
	TaskQueue[key] = task

	if time.Now().After(task.Time) {
		task.Callback()
	} else {
		time.AfterFunc(task.Time.Sub(time.Now()), task.Callback)
	}
}

// taskCallback 是任务的回调函数，打印消息
func taskCallback(tasktime time.Time, message string) func() {
	return func() {
		fmt.Println("Task executed:", message)
		mu.Lock()
		defer mu.Unlock()
		delete(TaskQueue, fmt.Sprintf("%d", tasktime.Unix()))
	}
}

// addTaskHandler 处理HTTP请求以添加任务
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method, use POST", http.StatusBadRequest)
		return
	}
	// 从http请求中提取查询参数
	query := r.URL.Query()
	timeStr := query.Get("timeStr")
	message := query.Get("message")

	const layout = "2006-01-02 15:04:05"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}

	// 从解析后的时间中提取年、月、日、时、分、秒
	year, month, day := t.Date()
	hour, minute, second := t.Clock()

	// 将解析后得到的时间解析成本地时区
	taskTime := time.Date(year, month, day, hour, minute, second, 0, time.Local)

	// 创建任务并添加到队列
	task := &Task{
		Time:     taskTime,
		Message:  message,
		Callback: taskCallback(taskTime, message),
	}
	addTask(task)

	fmt.Println("Task added successfully")
}

func main() {
	http.HandleFunc("/add-task", addTaskHandler)
	fmt.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
