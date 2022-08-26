package example

import (
	"context"
	"errors"
	"fmt"
	"go-micro.dev/v4"
	"skframe/app/protobuffs/Student"
	"skframe/pkg/rpc"
	"time"
)

type StudentManager struct {
}

func (*StudentManager) GetStudent(ctx context.Context, request *Student.StudentRequest, response *Student.Student) error {

	studentMap := map[string]Student.Student{
		"davie":  Student.Student{Name: "davie", Classes: "软件工程专业", Grade: 80},
		"steven": Student.Student{Name: "steven", Classes: "计算机科学与技术", Grade: 90},
		"tony":   Student.Student{Name: "tony", Classes: "计算机网络工程", Grade: 85},
		"jack":   Student.Student{Name: "jack", Classes: "工商管理", Grade: 96},
	}
	if request.Name == "" {
		return errors.New("请求参数错误，请重新请求。")
	}
	// 获取对应的student
	student := studentMap[request.Name]
	if student.Name != "" {
		fmt.Println(student.Name, student.Classes, student.Grade)
		*response = student
		return nil
	}
	return errors.New("未查询到相关学生信息")
}

func MicroServer() {
	rpc := rpc.Micro{}
	rpc.NewServer("127.0.0.1:8500","student.server", func(service micro.Service) {
		Student.RegisterStudentServiceHandler(service.Server(),new(StudentManager))

	})

}

func MicroClient() {
	rpc := rpc.Micro{}
	rpc.NewClient("127.0.0.1:8500", func(service micro.Service) {
		client := Student.NewStudentService("student.server",service.Client())
		test1,err := client.GetStudent(context.TODO(),&Student.StudentRequest{Name: "davie"})
		fmt.Println(err)
		fmt.Println(test1)
	})
}

func TestMicro()  {
	go func() {
		MicroServer()
		fmt.Println("xxx")
		time.Sleep(20 * time.Second)
	}()

	time.Sleep(5 *time.Second)
	MicroClient()
	time.Sleep(10 * time.Second)
	return
}
