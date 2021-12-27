package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (lC *lruCache) Set(key Key, value interface{}) bool {
	// формируем элемент кеша
	cItem := &cacheItem{key, value}

	// пытаемся получить элемент из словаря по ключу
	mapItem, ok := lC.items[key]

	if ok {
		// элемент присутствует в словаре, меняем значение и перемещаем в начало очереди
		mapItem.Value = cItem
		lC.queue.MoveToFront(mapItem)
	} else {
		// элемент отсутствует
		// проверка переполнения кеша
		if lC.queue.Len() >= lC.capacity {
			// не хватает места для нового элемента
			// получаем последний элемент очереди
			lastItem := lC.queue.Back()
			// удаляем его из очереди
			lC.queue.Remove(lastItem)
			// и из словаря
			delete(lC.items, lastItem.Value.(*cacheItem).key)
		}

		// добавляем в очередь и в словарь
		lC.items[key] = lC.queue.PushFront(cItem)
	}
	return ok
}

func (lC *lruCache) Get(key Key) (interface{}, bool) {
	// пытаемся получить элемент из словаря по ключу
	if mapItem, ok := lC.items[key]; ok {
		// элемент присутствует в словаре, перемещаем в начало очереди
		lC.queue.MoveToFront(mapItem)
		// и возвращаем
		return mapItem.Value.(*cacheItem).value, true
	}

	// элемент отсутствует
	return nil, false
}

func (lC *lruCache) Clear() {
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
