package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/tomahawk360/lab2sd/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:6000", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("No se puede conectar con el NameNode: %s\n", err)
	}

	serviceClient := pb.NewPersonaServiceClient(conn)

	/* 	person := pb.Persona{
	   		Nombre:   "Antonia",
	   		Apellido: "Ortiz",
	   	}

	   	serviceClient.Subir(
	   		context.Background(),
	   		&pb.SubirPersonaReq{
	   			Persona: &person,
	   			Estado:  true,
	   		},
	   	) */

	personas, err_b := serviceClient.Bajar(
		context.Background(),
		&pb.BajarPersonaReq{
			Estado: false,
		},
	)

	if err_b != nil {
		fmt.Printf("Error al Bajar en el lado de ONU: %s", err_b)
	}

	for i := 0; i < len(personas.GetPersona()); i++ {
		person := personas.GetPersona()[i]

		fmt.Printf("%s %s\n", person.GetNombre(), person.GetApellido())
	}
}
