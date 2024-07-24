package core

func ErrorHandle(err error) {
	// @TODO: Database handle and Log handle
	logHandle(err)
	jsonHandle(err)
}

func logHandle(err error) {
	//fmt.Println(err.Error())
}

func jsonHandle(err error) {}
