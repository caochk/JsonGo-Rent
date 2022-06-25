package handler

import (
	"JsonGo-Rent/web/model"
	GETAREA "JsonGo-Rent/web/proto/getAreaProto"
	"context"
	"encoding/json"
	"github.com/asim/go-micro/v3"
	"github.com/beego/beego/v2/core/logs"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const ServerName = "go.micro.srv.GetArea"

func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logs.Info("---------------- 获取地区请求客户端 url : api/v1.0/areas ----------------")

	// 创建服务
	service := micro.NewService()
	// 初始化服务
	service.Init()

	client := GETAREA.NewGetAreaService(ServerName, service.Client())
	rsp, err := client.GetArea(context.Background(), &GETAREA.Request{})
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	logs.Info("rsp:", rsp)

	// 声明返回的地域数据
	var areaList []model.Area
	// 循环读取服务返回的数据
	for _, value := range rsp.Data {
		area := model.Area{Id: int(value.Aid), Name: value.Aname, Houses: nil}
		areaList = append(areaList, area)
	}
	// 创建返回数据Map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   areaList,
	}
	w.Header().Set("Content-Type", "application/json")
	// 返回给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
}
