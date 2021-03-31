package consistenthash

import (
	"log"
	"strconv"
	"testing"
)

func TestMap_Add(t *testing.T) {
	m :=  New(9,nil)
	m.Add("emt")
	m.Add("suosi")
	m.Add("pandora")
	//println(len(m.hashMap))
}

func TestMap_Get(t *testing.T) {
	m :=  New(9,nil)
	m.Add("emt")
	m.Add("suosi")
	m.Add("pandora")
	//
	//v := m.Get("emt")
	//hash := int(m.hash([]byte("emt")))
	//fmt.Println(hash)
	//fmt.Println(v)
	//
	//for v1,v2:= range m.hashMap{
	//	fmt.Println(v1,"  ",v2)
	//}
}

func TestHash(t *testing.T){
	m := New(3, func(data []byte) uint32 {
		hash,_ := strconv.Atoi(string(data))
		return uint32(hash)
	})


	//2  4  6
	m.Add("2","4","6")
	testcases := map[string]string{
		"2":"2",
		"11":"2",
		"23":"4",
		"27":"2",
	}


	for k,v := range testcases{
		if v != m.Get(k){

			log.Fatal("error!"," ",k," ",m.Get(k))
		}
	}

	m.Add("8")
	if m.Get("27") != "8"{
		log.Fatal("error !" + "8")
	}


}
