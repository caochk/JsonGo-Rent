package handler

import (
	"GetImageCd/utils"
	"context"
	"encoding/json"
	"github.com/afocus/captcha"
	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/logs"
	"image/color"
	"time"

	pb "GetImageCd/proto"
	"github.com/asim/go-micro/v3/util/log"
)

type GetImageCd struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetImageCd) GetImageCd(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logs.Info("---------------- GET  /api/v1.0/imagecode/:uuid GetImage() ------------------")
	// 利用第三方库创建验证码结构体实例
	code := captcha.New()
	// 设置字体（若无字体文件可能失败）
	if err := code.SetFont("comic.ttf"); err != nil {
		logs.Info("无此字体文件")
		panic(err.Error())
	}
	// 设置验证码图片大小
	code.SetSize(90, 40)
	// 设置干扰的强度
	code.SetDisturbance(captcha.MEDIUM)
	// 设置前景色
	code.SetFrontColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	// 设置背景色（可多个颜色随机替换）
	code.SetBkgColor(color.RGBA{R: 255, A: 255}, color.RGBA{B: 255, A: 255}, color.RGBA{G: 153, A: 255})
	// 生成验证码图片
	img, str := code.Create(4, captcha.NUM)

	// 记录到日志
	logs.Info(str)

	b := *img // 解引用（取地址对应内容？）
	c := *(b.RGBA)

	/* 验证码制作完成，准备返回数据 */
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	rsp.Pix = c.Pix
	rsp.Stride = int64(c.Stride)
	rsp.Max = &pb.Response_Point{X: int64(c.Rect.Max.X), Y: int64(c.Rect.Max.Y)}
	rsp.Min = &pb.Response_Point{X: int64(c.Rect.Min.X), Y: int64(c.Rect.Min.Y)}

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
	/* 将该用户与其对应的验证码信息存入缓存，30分钟内有效 */
	_ = bm.Put(context.TODO(), req.Uuid, str, time.Second*1800)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetImageCd) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.GetImageCd_StreamStream) error {
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
func (e *GetImageCd) PingPong(ctx context.Context, stream pb.GetImageCd_PingPongStream) error {
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
