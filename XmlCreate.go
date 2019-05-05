package GoMybatis

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

var _XmlData = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <!--logic_enable 逻辑删除字段-->
    <!--logic_deleted 逻辑删除已删除字段-->
    <!--logic_undelete 逻辑删除 未删除字段-->
    <!--version_enable 乐观锁版本字段,支持int,int8,int16,int32,int64-->
    <resultMap id="BaseResultMap" tables="#{table}">
    #{resultMapBody}
    </resultMap>
</mapper>
`
var _XmlLogicEnable = `logic_enable="true" logic_undelete="1" logic_deleted="0"`
var _XmlVersionEnable = `version_enable="true"`
var _XmlIdItem = `<id column="id" property="id"/>`
var _ResultItem = `<result column="#{property}" property="#{property}" langType="#{langType}" #{version} #{logic} />`

/**
//例子

//GoMybatis当前是以xml内容为主gm:""注解只是生成xml的时候使用
//定义数据库模型, gm:"id"表示输出id的xml,gm:"version"表示为输出版本号的xml，gm:"logic"表示输出逻辑删除xml
type TestActivity struct {
	Id         string    `json:"id" gm:"id"`
	Uuid       string    `json:"uuid"`
	Name       string    `json:"name"`
	PcLink     string    `json:"pcLink"`
	H5Link     string    `json:"h5Link"`
	Remark     string    `json:"remark"`
	Version    int       `json:"version" gm:"version"`
	CreateTime time.Time `json:"createTime"`
	DeleteFlag int       `json:"deleteFlag" gm:"logic"`
}


func TestUserAddres(t *testing.T)  {
	var s=utils.CreateDefaultXml("biz_user_address",TestActivity{})//创建xml内容
	utils.OutPutXml("D:/GOPATH/src/dao/ActivityMapper.xml",[]byte(s))//写入磁盘
}
 */
//根据结构体 创建xml文件
func CreateXml(tableName string, bean interface{}) []byte {
	var content = ""
	var tv = reflect.TypeOf(bean)
	if tv.Kind() == reflect.Ptr {
		tv = tv.Elem()
	}
	for i := 0; i < tv.NumField(); i++ {
		var item = tv.Field(i)
		if item.Name == "id" {
			content += _XmlIdItem
			content += "\n"
		} else {
			var itemName = item.Name
			var itemJson = item.Tag.Get("json")
			if itemJson != "" {
				itemName = itemJson
			}

			var itemStr = strings.Replace(_ResultItem, "#{property}", itemName, -1)
			itemStr = strings.Replace(itemStr, "#{langType}", item.Type.Name(), -1)

			var gm = item.Tag.Get("gm")
			if gm != "" {
				if gm == "id" {
					content += _XmlIdItem
					content += "\n"
				}
				if gm == "version" {
					itemStr = strings.Replace(itemStr, "#{version}", _XmlVersionEnable, -1)
				}
				if gm == "logic" {
					itemStr = strings.Replace(itemStr, "#{logic}", _XmlLogicEnable, -1)
				}
			}
			//clean
			itemStr = strings.Replace(itemStr, "#{version}", "", -1)
			itemStr = strings.Replace(itemStr, "#{logic}", "", -1)

			content += "\t" + itemStr
			if i+1 < tv.NumField() {
				content += "\n"
			}
		}
	}
	var res = strings.Replace(_XmlData, "#{resultMapBody}", content, -1)
	res = strings.Replace(res, "#{table}", tableName, -1)
	return []byte(res)
}

//输出文件
func OutPutXml(fileName string, body []byte) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f.Write(body)
		if err != nil {
			println("写入文件失败：" + err.Error())
		} else {
			println("写入文件成功：" + fileName)
		}
	}
}
