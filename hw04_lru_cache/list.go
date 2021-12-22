package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length int
	front  *ListItem
	back   *ListItem
}

func (l *list) empty() {
	l.front = nil
	l.back = nil
	l.length = 0
}

func (l *list) FirstItem(item *ListItem) {
	l.back = item
	l.front = item

	item.Prev = nil
	item.Next = nil

	l.length++
}

func (l *list) pushingFront(item *ListItem) bool {
	return l.front == item
}

func (l *list) pushingBack(item *ListItem) bool {
	return l.back == item
}

func (l *list) PushingFrontItem(item *ListItem) {
	old := l.front

	if l.Len() == 0 {
		l.FirstItem(item)
		return
	}

	old.Prev = item
	item.Next = old

	l.front = item
	l.length++
}

func (l *list) PushingBackItem(item *ListItem) {
	old := l.back
	if l.Len() == 0 {
		l.FirstItem(item)
		return
	}

	old.Next = item
	item.Prev = old

	l.back = item
	l.length++
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  l.front,
	}

	l.PushingFrontItem(item)

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Prev:  l.back,
	}

	l.PushingBackItem(item)

	return item
}

func (l *list) Remove(item *ListItem) {
	if l.Len() == 1 {
		l.empty()
		return
	}

	l.length--

	if l.pushingFront(item) {
		l.front = item.Next
		l.front.Prev = nil
		return
	}

	if l.pushingBack(item) {
		l.back = item.Prev
		l.back.Next = nil
		return
	}

	item.Prev.Next = item.Next
	item.Next.Prev = item.Prev
}

func (l *list) MoveToFront(item *ListItem) {
	if l.Len() == 1 {
		return
	}

	l.Remove(item)
	l.PushingFrontItem(item)
}

func NewList() List {
	return new(list)
}
