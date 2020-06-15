package cmd

type Timeline struct {
	PagingToken string
	Moments     []Moment
}

type TimelineRequest struct {
}

type Moment struct {
}

func Serve(request TimelineRequest) (Timeline, error) {
	return Timeline{
		PagingToken: "someToken",
	}, nil
}
