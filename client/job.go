package client


type Job struct {
   Data []byte
}

type FutureResult struct {
	Success chan Job
	Failure chan error
}



