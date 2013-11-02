/**
 * Created with IntelliJ IDEA.
 * User: pfeairheller
 * Date: 10/31/13
 * Time: 8:38 PM
 * To change this template use File | Settings | File Templates.
 */
package gopi

import "errors"

type Worker struct {
	funcs JobFuncs
	servers []string
	IsRunning bool
}

type JobHandler func(*Job) error

type JobFunc func(*Job) ([]byte, error)

// The definition of the callback function.
type jobFunc struct {
    f JobFunc
    timeout uint32
}

// Map for added function.
type JobFuncs map[string]*jobFunc


func NewWorker() (w *Worker) {
	w = new(Worker)
	w.funcs = make(JobFuncs)
	w.IsRunning = false
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

func (w *Worker) AddFunc(funcname string, f JobFunc) (err error) {
	if _, ok := w.funcs[funcname]; ok {
		return errors.New("The function already exists: "+ funcname)
	}
	w.funcs[funcname] = &jobFunc{f: f, timeout: w}

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

