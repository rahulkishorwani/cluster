package main
import (
	//"fmt"
	"github.com/rahulkishorwani/cluster/myluster"
	"sync"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)
func main() {
   
	numberofservers:=0
	numberofservers,_=strconv.Atoi(os.Args[1])
	wg := new(sync.WaitGroup)
    	wg.Add(numberofservers)

	for i:=0;i<numberofservers;i++ {
		id:=strconv.Itoa(i)
		go cluster.serverfunc(id,20000,numberofservers,wg)	
	}
	
	wg.Wait()


	binary, lookErr := exec.LookPath("bash")
    	if lookErr != nil {
    	    panic(lookErr)
    	}
	

	args := []string{"bash", "mybash.sh"}
	env := os.Environ()
	
	
	execErr := syscall.Exec(binary, args, env)
    	if execErr != nil {
        	panic(execErr)
    	}

	
	
}
