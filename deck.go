package main

type Deck struct {
	mainLen  int
	extraLen int
	sideLen  int
	main     []*CardDataC
	extra    []*CardDataC
	side     []*CardDataC
}

func (d *Deck) Clear() {
	d.ClearMain()
	d.ClearExtra()
	d.ClearSide()
}
func (d *Deck) ClearMain() {
	d.mainLen = 0
}
func (d *Deck) ClearExtra() {
	d.extraLen = 0
}
func (d *Deck) ClearSide() {
	d.sideLen = 0
}
