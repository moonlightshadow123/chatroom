package myfile

import(
	"os"
	"fmt"
	"time"
    "runtime"
    "path/filepath"
)

var(
	_, b, _, _ = runtime.Caller(0)
    basepath,_   = filepath.Abs(filepath.Dir(os.Args[0]))
  	addr_suffix = "/static/files/"
    Statfile_path = filepath.Join(basepath, "static/files")
)

func GetPathAddr(filename string)(string, string){
	date := time.Now().Unix()
	newname := string(fmt.Sprint(date))+"_"+filename
	file_path := filepath.Join(Statfile_path, newname)
	addr := addr_suffix + newname
	return file_path, addr
}