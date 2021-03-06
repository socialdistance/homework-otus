package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func worker(done In, outChannel In) Out {
	bindChannel := make(Bi)
	go func() {
		defer close(bindChannel)
		for {
			select {
			case <-done:
				return
			case value, ok := <-outChannel:
				if !ok {
					return
				}
				bindChannel <- value
			}
		}
	}()

	return bindChannel
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outChannel := in
	for _, stage := range stages {
		outChannel = stage(worker(done, outChannel))
	}

	return outChannel
}
