package models

import (
	"github.com/snail007/go-activerecord/mysql"
	"time"
)

const Table_Message_Name = "message"



type MessageTree struct {
	Message_id string
	Parent_id string
	Space_id string
	Document_id string
	Content string
	Path string
	Create_user_id string
	To_user_id string
	Is_read string
	Is_good string
	Is_bad string
	Create_time string
	Children []MessageTree
}

type Message struct {
}
var MessageModel = Message{}

// get message by message_id
func (m *Message) GetMessageByMessageId(messageId string) (message map[string]string,err error){
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Message_Name).Where(map[string]interface{}{
		"message_id": messageId,
	}))
	if err != nil {
		return
	}
	message = rs.Row()
	return
}

// insert message
func (m *Message) Insert(MessageValue map[string]interface{}) (id int64, err error){
	db := G.DB()
	MessageValue["create_time"] = time.Now().Unix()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert(Table_Message_Name, MessageValue))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return id, nil
}

func (m *Message) GetMessagesByDocumentIdAndPath(documentId,path string) (message []MessageTree,err error){
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Message_Name).Where(map[string]interface{}{
		"path": path,
		"document_id":documentId,
	}))
	if err != nil {
		return
	}
	tmp_messages := rs.Rows()
	if len(tmp_messages) == 0{
                return
        }
	tmp_messageTree := make([]MessageTree,0,1000)
	var tmp_message MessageTree
	for i:= 0;i<len(tmp_messages);i++{
		tmp_message.Message_id = tmp_messages[i]["message_id"]
		tmp_message.Parent_id = tmp_messages[i]["parent_id"]
		tmp_message.Document_id = tmp_messages[i]["document_id"]
		tmp_message.Space_id = tmp_messages[i]["space_id"]
		tmp_message.Content = tmp_messages[i]["content"]
		tmp_message.Path = tmp_messages[i]["path"]
		tmp_message.Create_user_id = tmp_messages[i]["create_user_id"]
		tmp_message.Is_read = tmp_messages[i]["is_read"]
		tmp_message.Is_good = tmp_messages[i]["is_good"]
		tmp_message.Is_bad = tmp_messages[i]["is_bad"]
		tmp_message.To_user_id = tmp_messages[i]["to_user_id"]
		tmp_message.Create_time = tmp_messages[i]["create_time"]
		tmp_messageTree = append(tmp_messageTree,tmp_message)
	}
	message = tmp_messageTree
	return
}
func (m *Message) GetMessagesTree(documentId,path string) (message []MessageTree,err error){
        if path == ""{
                path = "0"
        }
        root_message,err := m.GetMessagesByDocumentIdAndPath(documentId,path)
        if err != nil {
                return
        }
        if len(root_message) == 0{
                return
        }
        for i := 0; i < len(root_message); i++ {
                tmp_path:=root_message[i].Path+","+root_message[i].Message_id
                tmp_children,_ := m.GetMessagesTree(documentId,tmp_path)
                root_message[i].Children = tmp_children
        }
        message = root_message
        return
}
