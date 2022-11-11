package make

import (
	"github.com/spf13/cobra"
)

var CmdMakeController = &cobra.Command{
	Use:   "ctl",
	Short: "Crate model file, example: make model user",
	Run:   runMakeController,
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func runMakeController(cmd *cobra.Command, args []string) {

	// 格式化模型名称，返回一个 Model 对象
	model := MakeModelFromString(args[0])
	CreateFileFromStub("app/controllers/"+model.PackageName+"_controller.go", "controller/controller", model)
	CreateFileFromStub("app/requests/"+model.PackageName+"_request.go", "controller/request", model)

}
