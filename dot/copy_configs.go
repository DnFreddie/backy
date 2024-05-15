package dot

import (
	"fmt"

	"github.com/DnFreddie/backy/utils"
)
	
func CreateSymlinksTemp(dotfiels []Dotfile){

	c_dir,err:= utils.ScanDir(".config")

//TODOOOOOO!
	// for  _, dot := range dotfiels{

	// }



	if err != nil {
		return 
	}
fmt.Println(c_dir)
}
