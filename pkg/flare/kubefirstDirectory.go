package flare

import (
	"log"
	"fmt"
	"os"
	"errors"
	)


//Verify the state of the kubefirst directory
//
// Output:
//   $PATH/.kubefirst
func CheckKubefirstDir(home string) string {
	k1sDir :=  fmt.Sprintf("%s/.kubefirst", home)
	if _, err := os.Stat(k1sDir); err == nil {
		// path/to/whatever exists
		log.Printf("\".kubefirst\" file found: %s", k1sDir)	
		log.Printf("	\".kubefirst\" will be generated by installation process, if exist means a installation may already be executed" )		  
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		log.Printf("\".kubefirst\" file not found: %s", k1sDir)				  
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		log.Printf("Unable to check is \".kubefirst\" if file exists" )		  
	}
	return k1sDir
}