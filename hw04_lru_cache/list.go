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

func NewListItem(v interface{}) *ListItem {
	return &ListItem{Value: v}
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) setFront(i *ListItem) {
	i.Prev = nil
	connectWith(i, l.front)

	l.front = i

	if l.back == nil {
		l.back = i
	}

	l.len++
}

func (l *list) setBack(i *ListItem) {
	i.Next = nil
	connectWith(l.back, i)

	l.back = i

	if l.front == nil {
		l.front = i
	}

	l.len++
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := NewListItem(v)

	l.setFront(i)

	return i
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := NewListItem(v)

	l.setBack(i)

	return i
}

func (l *list) Remove(i *ListItem) {
	connectWith(i.Prev, i.Next)

	if i.Prev == nil {
		l.front = i.Next
	}

	if i.Next == nil {
		l.back = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}

	l.Remove(i)
	l.setFront(i)
}

func NewList() List {
	return new(list)
}

func connectWith(prev, next *ListItem) {
	if prev != nil {
		prev.Next = next
	}

	if next != nil {
		next.Prev = prev
	}
}
