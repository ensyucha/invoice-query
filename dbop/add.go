package dbop

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"invoice/model"
)

// 新增查询结果
func AddResult(username string, results []*model.ResultItem) error {

	stmt, err := db.Prepare("INSERT INTO result_" + username + "(ensured, sealed, respCode, respMsg, qd, fpdm, fphm, " +
		"kprq, yzmSj, fpzt, fxqy, fplx, jqbm, jym, gfName, gfNsrsbh, gfAddressTel, gfBankZh, jshjL, sfName, sfNsrsbh, " +
		"sfAddressTel, sfBankZh, bz, jshjU, mxName, ggxh, unit, price, je, sl, se, totalJe, totalSe, queryTime, num) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);")

	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, item := range results {

		_, err = stmt.Exec(item.Ensured, item.Sealed, item.RespCode, item.RespMsg, item.Qd, item.Fpdm,
			item.Fphm, item.Kprq, item.YzmSj, item.Fpzt, item.Fxqy, item.Fplx, item.Jqbm, item.Jym,
			item.GfName, item.GfNsrsbh, item.GfAddressTel, item.GfBankZh, item.JshjL, item.SfName, item.SfNsrsbh,
			item.SfAddressTel, item.SfBankZh, item.Bz, item.JshjU, item.MxName, item.Ggxh, item.Unit, item.Price,
			item.Je, item.Sl, item.Se, item.TotalJe, item.TotalSe, item.QueryTime, item.Num)

		if err != nil {
			return err
		}
	}

	return nil
}

func AddDataToDB(queryArray *model.QueryArray, user *model.User) context.Map {

	tableName := "result_" + user.Username

	stmt, err := db.Prepare("INSERT INTO " + tableName + " (fpdm,fphm,kprq,je) VALUES(?,?,?,?);")

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "插入新数据失败：" + err.Error(),
		}
	}

	defer stmt.Close()

	for _, item := range queryArray.QueryArray {
		_, err = stmt.Exec(item.Fpdm, item.Fphm, item.Kprq, item.Je)

		if err != nil {
			return iris.Map{
				"status": "failed",
				"message": "插入新数据失败：" + err.Error(),
			}
		}
	}

	return iris.Map{
		"status": "success",
		"message": "插入新数据成功：",
	}
}