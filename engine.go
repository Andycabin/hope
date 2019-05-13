// engine
// 爬虫核心
// 控制整个爬虫流程：生成请求，下载，解析，存储
package engine

import (
	"log"

	"hope/http/request"
	"hope/model"
	"hope/tasker"
)

// 核心引擎
// 调度器Scheduler
type Engine struct {
	Tasker         Tasker
	ItemSaver      *model.DB
	GoroutineCount int
}

type Tasker interface {
	Submit(*request.Request)
	ExecutorReady(chan *request.Request)
	Run(chan *request.ParseResult, chan *request.Request)
	TaskExecutor(taskNum int)
}

// 启动方法
func (e *Engine) Run(seeds ...*request.Request) {
	// 传入初始Request，可一个，也可一个数组
	// 结果Chan，任务Chan
	resultChan := make(chan *request.ParseResult)
	taskChan := make(chan *request.Request)
	// 启动任务管理器，启动存储器
	e.Tasker.Run(resultChan, taskChan)
	e.ItemSaver.Run()

	// 注册模型
	// model.Register(new(model.Profile))
	// 根据GoroutineCount创建任务执行器
	for i := 1; i <= e.GoroutineCount; i++ {
		e.Tasker.TaskExecutor(i)
	}
	for _, r := range seeds {
		e.Tasker.Submit(r)
	}
	itemCount := 1
	for {
		result := <-resultChan
		for _, item := range result.Items {
			log.Printf("[Hope engine]: <Core> got #%d items, %s", itemCount, item)
			go func() { e.ItemSaver.ResultChan <- item }()
			itemCount++
		}
		for _, req := range result.Requests {
			e.Tasker.Submit(req)
		}
	}
}

func NewSpider(driverName, driverAddress string, goNum int) *Engine {
	return &Engine{
		Tasker:         &tasker.Tasker{},
		ItemSaver:      model.NewDB(driverName, driverAddress),
		GoroutineCount: goNum,
	}
}
