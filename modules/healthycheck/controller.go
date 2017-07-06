package healthycheck

import (
	iris "gopkg.in/kataras/iris.v4"
	"os/exec"
	"fmt"
)

func Checking (ctx *iris.Context){
	type Status struct{
		Postgre string `json:"postgre"`
	}

	// m := []Status{}

	cmd := exec.Command("pg_ctl","-D","/usr/local/var/postgres","status")
	fmt.Println(cmd)

	err := cmd.Start()
	fmt.Println(err)

	// if err != nil{
	// 	fmt.Println("idup")
	// }else{
	// 	// m.Postgre := "Postgre mati"
	// 	fmt.Println("mati")
	// }
	// ctx.JSON(iris.StatusOK, iris.Map{
	// 	"status": "success",
	// 	"data":   m,
	// })

}