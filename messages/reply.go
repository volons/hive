package messages

type Reply struct {
	result interface{}
	err    error
}

func newReply(result interface{}, err error) Reply {
	return Reply{
		result: result,
		err:    err,
	}
}

func (r Reply) Result() interface{} {
	return r.result
}

func (r Reply) Err() error {
	return r.err
}
