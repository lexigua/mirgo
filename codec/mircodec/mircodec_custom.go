package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
	"reflect"
)

func init() {
	codec.RegisterCodec(new(MirUserInformationCodec))
	codec.RegisterCodec(new(MirPlayerInspectCodec))
}

/*
MirUserInformationCodec
*/

// MirUserInformationCodec ...
type MirUserInformationCodec struct{}

// Name 返回名字
func (*MirUserInformationCodec) Name() string {
	return "MirUserInformationCodec"
}

// MimeType 我也不知道是干嘛的
func (*MirUserInformationCodec) MimeType() string {
	return "application/binary"
}

// Encode 将数据转换为字节数组
func (*MirUserInformationCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	//return encode(msgObj)
	var bytes []byte
	ui := msgObj.(*server.UserInformation)
	writer := &BytesWrapper{Bytes: &bytes}
	writer.Write(ui.ObjectID)
	writer.Write(ui.RealID)
	writer.Write(ui.Name)
	writer.Write(ui.GuildName)
	writer.Write(ui.GuildRank)
	writer.Write(ui.NameColor)
	writer.Write(ui.Class)
	writer.Write(ui.Gender)
	writer.Write(ui.Level)
	writer.Write(ui.Location.X)
	writer.Write(ui.Location.Y)
	writer.Write(ui.Direction)
	writer.Write(ui.Hair)
	writer.Write(ui.HP)
	writer.Write(ui.MP)
	writer.Write(ui.Experience)
	writer.Write(ui.MaxExperience)
	writer.Write(ui.LevelEffect)

	// Inventory
	hasInventory := true
	if ui.Inventory == nil || len(ui.Inventory) == 0 {
		hasInventory = false
	}
	writer.Write(hasInventory)
	if hasInventory {
		l := len(ui.Inventory)
		//l := 46
		writer.Write(int32(l))
		for i := 0; i < l; i++ {
			hasUserItem := !IsNull(ui.Inventory[i])
			writer.Write(hasUserItem)
			if !hasUserItem {
				continue
			}
			writer.Write(&ui.Inventory[i])
		}
	}

	// Equipment
	hasEquipment := true
	if ui.Equipment == nil || len(ui.Equipment) == 0 {
		hasEquipment = false
	}
	writer.Write(hasEquipment)
	if hasEquipment {
		l := len(ui.Equipment)
		//l := 14
		writer.Write(int32(l))
		for i := 0; i < l; i++ {
			hasUserItem := !IsNull(ui.Equipment[i])
			writer.Write(hasUserItem)
			if !hasUserItem {
				continue
			}
			writer.Write(&ui.Equipment[i])
		}
	}

	// QuestInventory
	hasQuestInventory := true
	if ui.QuestInventory == nil || len(ui.QuestInventory) == 0 {
		hasQuestInventory = false
	}
	writer.Write(hasQuestInventory)
	if hasQuestInventory {
		l := len(ui.QuestInventory)
		//l := 40
		writer.Write(int32(l))
		for i := 0; i < l; i++ {
			hasUserItem := !IsNull(ui.QuestInventory[i])
			writer.Write(hasUserItem)
			if !hasUserItem {
				continue
			}
			writer.Write(&ui.QuestInventory[i])
		}
	}
	writer.Write(ui.Gold)
	writer.Write(ui.Credit)
	return *writer.Bytes, nil
}

func IsNull(ui common.UserItem) bool {
	if ui.ID == 0 && ui.ItemID == 0 {
		return true
	}
	return false
}

