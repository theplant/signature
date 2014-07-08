package signature

import (
	"encoding/gob"
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

func (this *SerializableItem) ToMap(t reflect.Type) interface{} {
	if len(this.Keys) <= 0 {
		return nil
	}
	oMap := reflect.MakeMap(t)
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
				ki := k.Interface()
				vi := v.MapIndex(k).Interface()
				kis = append(kis, ki)
				vis = append(vis, vi)

				// Register Gob so that it can encoding customized struct
				kType := reflect.TypeOf(ki)
				kDefaultValue := reflect.New(kType).Elem().Interface()
				gob.Register(kDefaultValue)
				vType := reflect.TypeOf(vi)
				vDefaultValue := reflect.New(vType).Elem().Interface()
				gob.Register(vDefaultValue)
			}
			r = &SerializableItem{kis, vis}
		}
	}
	return
}
