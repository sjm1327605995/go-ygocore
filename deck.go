package main

type Deck struct {
	main  []*CardDataC
	extra []*CardDataC
	side  []*CardDataC
}

func (d *Deck) Clear() {
	d.main = nil
	d.extra = nil
	d.side = nil
}
