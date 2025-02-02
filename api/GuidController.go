package api

import (
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/pkg/engine/db"
	"github.com/1340691923/ElasticView/pkg/jwt"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/util"
	. "github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

// 引导控制器
type GuidController struct {
	BaseController
}

// 完成新手引导
func (this GuidController) Finish(ctx *Ctx) error {
	c, err := jwt.ParseToken(ctx.Get("X-Token"))
	if err != nil {
		return this.Error(ctx, err)
	}

	gmGuidModel := model.GmGuidModel{}
	err = ctx.BodyParser(&gmGuidModel)
	if err != nil {
		return this.Error(ctx, err)
	}
	_, err = db.SqlBuilder.
		Insert(gmGuidModel.TableName()).
		SetMap(util.Map{
			"uid":       c.ID,
			"guid_name": gmGuidModel.GuidName,
			"created":   time.Now().Format(util.TimeFormat),
		}).RunWith(db.Sqlx).Exec()

	if err != nil && (strings.Contains(err.Error(), "Error 1062") || strings.Contains(err.Error(), "UNIQUE")) {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

//是否完成新手引导
func (this GuidController) IsFinish(ctx *Ctx) error {
	c, err := jwt.ParseToken(ctx.Get("X-Token"))
	if err != nil {
		return this.Error(ctx, err)
	}

	gmGuidModel := model.GmGuidModel{}
	err = ctx.BodyParser(&gmGuidModel)
	if err != nil {
		return this.Error(ctx, err)
	}
	sql, args, err := db.SqlBuilder.
		Select("count(*)").
		From(gmGuidModel.TableName()).
		Where(db.Eq{
			"uid":       c.ID,
			"guid_name": gmGuidModel.GuidName,
		}).ToSql()

	if err != nil {
		return this.Error(ctx, err)
	}
	var count int
	err = db.Sqlx.Get(&count, sql, args...)
	if util.FilterMysqlNilErr(err) {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, count > 0)
}
