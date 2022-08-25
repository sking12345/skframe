package example

import (
	"context"
	"errors"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"skframe/app/protobuffs/Student"
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

func microServer() {
	consulReg := consul.NewRegistry( //新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
		registry.Addrs("127.0.0.1:8500"),
	)
	server := micro.NewService( //创建一个新的服务对象
		micro.Name("student"),
		micro.Version("v1.0.0"),
		micro.Registry(consulReg),
		micro.RegisterTTL(10*time.Second),
		micro.RegisterInterval(5*time.Second),
	)

	server.Init()
	Student.RegisterStudentServiceHandler(server.Server(), new(StudentManager))
	err := server.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func microClient() {
	consulReg := consul.NewRegistry( //新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
		registry.Addrs("127.0.0.1:8500"),
	)
	service := micro.NewService(
		micro.Name("student.client"),
		micro.Registry(consulReg),
	)
	service.Init()

	studentService := Student.NewStudentServiceClient("student_service", service.Client())
	res, err := studentService.GetStudent(context.TODO(), &Student.StudentRequest{Name: "davie"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Name)
	fmt.Println(res.Classes)
	fmt.Println(res.Grade)
	time.Sleep(50 * time.Second)

}
