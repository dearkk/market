package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
)

func main() {
	fmt.Print("123\n")
	str := make(map[string] interface{})
	f, err := os.OpenFile("D:\\bean\\master\\conf.yaml", os.O_RDONLY,0600)
	defer f.Close()
	if err !=nil {
		fmt.Println(err.Error())
	} else {
		contentByte,_ :=ioutil.ReadAll(f)
		yaml.Unmarshal(contentByte, str)
		//fmt.Printf("str: %+v", str)
		for _, v := range str {
			t := reflect.TypeOf(v).Kind().String()
			fmt.Printf("value type: %+v\n", t)
			if value, ok := v.([]interface{}); ok {
				for _, vv := range value {
					fmt.Printf("value type2 : %+v\n", vv)
					fmt.Printf("value type2: %+v\n", reflect.TypeOf(vv).Kind().String())
					if reflect.TypeOf(vv).Kind().String() == "map" {
						for k, vvv := range vv.(map[interface{}]interface{}) {
							fmt.Printf("value type3 : k: %s, %v\n", k, vvv)
						}
					}
				}
			}
		}
	}
}