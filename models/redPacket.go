package models

type RedPacket struct {
	id                       string
	enterprise_basic_id      string
	publisher_department     string //发布者所属部门
	publisher_nickname       string //发布者昵称
	redpacket_titie          string //红包标题
	redpacket_cover_plan     string //红包封面图
	redpacket_summary        string //红包摘要
	people_age_max           string //抢红包人年龄最大值',
	people_age_min           string //抢红包人年龄最小值',
	people_sex               string //抢红包人性别  1：男；0：女',
	people_edication         string //'抢红包人教育背景。1：小学及以下；2：初中；3：高中；4：学士学位；5：硕士学位；6：博士及以上。',
	people_skill             string //'抢红包人技能关键字',
	people_experience        string // '抢红包人经历关键字',
	people_profession_status string // '抢红包人职业状态。1：无业；2：在职；3：离职；4：退休。',
	money_amount             string // '红包总金额',
	redpacket_amount         string // '红包总数',
	redpacket_money_max      string // '单个红包金额最大值',
	redpacket_money_min      string // '单个红包金额最小值',
	distributed_location     string // '派发地点中心',
	distributed_radius       string //'派发半径',
	distributed_start_time   string // '派发开始时间',
	distributed_end_time     string // '派发结束时间',
	outside_chain            string // '外链',
	bottom_advert            string //底部广告
	five_minutes_advert      string //5秒广告
	publish_status           string //发布状态。1：未发布；2：进行中；3：已发布；4：失效
	create_time              string //创建时间
	update_time              string //更新时间
}
