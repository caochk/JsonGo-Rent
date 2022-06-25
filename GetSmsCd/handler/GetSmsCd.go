package handler

import (
	"GetSmsCd/model"
	"GetSmsCd/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"strconv"
	"strings"
	"time"

	pb "GetSmsCd/proto"
)

type GetSmsCd struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetSmsCd) GetSmsCd(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- GET smscd  api/v1.0/smscode/:id ----------------")

	// 通过手机号码进行注册时发现该用户已经注册
	o := orm.NewOrm()
	var user = model.User{Mobile: req.Mobile}
	if err := o.Read(&user); err != nil {
		logs.Info("该手机号已注册")
		logs.Info(err.Error())
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	/* 连接Redis。利用beego的cache模块（cache中的Redis实现好像是用的redigo库）*/
	// 以下是beego之cache模块连接Redis的固定格式
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
	// 获取Redis中存储的图片验证码
	imageCode, _ := bm.Get(context.TODO(), req.Id)
	if imageCode == nil {
		logs.Info("从Redis中获取图片验证码失败")
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(utils.RECODE_NODATA)
		return nil
	}
	// 获取到了图片验证码
	imageCodeString, _ := redis.String(imageCode, nil)
	// 验证用户输入图片验证码是否正确
	if req.Text != imageCodeString {
		logs.Info("用户输入图片验证码错误")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	/*生成随机数*/
	t := rand.New(rand.NewSource(time.Now().UnixNano()))
	sms := t.Intn(8999) + 1000 // 这样会创建出1000~9999的随机数
	//预先创建好的appid
	messageconfig["appid"] = "29672"
	//预先获得的app的key
	messageconfig["appkey"] = "89d90165cbea8cae80137d7584179bdb"
	//加密方式默认
	messageconfig["signtype"] = "md5"

	//创建短信发送的句柄（submail是一个云通信工具，已经无法使用，得换方案了！）
	messagexsend := submail.CreateMessageXSend()
	//设置发送短信的手机号
	submail.MessageXSendAddTo(messagexsend, req.Mobile)
	//短信主题
	submail.MessageXSendSetProject(messagexsend, "rent sms")
	//验证码
	submail.MessageXSendAddVar(messagexsend, "sms:", strconv.Itoa(sms))
	//发送短信
	send := submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig)

	//对短信的发送的验证码进行校验
	bo := strings.Contains(send, "success")
	if bo != true {
		fmt.Println("短信验证码发送失败")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//将sms和手机号存入redis
	err = bm.Put(context.TODO(), req.Mobile, strconv.Itoa(sms), time.Second*300)
	if err != nil {
		fmt.Println("sms存储失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetSmsCd) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.GetSmsCd_StreamStream) error {
	log.Infof("Received pb.Stream request with count: %d", req.Count)

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
func (e *GetSmsCd) PingPong(ctx context.Context, stream pb.GetSmsCd_PingPongStream) error {
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
