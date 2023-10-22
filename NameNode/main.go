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

// db conn
func Stub() *pb.PersonaServiceClient {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("No se puede conectar con el DataNode: %s\n", err)
	}

	serviceClient := pb.NewPersonaServiceClient(conn)

	return &serviceClient
}

// utils
func WriteText(data string) {
	file, err_a := os.OpenFile("NameNode/data.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err_a != nil {
		fmt.Printf("Error al abrir el archivo: %s", err_a)
		return
	}

	defer file.Close()

	_, err_w := fmt.Fprintf(file, data)
	if err_w != nil {
		fmt.Printf("Error al escribir sobre el archivo: %s", err_w)
		return
	}
}

func ReadStatusFromText(status string) ([]int, []int) {
	content, err := os.ReadFile("NameNode/data.txt")
	if err != nil {
		fmt.Printf("Error al leer el archivo de entrada: %s", err)
		return nil, nil
	}

	temp := strings.Split(string(content), "/n")

	db1 := []int{}
	db2 := []int{}

	for i := 0; i < len(temp); i++ {
		temp_2 := strings.Split(string(temp[i]), " ")

		if temp_2[1] == status {
			if temp_2[2] == "1" {
				temp_3, err_t := strconv.Atoi(temp_2[0])
				if err_t != nil {
					fmt.Printf("Error al transformar string a int: %s", err_t)
				}

				db1 = append(db1, temp_3)
			}

			if temp_2[2] == "2" {
				temp_3, err_t := strconv.Atoi(temp_2[0])
				if err_t != nil {
					fmt.Printf("Error al transformar string a int: %s", err_t)
				}

				db2 = append(db2, temp_3)
			}
		}
	}

	return db1, db2
}

func ReadLinesText() int {
	content, err := os.ReadFile("NameNode/data.txt")
	if err != nil {
		fmt.Printf("Error al leer el archivo de entrada: %s", err)
		return 0
	}

	temp := strings.Split(string(content), "/n")
	return len(temp) - 1

}

// RPCs
func (s *server) Subir(ctx context.Context, req *pb.SubirPersonaReq) (*pb.SubirPersonaRes, error) {
	persona := req.GetPersona()
	persona_estado := req.GetEstado()

	persona_id := ReadLinesText()
	persona_id = persona_id + 1

	temp := fmt.Sprint(persona_id) + " " + "1" + " " + persona_estado
	WriteText(temp)

	dn_client := *Stub()

	dn_client.Guardar(
		context.Background(),
		&pb.GuardarPersonaReq{
			Id:      int64(persona_id),
			Persona: persona,
		},
	)

	fmt.Printf("Solicitud de %s recibida, mensaje enviado: %s")
	return &pb.SubirPersonaRes{}, nil
}

func (s *server) Bajar(ctx context.Context, req *pb.BajarPersonaReq) (*pb.BajarPersonaRes, error) {
	//persona_estado := req.GetEstado()

	//db1, db2 := ReadStatusFromText(persona_estado)

	personas := []*pb.Persona{}

	for j := 0; j < 2; j++ {

		dn_client := *Stub()

		stream, err_stm := dn_client.Obtener(context.Background())
		if err_stm != nil {
			fmt.Printf("Error al abrir stream en Obtener del lado cliente: %s", err_stm)
		}

		for i := 0; i < 10; i++ { //Iterar en cada elementos de los array de ReadStatusFromText()
			if err_send := stream.Send(
				&pb.ObtenerPersonaReq{
					Id: 1, //Id del elemento x del array n de ReadStatusFromText()
				},
			); err_send != nil {
				fmt.Printf("Error al enviar en Obtener del lado cliente: %s", err_send)
			}
		}

		if err_csend := stream.CloseSend(); err_csend != nil {
			fmt.Printf("Error al enviar en Obtener del lado cliente: %s", err_csend)
		}

		for {
			res, err_res := stream.Recv()
			if err_res == io.EOF {
				break
			}
			if err_res != nil {
				fmt.Printf("Error al enviar en Obtener del lado cliente: %s", err_res)
			}

			personas = append(personas, res.GetPersona())
		}
	}

	return &pb.BajarPersonaRes{
		Persona: personas,
	}, nil
}

// main
func main() {
	listener, err_p := net.Listen("tcp", "localhost:6000")

	if err_p != nil {
		fmt.Printf("Error al conectar en el puerto 6000: %s\n", err_p)
	}

	serv := grpc.NewServer()

	pb.RegisterPersonaServiceServer(serv, &server{})

	if err_s := serv.Serve(listener); err_s != nil {
		fmt.Printf("Error al inicializar el server: %s\n", err_s)
	}
}
