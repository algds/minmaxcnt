package minmaxcnt

import "container/list"

// Interface represents all permitted operations
// that run in O(1) time: constant time operations.
type Interface interface {

	// Increment increases the count associated with the key by 1.
	// If the key doesn't exist, it is created with count of 1.
	Increment(key interface{})

	// Decrement decreases the count associated with the key by 1.
	// If the key doesn't exist, the result is no operation.
	// If as a result, the count drops to zero, it is removed.
	// Counts can never become negative hence the 'uint'.
	Decrement(key interface{})

	// Count returns the current count of the given key.
	// If the key does not exist, zero is returned.
	Count(key interface{}) uint

	// Max deterministically returns the key with the highest count.
	// Of the keys sharing the same maximum count, the oldest of them is returned.
	Max() (interface{}, uint)

	// Min deterministically returns the key with the lowest non-zero count.
	// Of the keys sharing the same minimum count, the newest of them is returned.
	Min() (interface{}, uint)
}

type data struct {
	entryList         *list.List
	lookupEntryByKey  map[interface{}]*list.Element
	lookupListByCount map[uint]*list.Element
}

type entry struct {
	count       uint
	values      *list.List
	lookupByKey map[interface{}]*list.Element
}

func (d *data) Increment(key interface{}) {
	eSource := d.lookupEntryByKey[key]

	if eSource == nil {
		vals := list.New()
		lookup := make(map[interface{}]*list.Element)
		lookup[key] = vals.PushFront(key)

		newEntry := &entry{
			count:       0,
			values:      vals,
			lookupByKey: lookup,
		}
		eSource = d.entryList.PushFront(newEntry)
	}

	sourceEntry := eSource.Value.(*entry)

	eTarget := d.lookupListByCount[sourceEntry.count+1]

	if eTarget == nil {
		newEntry := &entry{
			count:       sourceEntry.count + 1,
			values:      list.New(),
			lookupByKey: make(map[interface{}]*list.Element),
		}
		eTarget = d.entryList.InsertAfter(newEntry, eSource)
	}

	targetEntry := eTarget.Value.(*entry)

	d.lookupEntryByKey[key] = eTarget
	d.lookupListByCount[targetEntry.count] = eTarget
	targetEntry.lookupByKey[key] = targetEntry.values.PushFront(key)

	sourceEntryKey := sourceEntry.lookupByKey[key]
	sourceEntry.values.Remove(sourceEntryKey)
	delete(sourceEntry.lookupByKey, key)

	if len(sourceEntry.lookupByKey) == 0 {
		delete(d.lookupListByCount, sourceEntry.count)
		d.entryList.Remove(eSource)
	}
}

func (d *data) Decrement(key interface{}) {
	eSource := d.lookupEntryByKey[key]

	if eSource == nil {
		return
	}

	sourceEntry := eSource.Value.(*entry)

	eTarget := d.lookupListByCount[sourceEntry.count-1]

	if eTarget == nil {
		newEntry := &entry{
			count:       sourceEntry.count - 1,
			values:      list.New(),
			lookupByKey: make(map[interface{}]*list.Element),
		}
		eTarget = d.entryList.InsertBefore(newEntry, eSource)
	}

	targetEntry := eTarget.Value.(*entry)

	d.lookupEntryByKey[key] = eTarget
	d.lookupListByCount[targetEntry.count] = eTarget
	targetEntry.lookupByKey[key] = targetEntry.values.PushFront(key)

	sourceEntryKey := sourceEntry.lookupByKey[key]
	sourceEntry.values.Remove(sourceEntryKey)
	delete(sourceEntry.lookupByKey, key)

	if len(sourceEntry.lookupByKey) == 0 {
		delete(d.lookupListByCount, sourceEntry.count)
		d.entryList.Remove(eSource)
	}

	if targetEntry.count == 0 {
		delete(d.lookupEntryByKey, key)
		delete(d.lookupListByCount, targetEntry.count)
		d.entryList.Remove(eTarget)
	}
}

func (d *data) Count(key interface{}) uint {
	if e, found := d.lookupEntryByKey[key]; found {
		return e.Value.(*entry).count
	}

	return 0
}

func (d *data) Max() (interface{}, uint) {
	if d.entryList.Len() == 0 {
		return nil, 0
	}

	back := d.entryList.Back().Value.(*entry)
	val := back.values.Back().Value

	return val, back.count
}

func (d *data) Min() (interface{}, uint) {
	if d.entryList.Len() == 0 {
		return nil, 0
	}

	front := d.entryList.Front().Value.(*entry)
	val := front.values.Front().Value

	return val, front.count
}

// New returns a new instance of a BigOhOne data structure.
func New() Interface {
	return &data{
		entryList:         list.New(),
		lookupEntryByKey:  make(map[interface{}]*list.Element),
		lookupListByCount: make(map[uint]*list.Element),
	}
}
