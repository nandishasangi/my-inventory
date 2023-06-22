package main

func main() {
	app := App{}
	app.Initialise("139.59.16.135", 5000, "shipx", "bijapur", "inventory")
	app.Run("localhost:10000")
}
