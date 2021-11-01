package main

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"protobuf-go/pb/first"
	"protobuf-go/pb/second"
	"protobuf-go/pb/third"
)

func main() {

	// pb -> 二进制file
	//pm := NewPersonMessage()
	//err := WriteToFile("person.bin", pm)
	//if err != nil {
	//	return
	//}

	// 二进制file -> pb
	//pm2 := &first.PersonMessage{}
	//err := ReadFromFile("person.bin", pm2)
	//if err != nil {
	//	return
	//}
	//fmt.Printf("%#v\n", pm2)

	// pb -> json
	//json := toJSON(pm2)
	//fmt.Println(json) //{"id":1234, "isAdult":true, "name":"Leo", "luckyNumbers":[1, 2, 3, 4]}

	// json -> pb
	//pb3 := &first.PersonMessage{}
	//_ = fromJSON(json, pb3)
	//fmt.Printf("%#v\n", pb3)

	//em:=NewEnumMessage()
	//fmt.Println(em.GetGender())
	p := NewComplexMessage()
	fmt.Println(p.GetHobbies()[1].GetName())
}

func NewComplexMessage() *third.Person {
	p := &third.Person{
		Id: 1,
		Hobbies: []*third.Hobby{
			{Name: "basketball"},
			{Name: "running"},
		},
	}
	return p
}

func NewEnumMessage() *second.EnumMessage {
	em := &second.EnumMessage{
		Id:     345,
		Gender: second.Gender_FEMALE,
	}
	return em
}

func fromJSON(in string, pb proto.Message) error {
	err := protojson.Unmarshal([]byte(in), pb)
	return err
}

func toJSON(pb proto.Message) string {
	marshaller := protojson.MarshalOptions{
		Indent: "	",
	}
	dataBytes, _ := marshaller.Marshal(pb)
	return string(dataBytes)
}

func ReadFromFile(fileName string, pb proto.Message) error {
	dataBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("读取文件失败", err.Error())
	}
	err = proto.Unmarshal(dataBytes, pb)
	if err != nil {
		log.Fatalln("反序列化失败", err.Error())
	}
	return err
}

func WriteToFile(fileName string, pb proto.Message) error {
	dataBytes, err := proto.Marshal(pb)
	if err != nil {
		log.Fatalln("无法序列化")
	}
	err = ioutil.WriteFile(fileName, dataBytes, 0644)
	return err
}
func NewPersonMessage() *first.PersonMessage {
	pm := first.PersonMessage{
		Id:           1234,
		IsAdult:      true,
		Name:         "Leo",
		LuckyNumbers: []int32{1, 2, 3, 4},
	}
	//fmt.Println(pm)
	return &pm
}
