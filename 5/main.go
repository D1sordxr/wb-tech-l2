package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	// если для err задать тип *customError, то код выведет "ok"
	var err error // создается err с типом и значением nil
	err = test()  // err присваивается тип *customError => err != nil
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
