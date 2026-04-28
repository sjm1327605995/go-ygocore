package go_ygocore

import "github.com/sjm1327605995/go-ygocore/core"

// Duel 封装了决斗实例的句柄，对外隐藏底层 uintptr
type Duel struct {
	pduel  uintptr
	client *Client
}

// Client 对外提供的高层封装，隐藏了 Library 和 uintptr 细节
type Client struct {
	lib *core.Library
}

// NewClient 创建一个新的 Client，加载动态库并初始化
func NewClient(opts core.LibraryOptions) (*Client, error) {
	lib, err := core.LoadLibrary(opts)
	if err != nil {
		return nil, err
	}
	return &Client{lib: lib}, nil
}

// Close 关闭 Client，释放动态库资源
func (c *Client) Close() error {
	if c.lib != nil {
		return c.lib.Close()
	}
	return nil
}

// CreateDuel 创建决斗实例（使用单个 seed，已废弃，仅用于回放模式）
func (c *Client) CreateDuel(seed int32) *Duel {
	pduel := c.lib.CreateDuel(seed)
	return &Duel{pduel: pduel, client: c}
}

// CreateDuelV2 创建决斗实例（使用 seed 数组）
func (c *Client) CreateDuelV2(seedSequence [core.SEED_COUNT]uint32) *Duel {
	pduel := c.lib.CreateDuelV2(seedSequence)
	return &Duel{pduel: pduel, client: c}
}

// StartDuel 开始决斗
func (d *Duel) StartDuel(options uint32) {
	d.client.lib.StartDuel(d.pduel, options)
}

// EndDuel 结束决斗
func (d *Duel) EndDuel() {
	d.client.lib.EndDuel(d.pduel)
}

// SetPlayerInfo 设置玩家信息
func (d *Duel) SetPlayerInfo(playerid int32, lp int32, startcount int32, drawcount int32) {
	d.client.lib.SetPlayerInfo(d.pduel, playerid, lp, startcount, drawcount)
}

// GetLogMessage 获取日志消息
func (d *Duel) GetLogMessage(buf []byte) {
	d.client.lib.GetLogMessage(d.pduel, buf)
}

// GetMessage 获取消息
func (d *Duel) GetMessage(buf []byte) int32 {
	return d.client.lib.GetMessage(d.pduel, buf)
}

// Process 执行一个游戏 tick
func (d *Duel) Process() int32 {
	return d.client.lib.Process(d.pduel)
}

// NewCard 添加卡片到决斗状态
func (d *Duel) NewCard(code uint32, owner uint8, playerid uint8, location uint8, sequence uint8, position uint8) {
	d.client.lib.NewCard(d.pduel, code, owner, playerid, location, sequence, position)
}

// NewTagCard 添加卡片到 tag 池
func (d *Duel) NewTagCard(code uint32, owner uint8, location uint8) {
	d.client.lib.NewTagCard(d.pduel, code, owner, location)
}

// QueryCard 查询特定位置的卡片
func (d *Duel) QueryCard(playerid uint8, location uint8, sequence uint8, queryFlag uint32, buf []byte, useCache int32) int32 {
	return d.client.lib.QueryCard(d.pduel, playerid, location, sequence, queryFlag, buf, useCache)
}

// QueryFieldCount 查询特定区域卡片数量
func (d *Duel) QueryFieldCount(playerid uint8, location uint8) int32 {
	return d.client.lib.QueryFieldCount(d.pduel, playerid, location)
}

// QueryFieldCard 查询特定区域所有卡片
func (d *Duel) QueryFieldCard(playerid uint8, location uint8, queryFlag uint32, buf []byte, useCache int32) int32 {
	return d.client.lib.QueryFieldCard(d.pduel, playerid, location, queryFlag, buf, useCache)
}

// QueryFieldInfo 查询场地信息
func (d *Duel) QueryFieldInfo(buf []byte) int32 {
	return d.client.lib.QueryFieldInfo(d.pduel, buf)
}

// SetResponseI 设置整数响应
func (d *Duel) SetResponseI(value int32) {
	d.client.lib.SetResponseI(d.pduel, value)
}

// SetResponseB 设置字节数组响应
func (d *Duel) SetResponseB(buf []byte) {
	d.client.lib.SetResponseB(d.pduel, buf)
}

// PreloadScript 预加载脚本
func (d *Duel) PreloadScript(script string) int32 {
	return d.client.lib.PreloadScript(d.pduel, script)
}
