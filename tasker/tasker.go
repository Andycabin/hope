// tasker
package tasker

import (
	"hope/http/download"
	"hope/http/request"
	"log"
)

type Tasker struct {
	taskChan     chan *request.Request
	executorChan chan chan *request.Request
	resultChan   chan *request.ParseResult
}

// 向taskChan提交Task
func (t *Tasker) Submit(r *request.Request) {
	t.taskChan <- r
}

// 任务执行器创建自己的chan后传入任务执行器队列中
func (t *Tasker) ExecutorReady(p chan *request.Request) {
	t.executorChan <- p
}

// 启动任务管理器
func (t *Tasker) Run(resultChan chan *request.ParseResult, taskChan chan *request.Request) {
	// taskChan，调度器将Request传入任务管理器的chan
	t.taskChan = taskChan
	// performChan，执行器Chan（管理根据传入goroutineCount生成多个任务执行器）
	// 每个任务执行器自己一个chan
	t.executorChan = make(chan chan *request.Request)
	// resultChan，传输结果
	t.resultChan = resultChan
	go func() {
		// 任务队列，执行器队列
		var taskQ []*request.Request
		var executorQ []chan *request.Request
		for {
			// 将一个任务放入一个执行器中
			var activeTask *request.Request
			var activeExecutor chan *request.Request
			if len(taskQ) > 0 && len(executorQ) > 0 {
				activeTask = taskQ[0]
				activeExecutor = executorQ[0]
			}
			// 新任务放入任务队列，新任务执行器放入执行器队列
			select {
			case task := <-t.taskChan:
				taskQ = append(taskQ, task)
			case executor := <-t.executorChan:
				executorQ = append(executorQ, executor)
			case activeExecutor <- activeTask:
				taskQ = taskQ[1:]
				executorQ = executorQ[1:]
			}
		}
	}()
}

// 任务执行器
func (t *Tasker) TaskExecutor(taskNum int) {
	go func() {
		for {
			// 创建属于任务执行器得chan,传给任务管理器管理，同时自己也从自己的Chan里面读取Request
			myChan := make(chan *request.Request)
			t.ExecutorReady(myChan)
			req := <-myChan
			log.Printf("[Hope engine]: <Executor %d> run task <crawling %s>", taskNum, req.Url)
			resp, err := download.Download(req)
			if err != nil {
				log.Printf("[Hope engine]: <Executor %d> run task <crawling %s> fail, error %v", taskNum, req.Url, err)
				continue
			}
			log.Printf("[Hope engine]: <Executor %d> run task <crawling %s> success, status code %d", taskNum, req.Url, resp.StatusCode)
			log.Printf("[Hope engine]: <Executor %d> run task <parsing %s>", taskNum, req.Url)
			result := resp.ParseFunc(resp.Body)
			log.Printf("[Hope engine]: <Executor %d> run task <parsing %s> success, got %d requests %d items", taskNum, req.Url, len(result.Requests), len(result.Items))
			t.resultChan <- result
		}
	}()
}
