package handler

import (
	"GetArea/model"
	"GetArea/utils"
	"context"
	"encoding/json"
	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"time"

	pb "GetArea/proto"
	log "github.com/micro/micro/v3/service/logger"
)

type GetArea struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetArea) GetArea(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("=========GetArea  api/v1.0/areas=========")

	// 初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	// 连接Redis。利用beego的cache模块（cache中的Redis实现好像是用的redigo库）
	// 以下应该是beego之cache模块连接Redis的固定格式
	redisConfigMap := map[string]string{
		"key":   utils.G_server_name,
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	logs.Info("连接Redis：", redisConfigMap)
	// 将map转换为json，为连接Redis做准备
	redisConfig, _ := json.Marshal(redisConfigMap)
	// 连接Redis
	bm, err := cache.NewCache("redis", string((redisConfig)))
	if err != nil {
		logs.Info("连接Redis失败：", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return err
	}

	// 先尝试从缓存中获取数据
	areasInfoValue, _ := bm.Get(context.TODO(), "areas_info")
	if areasInfoValue != nil {
		logs.Info("获取到了缓存")
		// Redis中获取到的结果需要进行json解码，准备好用于存放解码结果的变量
		var areaInfo []map[string]interface{}
		_ = json.Unmarshal(areasInfoValue.([]byte), &areaInfo)
		// 将Redis中获取结果放入rsp结构体
		for _, value := range areaInfo {
			areaAddress := pb.Response_Address{Aid: value["aid"].(int32), Aname: value["aname"].(string)}
			rsp.Data = append(rsp.Data, &areaAddress)
		}
		return nil
	}

	// 没能从Redis中获取到，那就查询数据库
	o := orm.NewOrm()
	var areas []model.Area
	resNum, err := o.QueryTable("area").All(&areas)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	if resNum == 0 {
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(utils.RECODE_NODATA)
		return nil
	}

	// 数据库中查到数据，在Redis中留一份
	logs.Info("数据库数据写入缓存")
	areaInfoStr, _ := json.Marshal(areas)
	if err := bm.Put(context.TODO(), "areas_info", areaInfoStr, time.Second*3600); err != nil {
		logs.Debug("数据库数据写入缓存出错：", err)
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(utils.RECODE_NODATA)
		return nil
	}

	// 以上异常皆未发生，返回获取到的area信息
	for _, value := range areas {
		areaAddress := pb.Response_Address{Aid: int32(value.Id), Aname: value.Name}
		rsp.Data = append(rsp.Data, &areaAddress)
	}
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetArea) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.GetArea_StreamStream) error {
	log.Infof("Received GetArea.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&pb.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetArea) PingPong(ctx context.Context, stream pb.GetArea_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&pb.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
