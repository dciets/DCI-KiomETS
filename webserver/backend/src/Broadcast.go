package main

type Subject interface {
	register(observer Observer)
	deregister(observer Observer)
	notifyAll()
}
type Broadcast struct {
	observers []Observer
	data      string
}

func (b *Broadcast) register(o Observer) {
	b.observers = append(b.observers, o)
}

func removeFromslice(observerList []Observer, observerToRemove Observer) []Observer {
	observerListLength := len(observerList)
	for i, observer := range observerList {
		if observerToRemove.getID() == observer.getID() {
			observerList[observerListLength-1], observerList[i] = observerList[i], observerList[observerListLength-1]
			return observerList[:observerListLength-1]
		}
	}
	return observerList
}

func (b *Broadcast) deregister(o Observer) {
	b.observers = removeFromslice(b.observers, o)
}

func (b *Broadcast) notifyAll() {
	for _, observer := range b.observers {
		observer.update(b.data)
	}
}

func (b *Broadcast) setData(data string) {
	b.data = data
	b.notifyAll()
}

var gameBroadcast = &Broadcast{observers: []Observer{}, data: ""}
var scoreBroadcast = &Broadcast{observers: []Observer{}, data: ""}
