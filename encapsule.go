package signature

import (
	"reflect"
	"sort"
)

type SortableKeys []reflect.Value
type ByString struct{ SortableKeys }

func (this SortableKeys) Len() int {
	return len(this)
}
func (this SortableKeys) Swap(i int, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this ByString) Less(i int, j int) bool {
	return this.SortableKeys[i].String() < this.SortableKeys[j].String()
}

type SerializableItem struct {
	Keys   []interface{}
	Values []interface{}
}

func (this *SerializableItem) ToMap() interface{} {
	if len(this.Keys) <= 0 {
		return nil
	}
	mapType := reflect.MapOf(reflect.ValueOf(this.Keys[0]).Type(), reflect.ValueOf(this.Values[0]).Type())
	oMap := reflect.MakeMap(mapType)
	for i, k := range this.Keys {
		oMap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(this.Values[i]))
	}
	return oMap.Interface()
}

func InitWithMap(e interface{}) (r *SerializableItem) {
	v := reflect.Indirect(reflect.ValueOf(e))
	length := len(v.MapKeys())
	if length > 0 {
		// sort the keys
		ks := []reflect.Value{}
		kis := []interface{}{}
		vis := []interface{}{}
		for _, k := range v.MapKeys() {
			ks = append(ks, k)
		}
		sort.Sort(ByString{ks})

		// encapsule the maps to the struct with certain order
		if sort.IsSorted(ByString{ks}) {
			for _, k := range ks {
				kis = append(kis, k.Interface())
				vis = append(vis, v.MapIndex(k).Interface())
			}
			r = &SerializableItem{kis, vis}
		}
	}
	return
}
