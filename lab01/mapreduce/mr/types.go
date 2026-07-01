package mr

// KeyValue is a single key-value pair.
type KeyValue struct {
	Key   string
	Value string
}

// ByKey attaches the methods of sort.Interface to []KeyValue, so we can call
// sort.Sort(ByKey(pairs)) to order pairs alphabetically by Key.
type ByKey []KeyValue

// Len returns the length of the given data.
func (b ByKey) Len() int { return len(b) }

// Swap exchanges the values of two given pairs.
func (b ByKey) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// Less evaluates if ith-key is less than jth.
func (b ByKey) Less(i, j int) bool { return b[i].Key < b[j].Key }
