package core

func ErrorHandle(err error) {
	logHandle(err)
	jsonHandle(err)
}

func logHandle(err error) {
	//fmt.Println(err.Error())
}

func jsonHandle(err error) {

}
