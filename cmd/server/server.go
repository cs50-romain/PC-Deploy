/*
This package will be used as the server. It will listen for incoming connections from the client. It will receive logs and save them. We can send client coming from the controlcenter to the client
*/
package server

type Server struct {
	addr string
	port int
}
