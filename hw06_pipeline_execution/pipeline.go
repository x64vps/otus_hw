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
		out = stage(doneStage(out, done))
	}

	return doneStage(out, done)
}

func doneStage(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}

				select {
				case out <- v:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()

	return out
}
