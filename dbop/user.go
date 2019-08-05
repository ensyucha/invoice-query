package dbop

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"invoice/auth"
	"invoice/model"
	"log"

)

// 新增用户
func AddUser(user *model.User) context.Map {

	if user.Username == "admin" {
		return iris.Map{
			"status": "failed",
			"message": "新增用户失败: 不允许用户名为 admin",
		}
	}

	if user.Usage < 0 {
		return iris.Map{
			"status": "failed",
			"message": "预分配额度不能少于0",
		}
	}

	newUnusedUsage, ok, err := UCEnough(user.Usage)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "新增用户失败: " + err.Error(),
		}
	}

	// 额度判断能否分配
	if !ok { // 新增用户预分配额度大于闲置额度，无法创建用户
		return iris.Map{
			"status": "failed",
			"message": "新增用户失败: 预分配额度大于闲置额度",
		}
	}

	stmt, err := db.Prepare("INSERT INTO users(username, nickname, password, usages) VALUES (?,?,?,?);")

	if err != nil {

		return iris.Map{
			"status":  "failed",
			"message": "新增用户失败: " + err.Error(),
		}
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.NickName, user.Password, user.Usage)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "新增用户失败" + err.Error(),
		}
	}

	err = UCUpdateUnusedUsage(&model.System{UnusedUsage:newUnusedUsage})

	if err != nil {
		return iris.Map{
			"status": "ok",
			"message": "新增用户失败：" + err.Error(),
		}
	}

	err = createResultTable(user.Username)

	if err != nil {
		err2 := removeResultTable(user.Username)

		if err2 != nil {
			return iris.Map{
				"status":  "failed",
				"message": "新增用户失败: 无法创建结果表：" + err.Error() + err2.Error(),
			}
		} else {
			return iris.Map{
				"status":  "failed",
				"message": "新增用户失败: 无法创建结果表：" + err.Error(),
			}
		}
	}

	WriteLog("system", "新增用户：" + user.Username, "manager")

	return iris.Map{
		"status": "ok",
		"message": "新增用户成功",
	}
}

// 删除用户
func RemoveUser(user *model.User) context.Map {

	remainUsage, err := UCGetUserUsage(user)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "删除用户失败: " + err.Error(),
		}
	}

	nowUnusedUsage, err := UCGetUnusedUsage()

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "删除用户失败: " + err.Error(),
		}
	}

	err = UCUpdateUnusedUsage(&model.System{UnusedUsage:remainUsage + nowUnusedUsage})

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "删除用户失败: " + err.Error(),
		}
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE username=?;")

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "删除用户失败: " + err.Error(),
		}
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Username)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "删除用户失败" + err.Error(),
		}
	}

	auth.RemoveTokenByUsername(user.Username) // 删除对应的token

	err = removeResultTable(user.Username)

	if err != nil {
		return iris.Map{
			"status": "ok",
			"message": "删除用户数据表失败：" + err.Error(),
		}
	}

	WriteLog("system", "删除用户：" + user.Username, "manager")

	return iris.Map{
		"status": "ok",
		"message": "删除用户成功",
	}
}

// 更新用户
func UpdateUser(user *model.User) context.Map {

	if user.Usage < 0 {
		return iris.Map{
			"status": "failed",
			"message": "分配额度不能少于0",
		}
	}

	nowUserUsage, err := UCGetUserUsage(user)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "获取用户原始额度失败:" + err.Error(),
		}
	}

	diffUsage := user.Usage - nowUserUsage // 获取额度差值

	newUnusedUsage, ok, err := UCEnough(diffUsage)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "计算额度差值失败:" + err.Error(),
		}
	}

	if !ok {
		return iris.Map{
			"status": "failed",
			"message": "更新用户失败: 预分配额度大于未分配额度",
		}
	}

	stmt, err := db.Prepare("UPDATE users SET nickname=?, password=?, usages=? WHERE username=?;")

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "更新用户失败: " + err.Error(),
		}
	}

	defer stmt.Close()

	log.Println(user.NickName, user.Password, user.Usage, user.Username)

	_, err = stmt.Exec(user.NickName, user.Password, user.Usage, user.Username)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "更新用户失败" + err.Error(),
		}
	}

	err = UCUpdateUnusedUsage(&model.System{UnusedUsage:newUnusedUsage})

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "更新未分配额度失败:" + err.Error(),
		}
	}

	auth.RemoveTokenByUsername(user.Username) // 删除对应的token

	WriteLog("system", "更新用户：" + user.Username, "manager")

	return iris.Map{
		"status": "ok",
		"message": "更新用户成功",
	}
}

