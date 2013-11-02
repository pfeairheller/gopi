/**
 * Created with IntelliJ IDEA.
 * User: pfeairheller
 * Date: 10/31/13
 * Time: 8:03 PM
 * To change this template use File | Settings | File Templates.
 */
package gopi

import "encoding/json"


type Job struct {
	Data []byte
}

func (job *Job) Value(target interface{}) (err error) {
	return json.Unmarshal(job.Data, target)
}



