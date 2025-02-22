package main

func handle(err error) {
	if err == nil {
		return
	}
	panic(err)
}
