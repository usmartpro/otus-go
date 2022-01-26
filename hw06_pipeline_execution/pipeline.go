package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = stage(checkStage(done, out))
	}

	return out
}

func checkStage(done In, in In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case value, ok := <-in:
				if ok {
					// put to out chan
					out <- value
				} else {
					// empty chan
					return
				}
			case <-done:
				// terminate
				return
			}
		}
	}()

	return out
}
