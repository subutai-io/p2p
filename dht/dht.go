package dht

type DHTClient struct {
	id string
}

func (dht *DHTClient) SetId(id string) {
	dht.id = id
}
