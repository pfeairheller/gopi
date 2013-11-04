/**
 * Created with IntelliJ IDEA.
 * User: pfeairheller
 * Date: 10/31/13
 * Time: 8:38 PM
 * To change this template use File | Settings | File Templates.
 */
package gopi

import (
	"errors"
)

const (
    Unlimited = 0
    OneByOne = 1
	QueueSize = 8
    Immediately = 0
)

type Worker struct {
	funcs JobFuncs
	agents []*agent
	in chan *Job
	IsRunning bool
	ErrHandler func(error)
}

type JobHandler func(*Job) error

type JobFunc func(*Job) ([]byte, error)

// The definition of the callback function.
type jobFunc struct {
    f JobFunc
    numberOfWorkers int
	c chan []byte
}

// Map for added function.
type JobFuncs map[string]*jobFunc


func NewWorker() (w *Worker) {
	w = new(Worker)
	w.funcs = make(JobFuncs)
	w.IsRunning = false
	w.in = make(chan *Job, QueueSize)
	return w

}

func (w *Worker) AddServer(addr string) (err error) {
	agent, err := newAgent(addr, w)
	if err != nil {
		return
	}

	w.agents = append(w.agents, agent)

	return
}

func (w *Worker) AddFunc(funcname string, f JobFunc, numberOfWorkers int) (err error) {
	if _, ok := w.funcs[funcname]; ok {
		return errors.New("The function already exists: "+ funcname)
	}
	w.funcs[funcname] = &jobFunc{f: f, numberOfWorkers: numberOfWorkers, c: make(chan []byte)}

//  if w.running {
//      w.addFunc(funcname, timeout)
//  }
	return

}

func (w *Worker) RemoveFunc(funcname string) (err error) {
	if _, ok := w.funcs[funcname]; ok {
		return errors.New("The function does not exist: "+ funcname)
	}
	delete (w.funcs, funcname)

//	if worker.running {
//     worker.removeFunc(funcname)
//	}
	return

}

func (w *Worker) Work() {
	defer func() {
		for _, v := range w.agents {
			v.Close()
		}
	}()
	w.IsRunning = true
	for _, v := range w.agents {
		go v.Work()
	}

	ok := true
	for ok {
		var job *Job
		if job, ok = <-w.in; ok {
			go w.dealJob(job)
		}
	}
}

func (w *Worker) dealJob(job *Job) {
    defer func() {
        job.Close()
    }()

	if err := w.exec(job); err != nil {

		if w.JobHandler != nil {
			if err := w.JobHandler(job); err != nil {
				w.err(err)
			}
		}
	}

}

func (worker *Worker) exec(job *Job) (err error) {
	f, ok := worker.funcs[job.Fn]
	if !ok {
		return errors.New("The function does not exist: " + job.Fn)
	}

	f.f(job)

	return
}


