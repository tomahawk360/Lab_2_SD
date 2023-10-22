package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	pb "github.com/tomahawk360/lab2sd/proto"
)

type server struct {
	pb.UnimplementedPersonaServiceServer
}

// utils
func WriteText(data string) {
	file, err_a := os.OpenFile("DateNode/data.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err_a != nil {
		fmt.Printf("Error al abrir el archivo: %s/n", err_a)
		return
	}

	defer file.Close()

	_, err_w := fmt.Fprintln(file, data)
	if err_w != nil {
		fmt.Printf("Error al escribir sobre el archivo: %s/n", err_w)
		return
	}
}

func ReadLineText(id int) (string, string) {
	content, err := os.ReadFile("DateNode/data.txt")
	if err != nil {
		fmt.Printf("Error al leer el archivo de entrada: %s/n", err)
		return "", ""
	}

	temp := strings.Split(string(content), "/n")

	for i := 0; i < len(temp); i++ {
		temp_2 := strings.Split(string(temp[i]), " ")

		if temp_2[0] == fmt.Sprint(id) {
			return temp_2[1], temp_2[2]
		}
	}

	return "", ""
}

// RPCs
func (s *server) Guardar(ctx context.Context, req *pb.GuardarPersonaReq) (*pb.GuardarPersonaRes, error) {
	persona_id := req.GetId()
	persona_nombre := req.GetPersona().GetNombre()
	persona_apellido := req.GetPersona().GetApellido()

	temp := strconv.FormatInt(persona_id, 10) + " " + persona_nombre + " " + persona_apellido
	WriteText(temp)

	return &pb.GuardarPersonaRes{}, nil
}

func (s *server) Obtener(srv pb.PersonaService_ObtenerServer) error {
	for {
		req, err_r := srv.Recv()
		if err_r == io.EOF {
			return nil
		}
		if err_r != nil {
			fmt.Printf("Error al recibir en Obtener del lado servidor: %s", err_r)
			return nil
		}

		persona_id := req.GetId()

		persona_nombre, persona_apellido := ReadLineText(int(persona_id))

		if err_s := srv.Send(&pb.ObtenerPersonaRes{
			Persona: &pb.Persona{
				Nombre:   persona_nombre,
				Apellido: persona_apellido,
			},
		}); err_s != nil {
			fmt.Printf("Error al enviar en Obtener del lado servidor: %s", err_s)
			return nil
		}
	}
}

// Main
func main() {
	listener, err_p := net.Listen("tcp", "localhost:5000")

	if err_p != nil {
		fmt.Printf("Error al conectar en el puerto 5000: %s\n", err_p)
	}

	serv := grpc.NewServer()

	pb.RegisterPersonaServiceServer(serv, &server{})

	if err_s := serv.Serve(listener); err_s != nil {
		fmt.Printf("Error al inicializar el server: %s\n", err_s)
	}
}
