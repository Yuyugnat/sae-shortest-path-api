package data

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	c "sae-shortest-path/connection"
	bug "sae-shortest-path/debugging"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fastjson"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var p fastjson.Parser

type Neighbor struct {
	Gid    int     `json:"gid"`
	Length float64 `json:"length"`
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
}

type NeighborTable struct {
	debugger   *bug.Debug
	Table      map[int][]Neighbor
	dataLoaded bool
	dbConn     *c.PostgresConn
}

var instance *NeighborTable

func Load() {
	GetInstance()
}

func GetInstance() *NeighborTable {
	if instance == nil {
		conn, err := c.GetInstance()
		if err != nil {
			panic(err)
		}
		instance = &NeighborTable{
			Table:      make(map[int][]Neighbor),
			dataLoaded: false,
			dbConn:     conn,
		}
		instance.debugger = bug.NewDebug()
		go func() {
			instance.loadData()
		}()
	}
	return instance
}

func (f *NeighborTable) loadData() {
	query := `
		SELECT gid, voisins
		FROM voisins_jsonb
	`
	var rows *sql.Rows
	var err error

	f.Debug().GetTimeUsing("query (load data)", func() {
		rows, err = f.dbConn.DB.Query(query)
	})

	if err != nil {
		fmt.Println("Data not loaded (query) : ", err)
		f.dataLoaded = false
		return
	}

	var gid int
	var voisins string

	for rows.Next() {
		f.Debug().GetTimeUsing("scan (load data)", func() {
			err = rows.Scan(&gid, &voisins)
		})
		if err != nil {
			fmt.Println("Data not loaded (scan) : ", err)
			f.dataLoaded = false
			return
		}
		var neighbours []Neighbor
		f.Debug().GetTimeUsing("unmarshal (load data)", func() {
			err = json.Unmarshal([]byte(voisins), &neighbours)
		})

		if err != nil {
			fmt.Println("Data not loaded (unmarshal) : ", err)
			f.dataLoaded = false
			return
		}
		f.Debug().GetTimeUsing("insert in map (load data)", func() {
			f.Table[gid] = neighbours
		})
	}

	fmt.Println("Data loaded !")
	f.dataLoaded = true
	f.Debug().GetTimeUsing("get real size (load data)", func() {
		res, _ := getRealSizeOf(f.Table)
		fmt.Printf("Size of unsafed data %.3f MB\n", float64(res)/float64(1024*1024))
	})

	f.Debug().Print()
}

// func (f *NeighborTable) loadData2() {
// 	var err error
// 	var file []byte
// 	f.Debug().GetTimeUsing("read file (load data)", func() {
// 		file, err = os.ReadFile("data.json")
// 	})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	f.Debug().GetTimeUsing("unmarshal (load data)", func() {
// 		err = json.Unmarshal(file, &f.Table)
// 	})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	f.Debug().Print()
// 	f.dataLoaded = true

// 	// fmt.Println("Data loaded :", f.Table)
// 	res, _ := getRealSizeOf(f.Table)
// 	fmt.Printf("Size of data %f MB\n", float64(res)/float64(1024*1024))

// }

func getRealSizeOf(v interface{}) (int, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return 0, err
	}
	return b.Len(), nil
}

func (f *NeighborTable) Debug() *bug.Debug {
	return f.debugger
}

func (f *NeighborTable) Get(gid int) []Neighbor {
	return f.Table[gid]
}

func (f *NeighborTable) Ready() bool {
	return f.dataLoaded
}
