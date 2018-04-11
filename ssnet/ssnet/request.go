package ssnet

type GetRequest struct {
	Url    string
	SSlist SteppingStones
}

func NewGetRequest(u string, sslist SteppingStones) GetRequest {
	return GetRequest{
		Url:    u,
		SSlist: sslist,
	}
}
