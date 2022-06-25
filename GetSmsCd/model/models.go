package model

import (
	"GetSmsCd/utils"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

// User 用户表 表名：user
type User struct {
	Id            int           `json:"user_id"`                      // 用户ID
	Name          string        `orm:"size(32)" json:"name"`          // 用户名
	Password_hash string        `orm:"size(128)" json:"password"`     // 用户密码密文
	Mobile        string        `orm:"size(11);unique" json:"mobile"` // 手机号
	Real_name     string        `orm:"size(32)" json:"real_name"`     // 实名姓名
	Id_card       string        `orm:"size(20)" json:"id_card"`       // 实名身份证
	Avatar_url    string        `orm:"size(256)" json:"avatar_url"`   // 用户头像路径，图片通过fastDFS进行存储
	Houses        []*House      `orm:"reverse(many)" json:"houses"`   // 用户发布的房源信息，一对多关系反向，一人可发布多套房
	Orders        []*OrderHouse `orm:"reverse(many)" json:"orders"`   // 用户产生的订单信息，一对多关系反向，一人可下多单
}

// House 房屋信息表 表名：house
type House struct {
	Id              int           `json:"house_id"`                        // 房屋编号
	User            *User         `orm:"rel(fk)" json:"user_id"`           // 该房屋所属主人的ID，一对多关系
	Area            *Area         `orm:"rel(fk)" json:"area_id"`           // 该房屋归属地的区域ID，一对多关系
	Title           string        `orm:"size(64)" json:"title"`            // 房屋信息的标题
	Price           int           `orm:"default(0)" json:"price"`          // 房屋单价,单位:分   每次的价格要乘以100
	Address         string        `orm:"size(512)" json:"address"`         // 地址
	Room_count      int           `orm:"default(1)" json:"room_count"`     // 房间数目
	Acreage         int           `orm:"default(0)" json:"acreage"`        // 房屋总面积
	Unit            string        `orm:"size(32)" json:"unit"`             // 房屋单元,如 几室几厅
	Capacity        int           `orm:"default(1)" json:"capacity"`       // 房屋容纳的总人数
	Beds            string        `orm:"size(64)" json:"beds"`             // 房屋床铺配置
	Deposit         int           `orm:"default(0)" json:"deposit"`        // 押金
	Min_days        int           `orm:"default(1)" json:"min_days"`       // 最少入住的天数
	Max_days        int           `orm:"default(0)" json:"max_days"`       // 最多入住的天数 0表示不限制
	Order_count     int           `orm:"default(0)" json:"order_count"`    // 已预订完成的该房屋的订单数
	Index_image_url string        `orm:"size(256)" json:"index_image_url"` // 房屋主图片路径
	Facilities      []*Facility   `orm:"reverse(many)" json:"facilities"`  // 房屋设施，多对多关系反向
	Images          []*HouseImage `orm:"reverse(many)" json:"img_urls"`    // 房屋次要图片，一对多关系反向
	Orders          []*OrderHouse `orm:"reverse(many)" json:"orders"`      // 该房屋对应订单，一房屋被多次成交遂有多订单
	Ctime           time.Time     `orm:"auto_now_add;type(datetime)" json:"ctime"`
}

// Area 区域信息表 表名：area (区域信息是需要我们手动添加到数据库中)
type Area struct {
	Id     int      `json:"aid"`                        // 区域ID
	Name   string   `orm:"size(32)" json:"aname"`       // 区域名字
	Houses []*House `orm:"reverse(many)" json:"houses"` // 区域内所有房屋，一对多关系反向
}

// Facility 设施信息表 表名：facility （设施信息需提前手动添加）
type Facility struct {
	Id     int      `json:"fid"`     // 设施ID
	Name   string   `orm:"size(32)"` // 设施名字
	Houses []*House `orm:"rel(m2m)"` // 都有哪些房屋有此设施，多对多关系
}

// HouseImage 房屋次要图片表 表名：house_image
type HouseImage struct {
	Id    int    `json:"house_image_id"`         // 次要图片ID
	Url   string `orm:"size(256)" json:"url"`    // 次要图片url
	House *House `orm:"rel(fk)" json:"house_id"` // 次要图片所属房屋ID，一对多关系，一屋可有多张次要图片
}

// OrderHouse 订单表 表名：order
type OrderHouse struct {
	Id          int       `json:"order_id"`               // 订单ID
	User        *User     `orm:"rel(fk)" json:"user_id"`  // 下单用户的ID，一对多
	House       *House    `orm:"rel(fk)" json:"house_id"` // 预订房屋的ID
	Begin_date  time.Time `orm:"type(datetime)"`          // 预订的起始时间
	End_date    time.Time `orm:"type(datetime)"`          // 预订的结束时间
	Days        int       // 预订的总天数
	House_price int       // 房屋单价？？有点怪，感觉通过房屋ID即可得房屋价格
	Amount      int       // 订单总金额
	Status      string    `orm:"default(WAIT_ACCEPT)"`                 // 订单状态
	Comment     string    `orm:"size(512)"`                            // 订单评论
	Ctime       time.Time `orm:"auto_now;type(datetime)" json:"ctime"` // 每次更新此表，都会更新这个字段
	Credit      bool      // 表示个人征信情况 true表示良好
}

func init() {
	// 注册数据库驱动
	// mysql / sqlite3 / postgres / tidb 这几种是默认已经注册过的，所以可以无需设置
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	// 注册数据库
	var dsn = utils.G_mysql_user + ":" + utils.G_mysql_pass + "@tcp(" + utils.G_mysql_addr + ":" + utils.G_mysql_port + ")/" + utils.G_mysql_dbname + "?charset=utf8"
	err := orm.RegisterDataBase("default", "mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}
	// 注册模型（若要利用beego进行高级sql查询，此步为必须步骤）
	orm.RegisterModel(new(User), new(House), new(Area), new(Facility), new(HouseImage), new(OrderHouse))
	// 建表？
	if err := orm.RunSyncdb("default", false, true); err != nil {
		logs.Debug("[ERROR] sync db:", err)
	}
}
