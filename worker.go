/**
 * Created with IntelliJ IDEA.
 * User: pfeairheller
 * Date: 10/31/13
 * Time: 8:38 PM
 * To change this template use File | Settings | File Templates.
 */
package gopi

type Worker struct {


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


func (w *Worker) AddServer(addr string) (err error) {
	return nil
}

