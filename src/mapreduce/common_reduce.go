package mapreduce

import (
	"encoding/json"
	"os"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// TODO:
	// You will need to write this function.
	// You can find the intermediate file for this reduce task from map task number
	// m using reduceName(jobName, m, reduceTaskNumber).
	// Remember that you've encoded the values in the intermediate files, so you
	// will need to decode them. If you chose to use JSON, you can read out
	// multiple decoded values by creating a decoder, and then repeatedly calling
	// .Decode() on it until Decode() returns an error.
	//
	// You should write the reduced output in as JSON encoded KeyValue
	// objects to a file named mergeName(jobName, reduceTaskNumber). We require
	// you to use JSON here because that is what the merger than combines the
	// output from all the reduce tasks expects. There is nothing "special" about
	// JSON -- it is just the marshalling format we chose to use. It will look
	// something like this:
	//
	// enc := json.NewEncoder(mergeFile)
	// for key in ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()

	// Step 1: Read intermediate file from all map task
	kvlist := make(map[string][]string)
	for i := 0; i < nMap; i++ {
		name := reduceName(jobName, i, reduceTaskNumber)
		kvs := readFromFile(name)
		for _, kv := range kvs {
			if _, ok := kvlist[kv.Key]; !ok {
				kvlist[kv.Key] = []string{kv.Value}
			} else {
				kvlist[kv.Key] = append(kvlist[kv.Key], kv.Value)
			}
		}
	}

	// Step 2: Call reduceF and write to mergeFile
	mergeFilename := mergeName(jobName, reduceTaskNumber)
	outFile, err := os.Create(mergeFilename)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	enc := json.NewEncoder(outFile)
	for k, vlist := range kvlist {
		enc.Encode(KeyValue{k, reduceF(k, vlist)})
	}

}

func readFromFile(filename string) (kvs []KeyValue) {
	inFile, err := os.Open(filename)
	if err != nil {
		return
	}
	defer inFile.Close()

	dec := json.NewDecoder(inFile)
	var kv KeyValue
	for err := dec.Decode(&kv); err == nil; err = dec.Decode(&kv) {
		kvs = append(kvs, kv)
	}

	return
}
