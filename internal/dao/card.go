package dao

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"

	"github.com/88250/lute"
	"github.com/LinkinStars/dc/internal/base/db"
	"github.com/LinkinStars/dc/internal/model"
	"github.com/LinkinStars/dc/internal/val"
	"github.com/LinkinStars/go-scaffold/logger"
	"github.com/LinkinStars/go-scaffold/mistake"
)

func AddCard(content string) error {
	card := &model.Card{CreatedAt: time.Now(), OriginalText: content, ParsedText: Markdown2HTML(content)}
	_, err := db.Engine.InsertOne(card)
	if err != nil {
		return mistake.InternalServer("500", err.Error())
	}
	return nil
}

func UpdateCard(req *val.UpdateCardReq) error {
	card := &model.Card{
		OriginalText: req.Content,
		ParsedText:   Markdown2HTML(req.Content),
	}

	if len(req.CreatedAt) > 0 {
		date, err := time.Parse("2006-01-02", req.CreatedAt)
		if err == nil {
			card.CreatedAt = date
		} else {
			logger.Error(fmt.Sprintf("parse date error: %s", err.Error()))
		}
	}

	_, err := db.Engine.ID(req.ID).Update(card)
	if err != nil {
		return mistake.InternalServer("500", err.Error())
	}
	return nil
}

func AddCardPv(id int, pv int) error {
	card := &model.Card{}
	_, err := db.Engine.ID(id).Incr("pv", pv).Update(card)
	if err != nil {
		return mistake.InternalServer("500", err.Error())
	}
	return nil
}

func DeleteCard(id int) error {
	_, err := db.Engine.ID(id).Delete(&model.Card{})
	if err != nil {
		return mistake.InternalServer("500", err.Error())
	}
	return nil
}

func GetCardDetail(id int) (card *model.Card, err error) {
	card = &model.Card{}
	session := db.Engine.Desc("id")
	if id > 0 {
		session.ID(id)
	}
	_, err = session.Get(card)
	if err != nil {
		return nil, mistake.InternalServer("500", err.Error())
	}
	return card, nil
}

func GetFirstCard() (card *model.Card, exist bool, err error) {
	card = &model.Card{}
	exist, err = db.Engine.Asc("id").Get(card)
	if err != nil {
		return nil, false, mistake.InternalServer("500", err.Error())
	}
	return card, exist, nil
}

func GetCardDetailByOffset(id, offset int) (card *model.Card, err error) {
	cards := make([]*model.Card, 0)
	session := db.Engine.NewSession()
	if offset > 0 {
		session.Asc("id").Where("id > ?", id)
	} else {
		session.Desc("id").Where("id < ?", id)
	}
	err = session.Limit(1).Find(&cards)
	if err != nil {
		return nil, mistake.InternalServer("500", err.Error())
	}
	if len(cards) > 0 {
		return cards[0], nil
	}
	return nil, nil
}

func GetCards(page, pageSize int, keyword, date string) (cards []*model.Card, count int64, err error) {
	cards = make([]*model.Card, 0)
	startNum := (page - 1) * pageSize
	session := db.Engine.Desc("created_at")
	if len(keyword) > 0 {
		session.Where("original_text LIKE ?", "%"+keyword+"%")
	}
	if len(date) > 0 {
		dateTime, err := time.Parse(time.DateOnly, date)
		if err == nil {
			cur := now.New(dateTime)
			start := cur.BeginningOfDay()
			end := cur.EndOfDay()
			session.Where("created_at >= ?", start.Format(time.DateTime))
			session.Where("created_at <= ?", end.Format(time.DateTime))
		}
	}
	count, err = session.Limit(pageSize, startNum).FindAndCount(&cards)
	if err != nil {
		return nil, 0, mistake.InternalServer("500", err.Error())
	}
	return cards, count, nil
}

func GetCardsByTime(startTime, endTime time.Time) (cards []*model.Card, err error) {
	cards = make([]*model.Card, 0)
	session := db.Engine.Where("created_at >= ?", startTime)
	session.Where("created_at <= ?", endTime)
	err = session.Find(&cards)
	if err != nil {
		return nil, mistake.InternalServer("500", err.Error())
	}
	return cards, nil
}

func GetCardsRecord(startTime, endTime time.Time) (recordTime []string, err error) {
	recordTime = make([]string, 0)
	session := db.Engine.Table("card")
	session.Select("created_at")
	session.Where("deleted_at IS NULL")
	session.Where("created_at >= ?", startTime)
	session.Where("created_at <= ?", endTime)
	err = session.Find(&recordTime)
	if err != nil {
		return nil, mistake.InternalServer("500", err.Error())
	}
	return recordTime, nil
}

func CountCards() (count int64, err error) {
	count, err = db.Engine.Count(&model.Card{})
	if err != nil {
		return 0, mistake.InternalServer("500", err.Error())
	}
	return count, nil
}

func Markdown2HTML(source string) string {
	luteEngine := lute.New() // 默认已经启用 GFM 支持以及中文语境优化
	return luteEngine.MarkdownStr("demo", source)
	//mdConverter := goldmark.New(
	//	goldmark.WithExtensions(extension.GFM),
	//	goldmark.WithParserOptions(
	//		parser.WithAutoHeadingID(),
	//	),
	//	goldmark.WithRendererOptions(
	//		html.WithHardWraps(),
	//		html.WithUnsafe(),
	//	),
	//)
	//var buf bytes.Buffer
	//if err := mdConverter.Convert([]byte(source), &buf); err != nil {
	//	logger.Error(err)
	//	return source
	//}
	//return buf.String()
}
