package icache

import (
	"fmt"
	"reflect"
)

// SinkIf sink interface
type SinkIf interface {
	SetView(v View) error
	GetView() (View, error)
	SetBytes([]byte) error
	SetString(string) error
	SetObj(interface{}) error
	SetTTL(int32) error
}

////////////////////////////////////////////////////////
// stringSink

// StringSink new stringSink
func StringSink(sp *string) SinkIf {
	return &stringSink{sp: sp}
}

type stringSink struct {
	sp  *string
	ttl int32
}

// SetView set view
func (ss *stringSink) SetView(v View) error {
	switch v.v.(type) {
	case string:
		*ss.sp = v.v.(string)
	case []byte:
		*ss.sp = string(v.v.([]byte))
	default:
		return fmt.Errorf("not string view")
	}
	ss.ttl = v.ttl
	return nil
}

// GetView get view
func (ss *stringSink) GetView() (View, error) {
	v := View{
		v:   *ss.sp,
		ttl: ss.ttl,
	}
	return v, nil
}

// SetString set string
func (ss *stringSink) SetString(s string) error {
	*ss.sp = s
	return nil
}

// SetBytes set bytes
func (ss *stringSink) SetBytes(b []byte) error {
	*ss.sp = string(b)
	return nil
}

// SetObj set obj
func (ss *stringSink) SetObj(obj interface{}) error {
	return fmt.Errorf("not support interface obj")
}

// SetTTL set ttl
func (ss *stringSink) SetTTL(iTTL int32) error {
	ss.ttl = iTTL
	return nil
}

////////////////////////////////////////////////////////
// byteSink

// ByteSink new byteSink
func ByteSink(bp *[]byte) SinkIf {
	return &byteSink{bp: bp}
}

type byteSink struct {
	bp  *[]byte
	ttl int32
}

// SetView set view
func (bs *byteSink) SetView(v View) error {
	switch v.v.(type) {
	case string:
		*bs.bp = []byte(v.v.(string))
	case []byte:
		*bs.bp = cloneBytes(v.v.([]byte))
	default:
		return fmt.Errorf("not byte view")
	}
	bs.ttl = v.ttl
	return nil
}

// GetView get view
func (bs *byteSink) GetView() (View, error) {
	v := View{
		v:   *bs.bp,
		ttl: bs.ttl,
	}
	return v, nil
}

// SetBytes set bytes
func (bs *byteSink) SetBytes(b []byte) error {
	*bs.bp = cloneBytes(b)
	return nil
}

// SetString set string
func (bs *byteSink) SetString(s string) error {
	*bs.bp = []byte(s)
	return nil
}

// SetObj set obj
func (bs *byteSink) SetObj(obj interface{}) error {
	return fmt.Errorf("not support interface obj")
}

// SetTTL set ttl
func (bs *byteSink) SetTTL(iTTL int32) error {
	bs.ttl = iTTL
	return nil
}

////////////////////////////////////////////////////////
// objSink

// ObjSink new obj sink
func ObjSink(pObj interface{}) SinkIf {
	return &objSink{obj: pObj}
}

type objSink struct {
	obj interface{}
	ttl int32
}

// SetView set view
func (os *objSink) SetView(inView View) error {
	vType := reflect.TypeOf(inView.v)
	if vType.Kind() != reflect.Ptr {
		return fmt.Errorf("inView not ptr type")
	}
	objValue := reflect.ValueOf(os.obj)
	if objValue.Kind() != reflect.Ptr {
		return fmt.Errorf("obj not ptr type")
	}
	objValue.Elem().Set(reflect.ValueOf(inView.v).Elem())
	os.ttl = inView.ttl
	return nil
}

// GetView get view
func (os *objSink) GetView() (View, error) {
	v := View{
		v:   os.obj,
		ttl: os.ttl,
	}
	return v, nil
}

// SetBytes set bytes
func (os *objSink) SetBytes(b []byte) error {
	return fmt.Errorf("not support bytes")
}

// SetString set string
func (os *objSink) SetString(s string) error {
	return fmt.Errorf("not support string")
}

// SetObj set obj
func (os *objSink) SetObj(inObj interface{}) error {
	vType := reflect.TypeOf(inObj)
	if vType.Kind() != reflect.Ptr {
		return fmt.Errorf("inObj not ptr type")
	}
	objValue := reflect.ValueOf(os.obj)
	if objValue.Kind() != reflect.Ptr {
		return fmt.Errorf("obj not ptr type")
	}
	objValue.Elem().Set(reflect.ValueOf(inObj).Elem())
	return nil
}

// SetTTL set ttl
func (os *objSink) SetTTL(iTTL int32) error {
	os.ttl = iTTL
	return nil
}
