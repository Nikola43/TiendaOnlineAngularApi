package main

func main() {
	a := App{}
	a.Initialize("root", "", "angular_shop?parseTime=true")
	a.Run(":8080")
}
