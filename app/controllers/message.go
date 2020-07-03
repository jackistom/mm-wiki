package controllers

import (
	"regexp"
	"github.com/phachon/mm-wiki/app/models"
)



type MessageController struct {
	BaseController
}

type JSONStruct struct {
    DocumentId int
    ParentId int
    Msg  string
}
// add message
func (this *MessageController) Add() {
	documentId := this.GetString("document_id","0")
	parentId :=  this.GetString("parent_id","0")
	content := this.GetString("content","")
	var insertMessage map[string]interface{}
	if documentId == "0" {
		this.jsonError("没有选择文档！")
	}
	document, err := models.DocumentModel.GetDocumentByDocumentId(documentId)
	if err != nil {
		this.ErrorLog("添加文档评论 " + documentId + " 失败：" + err.Error())
		this.jsonError("添加评论失败！")
	}
	if len(document) == 0 {
		this.jsonError("评论文档不存在！")
	}
	spaceId := document["space_id"]
	space, err := models.SpaceModel.GetSpaceBySpaceId(spaceId)
	if err != nil {
		this.ErrorLog("获取空间信息失败：" + err.Error())
		this.jsonError("获取空间信息失败！")
	}
	if len(space) == 0 {
		this.jsonError("空间不存在！")
	}
	// check space document privilege
	_, isEditor, _ := this.GetDocumentPrivilege(space)
	if !isEditor {
		this.jsonError("您没有权限在该空间下评论！")
	}
	if content == "" {
		this.jsonError("文档名称不能为空！")
	}
	match, err := regexp.MatchString(`[\\\\/:*?\"<>、|]`, content)
	if err != nil {
		this.jsonError("评论格式不正确！")
	}
	if match {
		this.jsonError("评论称格式不正确！")
	}
	if parentId != "0"{
		parentMessage,err :=models.MessageModel.GetMessageByMessageId(parentId)
		if err != nil {
			this.ErrorLog("添加文档评论 " + documentId + " 失败：" + err.Error())
			this.ViewError("添加文档评论失败！")
		}
		if len(parentMessage) == 0 {
			this.jsonError("回复评论不存在！")
		}
		if parentMessage["space_id"] != spaceId {
			this.jsonError("回复评论不存在！")
		}
		insertMessage = map[string]interface{}{
			"parent_id":          parentMessage["message_id"],
			"document_id":        documentId,
			"space_id":           spaceId,
			"content":           content,
			"path":              parentMessage["path"]+","+parentId,
			"create_user_id":    this.UserId,
			"to_user_id":        parentMessage["create_user_id"],
		}
	} else {
		insertMessage = map[string]interface{}{
			"parent_id":          0,
			"document_id":        documentId,
			"space_id":           spaceId,
			"content":           content,
			"path":              0,
			"create_user_id":    this.UserId,
			"to_user_id":        document["create_user_id"],
		}
	}
	_, err = models.MessageModel.Insert(insertMessage)
	if err != nil {
		this.ErrorLog("评论失败：" + err.Error())
		this.jsonError("评论文档失败")
	}
	this.InfoLog(this.UserId +"评论 " + documentId + " 文档 " + " 成功")
	this.jsonSuccess("评论成功！", nil)
}

func (this *MessageController) GetMessage() {
	documentId := this.GetString("document_id","0")
	if documentId == "0" {
		this.jsonError("没有选择文档！")
	}
	document, err := models.DocumentModel.GetDocumentByDocumentId(documentId)
	if err != nil {
		this.ErrorLog("获取 " + documentId + " 评论失败：" + err.Error())
		this.jsonError("获取文档信息失败！")
	}
	if len(document) == 0 {
		this.jsonError("评论文档不存在！")
	}
	spaceId := document["space_id"]
	space, err := models.SpaceModel.GetSpaceBySpaceId(spaceId)
	if err != nil {
		this.ErrorLog("获取空间信息失败：" + err.Error())
		this.jsonError("获取空间信息失败！")
	}
	if len(space) == 0 {
		this.jsonError("空间不存在！")
	}
	// check space document privilege
	isVisit, _, _ := this.GetDocumentPrivilege(space)
	if !isVisit {
		this.jsonError("您没有权限访问该空间下评论！")
	}
	messageTree,err := models.MessageModel.GetMessagesTree(documentId,"0")
	if err != nil {
		this.jsonError(err.Error)
	}
	this.Data["messageTree"] = messageTree
	this.viewLayout("message/messagetree", "default")
//	this.jsonSuccess(messageTree, nil)

}