// 获取用户列表
func ListUser() context.Map {

	var group []model.User

	userResult, err := db.Query("SELECT * FROM users")

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "获取用户列表失败: " + err.Error(),
		}
	}

	defer userResult.Close()

	for userResult.Next() {

		var userItem model.User

		err = userResult.Scan(&userItem.UserId, &userItem.Username, &userItem.NickName, &userItem.Password,
			&userItem.Usage, &userItem.Total)

		if err != nil {
			return iris.Map{
				"status": "failed",
				"message": "获取用户列表失败: " + err.Error(),
			}
		}

		group = append(group, userItem)
	}

	return iris.Map{
		"status": "ok",
		"message": group,
	}
}

func CheckUser(user *model.User) (bool, string, error) {

	var systemItem model.System
	var userItem model.User
	var querySQL string

	if user.Username == "admin" {
		querySQL = "SELECT password FROM systems;"
	} else {
		querySQL = "SELECT nickname, password FROM users WHERE username='" + user.Username + "';"
	}

	queryResult, err := db.Query(querySQL)

	if err != nil {
		return false, "", err
	}

	defer queryResult.Close()

	if queryResult.Next() {

		if user.Username == "admin" {
			err = queryResult.Scan(&systemItem.Password)
		} else {
			err = queryResult.Scan(&userItem.NickName, &userItem.Password)
		}

		if err != nil {
			return false, "", err
		}
	}

	if (user.Username == "admin" && user.Password != systemItem.Password) ||
		(user.Username != "admin" && user.Password != userItem.Password) {
		return false, "", nil
	}

	return true, userItem.NickName, nil
}

func createResultTable(username string) error {

	sql := `CREATE TABLE IF NOT EXISTS result_` + username + ` (
	resultid INTEGER PRIMARY KEY AUTO_INCREMENT COMMENT '查询结果id',
	ensured varchar(300) DEFAULT '' COMMENT '确认状态',
	sealed char(1) DEFAULT '0' COMMENT '封存状态',
	respCode char(20) DEFAULT '' COMMENT '查询结果代号',
	respMsg char(100) DEFAULT '<span class="my-not-check">未查询</span>' COMMENT '查询结果',
	qd char(1) DEFAULT '' COMMENT '是否有清单',
	fpdm char(50) DEFAULT '' COMMENT '发票代码',
	fphm char(50) DEFAULT '' COMMENT '发票号码',
	kprq char(50) DEFAULT '' COMMENT '开票日期',
	yzmSj char(50) DEFAULT '' COMMENT '验证时间',
	fpzt char(50) DEFAULT '' COMMENT '发票状态',
	fxqy char(50) DEFAULT '' COMMENT '风险企业验证',
	fplx char(50) DEFAULT '' COMMENT '发票类型',
	jqbm char(100) DEFAULT '' COMMENT '机器编码',
	jym char(100) DEFAULT '' COMMENT '校验码',
	gfName varchar(300) DEFAULT '' COMMENT '供应方名称',
	gfNsrsbh varchar(300) DEFAULT '' COMMENT '供应方识别号',
	gfAddressTel varchar(300) DEFAULT '' COMMENT '供应方联系方式',
	gfBankZh varchar(300) DEFAULT '' COMMENT '供应方开户行',
	jshjL varchar(300) DEFAULT '' COMMENT '价税合计',
	sfName varchar(300) DEFAULT '' COMMENT '销售方名称',
	sfNsrsbh varchar(300) DEFAULT '' COMMENT '销售方识别号',
	sfAddressTel varchar(300) DEFAULT '' COMMENT '销售方联系方式',
	sfBankZh varchar(300) DEFAULT '' COMMENT '销售方开户行',
	bz varchar(300) DEFAULT '' COMMENT '备注信息',
	jshjU varchar(300) DEFAULT '' COMMENT '价税合计(大写)',
	mxName varchar(300) DEFAULT '' COMMENT '商品名',
	ggxh varchar(300) DEFAULT '' COMMENT '规格型号',
	unit char(100) DEFAULT '' COMMENT '单位',
	price varchar(300) DEFAULT '' COMMENT '单价',
	je varchar(300) DEFAULT '' COMMENT '金额',
	sl char(50) DEFAULT '' COMMENT '税率',
	se varchar(300) DEFAULT '' COMMENT '税额',
	totalJe varchar(300) DEFAULT '-1' COMMENT '总金额',
	totalSe varchar(300) DEFAULT '-1' COMMENT '总税额',
	queryTime char(100) DEFAULT '' COMMENT '查询时间',
	num varchar(300) DEFAULT '' COMMENT '数量'
);`

	_, err := db.Exec(sql)

	return err
}

func removeResultTable(username string) error {

	sql := `DROP TABLE result_` + username + `;`

	_, err := db.Exec(sql)

	return err
}