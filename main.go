package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	consts "rimor/pkg/consts"
	"rimor/pkg/engine/dictionary/xindex"
	index "rimor/pkg/engine/dictionary/xindex"
	"rimor/pkg/web/routes"
)

func main() {
	// data_present := flag.Bool("present", false, "this is mainly used for restoring of index when index has previously been calculated")
	// store := flag.Bool("storedata", false, "when this flag is set, index data will be stored in file")
	flag.Parse()



	router := routes.NewRouter()
	if err := router.Router.Listen(router.Port); err != nil {
		log.Fatal(err.Error())
	}
	

	// if *store && !(*data_present) {
	// 	fmt.Print("exporting\n\n\n")
	// 	ExportIndexToFile(indx)
	// }

	

	// for _, p := range indx.Records {
	// 	fmt.Print(p.GetTerm())
	// 	current := p.GetPostingList()
	// 	for current != nil {
	// 		fmt.Printf("%d ", current.GetDocID())
	// 		current = current.GetNextElem()
	// 	}
	// 	fmt.Print("\n")
	// }


}


func ExportIndexToFile(index *index.Xindex) error{
	f, err := os.Create(consts.INDEX_FILE_PATH)
	if err != nil {
		log.Fatalln(err.Error())
	}
	enc := json.NewEncoder(f)
	if err := enc.Encode(index.Records); err != nil {
		log.Fatalln(err.Error())
	}
	return nil
}

func ImportIndexToFile() *xindex.Xindex{
	f, err := os.Open(consts.INDEX_FILE_PATH)
	if err != nil {
		log.Fatalf("failed to restore index, err : %s", err.Error())
	}
	dec := json.NewDecoder(f)
	indx := xindex.Xindex{}
	if err := dec.Decode(&indx); err != nil {
		log.Fatalf("failed to decode index data into object, err : %s", err.Error())
	}
	return &indx	
}
