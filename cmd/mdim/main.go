package main

func main() {
	mdimCmd.AddCommand(qiniuCmd)
	mdimCmd.Execute()
}
