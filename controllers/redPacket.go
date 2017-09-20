package controllers

import (
	"encoding/json"
	"fmt"
	//"redPacket/models"
	"database/sql"
	"math"
	"math/rand"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

// Operations about Users
type RedPacketController struct {
	beego.Controller
}

type Location struct {
	access_token string
	latitude     string
	longitude    string
}

/***发红包****/
func (r *RedPacketController) SendPacket() {
	id := r.GetString("id")                                     //红包id
	enterprise_basic_id := r.GetString("enterprise_basic_id")   //公司id
	redpacket_title := r.GetString("redpacket_title")           //红包标题
	latitude := r.GetString("latitude")                         //派发中心纬度
	longitude := r.GetString("longitude")                       //派发中心经度
	radius := r.GetString("radius")                             //派发半径
	money_amount := r.GetString("money_amount")                 //红包金额
	redpacket_amount := r.GetString("redpacket_amount")         //红包数量
	distributed_location := r.GetString("distributed_location") //派发地点中心

	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()

	_, err = c.Do("HMSET", id, "enterprise_basic_id", enterprise_basic_id, "redpacket_title", redpacket_title, "latitude", latitude, "longitude", longitude, "radius", radius, "money_amount", money_amount, "redpacket_amount", redpacket_amount, "distributed_location", distributed_location)
	if err != nil {
		fmt.Println("redis set failed:", err)
	} else {
		/***将数据插入mysql**/
		db, _ := sql.Open("mysql", "root:521wangjiaxuan.@/sys?charset=utf8")
		stmt, _ := db.Prepare(`INSERT INTO gd_redbag (id,enterprise_basic_id,redpacket_title,latitude,longitude,distributed_radius,money_amount,redpacket_amount,distributed_location) values (?,?,?,?,?,?,?,?,?)`)
		res, _ := stmt.Exec(id, enterprise_basic_id, redpacket_title, latitude, longitude, radius, money_amount, redpacket_amount, distributed_location)
		fmt.Println(res)
		fmt.Println(distributed_location)
		response := make(map[string]string)
		response["code"] = "200"
		response["msg"] = "success"
		data, _ := json.Marshal(response)
		r.Data["json"] = string(data)
		r.ServeJSON()
	}

}

/****获取用户位置****/
func (r *RedPacketController) GetLocal() {

	c, err := redis.Dial("tcp", "127.0.0.1:6379") //链接redis
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()
	canGetPackets := make(map[int]map[string]string) //红包map
	/**获取请求数据**/
	user_latitude := r.GetString("latitude")
	user_longitude := r.GetString("longitude")
	redpackets, err := redis.Strings(c.Do("keys", "*")) //查询所有红包
	var i int
	i = 0
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	} else {
		for _, re := range redpackets {
			latitude, err := redis.String(c.Do("hget", re, "latitude"))
			longitude, err := redis.String(c.Do("hget", re, "longitude"))
			radius, err := redis.String(c.Do("hget", re, "radius"))
			if err != nil {
				fmt.Println("Select error", err)
				return
			} else {
				latitude, err := strconv.ParseFloat(latitude, 64)
				longitude, err := strconv.ParseFloat(longitude, 64)
				radius, err := strconv.ParseFloat(radius, 64)
				user_latitude, err := strconv.ParseFloat(user_latitude, 64)
				user_longitude, err := strconv.ParseFloat(user_longitude, 64)
				if err != nil {
					fmt.Println("转化失败")
				} else {
					check := (latitude-user_latitude)*(latitude-user_latitude) + (longitude-user_longitude)*(longitude-user_longitude)
					if check <= 1000000 {
						tempPecket := make(map[string]string)
						tempPecket["id"] = re
						tempPecket["latitude"] = strconv.FormatFloat(latitude, 'E', -1, 64)
						tempPecket["longitude"] = strconv.FormatFloat(longitude, 'E', -1, 64)
						tempPecket["distince"] = strconv.FormatFloat(math.Sqrt(check), 'E', -1, 64)
						canGetPackets[i] = tempPecket
						if check <= radius*radius {
							tempPecket["canGet"] = "0"
						} else {
							tempPecket["canGet"] = "1"
						}
						i++

					}
				}
			}
		}
	}
	data, _ := json.Marshal(canGetPackets)
	r.Data["json"] = string(data)
	r.ServeJSON()
}

/**抢红包******/
func (r *RedPacketController) Grad() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379") //链接redis
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()
	/***接收请求的用户token及红包id**/
	access_token := r.GetString("access_token")
	redpacket_id := r.GetString("id")
	/****获取红包剩余金额和数量****/
	cshare, err := redis.String(c.Do("hget", redpacket_id, "cshare"))
	cbalance, err := redis.String(c.Do("hget", redpacket_id, "cbalance"))
	if err != nil {
		fmt.Println(err)
	} else {
		/***数据类型转化****/
		cshare, _ := strconv.ParseFloat(cshare, 64)
		cbalance, _ := strconv.ParseFloat(cbalance, 64)
		/****判断剩余红包数****/
		if cshare > 0 {
			redpacket := make(map[string]float64) //判断
			if cshare != 1 {
				fmt.Println(11 - cshare)
				temp := (cbalance / cshare) * 2
				user_get := (rand.Float64() * temp) + 0.01
				cbalance = cbalance - user_get
				cshare--
				c.Do("hset", redpacket_id, "cshare", cshare)
				c.Do("hset", redpacket_id, "cbalance", cbalance)
				redpacket["money"] = user_get
				fmt.Println(cbalance)
				fmt.Println(user_get)
				/***将数据插入mysql**/
				db, _ := sql.Open("mysql", "root:521wangjiaxuan.@/sys?charset=utf8")
				stmt, _ := db.Prepare(`INSERT INTO gd_redbag_user (people_id,redpacket_id,redpacket_money) values (?,?,?)`)
				res, _ := stmt.Exec(access_token, redpacket_id, user_get)
				fmt.Println(res)
			} else {
				fmt.Println(11 - cshare)
				user_get := cbalance
				cbalance = cbalance - user_get
				cshare--
				c.Do("hset", redpacket_id, "cshare", cshare)
				c.Do("hset", redpacket_id, "cbalance", cbalance)
				redpacket["money"] = user_get
				fmt.Println(cbalance)
				fmt.Println(user_get)
				/***将数据插入mysql**/
				db, _ := sql.Open("mysql", "root:521wangjiaxuan.@/sys?charset=utf8")
				stmt, _ := db.Prepare(`INSERT INTO gd_redbag_user (people_id,redpacket_id,redpacket_money) values (?,?,?)`)
				res, _ := stmt.Exec(access_token, redpacket_id, user_get)
				fmt.Println(res)
			}
			data, _ := json.Marshal(redpacket)
			r.Data["json"] = string(data)
			r.ServeJSON()
		} else {
			c.Do("hset", redpacket_id, "cshare", 10)
			temp := (cbalance / cshare) * 2
			user_get := (rand.Float64() * temp) + 0.01
			cbalance = cbalance - user_get
			c.Do("hset", redpacket_id, "cbalance", 100)
			fmt.Println(cshare)
			fmt.Println(cbalance)
		}
	}
}
