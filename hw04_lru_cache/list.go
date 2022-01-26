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
	count     int
	backItem  *ListItem
	frontItem *ListItem
}

// Len длина списка.
func (l list) Len() int {
	return l.count
}

// Front первый элемент списка.
func (l list) Front() *ListItem {
	return l.frontItem
}

// Back последний элемент списка.
func (l list) Back() *ListItem {
	return l.backItem
}

// PushFront добавить значение в начало.
func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := &ListItem{Value: v}

	if l.count == 0 {
		// список пуст, newListItem также будет и последним элементом
		l.backItem = newListItem
	} else {
		// список не пуст, правим ссылки
		newListItem.Next = l.frontItem
		l.frontItem.Prev = newListItem
	}

	// заменаем первый элемент
	l.frontItem = newListItem
	l.count++

	return newListItem
}

// PushBack добавить значение в конец.
func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := &ListItem{Value: v}

	if l.count == 0 {
		// список пуст, newListItem также будет и первым элеметом
		l.frontItem = newListItem
	} else {
		// список не пуст, правим ссылки
		newListItem.Prev = l.backItem
		l.backItem.Next = newListItem
	}

	// заменаем первый элемент
	l.backItem = newListItem
	l.count++

	return newListItem
}

// Remove удалить элемент.
func (l *list) Remove(i *ListItem) {
	switch {
	case i.Prev == nil:
		// первый элемент
		i.Next.Prev = nil
		l.frontItem = i
	case i.Next == nil:
		// последний элемент
		i.Prev.Next = nil
		l.frontItem = i
	default:
		// не крайний элемент
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	l.count--
}

// MoveToFront переместить элемент в начало.
func (l *list) MoveToFront(i *ListItem) {
	switch {
	case i.Prev == nil:
		// первый элемент
		return
	case i.Next == nil:
		// последний элемент
		i.Prev.Next = nil
		l.backItem = i.Prev
	default:
		// не крайний элемент
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	// предыдущий первый элемент, меняем ссылки
	oldFirstItem := l.frontItem
	oldFirstItem.Prev = i
	i.Next = oldFirstItem

	// новый первый элемеент
	l.frontItem = i
}

func NewList() List {
	return new(list)
}
