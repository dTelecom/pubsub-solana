package pubsub

func (p *PubSub) Subscribe(topic string, handler Handler) {
	p.subscriptionsMu.Lock()
	defer p.subscriptionsMu.Unlock()

	p.subscriptions[topic] = append(p.subscriptions[topic], handler)
}
