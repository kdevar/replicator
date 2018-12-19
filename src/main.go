package main



func main(){

	config := NewConfig()
	streamer := NewStreamer(config)
	handler := NewHandler(streamer)
	server := NewIngestServer(config, handler)

	server.Run()
}