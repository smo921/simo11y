package main

func tee(done <-chan string, in <-chan map[string]interface{}, out chan<- map[string]interface{}) <-chan map[string]interface{} {
	forward := make(chan map[string]interface{})
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			default:
				message := <-in
				out <- message
				forward <- message
			}
		}
	}()
	return forward
}