// Decode 将字节数组转换为数据
func (*MirUserInformationCodec) Decode(data interface{}, msgObj interface{}) error {
	ui := msgObj.(*server.UserInformation)
	bytes := data.([]byte)
	reader := &BytesWrapper{Bytes: &bytes}
	ui.ObjectID = reader.ReadUInt32()
	ui.RealID = reader.ReadUInt32()
	ui.Name = reader.ReadString()
	ui.GuildName = reader.ReadString()
	ui.GuildRank = reader.ReadString()
	ui.NameColor = reader.ReadInt32()
	ui.Class = common.MirClass(reader.ReadByte())
	ui.Gender = common.MirGender(reader.ReadByte())
	ui.Level = reader.ReadUInt16()
	x := reader.ReadInt32()
	y := reader.ReadInt32()
	ui.Location = common.Point{X: uint32(x), Y: uint32(y)}
	ui.Direction = common.MirDirection(reader.ReadByte())
	ui.Hair = reader.ReadUInt8()
	ui.HP = reader.ReadUInt16()
	ui.MP = reader.ReadUInt16()
	ui.Experience = reader.ReadInt64()
	ui.MaxExperience = reader.ReadInt64()
	ui.LevelEffect = common.LevelEffects(reader.ReadUInt8())

	// Inventory
	if reader.ReadBoolean() {
		count := reader.ReadInt32()
		ui.Inventory = make([]common.UserItem, count)
		for i := 0; i < int(count); i++ {
			if reader.ReadBoolean() {
				last := reader.Last()
				item := &ui.Inventory[i]
				*reader.Bytes = decodeValue(reflect.ValueOf(item), last)
			}
		}
	}

	// Equipment
	if reader.ReadBoolean() {
		count := reader.ReadInt32()
		ui.Equipment = make([]common.UserItem, count)
		for i := 0; i < int(count); i++ {
			if reader.ReadBoolean() {
				last := reader.Last()
				item := &ui.Equipment[i]
				*reader.Bytes = decodeValue(reflect.ValueOf(item), last)
			}
		}
	}

	// QuestInventory
	if reader.ReadBoolean() {
		count := reader.ReadInt32()
		ui.QuestInventory = make([]common.UserItem, count)
		for i := 0; i < int(count); i++ {
			if reader.ReadBoolean() {
				last := reader.Last()
				item := &ui.QuestInventory[i]
				*reader.Bytes = decodeValue(reflect.ValueOf(item), last)
			}
		}
	}
	ui.Gold = reader.ReadUInt32()
	ui.Credit = reader.ReadUInt32()
	return nil
}

/*
MirPlayerInspectCodec
*/
type MirPlayerInspectCodec struct{}

// Name 编码器的名字
func (*MirPlayerInspectCodec) Name() string {
	return "MirPlayerInspectCodec"
}

// MimeType 兼容http类型
func (*MirPlayerInspectCodec) MimeType() string {
	return "application/binary"
}

// Encode 将数据转换为字节数组
func (*MirPlayerInspectCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	var bytes []byte
	pi := msgObj.(*server.PlayerInspect)
	writer := &BytesWrapper{Bytes: &bytes}
	writer.Write(pi.Name)
	writer.Write(pi.GuildName)
	writer.Write(pi.GuildRank)
	// Equipment
	l := len(pi.Equipment)
	if l != 14 {
		panic("equipment != 14")
	}
	//l := 14
	writer.Write(int32(l))
	for i := 0; i < l; i++ {
		hasUserItem := !IsNull(pi.Equipment[i])
		writer.Write(hasUserItem)
		if !hasUserItem {
			continue
		}
		writer.Write(&pi.Equipment[i])
	}

	writer.Write(pi.Class)
	writer.Write(pi.Gender)
	writer.Write(pi.Hair)
	writer.Write(pi.Level)
	writer.Write(pi.LoverName)
	return *writer.Bytes, nil
}

// Decode 将字节数组转换为数据
func (*MirPlayerInspectCodec) Decode(data interface{}, msgObj interface{}) error {
	pi := msgObj.(*server.PlayerInspect)
	bytes := data.([]byte)
	reader := &BytesWrapper{Bytes: &bytes}
	pi.Name = reader.ReadString()
	pi.GuildName = reader.ReadString()
	pi.GuildRank = reader.ReadString()
	count := reader.ReadInt32()
	pi.Equipment = make([]common.UserItem, count)
	for i := 0; i < int(count); i++ {
		if reader.ReadBoolean() {
			last := reader.Last()
			item := &pi.Equipment[i]
			*reader.Bytes = decodeValue(reflect.ValueOf(item), last)
		}
	}
	pi.Class = common.MirClass(reader.ReadByte())
	pi.Gender = common.MirGender(reader.ReadByte())
	pi.Hair = uint8(reader.ReadByte())
	pi.Level = reader.ReadUInt16()
	pi.LoverName = reader.ReadString()
	return nil
}
