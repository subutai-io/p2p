package ptp

// Cross-peer communication handlers

// commStatusReportHandler handles status reports from another peer
func commStatusReportHandler(data []byte) error {
	return nil
}

// commSubnetInfoHandler request/response of network subnet. Data format is as follows:
// hash[36] - subnet[?]
// If subnet is empty, that means that this is a request. Hash is a mandatory, but just for a sanity check
func commSubnetInfoHandler(data []byte) error {
	//hash := data[0:36]

	return nil
}
