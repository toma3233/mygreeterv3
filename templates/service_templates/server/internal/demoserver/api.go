package demoserver

import (
	pb "<<apiModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>/v1"
)

type Server struct {
	// When the Unimplemented<<.serviceInput.serviceName>>Server struct is embedded,
	// the generated method/implementation in .pb file will be associated with this struct.
	// If this struct doesn't implment some methods,
	// the .pb ones will be used. If this struct implement the methods, it will override the .pb ones.
	// The reason is that anonymous field's methods are promoted to the struct.
	//
	// When this struct is NOT embedded,, all methods have to be implemented to meet the interface requirement.
	// See https://go.dev/ref/spec#Struct_types.
	pb.Unimplemented<<.serviceInput.serviceName>>Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) init(options Options) {	
}
