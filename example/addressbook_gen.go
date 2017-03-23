// Generated by github.com/davyxu/gosproto/sprotogen
// DO NOT EDIT!

package example

import (
	"reflect"
	"github.com/davyxu/gosproto"
	"github.com/davyxu/goobjfmt"
	"github.com/davyxu/cellnet/codec/sproto"
)

type MyCar int32

const (
	MyCar_Monkey MyCar = 1
	MyCar_Monk   MyCar = 2
	MyCar_Pig    MyCar = 3
)

var MyCar_ValueByName = map[string]int32{
	"Monkey": 1,
	"Monk":   2,
	"Pig":    3,
}

var MyCar_NameByValue = map[int32]string{
	1: "Monkey",
	2: "Monk",
	3: "Pig",
}

func (self MyCar) String() string {
	return sproto.EnumName(MyCar_NameByValue, int32(self))
}

type PhoneNumber struct {
	Number string `sproto:"string,0,name=Number"`

	Type int32 `sproto:"integer,1,name=Type"`
}

func (self *PhoneNumber) String() string { return goobjfmt.CompactTextString(self) }

type Person struct {
	Name string `sproto:"string,0,name=Name"`

	Id int32 `sproto:"integer,1,name=Id"`

	Email string `sproto:"string,2,name=Email"`

	Phone []*PhoneNumber `sproto:"struct,3,array,name=Phone"`
}

func (self *Person) String() string { return goobjfmt.CompactTextString(self) }

type AddressBook struct {
	Person []*Person `sproto:"struct,0,array,name=Person"`
}

func (self *AddressBook) String() string { return goobjfmt.CompactTextString(self) }

type MyData struct {
	Name string `sproto:"string,0,name=Name"`

	Type MyCar `sproto:"integer,1,name=Type"`

	Int32 int32 `sproto:"integer,2,name=Int32"`

	Uint32 uint32 `sproto:"integer,3,name=Uint32"`

	Int64 int64 `sproto:"integer,4,name=Int64"`

	Uint64 uint64 `sproto:"integer,5,name=Uint64"`

	Bool bool `sproto:"boolean,6,name=Bool"`
}

func (self *MyData) String() string { return goobjfmt.CompactTextString(self) }

type MyProfile struct {
	NameField *MyData `sproto:"struct,0,name=NameField"`

	NameArray []*MyData `sproto:"struct,1,array,name=NameArray"`

	NameMap []*MyData `sproto:"struct,2,array,name=NameMap"`
}

func (self *MyProfile) String() string { return goobjfmt.CompactTextString(self) }

var SProtoStructs = []reflect.Type{

	reflect.TypeOf((*PhoneNumber)(nil)).Elem(), // 4271979557
	reflect.TypeOf((*Person)(nil)).Elem(),      // 1498745430
	reflect.TypeOf((*AddressBook)(nil)).Elem(), // 2618161298
	reflect.TypeOf((*MyData)(nil)).Elem(),      // 2244887298
	reflect.TypeOf((*MyProfile)(nil)).Elem(),   // 438153711
}

func init() {
	sprotocodec.AutoRegisterMessageMeta(SProtoStructs)
}
