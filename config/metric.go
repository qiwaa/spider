package config

type receiveMetric struct {
	Ping         uint16
	FindNode     uint16
	GetPeers     uint16
	AnnouncePeer uint16
}

type Metric struct {
	Receive *receiveMetric
}

func newReceiveMetric() (*receiveMetric) {
	return &receiveMetric{
		GetPeers:     0,
		AnnouncePeer: 0,
	}
}

func NewMetric() (*Metric) {
	return &Metric{
		Receive: newReceiveMetric(),
	}
}
