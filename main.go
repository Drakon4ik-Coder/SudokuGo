package main

func frame() bool {
	return true
}

func game() bool {
	frameContinue := frame()
	for frameContinue {
		frameContinue = frame()
		ClearConsole()
	}
	return true
}

func main() {
	game()
}
