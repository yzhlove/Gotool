package main

import (
	"fmt"
	"os"
)

//go:generate msgp -tests=false -io=false

type Ext struct {
	Count int
}

type A struct {
	Att [2]Ext
}

func main() {
	//writeFile()
	//readFile()

	a := A{
		Att: [2]Ext{
			Ext{
				Count: 1,
				//Times: 10,
			},
			Ext{
				Count: 2,
				//Times: 20,
			},
		},
	}

	for i, k := range a.Att {
		if i == 1 {
			a.Att[i].Count += 10
		}
		fmt.Println(k.Count, a.Att[i].Count)
	}

}

func writeFile() {
	a := A{
		Att: [2]Ext{
			Ext{
				Count: 1,
				//Times: 10,
			},
			Ext{
				Count: 2,
				//Times: 20,
			},
		},
	}

	data, err := a.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile("data.bin", data, 0644); err != nil {
		panic(err)
	}
}

func readFile() {
	data, err := os.ReadFile("data.bin")
	if err != nil {
		panic(err)
	}

	var a A
	_, err = a.UnmarshalMsg(data)
	if err != nil {
		panic(err)
	}

	println(a.Att[0].Count)
	println(a.Att[1].Count)
}
