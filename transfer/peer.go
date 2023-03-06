package transfer

type Peer struct{}

func NewPeer() *Peer {
	return new(Peer)
}

// LocalIPAddr return current peer local ipv4 and port, if in received peer, you should omit port filed.
func (p *Peer) LocalIPAddr() ([]string, error) {
	ip, _, err := localIPv4WithNetwork()
	if err != nil {
		return []string{}, err
	}
	port := servPort
	return []string{ip, port}, nil
}
