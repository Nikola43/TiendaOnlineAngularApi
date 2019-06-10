package main

func main() {
	a := App{}
	a.Initialize("root", "-Gdeuapmw18", "angular_shop?parseTime=true")
	a.Run(":8080")
}
