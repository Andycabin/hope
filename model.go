// model
package model

type Profile struct {
	Name       string `name:"name"`
	Gender     string `name:"gender"`
	Age        int    `name:"age"`
	Height     int    `name:"height"`
	Weight     int    `name:"weight"`
	Income     string `name:"income"`
	Marriage   string `name:"marriage"`
	Education  string `name:"education"`
	Occupation string `name:"occupation"`
	Hukou      string `name:"hukou"`
}
